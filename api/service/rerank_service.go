package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/meta"
	relay_model "github.com/songquanpeng/one-api/relay/model"
)

// BailianRerankService 处理百炼 rerank API 调用的服务
type BailianRerankService struct{}

// RerankRequest 与 controller 中的 RerankRequest 结构相同
type RerankRequest struct {
	Model           string   `json:"model" example:"gte-rerank-v2" binding:"required"`
	Query           string   `json:"query" example:"人工智能的发展历程" binding:"required"`
	Documents       []string `json:"documents" example:"[\"人工智能起源于1950年代，图灵提出了著名的图灵测试\",\"深度学习是机器学习的一个分支，使用神经网络进行学习\",\"自然语言处理是人工智能的重要应用领域之一\"]" binding:"required"`
	TopN            *int     `json:"top_n,omitempty" example:"3"`
	ReturnDocuments *bool    `json:"return_documents,omitempty" example:"true"`
}

// RerankResult represents a single rerank result
type RerankResult struct {
	Object         string          `json:"object" example:"rerank_result"` // Result object type
	Index          int             `json:"index" example:"0"`              // Original index in input documents
	RelevanceScore float64         `json:"relevance_score" example:"0.95"` // Relevance score (0-1)
	Document       *RerankDocument `json:"document,omitempty"`             // Document content (if return_documents=true)
}

// RerankDocument represents document content in rerank result
type RerankDocument struct {
	Text string `json:"text" example:"文档内容"` // Document text content
}

// RerankResponse 与 controller 中的 RerankResponse 结构相同
type RerankResponse struct {
	Object string         `json:"object" example:"list"`
	Data   []RerankResult `json:"data"`
	Model  string         `json:"model" example:"gte-rerank-v2"`
	Usage  RerankUsage    `json:"usage"`
}

// RerankUsage 与 controller 中的 RerankUsage 结构相同
type RerankUsage struct {
	TotalTokens int `json:"total_tokens" example:"150"`
}

// CallBailianRerankAPI 调用百炼 rerank API
func (s *BailianRerankService) CallBailianRerankAPI(ctx context.Context, req *RerankRequest, meta *meta.Meta) (*RerankResponse, *relay_model.Usage, error) {
	// 创建百炼适配器请求格式
	bailianReq := struct {
		Model      string   `json:"model"`
		Query      string   `json:"query"`
		Documents  []string `json:"documents"`
		TopN       *int     `json:"top_n,omitempty"`
		ReturnDocs *bool    `json:"return_documents,omitempty"`
	}{
		Model:      req.Model,
		Query:      req.Query,
		Documents:  req.Documents,
		TopN:       req.TopN,
		ReturnDocs: req.ReturnDocuments,
	}

	// 构建请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": bailianReq.Model,
		"input": map[string]interface{}{
			"query":     bailianReq.Query,
			"documents": bailianReq.Documents,
		},
		"parameters": func() map[string]interface{} {
			params := make(map[string]interface{})
			if bailianReq.TopN != nil {
				params["top_n"] = *bailianReq.TopN
			}
			if bailianReq.ReturnDocs != nil {
				params["return_documents"] = *bailianReq.ReturnDocs
			} else {
				params["return_documents"] = false
			}
			return params
		}(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	// 构建正确的 rerank API URL
	baseUrl := meta.BaseURL
	if baseUrl == "" {
		baseUrl = "https://dashscope.aliyuncs.com"
	}
	url := fmt.Sprintf("%s/api/v1/services/rerank/text-rerank/text-rerank", baseUrl)

	// 详细的请求日志
	logger.SysLogf("🚀 百炼Rerank API请求开始")
	logger.SysLogf("┌─────────────────────────────────────────────────────────────")
	logger.SysLogf("│ 📡 请求URL: %s", url)
	logger.SysLogf("│ 🔑 API Key: %s", maskAPIKey(meta.APIKey))
	logger.SysLogf("│ 🤖 模型名称: %s", req.Model)
	logger.SysLogf("│ 📝 请求方法: POST")
	logger.SysLogf("│ 📊 查询长度: %d 字符", len(req.Query))
	logger.SysLogf("│ 📚 文档数量: %d", len(req.Documents))
	if req.TopN != nil {
		logger.SysLogf("│ 🔢 TopN: %d", *req.TopN)
	}
	if req.ReturnDocuments != nil {
		logger.SysLogf("│ 📄 返回文档: %v", *req.ReturnDocuments)
	}
	logger.SysLogf("└─────────────────────────────────────────────────────────────")

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(requestBody))
	if err != nil {
		logger.SysErrorf("❌ 创建HTTP请求失败: %v", err)
		return nil, nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+meta.APIKey)

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		logger.SysErrorf("❌ 百炼Rerank请求失败: %v", err)
		return nil, nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	logger.SysLogf("✅ 百炼Rerank请求完成 - 状态码: %d", resp.StatusCode)

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.SysErrorf("❌ 百炼Rerank请求失败 - 状态码: %d, 响应: %s", resp.StatusCode, string(body))
		return nil, nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.SysErrorf("❌ 读取响应失败: %v", err)
		return nil, nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析百炼响应
	var bailianResponse map[string]interface{}
	if err := json.Unmarshal(responseBody, &bailianResponse); err != nil {
		logger.SysErrorf("❌ 解析响应失败: %v", err)
		return nil, nil, fmt.Errorf("解析响应失败: %v", err)
	}

	// 转换为标准格式
	return s.convertBailianRerankResponse(bailianResponse, req)
}

