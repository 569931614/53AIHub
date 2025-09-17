package coze

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	db_model "github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/logger"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
)

type WorkflowAdaptor struct {
	meta         *meta.Meta
	CustomConfig *custom.CustomConfig
}

func (a *WorkflowAdaptor) Init(meta *meta.Meta) {
	a.meta = meta
}

func (a *WorkflowAdaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	baseUrl, err := custom.GetBaseURL(meta.BaseURL)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/v1/workflow/run", baseUrl), nil
}

func (a *WorkflowAdaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	custom.SetupCommonRequestHeader(c, req, meta)
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)
	return nil
}

func (a *WorkflowAdaptor) ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	// 提取工作流ID
	workflowID := extractWorkflowIDFromModel(request.Model)
	if workflowID == "" {
		return nil, errors.New("workflow ID cannot be empty")
	}

	// 提取用户输入（取最后一条用户消息）
	var input string
	for i := len(request.Messages) - 1; i >= 0; i-- {
		if request.Messages[i].Role == "user" {
			input = request.Messages[i].StringContent()
			break
		}
	}

	if input == "" {
		return nil, errors.New("no user input found in messages")
	}

	// 构建工作流请求参数
	parameters := make(map[string]interface{})

	// 优先使用前端直接传入的工作流参数（从 relay 层传递）
	if c != nil {
		if workflowParams, exists := c.Get("workflow_parameters"); exists {
			if params, ok := workflowParams.(map[string]interface{}); ok {
				// 直接使用传入的参数，不做任何验证或转换
				parameters = params
				logger.SysLogf("Coze工作流请求 - 使用前端传入的参数: %+v", parameters)
			}
		}
	}

	// 如果没有直接传入的参数，使用默认的 input 参数（向后兼容）
	if len(parameters) == 0 {
		parameters["input"] = input
		logger.SysLogf("Coze工作流请求 - 使用默认参数映射: input=%s", input)
	}

	// 构建工作流请求
	workflowRequest := &WorkflowRequest{
		WorkflowID: workflowID,
		Parameters: parameters,
	}

	logger.SysLogf("Coze工作流请求 - WorkflowID: %s, Parameters: %+v", workflowID, parameters)

	return workflowRequest, nil
}

// ConvertWorkflowRequest 直接从工作流参数构造请求，简化参数传递
func (a *WorkflowAdaptor) ConvertWorkflowRequest(workflowID string, parameters map[string]interface{}) *WorkflowRequest {
	logger.SysLogf("Coze工作流直接请求 - WorkflowID: %s, Parameters: %+v", workflowID, parameters)

	// 处理参数中的文件上传
	processedParameters, err := a.processWorkflowParameters(parameters)
	if err != nil {
		logger.SysErrorf("处理工作流文件参数失败: %v", err)
		// 如果文件处理失败，使用原始参数继续执行
		processedParameters = parameters
	}

	return &WorkflowRequest{
		WorkflowID: workflowID,
		Parameters: processedParameters,
	}
}

// processWorkflowParameters 处理工作流参数中的文件上传
func (a *WorkflowAdaptor) processWorkflowParameters(parameters map[string]interface{}) (map[string]interface{}, error) {
	if a.meta == nil {
		return parameters, fmt.Errorf("meta is nil")
	}

	processedParams := make(map[string]interface{})

	for key, value := range parameters {
		processedValue, err := a.processParameterValue(value)
		if err != nil {
			logger.SysErrorf("处理参数 %s 失败: %v", key, err)
			// 如果单个参数处理失败，使用原始值
			processedParams[key] = value
		} else {
			processedParams[key] = processedValue
		}
	}

	return processedParams, nil
}

// processParameterValue 递归处理参数值，支持字符串、数组、对象
func (a *WorkflowAdaptor) processParameterValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		// 检查是否为 file_id: 格式
		return a.processFileIDString(v)
	case []interface{}:
		// 处理数组
		processedArray := make([]interface{}, len(v))
		for i, item := range v {
			processedItem, err := a.processParameterValue(item)
			if err != nil {
				processedArray[i] = item // 使用原始值
			} else {
				processedArray[i] = processedItem
			}
		}
		return processedArray, nil
	case map[string]interface{}:
		// 处理对象
		processedMap := make(map[string]interface{})
		for k, val := range v {
			processedVal, err := a.processParameterValue(val)
			if err != nil {
				processedMap[k] = val // 使用原始值
			} else {
				processedMap[k] = processedVal
			}
		}
		return processedMap, nil
	default:
		// 其他类型直接返回
		return value, nil
	}
}

