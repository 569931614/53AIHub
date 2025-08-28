package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/53AI/53AIHub/common/ctxkey"
	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/common/session"
	"github.com/53AI/53AIHub/common/utils/helper"
	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/middleware"
	"github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service"
	adaptor53AI "github.com/53AI/53AIHub/service/hub_adaptor/53AI"
	"github.com/53AI/53AIHub/service/hub_adaptor/coze"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/53AI/53AIHub/service/hub_adaptor/dify"
	"github.com/53AI/53AIHub/service/hub_adaptor/fastgpt"
	"github.com/53AI/53AIHub/service/hub_adaptor/n8n"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	oneapi_model "github.com/songquanpeng/one-api/model"
	"github.com/songquanpeng/one-api/monitor"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/apitype"
	billing_ratio "github.com/songquanpeng/one-api/relay/billing/ratio"
	"github.com/songquanpeng/one-api/relay/channeltype"
	"github.com/songquanpeng/one-api/relay/constant/role"
	"github.com/songquanpeng/one-api/relay/controller"
	"github.com/songquanpeng/one-api/relay/controller/validator"
	"github.com/songquanpeng/one-api/relay/meta"
	relay_meta "github.com/songquanpeng/one-api/relay/meta"
	relay_model "github.com/songquanpeng/one-api/relay/model"
	"github.com/songquanpeng/one-api/relay/relaymode"
)

type Message struct {
	Role    string `json:"role" example:"user"`
	Content string `json:"content" example:"who are you"`
}