// convertBailianRerankResponse 转换百炼 rerank 响应为标准格式
func (s *BailianRerankService) convertBailianRerankResponse(bailianResp map[string]interface{}, req *RerankRequest) (*RerankResponse, *relay_model.Usage, error) {
	// 解析输出数据
	output, ok := bailianResp["output"].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("响应格式错误：缺少 output 字段")
	}

	results, ok := output["results"].([]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("响应格式错误：缺少 results 字段")
	}

	// 转换结果
	var rerankResults []RerankResult
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		index, _ := resultMap["index"].(float64)
		score, _ := resultMap["relevance_score"].(float64)

		rerankResult := RerankResult{
			Object:         "rerank_result",
			Index:          int(index),
			RelevanceScore: score,
		}

		// 如果需要返回文档内容
		if req.ReturnDocuments != nil && *req.ReturnDocuments {
			if int(index) < len(req.Documents) {
				rerankResult.Document = &RerankDocument{
					Text: req.Documents[int(index)],
				}
			}
		}

		rerankResults = append(rerankResults, rerankResult)
	}

	// 计算 token 使用量
	usage := s.calculateRerankUsage(req, len(rerankResults))

	response := &RerankResponse{
		Object: "list",
		Data:   rerankResults,
		Model:  req.Model,
		Usage: RerankUsage{
			TotalTokens: usage.TotalTokens,
		},
	}

	logger.SysLogf("✅ 响应转换完成 - 结果数量: %d, Token使用: %d", len(rerankResults), usage.TotalTokens)

	return response, usage, nil
}

// calculateRerankUsage 计算 rerank 的 token 使用量
func (s *BailianRerankService) calculateRerankUsage(req *RerankRequest, resultCount int) *relay_model.Usage {
	// 计算输入 token（query + documents）
	queryTokens := openai.CountTokenText(req.Query, req.Model)

	documentsText := strings.Join(req.Documents, " ")
	documentsTokens := openai.CountTokenText(documentsText, req.Model)

	promptTokens := queryTokens + documentsTokens

	// rerank 通常没有生成内容，completion tokens 为 0
	completionTokens := 0

	totalTokens := promptTokens + completionTokens

	logger.SysLogf("📊 Token计算详情 - Query: %d, Documents: %d, Total: %d",
		queryTokens, documentsTokens, totalTokens)

	return &relay_model.Usage{
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
	}
}

// maskAPIKey 遮蔽API密钥的敏感部分
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "****"
	}
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}