// processFileIDString 处理 file_id: 格式的字符串
func (a *WorkflowAdaptor) processFileIDString(value string) (interface{}, error) {
	// 检查是否为 file_id: 格式
	if !strings.HasPrefix(value, "file_id:") {
		return value, nil
	}

	// 提取文件ID
	fileIDStr := strings.TrimPrefix(value, "file_id:")
	fileID, err := strconv.ParseInt(fileIDStr, 10, 64)
	if err != nil {
		logger.SysErrorf("解析文件ID失败: %s, error: %v", fileIDStr, err)
		return value, err
	}

	// 获取上传文件对象
	uploadFile, err := db_model.GetUploadFileByID(fileID)
	if err != nil {
		logger.SysErrorf("获取上传文件失败: ID=%d, error: %v", fileID, err)
		return value, err
	}

	// 获取渠道文件映射
	channelID := a.meta.ChannelId
	modelName := "workflow-" + strings.TrimPrefix(a.meta.ActualModelName, "workflow-")

	fileMapping := uploadFile.GetChannelFileMapping(channelID, modelName)
	if fileMapping == nil {
		// 创建新的文件映射
		fileMapping = &db_model.ChannelFileMapping{}
		err := CozeUploadFile(a.meta, uploadFile, fileMapping)
		if err != nil {
			logger.SysErrorf("上传文件到Coze失败: %v", err)
			return value, err
		}
		err = db_model.CreateChannelFileMapping(fileMapping)
		if err != nil {
			logger.SysErrorf("创建文件映射失败: %v", err)
			return value, err
		}
	} else if helper.GetTimestamp() > fileMapping.ExpirationTime {
		// 文件映射已过期，重新上传
		err := CozeUploadFile(a.meta, uploadFile, fileMapping)
		if err != nil {
			logger.SysErrorf("重新上传文件到Coze失败: %v", err)
			return value, err
		}
		err = db_model.UpdateChannelFileMapping(fileMapping)
		if err != nil {
			logger.SysErrorf("更新文件映射失败: %v", err)
			return value, err
		}
	}

	// 根据文件类型生成对应的格式
	var fileObject map[string]interface{}
	if strings.HasPrefix(uploadFile.MimeType, "image/") {
		fileObject = map[string]interface{}{
			"type":    "image",
			"file_id": fileMapping.ChannelFileID,
		}
	} else {
		fileObject = map[string]interface{}{
			"type":    "file",
			"file_id": fileMapping.ChannelFileID,
		}
	}

	// 将文件对象转换为 JSON 字符串，然后包装在数组中（Coze 工作流要求的格式）
	fileObjectJSON, err := json.Marshal(fileObject)
	if err != nil {
		logger.SysErrorf("序列化文件对象失败: %v", err)
		return value, err
	}

	// 返回包含 JSON 字符串的数组格式
	fileArray := []string{string(fileObjectJSON)}

	logger.SysLogf("工作流文件处理成功 - 原始ID: %d, 渠道文件ID: %s, 类型: %s, 数组格式: %v",
		fileID, fileMapping.ChannelFileID, uploadFile.MimeType, fileArray)

	return fileArray, nil
}

func (a *WorkflowAdaptor) ConvertImageRequest(request *model.ImageRequest) (any, error) {
	return nil, errors.New("workflow adaptor does not support image requests")
}

func (a *WorkflowAdaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	return custom.DoRequestHelper(a, c, meta, requestBody)
}