type ChatRequest struct {
	Messages         []Message `json:"messages"`
	Stream           bool      `json:"stream"`
	Model            string    `json:"model" example:"agent-6"`
	Temperature      float64   `json:"temperature,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	ConversationID   int64     `json:"conversation_id"`
}

// WorkflowRunRequest 工作流运行请求结构体
type WorkflowRunRequest struct {
	Parameters     map[string]interface{} `json:"parameters"`      // 工作流参数
	Stream         bool                   `json:"stream"`          // 是否流式响应（工作流不支持，会被忽略）
	Model          string                 `json:"model"`           // Agent模型
	ConversationID int64                  `json:"conversation_id"` // 会话ID
}

// @Summary Workflow Run
// @Description 工作流运行接口，返回标准格式的工作流执行结果
// @Tags Workflow
// @Accept json
// @Produce json
// @Param workflowRequest body WorkflowRunRequest true "WorkflowRunRequest"
// @Success 200 {object} model.CommonResponse{data=custom.WorkflowResponseData}
// @Router /v1/workflow/run [post]
// @Security BearerAuth
func WorkflowRun(c *gin.Context) {
	c.Set(ctxkey.Group, "vip")

	// 记录开始时间用于计算执行时长
	startTime := time.Now()
	c.Set("workflow_start_time", startTime)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, model.ParamError.ToResponse(errors.New("请求体读取失败")))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var workflowRequest WorkflowRunRequest
	if err := json.Unmarshal(body, &workflowRequest); err != nil {
		c.JSON(400, model.ParamError.ToResponse(errors.New("请求参数解析失败")))
		return
	}

	// 获取 agent 信息
	agent, err := GetSessionAgent(c)
	if err != nil {
		c.JSON(404, model.NotFound.ToResponse(errors.New("Agent 未找到")))
		return
	}

	// 验证是否为工作流类型的 agent
	if agent.AgentType != model.AgentTypeWorkflow {
		c.JSON(400, model.ParamError.ToResponse(errors.New("该 Agent 不是工作流类型")))
		return
	}

	// 检查流式响应参数
	if workflowRequest.Stream {
		logger.SysLogf("工作流请求设置了 stream=true，但工作流不支持流式响应，将忽略此参数")
	}

	logger.SysLogf("工作流运行请求 - Agent: %s, Stream: %v, Parameters: %+v",
		agent.Model, workflowRequest.Stream, workflowRequest.Parameters)

	// 执行工作流
	response, err := executeWorkflow(c, &workflowRequest, agent)
	if err != nil {
		logger.SysErrorf("工作流执行失败 - Agent: %s, Error: %v", agent.Model, err)

		// 根据错误类型返回不同的状态码
		statusCode := 500
		if strings.Contains(err.Error(), "参数") || strings.Contains(err.Error(), "输入") {
			statusCode = 400
		} else if strings.Contains(err.Error(), "未找到") || strings.Contains(err.Error(), "不存在") {
			statusCode = 404
		}

		// 根据状态码选择合适的响应码
		var responseCode model.ResponseCode
		switch statusCode {
		case 400:
			responseCode = model.ParamError
		case 404:
			responseCode = model.NotFound
		default:
			responseCode = model.SystemError
		}

		c.JSON(statusCode, responseCode.ToResponse(errors.New(err.Error())))
		return
	}

	logger.SysLogf("工作流执行成功 - Agent: %s, ExecuteID: %s",
		agent.Model, response.ExecuteID)

	// 保存工作流消息记录
	if err := saveWorkflowMessage(c, &workflowRequest, agent, response); err != nil {
		logger.SysErrorf("保存工作流消息失败: %v", err)
		// 不影响主流程，继续返回成功响应
	}

	// 使用标准响应格式返回
	c.JSON(200, model.Success.ToResponse(response))
}

func GetSessionAgent(c *gin.Context) (agent *model.Agent, err error) {
	sessionAgent, exists := c.Get(session.SESSION_AGENT)
	if !exists {
		return nil, errors.New("agent not found")
	}
	agent, ok := sessionAgent.(*model.Agent)
	if !ok {
		return nil, errors.New("agent not found")
	}
	return agent, nil
}

// extractWorkflowID 从 agent 配置中提取工作流ID
func extractWorkflowID(modelName, customConfig string) string {
	// 方法1: 如果模型名称已经是工作流ID格式，直接使用
	if strings.HasPrefix(modelName, "workflow-") {
		return strings.TrimPrefix(modelName, "workflow-")
	}

	// 方法2: 从 CustomConfig 中解析工作流ID
	if customConfig != "" {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(customConfig), &config); err == nil {
			if workflowID, ok := config["workflow_id"].(string); ok && workflowID != "" {
				return workflowID
			}
			// 也可能存储在其他字段中
			if workflowID, ok := config["bot_id"].(string); ok && workflowID != "" {
				return workflowID
			}
		}
	}

	// 方法3: 直接使用模型名称作为工作流ID（假设模型名称就是工作流ID）
	if modelName != "" {
		return modelName
	}

	return ""
}

func GetSessionConversation(c *gin.Context) (conversation *model.Conversation, err error) {
	sessionConversation, exists := c.Get(session.SESSION_CONVERSATION)
	if !exists {
		return nil, errors.New("conversation not found")
	}
	conversation, ok := sessionConversation.(*model.Conversation)
	if !ok {
		return nil, errors.New("conversation not found")
	}
	return conversation, nil
}

// @Summary Relay
// @Description Relay
// @Tags Relay
// @Accept json
// @Produce json
// @Param chatRequest body ChatRequest true "ChatRequest"
// @Success 500 {object} model.OpenAIErrorResponse
// @Router /v1/chat/completions [post]
// @Security BearerAuth
func Relay(c *gin.Context) {
	c.Set(ctxkey.Group, "vip")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, model.ParamError.ToOpenAIErrorRespone(nil))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 先获取 agent 信息来判断类型
	relayMode := relaymode.GetByPath(c.Request.URL.Path)
	agent, err := GetSessionAgent(c)
	if err != nil {
		c.JSON(404, model.NotFound.ToOpenAIErrorRespone(err))
		return
	}

	// 检查是否为工作流类型的 agent
	if agent.AgentType == model.AgentTypeWorkflow {
		// 工作流类型的 agent 不支持聊天接口，返回错误
		c.JSON(400, model.ParamError.ToOpenAIErrorRespone(errors.New("工作流类型的 Agent 请使用 /v1/workflow/run 接口")))
		return
	}

	// 处理普通聊天请求
	handleChatRequest(c, body, agent, relayMode)
}

// handleChatRequest 处理标准聊天请求
func handleChatRequest(c *gin.Context, body []byte, agent *model.Agent, relayMode int) {
	var chatRequest ChatRequest
	if err := json.Unmarshal(body, &chatRequest); err != nil {
		c.JSON(400, model.ParamError.ToOpenAIErrorRespone(nil))
		return
	}

	processChatRequest(c, &chatRequest, agent, relayMode)
}

// processChatRequest 处理聊天请求的通用逻辑
func processChatRequest(c *gin.Context, chatRequest *ChatRequest, agent *model.Agent, relayMode int) {
	ctx := c.Request.Context()

	// Handle object_string type messages
	// {"conversation_id":619,"frequency_penalty":0.5,"messages":[{"role":"user","content":"[{\"type\":\"text\",\"content\":\"解析这张图片\"},{\"type\":\"image\",\"content\":\"file_id:175\"}]"}],"model":"agent-56","presence_penalty":0.5,"stream":true,"temperature":0.2,"top_p":0.75}
	logger.SysLogf("Relay", "Relay", "RelayMode", relayMode, "Agent", agent)

	retryTimes := config.CHANNEL_RETRY_TIMES
	requestModel := agent.Model

	// 如果是工作流类型的 agent，需要转换模型名称格式
	if agent.AgentType == model.AgentTypeWorkflow {
		// 从 agent.Model 中提取工作流ID，或者使用 agent.Model 作为工作流ID
		// 格式转换：原模型名称 -> workflow-{workflow_id}
		workflowID := extractWorkflowID(agent.Model, agent.CustomConfig)
		if workflowID != "" {
			requestModel = fmt.Sprintf("workflow-%s", workflowID)
		}
		logger.SysLogf("工作流Agent - 原模型: %s, 转换后: %s", agent.Model, requestModel)
	}

	chatRequest.Model = requestModel

	// if 1o model, unset temperature, presence_penalty, frequency_penalty, top_p
	if agent.ChannelType == channeltype.OpenAI && strings.Contains(strings.ToLower(chatRequest.Model), "o1") {
		chatRequest.Temperature = 0
		chatRequest.PresencePenalty = 0
		chatRequest.FrequencyPenalty = 0
		chatRequest.TopP = 0
	}

	modifiedBody, err := json.Marshal(chatRequest)
	if err != nil {
		c.JSON(500, model.ParamError.ToOpenAIErrorRespone(nil))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
	logger.SysLogf("modifiedBody", string(modifiedBody))

	// bizErr := relayHelper(c, relayMode)
	// if bizErr == nil {
	// 	return
	// }

	var lastFailedChannelId int64
	for i := retryTimes; i > 0; i-- {
		// 使用新的服务函数获取渠道并检查/刷新token
		channel, err := service.GetChannelWithTokenRefresh(ctx, agent.Eid, agent.ChannelType, requestModel, lastFailedChannelId)
		if err != nil {
			logger.Errorf(ctx, "获取渠道失败: %s", err.Error())
			continue
		}

		middleware.SetupContextForSelectedChannel(c, channel, requestModel)
		logger.SysLogf("ChannelID", channel.ChannelID)
		channelId := c.GetInt64(ctxkey.ChannelId)
		lastFailedChannelId = channelId
		requestBody, err := common.GetRequestBody(c)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		bizErr := relayHelper(c, relayMode)
		if bizErr == nil {
			return
		}
		channelName := c.GetString(ctxkey.ChannelName)
		go processChannelRelayError(ctx, int(config.GetUserId(c)), int(channelId), channelName, *bizErr)
		// return error message
		c.JSON(500, model.OpenAIErrorResponse{
			Error: struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			}{
				Message: bizErr.Message,
				Type:    bizErr.Type,
			},
		})
		return
	}
	c.JSON(500, model.OpenAIErrorResponse{
		Error: struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		}{
			Message: "All channels are unavailable",
			Type:    "53aihub_error",
		},
	})
}

func relayHelper(c *gin.Context, relayMode int) *relay_model.ErrorWithStatusCode {
	var err *relay_model.ErrorWithStatusCode
	switch relayMode {
	case relaymode.ImagesGenerations:
		err = controller.RelayImageHelper(c, relayMode)
	// case relaymode.AudioSpeech:
	// 	fallthrough
	// case relaymode.AudioTranslation:
	// 	fallthrough
	// case relaymode.AudioTranscription:
	// 	err = controller.RelayAudioHelper(c, relayMode)
	// case relaymode.Proxy:
	// 	err = controller.RelayProxyHelper(c, relayMode)
	default:
		err = RelayTextHelper(c)
	}
	return err
}

func processChannelRelayError(ctx context.Context, userId int, channelId int, channelName string, err relay_model.ErrorWithStatusCode) {
	logger.Errorf(ctx, "relay error (channel id %d, user id: %d): %+v", channelId, userId, err.Error)
	if monitor.ShouldDisableChannel(&err.Error, err.StatusCode) {
		monitor.DisableChannel(channelId, channelName, err.Message)
	} else {
		monitor.Emit(channelId, false)
	}
}

func getAndValidateTextRequest(c *gin.Context, relayMode int) (*relay_model.GeneralOpenAIRequest, error) {
	textRequest := &relay_model.GeneralOpenAIRequest{}
	err := common.UnmarshalBodyReusable(c, textRequest)
	if err != nil {
		return nil, err
	}
	if relayMode == relaymode.Moderations && textRequest.Model == "" {
		textRequest.Model = "text-moderation-latest"
	}
	if relayMode == relaymode.Embeddings && textRequest.Model == "" {
		textRequest.Model = c.Param("model")
	}
	err = validator.ValidateTextRequest(textRequest, relayMode)
	if err != nil {
		return nil, err
	}
	return textRequest, nil
}

func getMappedModelName(modelName string, mapping map[string]string) (string, bool) {
	if mapping == nil {
		return modelName, false
	}
	mappedModelName := mapping[modelName]
	if mappedModelName != "" {
		return mappedModelName, true
	}
	return modelName, false
}

func setSystemPrompt(ctx context.Context, request *relay_model.GeneralOpenAIRequest, prompt string) (reset bool) {
	if prompt == "" {
		return false
	}
	if len(request.Messages) == 0 {
		return false
	}
	if request.Messages[0].Role == role.System {
		request.Messages[0].Content = prompt
		logger.Infof(ctx, "rewrite system prompt")
		return true
	}
	request.Messages = append([]relay_model.Message{{
		Role:    role.System,
		Content: prompt,
	}}, request.Messages...)
	logger.Infof(ctx, "add system prompt")
	return true
}

func getPromptTokens(textRequest *relay_model.GeneralOpenAIRequest, relayMode int) int {
	switch relayMode {
	case relaymode.ChatCompletions:
		return openai.CountTokenMessages(textRequest.Messages, textRequest.Model)
	case relaymode.Completions:
		return openai.CountTokenInput(textRequest.Prompt, textRequest.Model)
	case relaymode.Moderations:
		return openai.CountTokenInput(textRequest.Input, textRequest.Model)
	}
	return 0
}

func getRequestBody(c *gin.Context, meta *meta.Meta, textRequest *relay_model.GeneralOpenAIRequest, adaptor adaptor.Adaptor) (io.Reader, error) {
	if !config.EnforceIncludeUsage &&
		meta.APIType != apitype.OpenAI &&
		meta.OriginModelName == meta.ActualModelName &&
		meta.ChannelType != channeltype.Baichuan &&
		meta.SystemPrompt == "" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return nil, err
		}

		var requestData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
			return nil, err
		}

		delete(requestData, "conversation_id")
		modifiedBody, err := json.Marshal(requestData)
		if err != nil {
			return nil, err
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
		// return c.Request.Body, nil
	}

	// get request body
	var requestBody io.Reader
	convertedRequest, err := adaptor.ConvertRequest(c, meta.Mode, textRequest)
	if err != nil {
		logger.Debugf(c.Request.Context(), "converted request failed: %s\n", err.Error())
		return nil, err
	}

	if convertedRequestMap, ok := convertedRequest.(map[string]interface{}); ok {
		delete(convertedRequestMap, "conversation_id")
	}

	jsonData, err := json.Marshal(convertedRequest)
	if err != nil {
		logger.Debugf(c.Request.Context(), "converted request json_marshal_failed: %s\n", err.Error())
		return nil, err
	}
	logger.Debugf(c.Request.Context(), "converted request: \n%s", string(jsonData))
	requestBody = bytes.NewBuffer(jsonData)
	return requestBody, nil
}

func RelayTextHelper(c *gin.Context) *relay_model.ErrorWithStatusCode {
	ctx := c.Request.Context()
	user_id := config.GetUserId(c)
	meta := GetByContext(c)
	channelId := c.GetInt64(ctxkey.ChannelId)
	meta.ChannelId = int(channelId)
	meta.APIType = model.GetApiType(meta.ChannelType)
	startTime := time.Now()
	textRequest, err := getAndValidateTextRequest(c, meta.Mode)
	if err != nil {
		logger.Errorf(ctx, "getAndValidateTextRequest failed: %s", err.Error())
		return openai.ErrorWrapper(err, "invalid_text_request", http.StatusBadRequest)
	}
	meta.IsStream = textRequest.Stream

	if meta.IsStream {
		SetupStreamInterceptor(c)
	}

	// map model name
	meta.OriginModelName = textRequest.Model
	textRequest.Model, _ = getMappedModelName(textRequest.Model, meta.ModelMapping)
	meta.ActualModelName = textRequest.Model
	modelRatio := billing_ratio.GetModelRatio(textRequest.Model, meta.ChannelType)
	// groupRatio := billing_ratio.GetGroupRatio(meta.Group)
	groupRatio := 1.0
	ratio := modelRatio * groupRatio
	// set system prompt if not empty
	agent, err := GetSessionAgent(c)
	if err != nil {
		logger.Errorf(ctx, "getSessionAgent failed: %s", err.Error())
		return openai.ErrorWrapper(err, "invalid_text_request", http.StatusBadRequest)
	}
	systemPromptReset := false
	if agent.Prompt != "" {
		systemPromptReset = addAgentPrompt(ctx, textRequest, agent.Prompt, agent.ChannelType)
		modifiedBody, err := json.Marshal(textRequest)
		if err != nil {
			return openai.ErrorWrapper(err, "marshal_request_failed", http.StatusInternalServerError)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(modifiedBody))
	}

	promptTokens := getPromptTokens(textRequest, meta.Mode)
	meta.PromptTokens = promptTokens
	preConsumedQuota, bizErr := preConsumeQuota(ctx, textRequest, promptTokens, ratio, meta)
	if bizErr != nil {
		logger.Warnf(ctx, "preConsumeQuota failed: %+v", *bizErr)
		return bizErr
	}
	conversation, err := GetSessionConversation(c)
	if err != nil {
		logger.Errorf(ctx, "getSessionConversation failed: %s", err.Error())
		return openai.ErrorWrapper(err, "invalid_text_request", http.StatusBadRequest)
	}
	adaptor := service.GetAdaptor(meta.APIType)
	// set custom config
	if adaptor == nil {
		return openai.ErrorWrapper(fmt.Errorf("invalid api type: %d", meta.APIType), "invalid_api_type", http.StatusBadRequest)
	}
	adaptor.Init(meta)
	customConfig := &custom.CustomConfig{
		UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
		ConversationId:             conversation.ChannelConversationID,
		ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
		AIHubConversationId:        conversation.ConversationID,
	}

	err = service.SetCustomConfig(&adaptor, customConfig)

	if err != nil {
		return openai.ErrorWrapper(err, "convert_request_failed", http.StatusInternalServerError)
	}

	// get request body
	requestBody, err := getRequestBody(c, meta, textRequest, adaptor)
	if err != nil {
		return openai.ErrorWrapper(err, "convert_request_failed", http.StatusInternalServerError)
	}

	// do request
	resp, err := adaptor.DoRequest(c, meta, requestBody)
	if err != nil {
		logger.Errorf(ctx, "DoRequest failed: %s", err.Error())
		return openai.ErrorWrapper(err, "do_request_failed", http.StatusInternalServerError)
	}
	if isErrorHappened(meta, resp) {
		logger.SysErrorf("检测到错误响应 - StatusCode: %d, ContentType: %s, IsStream: %v, ChannelType: %d, ModelName: %s",
			resp.StatusCode, resp.Header.Get("Content-Type"), meta.IsStream, meta.ChannelType, meta.ActualModelName)
		//billing.ReturnPreConsumedQuota(ctx, preConsumedQuota, meta.TokenId)
		return controller.RelayErrorHandler(resp)
	}

	// do response
	usage, respErr := adaptor.DoResponse(c, resp, meta)
	logger.SysLogf("usage", usage)
	if respErr != nil {
		logger.Errorf(ctx, "respErr is not nil: %+v", respErr)
		//billing.ReturnPreConsumedQuota(ctx, preConsumedQuota, meta.TokenId)
		return respErr
	}

	responseContent, reasoningContent := GetResponseContent(c, meta.IsStream, resp)

	customConfig = service.GetCustomConfig(&adaptor)
	// post-consume quota
	go postConsumeQuota(c, agent, user_id, startTime, ctx, usage, meta,
		textRequest, ratio, preConsumedQuota, modelRatio, groupRatio,
		systemPromptReset, responseContent, reasoningContent, customConfig)
	return nil
}

func isErrorHappened(meta *meta.Meta, resp *http.Response) bool {
	if resp == nil {
		if meta.ChannelType == channeltype.AwsClaude {
			return false
		}
		return true
	}
	if resp.StatusCode != http.StatusOK &&
		// replicate return 201 to create a task
		resp.StatusCode != http.StatusCreated {
		return true
	}
	if meta.ChannelType == channeltype.DeepL {
		// skip stream check for deepl
		return false
	}

	// Coze 工作流特殊处理：工作流不支持流式响应，总是返回 JSON
	if meta.ChannelType == channeltype.Coze && strings.HasPrefix(meta.ActualModelName, "workflow-") {
		// 对于 Coze 工作流，即使设置了 stream=true，也会返回 JSON 格式，这是正常的
		return false
	}

	if meta.IsStream && strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") &&
		// Even if stream mode is enabled, replicate will first return a task info in JSON format,
		// requiring the client to request the stream endpoint in the task info
		meta.ChannelType != channeltype.Replicate {
		return true
	}
	return false
}

func addAgentPrompt(ctx context.Context, textRequest *relay_model.GeneralOpenAIRequest, agentPrompt string, channelType int) bool {
	if agentPrompt == "" {
		return false
	}

	if channelType == channeltype.OpenAI && strings.Contains(strings.ToLower(textRequest.Model), "o1") {
		textRequest.Messages = append([]relay_model.Message{{
			Role:    "user",
			Content: agentPrompt,
		}}, textRequest.Messages...)
		logger.Infof(ctx, "add agent prompt for o1 model")
		return true
	}

	if len(textRequest.Messages) > 0 && textRequest.Messages[0].Role == role.System {
		textRequest.Messages[0].Content = fmt.Sprintf("%s\n%s", agentPrompt, textRequest.Messages[0].Content)
		logger.Infof(ctx, "append agent system prompt")
	} else {
		textRequest.Messages = append([]relay_model.Message{{
			Role:    role.System,
			Content: agentPrompt,
		}}, textRequest.Messages...)
		logger.Infof(ctx, "add agent system prompt")
	}
	return true
}

func postConsumeQuota(c *gin.Context, agent *model.Agent, user_id int64, startTime time.Time,
	ctx context.Context, usage *relay_model.Usage, meta *meta.Meta, textRequest *relay_model.GeneralOpenAIRequest,
	ratio float64, preConsumedQuota int64, modelRatio float64,
	groupRatio float64, systemPromptReset bool, responseContent string, reasoningContent string, customConfig *custom.CustomConfig) {
	if usage == nil {
		logger.Error(ctx, "usage is nil, which is unexpected")
		return
	}
	var quota int64
	completionRatio := billing_ratio.GetCompletionRatio(textRequest.Model, meta.ChannelType)
	promptTokens := usage.PromptTokens
	completionTokens := usage.CompletionTokens
	quota = int64(math.Ceil((float64(promptTokens) + float64(completionTokens)*completionRatio) * ratio))
	if ratio != 0 && quota <= 0 {
		quota = 1
	}
	totalTokens := promptTokens + completionTokens
	if totalTokens == 0 {
		// in this case, must be some error happened
		// we cannot just return, because we may have to return the pre-consumed quota
		quota = 0
	}
	quotaDelta := quota - preConsumedQuota

	logContent := fmt.Sprintf("倍率：%.2f × %.2f × %.2f", modelRatio, groupRatio, completionRatio)

	messageJSON, err := json.Marshal(textRequest.Messages)
	if err != nil {
		logger.Errorf(ctx, "marshal messages failed: %s", err.Error())
		messageJSON = []byte("[]")
	}
	requestId := helper.GetRequestID(ctx)
	conversationId := c.GetInt64(session.SESSION_CONVERSATION_ID)
	message := &model.Message{
		Eid:              agent.Eid,
		UserID:           user_id,
		ConversationID:   conversationId,
		AgentID:          agent.AgentID,
		Message:          string(messageJSON),
		Answer:           responseContent,
		ReasoningContent: reasoningContent,
		ModelName:        textRequest.Model,
		Quota:            int(quotaDelta),
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokens,
		ChannelId:        int(meta.ChannelId),
		RequestId:        requestId,
		ElapsedTime:      helper.CalcElapsedTime(startTime),
		IsStream:         meta.IsStream,
		QuotaContent:     logContent,
	}
	model.CreateMessage(message)
	// conversation update
	if conversationId != 0 {
		conversation, err := model.GetConversationByIdAndUserId(agent.Eid, conversationId, user_id)
		if err != nil {
			logger.Errorf(ctx, "get conversation by id and user id failed: %s", err.Error())
		} else {
			lastMessage, _ := json.Marshal(map[string]string{
				"question": string(messageJSON),
				"answer":   responseContent,
			})

			conversation.Quota += int(quotaDelta)
			conversation.TotalTokens += totalTokens
			conversation.LastMessage = string(lastMessage)
			if customConfig != nil {
				if customConfig.ConversationId != "" {
					conversation.ChannelConversationID = customConfig.ConversationId
				}
				if customConfig.ConversationExpirationTime != 0 {
					conversation.ChannelConversationExpirationTime = customConfig.ConversationExpirationTime
				}
			}

			model.UpdateConversation(conversation)
		}
	}
}

func preConsumeQuota(ctx context.Context, textRequest *relay_model.GeneralOpenAIRequest, promptTokens int, ratio float64, meta *meta.Meta) (int64, *relay_model.ErrorWithStatusCode) {
	preConsumedQuota := getPreConsumedQuota(textRequest, promptTokens, ratio)
	return preConsumedQuota, nil
}

func getPreConsumedQuota(textRequest *relay_model.GeneralOpenAIRequest, promptTokens int, ratio float64) int64 {
	preConsumedTokens := config.PreConsumedQuota + int64(promptTokens)
	if textRequest.MaxTokens != 0 {
		preConsumedTokens += int64(textRequest.MaxTokens)
	}
	return int64(float64(preConsumedTokens) * ratio)
}

// executeWorkflow 执行工作流并返回标准响应数据
func executeWorkflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent) (*custom.WorkflowResponseData, error) {
	// 验证 parameters 参数
	if len(workflowRequest.Parameters) == 0 {
		return nil, errors.New("工作流参数不能为空")
	}

	logger.SysLogf("工作流执行开始 - Model: %s, ConversationID: %d, Parameters: %+v",
		workflowRequest.Model, workflowRequest.ConversationID, workflowRequest.Parameters)

	modelName := agent.Model
	// 获取渠道并检查/刷新token
	logger.SysLogf("工作流执行 - 开始获取渠道，Eid: %d, ChannelType: %d, Model: %s",
		agent.Eid, agent.ChannelType, modelName)

	// 使用新的服务函数获取渠道并检查/刷新token
	ctx := c.Request.Context()
	channel, err := service.GetChannelWithTokenRefresh(ctx, agent.Eid, agent.ChannelType, modelName, 0)
	if err != nil {
		// 如果是Coze渠道，尝试使用备用方法获取渠道
		if agent.ChannelType == channeltype.Coze {
			channel, err = model.GetFirstChannelByEidAndProviderType(agent.Eid, channeltype.Coze)
			if err != nil || channel == nil {
				return nil, fmt.Errorf("provider channel error")
			}
			channel.Models = channel.GetAddModelString(modelName)
			err := model.DB.Updates(channel).Error
			if err != nil {
				return nil, fmt.Errorf("update channel error")
			}
		} else {
			return nil, fmt.Errorf("获取渠道失败，Eid: %d, ChannelType: %d, Model: %s, Error: %v",
				agent.Eid, agent.ChannelType, modelName, err)
		}
	}

	logger.SysLogf("工作流执行 - 成功获取渠道，ChannelID: %d, BaseURL: %s",
		channel.ChannelID, channel.BaseURL)

	// 设置渠道上下文
	middleware.SetupContextForSelectedChannel(c, channel, modelName)

	// 直接调用工作流适配器执行
	return executeWorkflowDirect(c, workflowRequest, agent, channel, modelName)
}

// executeWorkflowDirect 直接执行工作流，简化参数传递
func executeWorkflowDirect(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 根据渠道类型选择对应的工作流适配器
	if channel.Type == channeltype.Coze || channel.Type == model.ChannelApiTypeCozeStudio {
		return executeCozeWorkflow(c, workflowRequest, agent, channel, modelName)
	}

	if channel.Type == model.ChannelApiDify {
		return executeDifyWorkflow(c, workflowRequest, agent, channel, modelName)
	}

	if channel.Type == channeltype.FastGPT || channel.Type == model.ChannelApiTypeFastGpt {
		// 这个 fastgpt 是因为 适配器的默认给 0 了，所以这里要手动设置一下，实际上数据库里面不会存 1007
		return executeFastGPTWorkflow(c, workflowRequest, agent, channel, modelName)
	}

	if channel.Type == model.ChannelApi53AI {
		return executeAI53Workflow(c, workflowRequest, agent, channel, modelName)
	}

	if channel.Type == model.ChannelApiTypeN8n {
		return executeN8nWorkflow(c, workflowRequest, agent, channel, modelName)
	}

	return nil, fmt.Errorf("不支持的渠道类型: %d", channel.Type)
}

// handleWorkflowError 处理工作流HTTP错误响应
func handleWorkflowError(resp *http.Response, workflowType string) error {
	// 读取错误响应内容
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	logger.SysErrorf("❌ %s工作流请求失败 - 状态码: %d, 响应: %s",
		workflowType, resp.StatusCode, string(body))

	// 尝试解析错误响应获取详细信息
	var detailMsg string
	if len(body) > 0 {
		// 尝试解析 JSON 错误响应
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			// 优先获取 message 字段
			if msg, ok := errorResp["message"].(string); ok && msg != "" {
				detailMsg = msg
			} else if msg, ok := errorResp["error"].(string); ok && msg != "" {
				detailMsg = msg
			} else if msg, ok := errorResp["detail"].(string); ok && msg != "" {
				detailMsg = msg
			}
		}

		// 如果无法解析 JSON，使用原始响应内容
		if detailMsg == "" {
			detailMsg = string(body)
		}
	}

	// 构建错误消息
	var errorMsg string
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		if detailMsg != "" {
			errorMsg = fmt.Sprintf("%s工作流请求参数错误：%s", workflowType, detailMsg)
		} else {
			errorMsg = fmt.Sprintf("%s工作流请求参数错误，状态码: %d", workflowType, resp.StatusCode)
		}
	} else {
		if detailMsg != "" {
			errorMsg = fmt.Sprintf("%s工作流请求失败：%s", workflowType, detailMsg)
		} else {
			errorMsg = fmt.Sprintf("%s工作流请求失败，状态码: %d", workflowType, resp.StatusCode)
		}
	}

	return errors.New(errorMsg)
}

// executeCozeWorkflow 执行 Coze 工作流
func executeCozeWorkflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 获取元数据
	meta := GetByContext(c)
	meta.APIType = model.GetApiType(channel.Type)
	meta.OriginModelName = modelName
	meta.ChannelId = int(channel.ChannelID)
	if channel.BaseURL != nil {
		meta.BaseURL = *channel.BaseURL
	}
	meta.APIKey = channel.Key

	// 应用模型映射
	mappedModel, _ := getMappedModelName(modelName, meta.ModelMapping)
	meta.ActualModelName = mappedModel

	logger.SysLogf("Coze工作流执行 - 模型映射，OriginModel: %s, ActualModel: %s",
		meta.OriginModelName, meta.ActualModelName)

	// 创建工作流适配器
	workflowAdaptor := &coze.WorkflowAdaptor{}
	workflowAdaptor.Init(meta)

	// 设置自定义配置
	user_id := config.GetUserId(c)
	conversation, err := GetSessionConversation(c)
	if err == nil {
		customConfig := &custom.CustomConfig{
			UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
			ConversationId:             conversation.ChannelConversationID,
			ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
			AIHubConversationId:        conversation.ConversationID,
		}
		workflowAdaptor.CustomConfig = customConfig
	}

	// 构建工作流请求
	workflowID := extractWorkflowID(agent.Model, agent.CustomConfig)
	if workflowID == "" {
		return nil, fmt.Errorf("无法提取工作流ID")
	}

	// 使用简化的方法直接构造请求
	cozeRequest := workflowAdaptor.ConvertWorkflowRequest(workflowID, workflowRequest.Parameters)

	// 序列化请求
	requestBody, err := json.Marshal(cozeRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化工作流请求失败: %v", err)
	}

	// 执行请求
	resp, err := workflowAdaptor.DoRequest(c, meta, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("执行工作流请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		return nil, handleWorkflowError(resp, "Coze")
	}

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取工作流响应失败: %v", err)
	}

	logger.SysLogf("Coze工作流原始响应 - StatusCode: %d, 响应长度: %d bytes",
		resp.StatusCode, len(responseBody))

	// 转换响应
	workflowResponse, err := workflowAdaptor.ConvertToWorkflowResponseData(responseBody)
	if err != nil {
		return nil, fmt.Errorf("转换工作流响应失败: %v", err)
	}

	// 设置响应信息
	workflowResponse.ChannelID = int(channel.ChannelID)
	workflowResponse.ModelName = agent.Model

	if len(workflowResponse.WorkflowOutputData) == 0 {
		logger.SysLogf("⚠️ Coze工作流执行成功但输出字段为空 - ExecuteID: %s", workflowResponse.ExecuteID)
		logger.SysLogf("🔍 Coze工作流详细输出数据: %+v", workflowResponse)
	} else {
		logger.SysLogf("✅ Coze工作流执行成功 - ExecuteID: %s, 输出字段数: %d",
			workflowResponse.ExecuteID, len(workflowResponse.WorkflowOutputData))
	}

	return workflowResponse, nil
}

// executeDifyWorkflow 执行 DIFY 工作流
func executeDifyWorkflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 获取元数据
	meta := GetByContext(c)
	meta.APIType = model.GetApiType(channel.Type)
	meta.OriginModelName = modelName
	meta.ChannelId = int(channel.ChannelID)
	if channel.BaseURL != nil {
		meta.BaseURL = *channel.BaseURL
	}
	meta.APIKey = channel.Key

	// 应用模型映射
	mappedModel, _ := getMappedModelName(modelName, meta.ModelMapping)
	meta.ActualModelName = mappedModel

	logger.SysLogf("DIFY工作流执行 - 模型映射，OriginModel: %s, ActualModel: %s",
		meta.OriginModelName, meta.ActualModelName)

	// 创建工作流适配器
	workflowAdaptor := &dify.DifyWorkflowAdaptor{}
	workflowAdaptor.Init(meta)

	// 设置自定义配置
	user_id := config.GetUserId(c)
	conversation, err := GetSessionConversation(c)
	if err == nil {
		customConfig := &custom.CustomConfig{
			UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
			ConversationId:             conversation.ChannelConversationID,
			ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
			AIHubConversationId:        conversation.ConversationID,
		}
		workflowAdaptor.CustomConfig = customConfig
	}

	// 构建工作流请求
	workflowID := extractWorkflowID(agent.Model, agent.CustomConfig)
	if workflowID == "" {
		return nil, fmt.Errorf("无法提取工作流ID")
	}

	// 使用简化的方法直接构造请求
	difyRequest, err := workflowAdaptor.ConvertWorkflowRequest(workflowID, workflowRequest.Parameters)
	if err != nil {
		return nil, fmt.Errorf("构造DIFY工作流请求失败: %v", err)
	}

	// 序列化请求
	requestBody, err := json.Marshal(difyRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化DIFY工作流请求失败: %v", err)
	}

	// 执行请求
	resp, err := workflowAdaptor.DoRequest(c, meta, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("执行DIFY工作流请求失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		return nil, handleWorkflowError(resp, "DIFY")
	}

	// 处理流式响应
	workflowResponse, err := workflowAdaptor.ProcessStreamingResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("处理DIFY工作流响应失败: %v", err)
	}

	// 设置响应信息
	workflowResponse.ChannelID = int(channel.ChannelID)
	workflowResponse.ModelName = agent.Model

	if len(workflowResponse.WorkflowOutputData) == 0 {
		logger.SysLogf("⚠️ DIFY工作流执行成功但输出字段为空 - ExecuteID: %s", workflowResponse.ExecuteID)
		logger.SysLogf("🔍 DIFY工作流详细输出数据: %+v", workflowResponse)
	} else {
		logger.SysLogf("✅ DIFY工作流执行成功 - ExecuteID: %s, 输出字段数: %d",
			workflowResponse.ExecuteID, len(workflowResponse.WorkflowOutputData))
	}

	return workflowResponse, nil
}

// executeFastGPTWorkflow 执行 FastGPT 工作流
func executeFastGPTWorkflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 检查 Agent 类型是否为工作流类型
	if agent.AgentType != model.AgentTypeWorkflow {
		return nil, fmt.Errorf("Agent 类型不是工作流类型，当前类型: %d", agent.AgentType)
	}

	// 获取元数据
	meta := GetByContext(c)
	meta.APIType = model.GetApiType(channel.Type)
	meta.OriginModelName = modelName
	meta.ChannelId = int(channel.ChannelID)
	if channel.BaseURL != nil {
		meta.BaseURL = *channel.BaseURL
	}
	meta.APIKey = channel.Key

	// 应用模型映射
	mappedModel, _ := getMappedModelName(modelName, meta.ModelMapping)
	meta.ActualModelName = mappedModel

	logger.SysLogf("FastGPT工作流执行 - 模型映射，OriginModel: %s, ActualModel: %s",
		meta.OriginModelName, meta.ActualModelName)

	// 创建工作流适配器
	workflowAdaptor := &fastgpt.FastGPTWorkflowAdaptor{}
	workflowAdaptor.Init(meta)

	// 设置自定义配置
	user_id := config.GetUserId(c)
	conversation, err := GetSessionConversation(c)
	if err == nil {
		customConfig := &custom.CustomConfig{
			UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
			ConversationId:             conversation.ChannelConversationID,
			ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
			AIHubConversationId:        conversation.ConversationID,
		}
		workflowAdaptor.CustomConfig = customConfig
	}

	// 转换工作流请求为 FastGPT 工作流请求
	fastgptRequest, err := workflowAdaptor.ConvertWorkflowToWorkflowRequest(workflowRequest.Parameters)
	if err != nil {
		return nil, fmt.Errorf("转换FastGPT工作流请求失败: %v", err)
	}

	// 序列化请求
	requestBody, err := json.Marshal(fastgptRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化FastGPT工作流请求失败: %v", err)
	}

	// 执行请求
	resp, err := workflowAdaptor.DoRequest(c, meta, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("执行FastGPT工作流请求失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		return nil, handleWorkflowError(resp, "FastGPT")
	}

	// 处理响应
	workflowResponse, err := workflowAdaptor.ProcessWorkflowResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("处理FastGPT工作流响应失败: %v", err)
	}

	// 设置响应信息
	workflowResponse.ChannelID = int(channel.ChannelID)
	workflowResponse.ModelName = agent.Model

	if len(workflowResponse.WorkflowOutputData) == 0 {
		logger.SysLogf("⚠️ FastGPT工作流执行成功但输出字段为空 - ExecuteID: %s", workflowResponse.ExecuteID)
		logger.SysLogf("🔍 FastGPT工作流详细输出数据: %+v", workflowResponse)
	} else {
		logger.SysLogf("✅ FastGPT工作流执行成功 - ExecuteID: %s, 输出字段数: %d",
			workflowResponse.ExecuteID, len(workflowResponse.WorkflowOutputData))
	}

	return workflowResponse, nil
}

// executeAI53Workflow 执行 53AI 工作流
func executeAI53Workflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 检查 Agent 类型是否为工作流类型
	if agent.AgentType != model.AgentTypeWorkflow {
		return nil, fmt.Errorf("Agent 类型不是工作流类型，当前类型: %d", agent.AgentType)
	}

	// 获取元数据
	meta := GetByContext(c)
	meta.APIType = model.GetApiType(channel.Type)
	meta.OriginModelName = modelName
	meta.ChannelId = int(channel.ChannelID)
	if channel.BaseURL != nil {
		meta.BaseURL = *channel.BaseURL
	}
	meta.APIKey = channel.Key

	// 应用模型映射
	mappedModel, _ := getMappedModelName(modelName, meta.ModelMapping)
	meta.ActualModelName = mappedModel

	logger.SysLogf("53AI工作流执行 - 模型映射，OriginModel: %s, ActualModel: %s",
		meta.OriginModelName, meta.ActualModelName)

	// 创建工作流适配器
	workflowAdaptor := &adaptor53AI.AI53WorkflowAdaptor{}
	workflowAdaptor.Init(meta)

	// 设置自定义配置
	user_id := config.GetUserId(c)
	conversation, err := GetSessionConversation(c)
	if err == nil {
		customConfig := &custom.CustomConfig{
			UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
			ConversationId:             conversation.ChannelConversationID,
			ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
			AIHubConversationId:        conversation.ConversationID,
		}
		workflowAdaptor.CustomConfig = customConfig
	}

	// 转换工作流请求为 53AI 工作流请求
	ai53Request, err := workflowAdaptor.ConvertWorkflowToAI53Request(workflowRequest.Parameters)
	if err != nil {
		return nil, fmt.Errorf("转换53AI工作流请求失败: %v", err)
	}

	// 序列化请求
	requestBody, err := json.Marshal(ai53Request)
	if err != nil {
		return nil, fmt.Errorf("序列化53AI工作流请求失败: %v", err)
	}

	// 执行请求
	resp, err := workflowAdaptor.DoRequest(c, meta, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("执行53AI工作流请求失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		return nil, handleWorkflowError(resp, "53AI")
	}

	// 处理响应
	workflowResponse, err := workflowAdaptor.ProcessAI53WorkflowResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("处理53AI工作流响应失败: %v", err)
	}

	// 设置响应信息
	workflowResponse.ChannelID = int(channel.ChannelID)
	workflowResponse.ModelName = agent.Model

	if len(workflowResponse.WorkflowOutputData) == 0 {
		logger.SysLogf("⚠️ 53AI工作流执行成功但输出字段为空 - TaskID: %s", workflowResponse.ExecuteID)
		logger.SysLogf("🔍 53AI工作流详细输出数据: %+v", workflowResponse)
	} else {
		logger.SysLogf("✅ 53AI工作流执行成功 - TaskID: %s, 输出字段数: %d",
			workflowResponse.ExecuteID, len(workflowResponse.WorkflowOutputData))
	}

	return workflowResponse, nil
}

// saveWorkflowMessage 保存工作流消息记录
func saveWorkflowMessage(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, response *custom.WorkflowResponseData) error {
	ctx := c.Request.Context()

	// 获取用户信息
	userId := config.GetUserId(c)
	if userId == 0 {
		return errors.New("用户ID获取失败")
	}

	// 获取会话ID
	conversationId := workflowRequest.ConversationID
	if conversationId == 0 {
		conversationId = c.GetInt64(session.SESSION_CONVERSATION_ID)
	}

	// 序列化工作流参数作为 message 内容
	parametersJSON, err := json.Marshal(workflowRequest.Parameters)
	if err != nil {
		logger.SysErrorf("序列化工作流参数失败: %v", err)
		parametersJSON = []byte("{}")
	}

	// 序列化工作流输出数据作为 answer 内容
	outputDataJSON, err := json.Marshal(response.WorkflowOutputData)
	if err != nil {
		logger.SysErrorf("序列化工作流输出数据失败: %v", err)
		outputDataJSON = []byte("{}")
	}

	// 获取请求ID
	requestId := helper.GetRequestID(ctx)
	if requestId == "" {
		requestId = response.ExecuteID // 使用 execute_id 作为 request_id
	}

	// 计算执行时间（如果有开始时间记录）
	var elapsedTime int64 = 0
	if startTimeValue, exists := c.Get("workflow_start_time"); exists {
		if startTime, ok := startTimeValue.(time.Time); ok {
			elapsedTime = helper.CalcElapsedTime(startTime)
		}
	}

	// 计算 token 消耗
	promptTokens, completionTokens, totalTokens := calculateWorkflowTokens(workflowRequest, response)

	// 获取费率信息（复用 chat 的费率计算逻辑）
	channelType := getWorkflowChannelType(response)
	modelRatio := billing_ratio.GetModelRatio(workflowRequest.Model, channelType)
	groupRatio := 1.0 // 与 chat 保持一致
	completionRatio := billing_ratio.GetCompletionRatio(workflowRequest.Model, channelType)
	ratio := modelRatio * groupRatio

	// 计算配额（复用 chat 的配额计算公式）
	quota := int64(math.Ceil((float64(promptTokens) + float64(completionTokens)*completionRatio) * ratio))
	if ratio != 0 && quota <= 0 {
		quota = 1 // 边界情况处理，与 chat 保持一致
	}

	// 生成配额内容记录（复用 chat 的格式）
	quotaContent := fmt.Sprintf("倍率：%.2f × %.2f × %.2f", modelRatio, groupRatio, completionRatio)

	// 创建消息记录
	message := &model.Message{
		Eid:               agent.Eid,
		UserID:            userId,
		ConversationID:    conversationId,
		AgentID:           agent.AgentID,
		Message:           string(parametersJSON), // 存储 parameters 的 JSON
		Answer:            string(outputDataJSON), // 存储 workflow_output_data 的 JSON
		ReasoningContent:  "",                     // 工作流暂不支持推理内容
		ModelName:         response.ModelName,
		Quota:             int(quota),
		PromptTokens:      promptTokens,
		CompletionTokens:  completionTokens,
		TotalTokens:       totalTokens,
		ChannelId:         response.ChannelID,
		RequestId:         requestId,
		ElapsedTime:       elapsedTime,
		IsStream:          false, // 工作流不支持流式
		QuotaContent:      quotaContent,
		AgentCustomConfig: agent.CustomConfig, // 历史记录
	}

	// 保存消息到数据库
	if err := model.CreateMessage(message); err != nil {
		return fmt.Errorf("创建消息记录失败: %v", err)
	}

	logger.SysLogf("工作流消息保存成功 - MessageID: %d, ExecuteID: %s", message.ID, response.ExecuteID)

	// 更新会话的最后消息（如果有会话ID）
	if conversationId != 0 {
		if err := updateConversationLastMessage(agent.Eid, conversationId, userId, string(parametersJSON), string(outputDataJSON), int(quota), totalTokens); err != nil {
			logger.SysErrorf("更新会话最后消息失败: %v", err)
			// 不返回错误，不影响主流程
		}
	}

	return nil
}

// updateConversationLastMessage 更新会话的最后消息和配额统计
func updateConversationLastMessage(eid, conversationId, userId int64, question, answer string, quota, totalTokens int) error {
	conversation, err := model.GetConversationByIdAndUserId(eid, conversationId, userId)
	if err != nil {
		return fmt.Errorf("获取会话失败: %v", err)
	}

	// 构造最后消息的 JSON 格式（与 chat 类型保持一致）
	lastMessage, err := json.Marshal(map[string]string{
		"question": question,
		"answer":   answer,
	})
	if err != nil {
		return fmt.Errorf("序列化最后消息失败: %v", err)
	}

	// 更新会话（复用 chat 的会话统计逻辑）
	conversation.LastMessage = string(lastMessage)
	conversation.Quota += quota
	conversation.TotalTokens += totalTokens
	if err := model.UpdateConversation(conversation); err != nil {
		return fmt.Errorf("更新会话失败: %v", err)
	}

	return nil
}

// calculateWorkflowTokens 计算工作流的 token 消耗
func calculateWorkflowTokens(workflowRequest *WorkflowRunRequest, response *custom.WorkflowResponseData) (promptTokens, completionTokens, totalTokens int) {
	// 计算输入 token（基于 parameters）
	parametersText := ""
	for key, value := range workflowRequest.Parameters {
		if strValue, ok := value.(string); ok {
			parametersText += key + ":" + strValue + " "
		} else {
			// 对于非字符串类型，序列化为 JSON
			if jsonBytes, err := json.Marshal(value); err == nil {
				parametersText += key + ":" + string(jsonBytes) + " "
			}
		}
	}

	// 使用现有的 token 计算逻辑
	promptTokens = openai.CountTokenText(parametersText, workflowRequest.Model)

	// 计算输出 token（基于 workflow_output_data）
	outputText := ""
	for key, value := range response.WorkflowOutputData {
		if strValue, ok := value.(string); ok {
			outputText += key + ":" + strValue + " "
		} else {
			// 对于非字符串类型，序列化为 JSON
			if jsonBytes, err := json.Marshal(value); err == nil {
				outputText += key + ":" + string(jsonBytes) + " "
			}
		}
	}

	completionTokens = openai.CountTokenText(outputText, workflowRequest.Model)
	totalTokens = promptTokens + completionTokens

	logger.SysLogf("工作流 Token 计算 - Model: %s, PromptTokens: %d, CompletionTokens: %d, TotalTokens: %d",
		workflowRequest.Model, promptTokens, completionTokens, totalTokens)

	return promptTokens, completionTokens, totalTokens
}

// executeN8nWorkflow 执行 n8n 工作流
func executeN8nWorkflow(c *gin.Context, workflowRequest *WorkflowRunRequest, agent *model.Agent, channel *model.Channel, modelName string) (*custom.WorkflowResponseData, error) {
	// 检查 Agent 类型是否为工作流类型
	if agent.AgentType != model.AgentTypeWorkflow {
		return nil, fmt.Errorf("Agent 类型不是工作流类型，当前类型: %d", agent.AgentType)
	}

	// 获取元数据
	meta := GetByContext(c)
	meta.APIType = model.GetApiType(channel.Type)
	meta.OriginModelName = modelName
	meta.ChannelId = int(channel.ChannelID)
	if channel.BaseURL != nil {
		meta.BaseURL = *channel.BaseURL
	}
	meta.APIKey = channel.Key

	// 应用模型映射
	mappedModel, _ := getMappedModelName(modelName, meta.ModelMapping)
	meta.ActualModelName = mappedModel

	logger.SysLogf("n8n工作流执行 - 模型映射，OriginModel: %s, ActualModel: %s",
		meta.OriginModelName, meta.ActualModelName)

	// 创建工作流适配器
	workflowAdaptor := &n8n.N8nWorkflowAdaptor{}
	workflowAdaptor.Init(meta)

	// 设置自定义配置
	user_id := config.GetUserId(c)
	conversation, err := GetSessionConversation(c)
	if err == nil {
		customConfig := &custom.CustomConfig{
			UserId:                     "angethub_u" + fmt.Sprintf("%d", user_id),
			ConversationId:             conversation.ChannelConversationID,
			ConversationExpirationTime: conversation.ChannelConversationExpirationTime,
			AIHubConversationId:        conversation.ConversationID,
		}
		workflowAdaptor.CustomConfig = customConfig
	}

	// 构建工作流请求
	workflowID := extractWorkflowID(agent.Model, agent.CustomConfig)
	if workflowID == "" {
		return nil, fmt.Errorf("无法提取工作流ID")
	}

	// 转换工作流请求为 n8n 工作流请求
	n8nRequest, err := workflowAdaptor.ConvertWorkflowRequest(workflowID, workflowRequest.Parameters)
	if err != nil {
		return nil, fmt.Errorf("转换n8n工作流请求失败: %v", err)
	}

	// 序列化请求
	requestBody, err := json.Marshal(n8nRequest)
	if err != nil {
		return nil, fmt.Errorf("序列化n8n工作流请求失败: %v", err)
	}

	// 执行请求
	resp, err := workflowAdaptor.DoRequest(c, meta, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("执行n8n工作流请求失败: %v", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode >= 400 {
		return nil, handleWorkflowError(resp, "n8n")
	}

	// 处理响应
	workflowResponse, err := workflowAdaptor.ProcessResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("处理n8n工作流响应失败: %v", err)
	}

	// 设置响应信息
	workflowResponse.ChannelID = int(channel.ChannelID)
	workflowResponse.ModelName = agent.Model

	if len(workflowResponse.WorkflowOutputData) == 0 {
		logger.SysLogf("⚠️ n8n工作流执行成功但输出字段为空 - ExecuteID: %s", workflowResponse.ExecuteID)
		logger.SysLogf("🔍 n8n工作流详细输出数据: %+v", workflowResponse)
	} else {
		logger.SysLogf("✅ n8n工作流执行成功 - ExecuteID: %s, 输出字段数: %d",
			workflowResponse.ExecuteID, len(workflowResponse.WorkflowOutputData))
	}

	return workflowResponse, nil
}

// getWorkflowChannelType 获取工作流的渠道类型
func getWorkflowChannelType(response *custom.WorkflowResponseData) int {
	// 从响应中获取渠道ID，然后查询渠道类型
	if response.ChannelID > 0 {
		if channel, err := model.GetChannelByID(int64(response.ChannelID)); err == nil {
			return channel.Type
		}
	}

	// 默认返回 Coze 类型（当前主要支持的工作流类型）
	return channeltype.Coze
}

func GetByContext(c *gin.Context) *relay_meta.Meta {
	meta := relay_meta.Meta{
		Mode:            relaymode.GetByPath(c.Request.URL.Path),
		ChannelType:     c.GetInt(ctxkey.Channel),
		ChannelId:       c.GetInt(ctxkey.ChannelId),
		TokenId:         c.GetInt(ctxkey.TokenId),
		TokenName:       c.GetString(ctxkey.TokenName),
		UserId:          c.GetInt(ctxkey.Id),
		Group:           c.GetString(ctxkey.Group),
		ModelMapping:    c.GetStringMapString(ctxkey.ModelMapping),
		OriginModelName: c.GetString(ctxkey.RequestModel),
		BaseURL:         c.GetString(ctxkey.BaseURL),
		APIKey:          strings.TrimPrefix(c.Request.Header.Get("Authorization"), "Bearer "),
		RequestURLPath:  c.Request.URL.String(),
		SystemPrompt:    c.GetString(ctxkey.SystemPrompt),
	}
	cfg, ok := c.Get(ctxkey.Config)
	if ok {
		meta.Config = cfg.(oneapi_model.ChannelConfig)
	}
	if meta.BaseURL == "" {
		if meta.ChannelType >= 0 && meta.ChannelType < len(channeltype.ChannelBaseURLs) {
			meta.BaseURL = channeltype.ChannelBaseURLs[meta.ChannelType]
		}
	}
	meta.APIType = channeltype.ToAPIType(meta.ChannelType)
	return &meta
}
