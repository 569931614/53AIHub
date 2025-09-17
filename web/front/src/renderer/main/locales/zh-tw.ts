export default {
  // 模塊導航
  module: {
    index: '首頁',
    toolbox: 'AI工具',
    tool: 'AI工具',
    agent: '智能體',
    chat: '对话',
    app: '應用',
    find: '發現',
    library: '知識庫',
    crop_space: '企業空間',
    prompt: '提示詞'
  },

  // 瀏覽器相關
  browser: {
    newTab: '新標籤頁',
    search_placeholder: '在百度AI中搜索，或者輸入一個網址',
    knowledge_clean: '知識清洗',
    load_error: '載入失敗',
    reload: '重新載入',
    not_logged_in: '未登入',
    login: '立即登入',
    window_position: '視窗位置',
    leftWindow: '左側小窗',
    centerWindow: '居中大窗',
    rightWindow: '右側小窗',
    alwaysOnTop: '視窗置頂'
  },

  // 清洗功能
  clean: {
    content: '內容抽取',
    knowledge: '知識整理',
    summarize: '知識萃取'
  },

  // 更新相關
  updater: {
    newVersion: '發現新版本，是否立即更新？',
    newVersionTitle: '更新提示',
    downloading: '正在下載更新',
    completed: '更新完成，是否立即重啟？',
    completedTitle: '更新提示',
    noUpdate: '當前已是最新版本',
    currentVersion: '當前版本 v{version}',
    checkUpdate: '檢查更新'
  },

  // 通用文本
  common: {
    new_friend: '新手指引',
    help_feedback: '幫助與反饋',
    attach: '附件',
    screenshot: '截圖',
    title: '標題',
    desc: '描述',
    tip: '提示',
    create_time: '創建時間',
    update_time: '更新時間',
    recently_updated: '最近更新',
    loading: '載入中',
    // 加載更多
    load_more: '載入更多',
    load_new: '載入新訊息',
    no_data: '暫無數據',
    unknown: '未知',
    usual: '常用',

    render: '渲染',
    preview: '預覽',
    code: '代碼',
    split: '分屏',
    all: '全部',
    recently_visit: '最近訪問',

    confirm_delete: '確定刪除該項？',
    copied: '已複製',
    image_cropper: '圖片裁剪',
    back: '返回',
    go_admin: '去後台',
    related_prompt: '相關提示詞',
    related_agent: '相關智能體',
    view_more: '查看更多',
    expand_more: '展開更多',
    collapse_more: '收起更多',
    or: '或',
    later: '稍後再說',
    allow_to: '已複製提示詞內容，訪問"{name}"站點',
    ai_generated: '內容由AI生成，僅供參考',
    back_home: '返回站點首頁'
  },

  // 狀態提示
  status: {
    error: '錯誤',
    warning: '警告',
    info: '信息',
    success: '成功',
    updated: '已更新',
    sent: '已發送',
    login_success: '已登入',
    logout_success: '已退出登入',
    not_found_account: '未找到賬號,請先註冊後再登入',
    existing_account: '已有帳號',
    approve_success: '已點讚',
    approve_cancel: '已取消',
    register_success: '註冊成功',
    update_success: '更改成功',
    copy_success: ' 複製成功 ',
    copy_error: ' 複製失敗'
  },

  // 操作按鈕
  action: {
    login: '登入',
    login_directly: '直接登入',
    logout: '退出登入',
    register: '註冊',
    user_register: '用戶註冊',
    forget_password: '忘記密碼？',
    update_password: '更新密碼',
    search: '搜索',
    setting: '設置',
    confirm: '確認',
    cancel: '取消',
    close: '關閉',
    open: '打開',
    ok: '確定',
    allow: '允許',
    and: '和',
    copy: '複製',
    upload: '上傳',
    download: '下載',
    del: '刪除',
    add: '添加',
    create: '創建',
    edit: '編輯',
    save: '保存',
    fixed: '固定',
    unfixed: '取消固定',
    collapse: '收起',
    expand: '展開',
    back: '返回',
    share: '分享',
    favorite: '收藏',
    unfavorite: '取消收藏',
    more: '更多',
    refresh: '刷新',
    apply: '應用',
    reset: '重置',
    submit: '提交',
    manage: '管理',
    rename: '重命名',
    find: '發現',
    preview: '預覽',
    reupload: '重新上傳',
    send: '發送',
    stop: '停止',
    pay: '去支付',
    ok_v2: '好的',
    view_more: '查看更多',
    click_upload: '點擊上傳',
    view: '查看',
    delete: '刪除',
    select_all: '全選'
  },

  file: {
    file_size: '單個文件大小不超過{size}MB',
    file_format: '支持格式：{format}',
    file_empty: '文件內容不能為空',
    file_exceed: '文件上傳大小不能超過{size}MB({name}大小已超出)',
    file_exceed_limit: '最多上傳{limit}個文件'
  },

  // 表單相關
  form: {
    icon: '圖標',
    name: '名稱',
    desc: '描述',
    input_placeholder: '請輸入',
    select_placeholder: '請選擇',
    icon_placeholder: '為保證顯示效果，請上傳寬高比 1:1 的圖片',
    required: '必填項',
    optional: '選填項',

    // 知識庫公開範圍
    library_view_scope: '知識庫公開範圍',
    library_view_public: '公開訪問',
    library_view_team: '團隊空間成員可訪問',
    library_view_member: '知識庫成員可訪問',

    avatar: '頭像',
    nickname: '暱稱',
    select_verify: '請選擇驗證方式',
    email_verify: '郵箱驗證',
    mobile_verify: '手機驗證',
    account: '賬號',
    remark: ' 備註 ',
    account_alias: '手機',
    password: '密碼',
    password_placeholder: '請輸入密碼',
    email: '郵箱',
    new_email: '新郵箱',
    mobile: '手機號',
    new_mobile: '新手機號',
    verify_code: '驗證碼',
    get_verify_code: '獲取驗證碼',
    verify_old_mobile: '驗證你的舊手機號碼',
    reset_password: '重置密碼',
    reset_password_method: '選擇重置密碼的方式',
    new_password: '新密碼',
    new_password_placeholder: '請輸入新密碼',
    new_password_confirm: '確認新密碼',
    new_password_confirm_placeholder: '請再次輸入新密碼',
    mobile_format: '請輸入正確手機號',
    email_format: '請輸入正確的郵箱',
    account_format: '請輸入正確的手機或郵箱',
    verify_code_format: '請輸入正確的驗證碼',
    password_length: '密碼長度必須在 8-20 個字符之間',
    password_format: '密碼必須包含大小寫字母和數字',
    password_not_match: '兩次輸入的密碼不一致',
    password_no_chinese: '密碼不能包含中文',
    change: '修改',
    existing_mobile: '該手機號已註冊，請使用該帳號登入。',
    existing_email: '該郵箱已註冊，請使用該帳號登入。'
  },

  // 首頁相關
  index: {
    agent_recommend: '精選智能體',
    agent_recommend_desc: '打造您的第二生產力引擎',
    prompt_recommend: '精選提示詞',
    prompt_recommend_desc: '用好提示詞才能始終領先一步',
    toolbox_recommend: '精選AI工具',
    toolbox_recommend_desc: '用好 AI 工具才能始終領先一步',
    use_history: '{count}人聊過',
    banner_title: '“{name} 專家”  效率倍增的秘密武器',
    banner_desc: '精選工作場景 AI 資源，讓 AI 成為你的生產力助手',
    search_placeholder: '輸入您想查找的AI資源',
    hot_search: '熱門搜索'
  },

  // 登入相關
  login: {
    agree: '登入即視您同意',
    terms_of_service: '條款和條件',
    privacy_policy: '隱私政策',
    password_login: '密碼登入',
    password_login_title: '密碼登入',
    message_login: '簡訊登入',
    message_login_title: '簡訊登入',
    wechat_login: '微信登入',
    wechat_login_title: '微信登入',
    wecom_login: '企業微信',
    wecom_login_title: '企業微信登入',
    unregistered_account_desc: '未註冊的賬號，請先註冊後再自動登入',
    unregistered_account_confirm: '未註冊手機號驗證後自動登錄',
    quit: '退出',
    login_by_wechat: '請使用微信登入',
    login_by_password: '使用密碼登入',
    login_by_mobile: '使用手機登入',
    bind_mobile: '綁定手機號',
    other_login_way: '其他登入方式',
    login_by_wecom: '請使用企業微信登入',
    immediate_login: '立即登入',
    wecom_login_error: '站點綁定的企業微信並非當前登錄的企業微信!\n請用正確的企微登入。'
  },

  // 註冊相關
  register: {
    agree: '註冊即視您同意',
    terms_of_service: '條款和條件',
    privacy_policy: '隱私政策',
    unregistered: '未註冊'
  },

  // 聊天相關
  chat: {
    search_placholder: '搜索智能體',
    collapse_side_bar: '收起側邊欄',
    expand_side_bar: '展開側邊欄',
    history: '對話記錄',
    new_conversation: '新對話',
    edit_conversation: '編輯對話名稱',

    conversation_confirm_delete: '刪除後，聊天記錄將不可恢復。',
    input: '輸入',
    output: '輸出',
    start_generate: '開始生成',
    regenerate: '重新生成',
    like: '點讚',
    dislike: '取消點讚',
    input_placeholder: '發送訊息',
    usage_guide: '使用指引',
    completion_empty_desc: '採用AI大模型智能生成內容，輸入需求即可一鍵成文，快去試試吧',
    usage_case: '使用案例',
    usage_scene: '使用場景',
    online_search: '聯網搜索',
    completion_generating: '內容生成中，請稍候...',
    completion_rendering: '內容渲染中...',

    // 暫無可用的智能體
    no_available_agent: '暫無可用的智能體',
    completion_completed: '已完成深度思考',
    completion_thinking: '深度思考中...',
    completion_scene: '關聯場景：',
    completion_next_action: '關聯場景：下一步操作',

    not_found_url: '無法找到 URL，請檢查智能體是否設定正確',
    completion_share_link: '對話鏈接已複製，快去分享吧~',
    no_available_agent_desc: '很遺憾，你要訪問的智能體已不存在，去看看其他有趣的智能體吧～'
  },

  // 知識庫相關
  library: {
    home: '首頁',
    name: '知識庫',
    create: '創建知識庫',
    all_spaces: '全員空間',
    space_name: '空間名稱',
    chat: '知識庫問答',
    file_count: '{count}個文件',
    upload_file: '上傳文件',
    manage: '管理知識庫',
    search: '搜索知識庫',
    all_libraries: '全部知識庫',
    empty_desc: '創建企業知識庫，沉澱團隊知識，提升協作效率',

    create_md: '新建知識',
    local_upload: '本地上傳',
    create_folder: '新建文件夾',
    tmpl_create: '模板創建',
    github_import: '從Github導入',
    docs_qq_import: '從騰訊文檔導入',
    docs_feishu_import: '從釘釘文檔導入',
    docs_dingtalk_import: '從飛書文檔導入'
  },

  agent: {
    no_data: '管理員還未發布智能體',
    dialogue_type: '對話型',
    applied_type: '應用型'
  },

  prompt: {
    default_sort: '默認排序',
    likes_sort: '按最多點讚排序',
    views_sort: '按最多瀏覽排序',
    content: '提示詞內容',
    auth_tip: '您所在的分组无权限',
    let_use_prompt: '去使用这个提示词'
  },

  toolbox: {
    // AI聚合问答
    name: 'AI聚合问答',
    title: '準備好了嗎？開始體驗百寶箱吧！',
    input_placeholder: '百寶箱可以讓你快速找到你需要的答案，Shift+Enter换行，Enter发送',
    max_agent_count: '最多只能选择{count}个模型',
    search_placeholder: '搜索網站',
    direct_access: ' 直接訪問 ',
    click_access: ' 點擊訪問 ',
    account_access: ' 共享帳號 ',
    account_text: ' 以下是共享的帳號與密碼，僅限本人使用',
    account_text2: '您可以複製下列共享的帳號密碼，訪問目標站點'
  },

  // 空間相關
  space: {
    create: '創建空間',
    join: '加入空間',
    manage: '管理空間',
    member: '成員管理',
    settings: '空間設置',
    space_usage: '已用 {used} / {total}',
    expand: '擴容'
  },

  // 時間相關
  time: {
    just_now: '剛剛',
    minutes_ago: '{count}分鐘前',
    hours_ago: '{count}小時前',
    days_ago: '{count}天前',
    yesterday: '昨天',
    today: '今天',
    week_last: '一週內',
    week_ago: '一週前'
  },

  profile: {
    info: '個人信息',
    profile: '個人資料',
    change_password: '修改密碼',
    login_password: '登入密碼',
    bind_accounts: '賬戶信息',
    subscription_info: '訂閱信息',
    order_info: '訂單記錄',
    unbind_account: '未綁定',
    bind_mobile: '手機號',
    bind_email: '電子郵箱',
    bind_google: '谷歌賬號',
    bind_wechat: '微信OpenId',
    bind: '綁定',
    unbind: '解綁',
    change: '更換',
    logined_devices: '已登入的設備',
    used_devices: '曾用設備',
    current_device: '當前設備',
    unbind_wechat_confirm_desc: '解綁後，將無法使用微信登入，是否繼續？',
    unbind_wechat_confirm_title: '解綁微信賬號綁定',
    unbind_wechat_confirm_cancel: '暫不解綁',
    unbind_wechat_confirm_ok: '確認解綁',
    renew: '續費',
    expire_time_desc: '{group_name}將於{day}天後({expire_time})到期，馬上續費延續會員權益',
    expired_time_desc: '{group_name}已於{expire_time}到期，馬上續費延續會員權益'
  },
  response_status: {
    '400': '請求參數錯誤',
    '401': '未授權，請登入',
    '403': '拒絕訪問',
    '404': '請求資源不存在',
    '500': '伺服器錯誤',
    '502': '網關錯誤',
    '504': '網關超時',
    '503': '服務不可用'
  },
  response_code: {
    success: '成功',
    param_error: '參數錯誤',
    database_error: '數據庫錯誤',
    network_error: '網絡錯誤，請重新嘗試',
    system_error: '系統錯誤',
    auth_error: '認證錯誤',
    not_found_error: '未找到錯誤',
    unauthorized_error: '身份認證錯誤',
    file_error: '文件錯誤',
    forbidden_error: '用戶已禁用，請聯繫管理員',
    agent_error: '智能體錯誤',
    token_expired_error: '身份過期錯誤',
    verification_code_error: '驗證碼錯誤'
  },
  response_message: {
    unknown_error: '未知錯誤',
    user_not_found: '用戶不存在',
    username_or_password_is_incorrect: '用戶名或密碼錯誤',
    username_already_exists: '用戶名已存在',
    mobile_already_bind: '手機號已被其他用戶綁定',
    email_already_bind: '郵箱已被其他用戶綁定',
    model_changed: '管理員已修改模型，請新建對話',
    wechat_already_bind: '微信號已被其他用戶綁定'
  },
  subscription: {
    free: '免費版',
    upgrade: '升級',
    version_title: '選擇版本',
    time_title: '選擇購買時長',
    month: '月',
    year: '年',
    credit_month_amount: '每月{amount}積分',
    agent_bots_title: '可使用智能體',
    ai_assistant_title: '跨平台的 AI 助手',
    time_unit_month: '1個月',
    time_unit_year: '1年',
    aside_title: '新訂單購買詳情',
    aside_desc: '為以下用戶購買',
    payment: '支付方式',
    wechat_pay: '微信',
    alipay: '支付寶',
    manual_pay: '手動轉賬',
    paypal: 'PayPal',
    total: '總計',
    order_title: '訂單支付',
    pay_amount: '支付金額',
    pay_by_wechat: '使用 微信掃碼 支付',
    pay_policy: '支付即視為你同意相關 {policy}',
    pay_success: '支付成功',
    pay_confirm: '確認支付完成，請點擊',
    payed: '已支付',
    manual_pay_success_title: '訂單已生成',
    manual_pay_success_desc: '請等待或聯繫管理員確認支付信息',
    expire_time: '到期時間'
  },
  authority: {
    group_not_permission: '您所在的分組無權限',
    payment_not_setting: '站點未配置收費方式，請直接聯繫管理員',
    agent_not_permission: '無權限使用此智能體，請聯繫管理員',
    login_not_permission: '請先登入',
    use_range: '使用范围'
  },
  upgrade_dialog: {
    title: '站點註冊用戶已達上限，請聯繫管理員~~',
    cancel: '關閉註冊'
  },
  order: {
    id: '訂單ID',
    subscription: '訂閱服務',
    amount: '訂單金額',
    pay_type: '支付方式',
    create_time: '下單時間',
    status: '狀態',
    action: '操作',
    search_placeholder: '請輸入訂單ID',

    payment: '支付',
    cancel: '取消訂單',
    cancel_confirm: '確認取消該訂單？',

    status_list: {
      all: '全部',
      not_confirm: '待確定',
      pending: '待支付',
      paid: '已支付',
      expired: '已超時',
      cancel: '已取消'
    }
  },
  filter: {
    start_time: '開始時間',
    end_time: '結束時間',

    date_range: {
      today: '今天',
      last_7_days: '過去7天',
      last_4_weeks: '過去4周',
      last_3_months: '過去3月',
      last_12_months: '過去12月',
      this_month: '本月至今',
      this_quarter: '本季度至今',
      this_year: '本年至今',
      all_time: '所有時間'
    }
  },
  table: {
    footer_text: '共有 {total} 條'
  },
  guide: {
    title: '歡迎使用53AI Hub',
    description: '當前為初始化環節，需錄入站點等基礎資訊，用於系統基礎配置搭建',
    website_info: '站點資訊',
    website_setting: '帳號設置',
    website_success: '完成安裝',
    website_info_name: '站點名稱',
    website_info_name_placeholder: '請輸入站點的名稱',
    website_info_logo: '站點logo',
    website_info_logo_placeholder: '請上傳Logo',
    website_info_logo_change: '修改',
    website_info_logo_upload: '上傳',
    website_info_logo_tip:
      '請上傳1:1比例的方形Logo（建議尺寸：512×512px）支援PNG/JPG格式，文件大小不超過2MB',
    website_style: '站點風格',
    website_info_language: '預設語言',
    website: '網站風格',
    software: '軟體風格',
    init_success: '初始化成功！',
    jump_now: '立即跳轉',
    jump_tip: '將在 {count} 秒後自動跳轉到首頁...',
    next: '下一步',
    init: '初始化',
    confirm_password: '確認密碼',
    confirm_password_placeholder: '請再次確認密碼'
  }
}