func (a *WorkflowAdaptor) DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode) {
	// 读取响应体
	responseBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, openai.ErrorWrapper(readErr, "read_response_body_failed", http.StatusInternalServerError)
	}

	closeErr := resp.Body.Close()
	if closeErr != nil {
		return nil, openai.ErrorWrapper(closeErr, "close_response_body_failed", http.StatusInternalServerError)
	}

	// 解析工作流响应
	var workflowResp WorkflowResponse
	if unmarshalErr := json.Unmarshal(responseBody, &workflowResp); unmarshalErr != nil {
		logger.SysErrorf("Coze工作流响应解析失败: %v", unmarshalErr)
		return nil, openai.ErrorWrapper(unmarshalErr, "unmarshal_response_body_failed", http.StatusInternalServerError)
	}

	// 检查工作流执行状态
	if workflowResp.Code != 0 {
		logger.SysErrorf("Coze工作流执行失败 - Code: %d, Msg: %s", workflowResp.Code, workflowResp.Msg)
		return nil, &model.ErrorWithStatusCode{
			Error: model.Error{
				Message: workflowResp.Msg,
				Code:    workflowResp.Code,
			},
			StatusCode: resp.StatusCode,
		}
	}

	// 转换为OpenAI格式响应
	openaiResponse := a.convertToOpenAIResponse(&workflowResp, meta.ActualModelName)

	// 序列化响应
	jsonResponse, marshalErr := json.Marshal(openaiResponse)
	if marshalErr != nil {
		return nil, openai.ErrorWrapper(marshalErr, "marshal_response_body_failed", http.StatusInternalServerError)
	}

	// 写入响应
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, writeErr := c.Writer.Write(jsonResponse)
	if writeErr != nil {
		logger.SysErrorf("写入响应失败: %v", writeErr)
	}

	// 计算使用量
	var responseText string
	if len(openaiResponse.Choices) > 0 {
		responseText = openaiResponse.Choices[0].Message.StringContent()
	}

	usage = &model.Usage{
		PromptTokens:     openaiResponse.Usage.PromptTokens,
		CompletionTokens: openaiResponse.Usage.CompletionTokens,
		TotalTokens:      openaiResponse.Usage.TotalTokens,
	}

	logger.SysLogf("Coze工作流响应成功 - Token: %d, 输出长度: %d", workflowResp.Token, len(responseText))

	return usage, nil
}

func (a *WorkflowAdaptor) GetModelList() []string {
	return []string{} // 工作流模型列表由渠道配置决定
}

func (a *WorkflowAdaptor) GetChannelName() string {
	return "coze-workflow"
}

// extractWorkflowIDFromModel 从模型名称中提取工作流ID
func extractWorkflowIDFromModel(modelName string) string {
	// 方法1: 如果是标准的 workflow-{id} 格式
	if strings.HasPrefix(modelName, "workflow-") {
		return strings.TrimPrefix(modelName, "workflow-")
	}

	// 方法2: 直接使用模型名称作为工作流ID
	// 这适用于手动创建的 agent，其模型名称就是工作流ID
	return modelName
}

// convertToOpenAIResponse 将工作流响应转换为OpenAI格式
func (a *WorkflowAdaptor) convertToOpenAIResponse(workflowResp *WorkflowResponse, modelName string) *openai.TextResponse {
	// 解析 data 字段（可能是字符串或对象）
	var outputData map[string]interface{}
	var content string

	// 处理 data 字段
	if dataStr, ok := workflowResp.Data.(string); ok {
		// data 是字符串，需要解析 JSON
		if err := json.Unmarshal([]byte(dataStr), &outputData); err != nil {
			logger.SysErrorf("解析工作流 data 字段失败: %v", err)
			content = dataStr // 如果解析失败，直接使用原字符串
		} else {
			// 成功解析，按照统一标准提取结果
			content = extractUnifiedOutput(outputData)
		}
	} else if dataMap, ok := workflowResp.Data.(map[string]interface{}); ok {
		// data 是对象
		outputData = dataMap
		content = extractUnifiedOutput(outputData)
	} else {
		// 其他情况，转换为字符串
		if dataBytes, err := json.Marshal(workflowResp.Data); err == nil {
			content = string(dataBytes)
		} else {
			content = "工作流执行完成，但输出格式无法解析"
		}
	}

	// 使用 Coze 返回的 token 数量，如果没有则估算
	var promptTokens, completionTokens, totalTokens int
	if workflowResp.Token > 0 {
		totalTokens = workflowResp.Token
		promptTokens = totalTokens / 3         // 估算输入占 1/3
		completionTokens = totalTokens * 2 / 3 // 估算输出占 2/3
	} else {
		// 简单估算
		promptTokens = len(content) / 4
		completionTokens = len(content) / 4
		totalTokens = promptTokens + completionTokens
	}

	return &openai.TextResponse{
		Id:      fmt.Sprintf("workflow-exec-%s", workflowResp.ExecuteID),
		Object:  "chat.completion",
		Model:   modelName,
		Created: helper.GetTimestamp(),
		Choices: []openai.TextResponseChoice{{
			Index: 0,
			Message: model.Message{
				Role:    "assistant",
				Content: content,
			},
			FinishReason: "stop",
		}},
		Usage: model.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}
}

