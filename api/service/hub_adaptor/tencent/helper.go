package tencent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/service/hub_adaptor/openai"
	"github.com/songquanpeng/one-api/relay/model"
)

// ConvertRequest 将OpenAI格式请求转换为腾讯云格式
func ConvertRequest(request model.GeneralOpenAIRequest, conversationID int64, UserId string, botAppKey string) *TencentRequest {
	// 构建内容和系统角色
	var content string
	var systemRole string

	for _, message := range request.Messages {
		messageContent := message.StringContent()
		if message.Role == "user" {
			content = messageContent
		} else if message.Role == "system" {
			systemRole = messageContent
		}
	}

	// 生成会话ID和访客ID
	sessionID := fmt.Sprintf("%d", conversationID)
	visitorBizID := UserId

	tencentReq := &TencentRequest{
		Content:      content,
		SessionID:    sessionID,
		BotAppKey:    botAppKey,
		VisitorBizID: visitorBizID,
		SystemRole:   systemRole,
		Incremental:  true,
	}

	// 设置流式传输
	if request.Stream {
		tencentReq.Stream = "enable"
	} else {
		tencentReq.Stream = "disable"
	}

	// 设置模型名称
	if request.Model != "" {
		tencentReq.ModelName = request.Model
	}

	logger.SysLogf("tencent request: %+v", tencentReq)
	return tencentReq
}

func generateUUID() string {
	// 简单的UUID生成（实际应用中应使用更完善的UUID库）
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rand.Uint32(),
		rand.Uint32()&0xffff,
		rand.Uint32()&0xffff,
		rand.Uint32()&0xffff,
		rand.Uint64()&0xffffffffffff)
}

// ConvertResponse 将腾讯云响应转换为OpenAI格式
func ConvertResponse(tencentResp *TencentResponse, modelName string) *openai.TextResponse {
	choice := openai.TextResponseChoice{
		Index: 0,
		Message: model.Message{
			Role:    "assistant",
			Content: tencentResp.Payload.Content,
		},
		FinishReason: "stop",
	}

	usage := model.Usage{
		PromptTokens:     0, // 腾讯云响应中没有token统计
		CompletionTokens: 0,
		TotalTokens:      0,
	}

	return &openai.TextResponse{
		Id:      tencentResp.Payload.RequestID,
		Object:  "chat.completion",
		Created: tencentResp.Payload.Timestamp,
		Model:   modelName,
		Choices: []openai.TextResponseChoice{choice},
		Usage:   usage,
	}
}

// ConvertStreamResponse 将腾讯云流式响应转换为OpenAI格式
func ConvertStreamResponse(data string, modelName string, previousContent string) *openai.ChatCompletionsStreamResponse {
	// 解析腾讯云流式数据
	var tencentResp TencentResponse
	if err := json.Unmarshal([]byte(data), &tencentResp); err != nil {
		logger.SysError("failed to parse tencent stream data: " + err.Error())
		return nil
	}

	// 只处理reply类型的响应
	if tencentResp.Type != "reply" {
		return nil
	}

	if tencentResp.Payload.IsFromSelf {
		// 不需要再重复一次自己的内容
		return nil
	}

	// 计算增量内容
	currentContent := tencentResp.Payload.Content
	var deltaContent string

	if len(currentContent) > len(previousContent) && strings.HasPrefix(currentContent, previousContent) {
		// 计算增量部分
		deltaContent = currentContent[len(previousContent):]
	} else {
		// 如果无法计算增量，返回完整内容（第一次或内容重置）
		deltaContent = currentContent
	}
	choice := openai.ChatCompletionsStreamResponseChoice{
		Index: 0,
		Delta: model.Message{
			Content: deltaContent,
			Role:    "assistant",
		},
	}

	// 判断是否结束
	emptyStr := ""
	var finishReason *string
	finishReason = &emptyStr
	if tencentResp.Payload.IsFinal {
		stopReason := "stop"
		finishReason = &stopReason
	}
	choice.FinishReason = finishReason

	// 使用session_id作为响应ID
	responseID := tencentResp.Payload.SessionID
	if responseID == "" {
		responseID = tencentResp.Payload.RequestID
	}

	return &openai.ChatCompletionsStreamResponse{
		Id:      responseID,
		Object:  "chat.completion.chunk",
		Created: tencentResp.Payload.Timestamp,
		Model:   modelName,
		Choices: []openai.ChatCompletionsStreamResponseChoice{choice},
	}
}

// IsStreamEnd 判断是否为流式响应结束
func IsStreamEnd(event string) bool {
	return strings.Contains(event, "finish") || strings.Contains(event, "done") || strings.Contains(event, "end")
}
