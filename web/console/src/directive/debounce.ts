// 定义类型
interface VNodeProps {
  onClick?: (...args: unknown[]) => unknown
  _onClick?: (...args: unknown[]) => unknown
  loading?: boolean
  [key: string]: unknown
}

interface VNode {
  ctx: {
    type: { name: string }
    vnode: {
      props: VNodeProps
      loading?: boolean
      key?: string
    }
    props: VNodeProps
    proxy?: {
      $forceUpdate?: () => void
    }
  }
}

interface HTMLElementWithListener extends HTMLElement {
  dListener?: (ev: Event) => void
  disabled?: boolean
}

/**
 * AOP函数 - 在原函数执行前后添加额外逻辑
 */
function AOP(
  func: (...args: unknown[]) => unknown,
  beforeFn: (...args: unknown[]) => void,
  afterFn: (...args: unknown[]) => void
) {
  return function (this: unknown, ...args: unknown[]) {
    beforeFn.apply(this, args)
    const ret = func.apply(this, args)

    if (ret && typeof ret === 'object' && 'then' in ret && 'catch' in ret) {
      // 处理Promise返回值
      ;(ret as Promise<unknown>).finally(() => {
        afterFn.apply(this, args)
      })
    } else {
      // 非Promise返回值，使用setTimeout模拟异步
      setTimeout(() => {
        afterFn.apply(this, args)
      }, 1000)
    }

    return ret
  }
}

/**
 * 创建函数副本，避免直接修改原函数
 */
function createFunctionCopy(originalFn: (...args: unknown[]) => unknown) {
  return function (this: unknown, ...args: unknown[]) {
    return originalFn.apply(this, args)
  }
}

/**
 * 创建防抖函数
 */
function debounce(fn: (...args: unknown[]) => unknown, delay = 1000, immediate = true) {
  let timer: number | null = null
  let hasExecuted = false

  return function (this: unknown, ...args: unknown[]) {
    if (timer) clearTimeout(timer)

    if (immediate && !hasExecuted) {
      // 立即执行
      fn.apply(this, args)
      hasExecuted = true
    }

    timer = window.setTimeout(() => {
      if (!immediate) {
        // 延迟执行
        fn.apply(this, args)
      }
      timer = null
      hasExecuted = false
    }, delay)
  }
}

/**
 * 指令处理函数
 */
const handler = (el: HTMLElement, binding: { value?: number }, vnode: VNode) => {
  const { ctx } = vnode
  const delay = binding.value || 1000

  // 处理ElButton组件
  if (ctx.type.name === 'ElButton') {
    const click = ctx.vnode.props?.onClick

    // 如果没有保存原始点击事件，则保存
    if (!ctx.vnode.props._onClick && click) {
      ctx.vnode.props._onClick = click
      ctx.vnode.key = `debounce_${Math.random().toString(36).substr(2, 9)}`
    }

    // 设置初始loading状态
    ctx.props.loading = ctx.vnode.loading || false

    // 使用防抖函数包装原始点击事件
    if (ctx.vnode.props._onClick) {
      const debouncedClick = debounce(createFunctionCopy(ctx.vnode.props._onClick), delay, true)

      // 使用AOP包装防抖后的点击事件，添加loading效果
      ctx.vnode.props.onClick = AOP(
        debouncedClick,
        () => {
          // 点击前立即设置loading状态
          ctx.props.loading = true
          ctx.vnode.loading = true
          // 强制更新组件状态
          if (ctx.proxy && typeof ctx.proxy.$forceUpdate === 'function') ctx.proxy.$forceUpdate()
        },
        () => {
          // 操作完成后取消loading状态
          ctx.props.loading = false
          ctx.vnode.loading = false
          // 强制更新组件状态
          if (ctx.proxy && typeof ctx.proxy.$forceUpdate === 'function') ctx.proxy.$forceUpdate()
        }
      )
    }
  } else {
    // 处理普通元素
    const elementWithListener = el as HTMLElementWithListener

    // 移除旧的事件监听器
    if (elementWithListener.dListener) {
      el.removeEventListener('click', elementWithListener.dListener)
    }

    // 创建新的防抖事件监听器
    elementWithListener.dListener = (_ev: Event) => {
      if (elementWithListener.disabled) return

      ctx.props.loading = true
      elementWithListener.disabled = true

      setTimeout(() => {
        elementWithListener.disabled = false
        ctx.props.loading = false
      }, delay)
    }

    // 添加事件监听
    el.addEventListener(
      'click',
      debounce(
        (...args: unknown[]) => {
          const ev = args[0] as Event
          elementWithListener.dListener?.(ev)
        },
        delay,
        true
      )
    )
  }
}

export default {
  mounted: handler,
  updated: handler,
}
