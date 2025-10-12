package n8n

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/render"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/model"
)

// ProcessN8nWorkflowResponse 处理 n8n 工作流响应
func ProcessN8nWorkflowResponse(c *gin.Context, resp *http.Response, meta interface{}) (*custom.WorkflowResponseData, error) {
	workflowAdaptor := &N8nWorkflowAdaptor{}

	// 处理响应
	workflowResponse, err := workflowAdaptor.ProcessResponse(resp)
	if err != nil {
		logger.SysErrorf("处理n8n工作流响应失败: %v", err)
		return nil, err
	}

	return workflowResponse, nil
}

// StreamHandler 处理 n8n 流式响应（如果需要的话）
func StreamHandler(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, *string, string) {
	var responseText string
	scanner := bufio.NewScanner(resp.Body)

	// 设置更大的缓冲区以处理大型响应
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanLines)
	common.SetEventStreamHeaders(c)

	channelConversationId := ""

	for scanner.Scan() {
		data := scanner.Text()
		if len(data) < 5 {
			continue
		}

		// n8n 可能不支持标准的 SSE 格式，这里预留处理逻辑
		// logger.SysLogf("n8n stream data: %s", data)

		// 如果 n8n 支持流式响应，在这里处理
		// 目前 n8n webhook 通常返回完整响应，不是流式的
	}

	if err := scanner.Err(); err != nil {
		logger.SysError("error reading n8n stream: " + err.Error())
	}

	render.Done(c)

	err := resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil, channelConversationId
	}

	return nil, &responseText, channelConversationId
}

// Handler 处理 n8n 阻塞式响应
func Handler(c *gin.Context, resp *http.Response, promptTokens int, modelName string) (*model.ErrorWithStatusCode, *string, string) {
	channelConversationId := ""

	// 使用工作流适配器处理响应
	workflowAdaptor := &N8nWorkflowAdaptor{}
	workflowResponse, err := workflowAdaptor.ProcessResponse(resp)

	if err != nil {
		return openai.ErrorWrapper(err, "process_n8n_response_failed", http.StatusInternalServerError), nil, channelConversationId
	}

	// 转换为 OpenAI 格式响应
	fullTextResponse := ResponseN8nToOpenAI(workflowResponse, modelName)
	jsonResponse, err := json.Marshal(fullTextResponse)
	if err != nil {
		return openai.ErrorWrapper(err, "marshal_response_body_failed", http.StatusInternalServerError), nil, channelConversationId
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = c.Writer.Write(jsonResponse)

	var responseText string
	if len(fullTextResponse.Choices) > 0 {
		responseText = fullTextResponse.Choices[0].Message.StringContent()
	}

	// n8n 没有会话ID概念，使用执行ID
	channelConversationId = workflowResponse.ExecuteID

	return nil, &responseText, channelConversationId
}

// ResponseN8nToOpenAI 将 n8n 工作流响应转换为 OpenAI 格式
func ResponseN8nToOpenAI(workflowResponse *custom.WorkflowResponseData, modelName string) *openai.TextResponse {
	var responseText string

	// 从工作流输出数据中提取主要内容
	if output, exists := workflowResponse.WorkflowOutputData["output"]; exists {
		if outputStr, ok := output.(string); ok {
			responseText = outputStr
		} else {
			// 如果不是字符串，转换为 JSON
			if jsonBytes, err := json.Marshal(output); err == nil {
				responseText = string(jsonBytes)
			} else {
				responseText = fmt.Sprintf("%v", output)
			}
		}
	} else {
		// 如果没有 output 字段，将整个输出数据转换为 JSON
		if jsonBytes, err := json.Marshal(workflowResponse.WorkflowOutputData); err == nil {
			responseText = string(jsonBytes)
		} else {
			responseText = "n8n workflow executed successfully"
		}
	}

	choice := openai.TextResponseChoice{
		Index: 0,
		Message: model.Message{
			Role:    "assistant",
			Content: responseText,
			Name:    nil,
		},
		FinishReason: "stop",
	}

	fullTextResponse := openai.TextResponse{
		Id:      fmt.Sprintf("n8n-exec-%s", workflowResponse.ExecuteID),
		Model:   modelName,
		Object:  "chat.completion",
		Created: helper.GetTimestamp(),
		Choices: []openai.TextResponseChoice{choice},
	}

	return &fullTextResponse
}

// ValidateN8nWorkflowID 验证 n8n 工作流ID格式
func ValidateN8nWorkflowID(workflowID string) bool {
	// n8n 工作流ID通常是 UUID 格式
	// 例如: d3c8ffe7-f052-4f51-a7b3-2e226d0b4d0c
	if len(workflowID) == 0 {
		return false
	}

	// 简单验证：包含连字符且长度合理
	return strings.Contains(workflowID, "-") && len(workflowID) >= 10
}
