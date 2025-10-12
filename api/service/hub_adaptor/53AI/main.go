package adaptor53AI

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/render"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/model"
)

func StreamHandler(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, *string, string) {
	var responseText string
	createdTime := helper.GetTimestamp()
	scanner := bufio.NewScanner(resp.Body)
	// 设置更大的缓冲区以处理大型响应 (1MB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanLines)
	common.SetEventStreamHeaders(c)
	var modelName string

	channelConversationId := ""
	for scanner.Scan() {
		data := scanner.Text()
		if len(data) < 5 || !strings.HasPrefix(data, "data:") {
			continue
		}
		data = strings.TrimPrefix(data, "data:")
		data = strings.TrimSuffix(data, "\r")
		// logger.SysLogf("53AI Stream Data: %s\n", data)

		var ai53Response StreamResponse
		err := json.Unmarshal([]byte(data), &ai53Response)
		if err != nil {
			logger.SysError("error unmarshalling stream response: " + err.Error())
			continue
		}

		response, _ := StreamResponseai53OpenAI(&ai53Response)
		if response == nil {
			continue
		}

		response.Model = modelName
		response.Created = createdTime

		err = render.ObjectData(c, response)
		if err != nil {
			logger.SysError(err.Error())
		}
		channelConversationId = ai53Response.ConversationID
	}

	if err := scanner.Err(); err != nil {
		logger.SysError("error reading stream: " + err.Error())
	}

	render.Done(c)

	err := resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil, channelConversationId
	}

	return nil, &responseText, channelConversationId
}

func StreamResponseai53OpenAI(ai53Response *StreamResponse) (*openai.ChatCompletionsStreamResponse, *Response) {
	var response *Response
	var stopReason string
	var choice openai.ChatCompletionsStreamResponseChoice

	// Build content from answer and append_content for normal and error events alike.
	answerStr := ""
	if ai53Response.Answer != "" {
		answerStr = ai53Response.Answer
	}
	if ai53Response.AppendContents != nil && len(ai53Response.AppendContents) > 0 {
		for _, appendContent := range ai53Response.AppendContents {
			if appendContent.Type == "text" {
				answerStr += appendContent.Content
			}
		}
	}

	// If this is an error event but upstream didn't provide textual answer,
	// still emit a minimal content to ensure front-end receives a normal data chunk.
	if ai53Response.Event == "error" && strings.TrimSpace(answerStr) == "" {
		// keep it minimal and non-structured, as requested: treat as normal data
		answerStr = "Upstream Error: \n" + ai53Response.Message
	}

	choice.Delta.Content = answerStr
	choice.Delta.Role = "assistant"

	finishReason := stopReasonAi53OpenAI(&stopReason)
	if finishReason != "null" {
		choice.FinishReason = &finishReason
	}

	var openaiResponse openai.ChatCompletionsStreamResponse
	openaiResponse.Object = "chat.completion.chunk"
	openaiResponse.Choices = []openai.ChatCompletionsStreamResponseChoice{choice}
	openaiResponse.Id = ai53Response.ConversationID
	return &openaiResponse, response
}

func BlockingResponseai53OpenAI(ai53Response *BlockResponse) (*openai.ChatCompletionsStreamResponse, *Response) {
	var response *Response
	var choice openai.ChatCompletionsStreamResponseChoice

	if ai53Response.Answer != "" {
		choice.Delta.Content = ai53Response.Answer
	}
	var openaiResponse openai.ChatCompletionsStreamResponse
	openaiResponse.Object = "chat.completion.chunk"
	openaiResponse.Choices = []openai.ChatCompletionsStreamResponseChoice{choice}
	openaiResponse.Id = ai53Response.ConversationID
	return &openaiResponse, response
}

func stopReasonAi53OpenAI(reason *string) string {
	if reason == nil {
		return ""
	}
	switch *reason {
	case "end_turn":
		return "stop"
	case "stop_sequence":
		return "stop"
	case "max_tokens":
		return "length"
	default:
		return *reason
	}
}

const (
	TransferMethodRemoteUrl = "remote_url"
	TransferMethodLocalFile = "local_file"
)

type File struct {
	TransferMethod string `json:"transfer_method"`
	Type           string `json:"type"`
	Url            string `json:"url"`
	UploadFileID   string `json:"upload_file_id"`
}

const (
	ResponseModeStream = "streaming"
	ResponseModeBlock  = "blocking"
)

type Request struct {
	ConversationId string      `json:"conversation_id"`
	User           string      `json:"user"`
	Query          string      `json:"query"`
	ResponseMode   string      `json:"response_mode"`
	Files          []File      `json:"files"`
	Inputs         interface{} `json:"inputs"`
}
