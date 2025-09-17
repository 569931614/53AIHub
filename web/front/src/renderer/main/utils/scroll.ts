/**
 * 通用滚动工具函数
 */

/**
 * 查找最近的滚动容器
 * @param element 起始元素
 * @returns 滚动容器元素
 */
export const findScrollContainer = (element: Element): Element | null => {
  let parent = element.parentElement
  while (parent) {
    const style = window.getComputedStyle(parent)
    if (style.overflowY === 'auto' || style.overflowY === 'scroll') {
      return parent
    }
    parent = parent.parentElement
  }
  return document.documentElement
}

/**
 * 滚动到指定元素
 * @param elementId 目标元素ID
 * @param offset 偏移量，默认150px
 * @param behavior 滚动行为，默认'smooth'
 */
export const scrollToElement = (
  elementId: string,
  offset = 150,
  behavior: ScrollBehavior = 'smooth'
) => {
  const targetElement = document.querySelector(elementId)
  if (!targetElement) return

  const scrollContainer = findScrollContainer(targetElement)
  if (!scrollContainer) return

  const containerRect = scrollContainer.getBoundingClientRect()
  const targetRect = targetElement.getBoundingClientRect()
  const scrollTop = scrollContainer.scrollTop + targetRect.top - containerRect.top - offset

  scrollContainer.scrollTo({
    top: Math.max(0, scrollTop),
    behavior
  })
}

/**
 * 滚动到指定元素（Promise版本）
 * @param elementId 目标元素ID
 * @param offset 偏移量，默认150px
 * @param behavior 滚动行为，默认'smooth'
 * @returns Promise，滚动完成后resolve
 */
export const scrollToElementAsync = (
  elementId: string,
  offset = 150,
  behavior: ScrollBehavior = 'smooth'
): Promise<void> => {
  return new Promise((resolve) => {
    scrollToElement(elementId, offset, behavior)

    // 等待滚动动画完成
    const scrollContainer = findScrollContainer(document.querySelector(elementId)!)
    if (scrollContainer) {
      const handleScrollEnd = () => {
        scrollContainer.removeEventListener('scroll', handleScrollEnd)
        resolve()
      }
      scrollContainer.addEventListener('scroll', handleScrollEnd)

      // 设置超时，防止滚动事件不触发
      setTimeout(resolve, 500)
    } else {
      resolve()
    }
  })
}
