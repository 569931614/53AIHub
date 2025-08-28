package n8n

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/common/storage"
	db_model "github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/relay/meta"
)

type N8nWorkflowAdaptor struct {
	meta         *meta.Meta
	CustomConfig *custom.CustomConfig
}

func (a *N8nWorkflowAdaptor) Init(meta *meta.Meta) {
	a.meta = meta
}

func (a *N8nWorkflowAdaptor) GetRequestURL(meta *meta.Meta) (string, error) {
	baseURL := strings.TrimSuffix(meta.BaseURL, "/")

	if strings.Contains(baseURL, "/webhook-test/") || strings.Contains(baseURL, "/webhook/") {
		return baseURL, nil
	}

	// 从模型名称中提取工作流ID
	workflowID := extractWorkflowID(meta.ActualModelName)
	if workflowID == "" {
		return "", fmt.Errorf("无法从模型名称中提取工作流ID: %s", meta.ActualModelName)
	}

	// 智能判断使用 webhook 还是 webhook-test
	var webhookPath string
	if strings.HasSuffix(baseURL, "/webhook-test") {
		// 如果 baseURL 以 /webhook-test 结尾，使用 test 模式
		webhookPath = "webhook-test"
		// 去掉 baseURL 中的 /webhook-test 部分
		baseURL = strings.TrimSuffix(baseURL, "/webhook-test")
	} else if strings.HasSuffix(baseURL, "/webhook") {
		// 如果 baseURL 以 /webhook 结尾，使用正式模式
		webhookPath = "webhook"
		// 去掉 baseURL 中的 /webhook 部分
		baseURL = strings.TrimSuffix(baseURL, "/webhook")
	} else {
		// 默认使用正式模式
		webhookPath = "webhook"
	}

	// 构造 n8n webhook URL
	url := fmt.Sprintf("%s/%s/%s", baseURL, webhookPath, workflowID)

	logger.SysLogf("n8n工作流URL构造 - BaseURL: %s, WebhookPath: %s, WorkflowID: %s, 完整URL: %s",
		meta.BaseURL, webhookPath, workflowID, url)

	return url, nil
}

