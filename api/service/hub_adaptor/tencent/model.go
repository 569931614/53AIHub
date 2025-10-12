package tencent

// TencentRequest 腾讯云请求结构体
type TencentRequest struct {
	RequestID         string            `json:"request_id,omitempty"`
	Content           string            `json:"content"`
	SessionID         string            `json:"session_id"`
	BotAppKey         string            `json:"bot_app_key"`
	VisitorBizID      string            `json:"visitor_biz_id"`
	StreamingThrottle int32             `json:"streaming_throttle,omitempty"`
	CustomVariables   map[string]string `json:"custom_variables,omitempty"`
	SystemRole        string            `json:"system_role,omitempty"`
	Incremental       bool              `json:"incremental,omitempty"`
	SearchNetwork     string            `json:"search_network,omitempty"`
	ModelName         string            `json:"model_name,omitempty"`
	Stream            string            `json:"stream,omitempty"`
	WorkflowStatus    string            `json:"workflow_status,omitempty"`
	VisitorLabels     []VisitorLabel    `json:"visitor_labels,omitempty"`
	FileInfos         []FileInfo        `json:"file_infos,omitempty"`
	TcadpUserID       string            `json:"tcadp_user_id,omitempty"`
}

type VisitorLabel struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type FileInfo struct {
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	FileURL  string `json:"file_url"`
	FileType string `json:"file_type"`
	DocID    string `json:"doc_id"`
}

// TencentResponse 腾讯云响应结构体
type TencentResponse struct {
	Type    string                 `json:"type"`
	Payload TencentResponsePayload `json:"payload"`
}

// TencentResponsePayload 腾讯云响应载荷
type TencentResponsePayload struct {
	RequestID       string             `json:"request_id"`
	Content         string             `json:"content"`
	RecordID        string             `json:"record_id"`
	RelatedRecordID string             `json:"related_record_id"`
	SessionID       string             `json:"session_id"`
	IsFromSelf      bool               `json:"is_from_self"`
	CanRating       bool               `json:"can_rating"`
	Timestamp       int64              `json:"timestamp"`
	IsFinal         bool               `json:"is_final"`
	IsEvil          bool               `json:"is_evil"`
	IsLLMGenerated  bool               `json:"is_llm_generated"`
	ReplyMethod     uint8              `json:"reply_method"`
	Knowledge       []TencentKnowledge `json:"knowledge"`
	OptionCards     []string           `json:"option_cards"`
	CustomParams    []string           `json:"custom_params"`
	TaskFlow        interface{}        `json:"task_flow"`
	WorkFlow        interface{}        `json:"work_flow"`
	QuoteInfos      []TencentQuoteInfo `json:"quote_infos"`
}

// TencentKnowledge 知识结构
type TencentKnowledge struct {
	ID    string `json:"id"`
	Type  uint32 `json:"type"`
	SegID string `json:"seg_id"`
}

// TencentQuoteInfo 引用信息
type TencentQuoteInfo struct {
	Index    int `json:"index"`
	Position int `json:"position"`
}

// TencentStreamResponse 腾讯云流式响应
type TencentStreamResponse struct {
	Event string `json:"event"` // 事件类型
	Data  string `json:"data"`  // 数据内容
}

// TencentErrorResponse 腾讯云错误响应
type TencentErrorResponse struct {
	Type  string           `json:"type"`
	Error TencentErrorInfo `json:"error"`
}

type TencentErrorInfo struct {
	Code    uint32 `json:"code"`
	Message string `json:"message"`
}
