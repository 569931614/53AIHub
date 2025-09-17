export default {
  // 模块导航
  module: {
    index: '首页',
    toolbox: 'AI工具',
    tool: 'AI工具',
    agent: '智能体',
    chat: '对话',
    app: '应用',
    find: '发现',
    library: '知识库',
    crop_space: '企业空间',
    prompt: '提示词'
  },

  // 浏览器相关
  browser: {
    newTab: '新标签页',
    search_placeholder: '在百度AI中搜索，或者输入一个网址',
    knowledge_clean: '知识清洗',
    load_error: '加载失败',
    reload: '重新加载',
    not_logged_in: '未登录',
    login: '立即登录',
    window_position: '窗口位置',
    leftWindow: '左侧小窗',
    centerWindow: '居中大窗',
    rightWindow: '右侧小窗',
    alwaysOnTop: '窗口置顶'
  },

  // 清洗功能
  clean: {
    content: '内容抽取',
    knowledge: '知识整理',
    summarize: '知识萃取'
  },

  // 更新相关
  updater: {
    newVersion: '发现新版本，是否立即更新？',
    newVersionTitle: '更新提示',
    downloading: '正在下载更新',
    completed: '更新完成，是否立即重启？',
    completedTitle: '更新提示',
    noUpdate: '当前已是最新版本',
    currentVersion: '当前版本 v{version}',
    checkUpdate: '检查更新'
  },

  // 通用文本
  common: {
    new_friend: '新手指引',
    help_feedback: '帮助与反馈',
    attach: '附件',
    screenshot: '截图',
    title: '标题',
    desc: '描述',
    tip: '提示',
    create_time: '创建时间',
    update_time: '更新时间',
    recently_updated: '最近更新',
    loading: '加载中',
    // 加载更多
    load_more: '加载更多',
    load_new: '加载新消息',
    no_data: '暂无数据',
    unknown: '未知',
    usual: '常用',

    render: '渲染',
    preview: '预览',
    code: '代码',
    split: '分屏',
    all: '全部',
    recently_visit: '最近访问',

    confirm_delete: '确定删除该项？',
    copied: '已复制',
    image_cropper: '图片裁剪',
    back: '返回',
    go_admin: '去后台',
    related_prompt: '相关提示词',
    related_agent: '相关智能体',

    view_more: '查看更多',

    expand_more: '展开更多',
    collapse_more: '收起更多',
    or: '或',
    later: '稍后再说',
    allow_to: '已复制提示词内容，访问"{name}"站点',
    ai_generated: '内容由AI生成，仅供参考',
    back_home: '返回站点首页'
  },

  // 状态提示
  status: {
    error: '错误',
    warning: '警告',
    info: '信息',
    success: '成功',
    updated: '已更新',
    sent: '已发送',
    login_success: '已登录',
    logout_success: '已退出登录',
    not_found_account: '未找到账号,请先注册后再登录',
    existing_account: '已有帐号',
    approve_success: '已点赞',
    approve_cancel: '已取消',
    register_success: '注册成功',
    update_success: '更改成功',
    copy_success: '复制成功',
    copy_error: '复制失败'
  },

  // 操作按钮
  action: {
    login: '登录',
    login_directly: '直接登录',
    logout: '退出登录',
    register: '注册',
    user_register: '用户注册',
    forget_password: '忘记密码？',
    update_password: '更新密码',
    search: '搜索',
    setting: '设置',
    confirm: '确认',
    cancel: '取消',
    close: '关闭',
    open: '打开',
    ok: '确定',
    allow: '允许',
    and: '和',
    copy: '复制',
    upload: '上传',
    download: '下载',
    del: '删除',
    add: '添加',
    create: '创建',
    edit: '编辑',
    save: '保存',
    fixed: '固定',
    unfixed: '取消固定',
    collapse: '收起',
    expand: '展开',
    back: '返回',
    share: '分享',
    favorite: '收藏',
    unfavorite: '取消收藏',
    more: '更多',
    refresh: '刷新',
    apply: '应用',
    reset: '重置',
    submit: '提交',
    manage: '管理',
    rename: '重命名',
    find: '发现',
    preview: '预览',
    reupload: '重新上传',
    send: '发送',
    stop: '停止',
    pay: '去支付',
    ok_v2: '好的',
    view_more: '查看更多',
    click_upload: '点击上传',
    view: '查看',
    delete: '删除',
    select_all: '全选',
    unselect_all: '取消全选',
    copy_link: '复制链接'
  },

  file: {
    support_format: '目前仅支持{format}格式文件',
    file_size: '单个文件大小不超过{size}MB',
    file_format: '支持格式：{format}',
    file_empty: '文件内容不能为空',
    file_exceed: '文件上传大小不能超过{size}MB({name}大小已超出)',
    file_exceed_limit: '最多上传{limit}个文件'
  },

  // 表单相关
  form: {
    icon: '图标',
    name: '名称',
    desc: '描述',
    input_placeholder: '请输入',
    select_placeholder: '请选择',
    icon_placeholder: '为保证显示效果，请上传宽高比 1:1 的图片',
    required: '必填项',
    optional: '选填项',

    // 知识库公开范围
    library_view_scope: '知识库公开范围',
    library_view_public: '公开访问',
    library_view_team: '团队空间成员可访问',
    library_view_member: '知识库成员可访问',

    avatar: '头像',
    nickname: '昵称',
    select_verify: '请选择验证方式',
    email_verify: '邮箱验证',
    mobile_verify: '手机验证',
    account: '账号',
    remark: '备注',
    account_alias: '手机',
    password: '密码',
    password_placeholder: '请输入密码',
    email: '邮箱',
    new_email: '新邮箱',
    mobile: '手机号',
    new_mobile: '新手机号',
    verify_code: '验证码',
    get_verify_code: '获取验证码',
    verify_old_mobile: '验证你的旧手机号码',
    reset_password: '重置密码',
    reset_password_method: '选择重置密码的方式',
    new_password: '新密码',
    new_password_placeholder: '请输入新密码',
    new_password_confirm: '重复新密码',
    new_password_confirm_placeholder: '请再次输入新密码',
    email_format: '请输入正确的邮箱',
    mobile_format: '请输入正确手机号',
    account_format: '请输入正确的手机或邮箱',
    verify_code_format: '请输入正确的验证码',
    username_length: '账号长度必须在 5-20 个字符之间',
    password_length: '密码长度必须在 8-20 个字符之间',
    password_format: '密码必须包含大小写字母和数字',
    password_not_match: '两次输入的密码不一致',
    password_no_chinese: '密码不能包含中文',
    change: '修改',
    existing_mobile: '该手机号已注册，请使用该账号登录。',
    existing_email: '该邮箱已注册，请使用该账号登录。'
  },

  // 首页相关
  index: {
    agent_recommend: '精选智能体',
    agent_recommend_desc: '打造您的第二生产力引擎',
    prompt_recommend: '精选提示词',
    prompt_recommend_desc: '掌握用大模型解决问题的方法',
    toolbox_recommend: '精选AI工具',
    toolbox_recommend_desc: '用好 AI 工具才能始终领先一步',
    use_history: '{count}人聊过',
    banner_title: '“{name} 专家”  效率倍增的秘密武器',
    banner_desc: '精选工作场景 AI 资源，让 AI 成为你的生产力助手',
    search_placeholder: '输入您想查找的AI资源',
    hot_search: '热门搜索'
  },

  // 登录相关
  login: {
    agree: '登录即视您同意',
    terms_of_service: '条款和条件',
    privacy_policy: '隐私政策',
    password_login: '账号密码',
    password_login_title: '密码登录',
    message_login: '短信登录',
    message_login_title: '短信登录',
    wechat_login: '微信登录',
    wechat_login_title: '微信登录',
    wecom_login: '企业微信',
    wecom_login_title: '企业微信登录',
    unregistered_account_desc: '未注册的账号，请先注册后再登录',
    unregistered_account_confirm: '未注册手机号验证后自动登录，',
    quit: '退出',
    login_by_wechat: '请使用微信扫码登录',
    login_by_password: '使用密码登录',
    login_by_mobile: '使用手机登录',
    bind_mobile: '绑定手机号',
    other_login_way: '其他登录方式',
    login_by_wecom: '请使用企业微信扫码登录',
    immediate_login: '立即登录',
    wecom_login_error: '站点绑定的企业微信并非当前登录的企业微信!\n请用正确的企微登录。'
  },

  // 注册相关
  register: {
    agree: '注册即视您同意',
    unregistered: '未注册',
    terms_of_service: '条款和条件',
    privacy_policy: '隐私政策'
  },

  // 聊天相关
  chat: {
    search_placholder: '搜索智能体',
    collapse_side_bar: '收起侧边栏',
    expand_side_bar: '展开侧边栏',
    history: '历史会话',
    new_conversation: '新对话',
    edit_conversation: '编辑对话名称',

    conversation_confirm_delete: '删除后，聊天记录将不可恢复。',
    input: '输入',
    output: '输出',
    start_generate: '开始生成',
    regenerate: '重新生成',
    like: '点赞',
    dislike: '取消点赞',
    input_placeholder: '发送消息',
    usage_guide: '使用指引',
    completion_empty_desc: '采用AI大模型智能生成内容，输入需求即可一键成文，快去试试吧',
    usage_case: '使用案例',
    usage_scene: '使用场景',
    online_search: '联网搜索',
    completion_generating: '内容生成中，请稍候...',
    completion_rendering: '内容渲染中...',

    // 暂无可用的智能体
    no_available_agent: '暂无可用的智能体',
    completion_completed: '已完成深度思考',
    completion_thinking: '深度思考中...',
    completion_scene: '关联场景：',
    completion_next_action: '关联场景：下一步操作',

    not_found_url: '无法找到URL，请检查智能体是否设置正确',
    completion_share_link: '对话链接已复制，快去分享吧~',
    no_available_agent_desc: '很遗憾，你要访问的智能体已不存在，去看看其他有趣的智能体吧～'
  },

  // 知识库相关
  library: {
    home: '首页',
    name: '知识库',
    create: '创建知识库',
    all_spaces: '全员空间',
    space_name: '空间名称',
    chat: '知识库问答',
    file_count: '{count}个文件',
    upload_file: '上传文件',
    manage: '管理知识库',
    search: '搜索知识库',
    all_libraries: '全部知识库',
    empty_desc: '创建企业知识库，沉淀团队知识，提升协作效率',

    create_md: '新建知识',
    local_upload: '本地上传',
    create_folder: '新建文件夹',
    tmpl_create: '模板创建',
    github_import: '从Github导入',
    docs_qq_import: '从腾讯文档导入',
    docs_feishu_import: '从钉钉文档导入',
    docs_dingtalk_import: '从飞书文档导入'
  },

  agent: {
    no_data: '管理员还未发布智能体',
    dialogue_type: '对话型',
    applied_type: '应用型'
  },

  prompt: {
    default_sort: '默认排序',
    likes_sort: '按最多点赞排序',
    views_sort: '按最多浏览排序',
    content: '提示词内容',

    auth_tip: '你所在的分组无权限',
    let_use_prompt: '去使用这个提示词'
  },

  toolbox: {
    name: 'AI聚合问答',
    title: '准备好了吗？开始体验AI聚合问答啦！',
    input_placeholder: '多个AI解答你输入的问题，Shift+Enter换行，Enter发送',
    max_agent_count: '最多只能选择{count}个模型',
    search_placeholder: '搜索网站',
    direct_access: '直接访问',
    click_access: '点击访问',
    account_access: '共享账号',
    account_text: '以下是共享的账号与密码，仅限本人使用',
    account_text2: '您可以复制下列共享的账号密码，访问目标站点'
  },

  // 空间相关
  space: {
    create: '创建空间',
    join: '加入空间',
    manage: '管理空间',
    member: '成员管理',
    settings: '空间设置',
    space_usage: '已用 {used} / {total}',
    expand: '扩容'
  },

  // 时间相关
  time: {
    just_now: '刚刚',
    minutes_ago: '{count}分钟前',
    hours_ago: '{count}小时前',
    days_ago: '{count}天前',
    yesterday: '昨天',
    today: '今天',
    week_last: '一周内',
    week_ago: '一周前'
  },

  profile: {
    info: '个人信息',
    profile: '个人资料',
    change_password: '修改密码',
    login_password: '登录密码',
    // 已绑定的账户
    bind_accounts: '账户信息',
    // 订阅信息
    subscription_info: '订阅信息',
    order_info: '订单记录',

    unbind_account: '未绑定',
    bind_mobile: '手机号',
    bind_email: '电子邮箱',
    bind_google: '谷歌账号',
    bind_wechat: '微信OpenId',
    bind: '绑定',
    bind_success: '已绑定',
    unbind: '解绑',
    unbind_success: '已解绑',
    change: '更换',
    change_success: '已更换',
    logined_devices: '已登录的设备',
    used_devices: '曾用设备',
    current_device: '当前设备',
    unbind_wechat_confirm_desc: '解绑后将无法使用该微信账号登录此账号，请谨慎操作！',
    unbind_wechat_confirm_title: '解绑微信账号绑定',
    unbind_wechat_confirm_cancel: '暂不解绑',
    unbind_wechat_confirm_ok: '确认解绑',
    bind_wechat_title: '微信扫码完成绑定'
  },
  response_status: {
    '400': '请求参数错误',
    '401': '未授权，请登录',
    '403': '拒绝访问',
    '404': '请求资源不存在',
    '500': '服务器错误',
    '502': '网关错误',
    '504': '网关超时',
    '503': '服务不可用'
  },
  response_code: {
    success: '成功',
    param_error: '参数错误',
    database_error: '数据库错误',
    network_error: '网络错误，请重新尝试',
    system_error: '系统错误',
    auth_error: '认证错误',
    not_found_error: '未找到错误',
    unauthorized_error: '身份认证错误',
    file_error: '文件错误',
    forbidden_error: '用户已禁用，请联系管理员 ',
    agent_error: '智能体错误',
    token_expired_error: '身份过期错误',
    verification_code_error: '验证码错误'
  },
  response_message: {
    unknown_error: '未知错误',
    user_not_found: '用户不存在',
    username_or_password_is_incorrect: '用户名或密码错误',
    username_already_exists: '用户名已存在',
    mobile_already_bind: '手机号已被其他用户绑定',
    email_already_bind: '邮箱已被其他用户绑定',
    model_changed: '管理员已修改模型，请新建对话',
    wechat_already_bind: '微信号已被其他用户绑定'
  },
  subscription: {
    free: '免费版',
    upgrade: '升级',
    version_title: '选择版本',
    time_title: '选择购买时长',
    month: '月',
    year: '年',
    credit_month_amount: '每月{amount}积分',
    agent_bots_title: '可使用智能体',
    ai_assistant_title: '跨平台的 AI 助手',
    time_unit_month: '1个月',
    time_unit_year: '1年',
    aside_title: '新订单购买详情',
    aside_desc: '为以下用户购买',
    payment: '支付方式',
    wechat_pay: '微信',
    alipay: '支付宝',
    manual_pay: '手动转账',
    paypal: 'PayPal',
    total: '总计',
    order_title: '订单支付',
    pay_amount: '支付金额',
    pay_by_wechat: '使用 微信扫码 支付',
    pay_policy: '支付即视为你同意相关 {policy}',
    pay_success: '支付成功',
    pay_confirm: '确认支付完成，请点击',
    payed: '已支付',
    manual_pay_success_title: '订单已生成',
    manual_pay_success_desc: '请等待或联系管理员确认支付信息',
    expire_time: '到期时间',
    renew: '续费',
    expire_time_desc: '{group_name}将于{day}天后({expire_time})到期，马上续费延续会员权益',
    expired_time_desc: '{group_name}已于{expire_time}到期，马上续费延续会员权益'
  },
  authority: {
    group_not_permission: '您所在的分组无权限',
    payment_not_setting: '站点未配置收费方式，请直接联系管理员',
    agent_not_permission: '无权限使用此智能体，请联系管理员',
    login_not_permission: '请先登录',
    use_range: '使用范围'
  },
  upgrade_dialog: {
    title: '站点注册用户已达上限，请联系管理员~~',
    cancel: '关闭注册'
  },
  order: {
    id: '订单ID',
    subscription: '订阅服务',
    amount: '订单金额',
    pay_type: '支付方式',
    create_time: '下单时间',
    status: '状态',
    action: '操作',
    search_placeholder: '请输入订单ID',

    payment: '支付',
    cancel: '取消订单',
    cancel_confirm: '确认取消该订单？',

    status_list: {
      all: '全部',
      not_confirm: '待确定',
      pending: '待支付',
      paid: '已支付',
      expired: '已超时',
      cancel: '已取消'
    }
  },
  filter: {
    start_time: '开始时间',
    end_time: '结束时间',

    date_range: {
      today: '今天',
      last_7_days: '过去7天',
      last_4_weeks: '过去4周',
      last_3_months: '过去3月',
      last_12_months: '过去12月',
      this_month: '本月至今',
      this_quarter: '本季度至今',
      this_year: '本年至今',
      all_time: '所有时间'
    }
  },
  table: {
    footer_text: '共有 {total} 条'
  },
  guide: {
    title: '欢迎使用53AI Hub',
    description: '当前为初始化环节，需录入站点等基础信息，用于系统基础配置搭建',
    website_info: '站点信息',
    website_setting: '账号设置',
    website_success: '完成安装',
    website_info_name: '站点名称',
    website_info_name_placeholder: '请输入站点的名称',
    website_info_logo: '站点logo',
    website_info_logo_placeholder: '请上传Logo',
    website_info_logo_change: '修改',
    website_info_logo_upload: '上传',
    website_info_logo_tip:
      '请上传1:1比例的方形Logo（建议尺寸：512×512px）支持PNG/JPG格式，文件大小不超过2MB',
    website_style: '站点风格',
    website_info_language: '默认语言',
    website: '网站风格',
    software: '软件风格',
    init_success: '初始化成功！',
    jump_now: '立即跳转',
    jump_tip: '将在 {count} 秒后自动跳转到首页...',
    next: '下一步',
    init: '初始化',
    confirm_password: '确认密码',
    confirm_password_placeholder: '请再次确认密码'
  }
}