func (a *N8nWorkflowAdaptor) SetupRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) error {
	// n8n 特殊认证：直接使用 API Key，不加 Bearer 前缀
	req.Header.Set("Authorization", meta.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logger.SysLogf("n8n工作流请求头设置 - Authorization: %s", maskAPIKey(meta.APIKey))

	return nil
}

// ConvertWorkflowRequest 将工作流参数转换为 n8n 请求格式
func (a *N8nWorkflowAdaptor) ConvertWorkflowRequest(workflowID string, parameters map[string]interface{}) (*N8nWorkflowRequest, error) {
	logger.SysLogf("n8n工作流请求转换 - WorkflowID: %s, Parameters: %+v", workflowID, parameters)

	// 处理参数中的文件上传
	processedParameters, err := a.processWorkflowParameters(parameters)
	if err != nil {
		logger.SysErrorf("处理n8n工作流文件参数失败: %v", err)
		// 如果文件处理失败，使用原始参数继续执行
		processedParameters = parameters
	}

	// n8n 支持任意 JSON 参数，传递处理后的参数
	request := N8nWorkflowRequest(processedParameters)

	return &request, nil
}

// processWorkflowParameters 处理工作流参数中的文件上传
func (a *N8nWorkflowAdaptor) processWorkflowParameters(parameters map[string]interface{}) (map[string]interface{}, error) {
	if a.meta == nil {
		return parameters, fmt.Errorf("meta is nil")
	}

	processedParams := make(map[string]interface{})

	for key, value := range parameters {
		processedValue, err := a.processParameterValue(value)
		if err != nil {
			logger.SysErrorf("处理n8n参数 %s 失败: %v", key, err)
			// 如果单个参数处理失败，使用原始值
			processedParams[key] = value
		} else {
			processedParams[key] = processedValue
		}
	}

	return processedParams, nil
}

// processParameterValue 递归处理参数值，支持字符串、数组、对象
func (a *N8nWorkflowAdaptor) processParameterValue(value interface{}) (interface{}, error) {
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

// processFileIDString 处理 file_id: 格式的字符串，转换为 n8n base64 格式
func (a *N8nWorkflowAdaptor) processFileIDString(value string) (interface{}, error) {
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

	// 检查文件大小限制（n8n 建议限制在 10MB 以内）
	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if uploadFile.Size > maxFileSize {
		logger.SysErrorf("文件过大，n8n不支持: ID=%d, Size=%d bytes, 限制=%d bytes",
			fileID, uploadFile.Size, maxFileSize)
		return value, fmt.Errorf("文件过大，n8n限制为10MB，当前文件大小: %d bytes", uploadFile.Size)
	}

	// 读取文件内容
	fileContent, err := storage.StorageInstance.Load(uploadFile.Key)
	if err != nil {
		logger.SysErrorf("读取文件内容失败: ID=%d, Key=%s, error: %v", fileID, uploadFile.Key, err)
		return value, fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 转换为 base64
	base64Content := base64.StdEncoding.EncodeToString(fileContent)

	// 构造 n8n 文件对象格式
	fileObject := map[string]interface{}{
		"filename":   uploadFile.FileName,
		"mimeType":   uploadFile.MimeType,
		"data":       base64Content,
		"size":       uploadFile.Size,
		"extension":  uploadFile.Extension,
		"originalId": fileID, // 保留原始ID用于调试
	}

	logger.SysLogf("✅ n8n工作流文件处理成功 - 原始ID: %d, 文件名: %s, 大小: %d bytes, MIME: %s",
		fileID, uploadFile.FileName, uploadFile.Size, uploadFile.MimeType)

	return fileObject, nil
}

// DoRequest 执行 n8n 工作流请求
func (a *N8nWorkflowAdaptor) DoRequest(c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	// 读取请求体用于日志输出
	bodyBytes, err := io.ReadAll(requestBody)
	if err != nil {
		return nil, fmt.Errorf("读取请求体失败: %v", err)
	}

	// 获取请求URL
	url, err := a.GetRequestURL(meta)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	err = a.SetupRequestHeader(c, req, meta)
	if err != nil {
		return nil, err
	}

	// 格式化输出请求信息
	logger.SysLogf("🚀 n8n工作流请求开始")
	logger.SysLogf("┌─────────────────────────────────────────────────────────────")
	logger.SysLogf("│ 📡 请求URL: %s", url)
	logger.SysLogf("│ 🔑 API Key: %s", maskAPIKey(meta.APIKey))
	logger.SysLogf("│ 📝 请求方法: POST")
	logger.SysLogf("│ 📋 Content-Type: application/json")
	logger.SysLogf("├─────────────────────────────────────────────────────────────")
	logger.SysLogf("│ 📦 请求参数:")

	// 格式化输出请求体
	var requestData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &requestData); err == nil {
		prettyJSON, _ := json.MarshalIndent(requestData, "│   ", "  ")
		logger.SysLogf("│   %s", string(prettyJSON))
	} else {
		logger.SysLogf("│   %s", string(bodyBytes))
	}
	logger.SysLogf("└─────────────────────────────────────────────────────────────")

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.SysErrorf("❌ n8n工作流请求失败: %v", err)
		return nil, err
	}

	logger.SysLogf("✅ n8n工作流请求成功 - 状态码: %d", resp.StatusCode)
	return resp, nil
}

// ProcessResponse 处理 n8n 工作流响应
func (a *N8nWorkflowAdaptor) ProcessResponse(resp *http.Response) (*custom.WorkflowResponseData, error) {
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取n8n工作流响应失败: %v", err)
	}

	logger.SysLogf("📡 n8n工作流原始响应 - StatusCode: %d, 响应长度: %d bytes",
		resp.StatusCode, len(responseBody))

	// 解析 n8n 响应（数组格式）
	var n8nResponse N8nWorkflowResponse
	if err := json.Unmarshal(responseBody, &n8nResponse); err != nil {
		logger.SysErrorf("解析n8n工作流响应失败: %v, 响应内容: %s", err, string(responseBody))
		return nil, fmt.Errorf("解析n8n工作流响应失败: %v", err)
	}

	// 转换为统一的工作流响应格式
	workflowOutputData := a.convertN8nResponseToOutputData(n8nResponse)

	// 生成执行ID（使用时间戳或其他唯一标识）
	executeID := fmt.Sprintf("n8n-exec-%d", len(responseBody))

	workflowResponse := &custom.WorkflowResponseData{
		WorkflowOutputData: workflowOutputData,
		ExecuteID:          executeID,
		ChannelID:          a.meta.ChannelId,
		ModelName:          a.meta.OriginModelName,
	}

	logger.SysLogf("✅ n8n工作流处理完成 - ExecuteID: %s, 输出字段数: %d",
		executeID, len(workflowOutputData))

	return workflowResponse, nil
}

// convertN8nResponseToOutputData 将 n8n 数组响应转换为标准输出格式
func (a *N8nWorkflowAdaptor) convertN8nResponseToOutputData(n8nResponse N8nWorkflowResponse) map[string]interface{} {
	outputData := make(map[string]interface{})

	if len(n8nResponse) == 0 {
		logger.SysLogf("⚠️ n8n工作流返回空数组")
		outputData["output"] = ""
		return outputData
	}

	// 如果只有一个响应项，直接展开其字段
	if len(n8nResponse) == 1 {
		for key, value := range n8nResponse[0] {
			outputData[key] = value
		}
		logger.SysLogf("n8n工作流单项响应转换 - 字段数: %d", len(outputData))
		return outputData
	}

	// 如果有多个响应项，将整个数组作为 items 字段
	outputData["items"] = n8nResponse
	outputData["count"] = len(n8nResponse)

	// 尝试提取第一项的 output 字段作为主要输出
	if firstItem := n8nResponse[0]; firstItem != nil {
		if output, exists := firstItem["output"]; exists {
			outputData["output"] = output
		}
	}

	logger.SysLogf("n8n工作流多项响应转换 - 项目数: %d, 输出字段数: %d",
		len(n8nResponse), len(outputData))

	return outputData
}

// extractWorkflowID 从模型名称中提取工作流ID
func extractWorkflowID(modelName string) string {
	// 支持 workflow-{id} 格式
	if strings.HasPrefix(modelName, "workflow-") {
		return strings.TrimPrefix(modelName, "workflow-")
	}

	// 直接使用模型名称作为工作流ID
	return modelName
}

// maskAPIKey 遮蔽 API Key 用于日志输出
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return strings.Repeat("*", len(apiKey))
	}
	return apiKey[:4] + strings.Repeat("*", len(apiKey)-8) + apiKey[len(apiKey)-4:]
}