// extractUnifiedOutput 按照统一标准提取输出结果
// 统一标准：优先查找 "结果" 字段，然后是其他常见字段
func extractUnifiedOutput(outputData map[string]interface{}) string {
	// 按优先级查找输出字段
	priorityFields := []string{"结果", "output", "result", "answer", "response", "content", "text"}

	for _, field := range priorityFields {
		if value, exists := outputData[field]; exists {
			if strValue, ok := value.(string); ok && strValue != "" {
				return strValue
			}
		}
	}

	// 如果没有找到标准字段，将整个输出序列化为JSON
	if outputBytes, err := json.Marshal(outputData); err == nil {
		return string(outputBytes)
	}

	return "工作流执行完成，但输出格式无法解析"
}

// ConvertToWorkflowResponseData 将 Coze 工作流响应转换为统一的 WorkflowResponseData 格式
func (a *WorkflowAdaptor) ConvertToWorkflowResponseData(responseBody []byte) (*custom.WorkflowResponseData, error) {
	// 解析工作流响应
	var workflowResp WorkflowResponse
	if err := json.Unmarshal(responseBody, &workflowResp); err != nil {
		logger.SysErrorf("Coze工作流响应解析失败: %v", err)
		return nil, fmt.Errorf("响应解析失败: %v", err)
	}

	// 检查工作流执行状态
	if workflowResp.Code != 0 {
		logger.SysErrorf("Coze工作流执行失败 - Code: %d, Msg: %s", workflowResp.Code, workflowResp.Msg)
		return nil, fmt.Errorf("工作流执行失败 - Code: %d, Msg: %s", workflowResp.Code, workflowResp.Msg)
	}

	// 转换工作流输出数据
	workflowOutputData := a.parseWorkflowOutputData(workflowResp.Data)

	// 构造统一的响应格式
	responseData := &custom.WorkflowResponseData{
		WorkflowOutputData: workflowOutputData,
		ExecuteID:          workflowResp.ExecuteID,
	}

	logger.SysLogf("Coze工作流响应转换完成 - ExecuteID: %s, 输出字段数: %d",
		workflowResp.ExecuteID, len(workflowOutputData))

	return responseData, nil
}

// parseWorkflowOutputData 解析工作流输出数据为 key-value 格式
func (a *WorkflowAdaptor) parseWorkflowOutputData(rawData interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	if rawData == nil {
		logger.SysLogf("Coze工作流数据为空")
		result["output"] = ""
		return result
	}

	// 处理字符串格式的 JSON 数据
	if dataStr, ok := rawData.(string); ok {
		var parsedData map[string]interface{}
		if err := json.Unmarshal([]byte(dataStr), &parsedData); err != nil {
			logger.SysLogf("Coze工作流数据 JSON 解析失败: %v, 原始数据长度: %d", err, len(dataStr))
			// 如果解析失败，将原始字符串作为 "output" 字段返回
			result["output"] = dataStr
			return result
		}
		// 成功解析，直接返回解析后的数据
		return parsedData
	}

	// 处理已经是 map 格式的数据
	if dataMap, ok := rawData.(map[string]interface{}); ok {
		return dataMap
	}

	// 处理数组类型的数据
	if dataArray, ok := rawData.([]interface{}); ok {
		result["output"] = dataArray
		result["count"] = len(dataArray)
		logger.SysLogf("Coze工作流返回数组数据，长度: %d", len(dataArray))
		return result
	}

	// 处理其他类型的数据，尝试序列化后再解析
	if dataBytes, err := json.Marshal(rawData); err == nil {
		var parsedData map[string]interface{}
		if err := json.Unmarshal(dataBytes, &parsedData); err == nil {
			return parsedData
		}
		// 如果解析失败，将序列化后的字符串作为 "output" 字段返回
		result["output"] = string(dataBytes)
	} else {
		// 如果序列化也失败，将原始数据转换为字符串
		result["output"] = fmt.Sprintf("%v", rawData)
	}

	logger.SysLogf("Coze工作流数据转换为默认格式，输出字段数: %d", len(result))
	return result
}
