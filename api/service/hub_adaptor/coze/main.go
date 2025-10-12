package coze

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/53AI/53AIHub/common/storage"
	db_model "github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/gin-gonic/gin"
	"github.com/songquanpeng/one-api/common"
	"github.com/songquanpeng/one-api/common/conv"
	"github.com/songquanpeng/one-api/common/helper"
	"github.com/songquanpeng/one-api/common/logger"
	"github.com/songquanpeng/one-api/common/render"
	"github.com/songquanpeng/one-api/relay/adaptor/coze/constant/messagetype"
	"github.com/songquanpeng/one-api/relay/adaptor/openai"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
)

// https://www.coze.com/open

func stopReasonCoze2OpenAI(reason *string) string {
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

func ConvertRequest(textRequest model.GeneralOpenAIRequest, meta *meta.Meta, customConfig *custom.CustomConfig) *Request {
	modelName := "bot-" + strings.TrimPrefix(meta.ActualModelName, "bot-")
	channelID := meta.ChannelId
	cozeRequest := Request{
		Stream: textRequest.Stream,
		UserID: customConfig.UserId,
		BotId:  strings.TrimPrefix(textRequest.Model, "bot-"),
	}
	for _, message := range textRequest.Messages {
		typeStr := TypeQuestion
		contentType := ContentTypeText
		if message.Role == "assistant" {
			typeStr = TypeAnswer
		}
		contentStr := message.StringContent()
		if contentStr == "" {
			continue
		}
		var contentObjs []db_model.ObjectStringContent
		if err := json.Unmarshal([]byte(contentStr), &contentObjs); err == nil {
			if len(contentObjs) == 0 {
				continue
			}
			var mergedContent []map[string]interface{}
			for _, contentObj := range contentObjs {
				if contentObj.Type == "text" {
					mergedContent = append(mergedContent, map[string]interface{}{
						"type": "text",
						"text": contentObj.Content,
					})
					continue
				}
				uoloadFile := contentObj.GetUploadFile()
				if uoloadFile == nil {
					logger.SysError("file not found")
					continue
				}
				fileMapping := uoloadFile.GetChannelFileMapping(channelID, modelName)
				if fileMapping == nil {
					fileMapping = &db_model.ChannelFileMapping{}
					err := CozeUploadFile(meta, uoloadFile, fileMapping)
					if err != nil {
						logger.SysError(fmt.Sprintf("upload file failed: %v", err))
						continue
					}
					err = db_model.CreateChannelFileMapping(fileMapping)
					if err != nil {
						logger.SysError(fmt.Sprintf("create file mapping failed: %v", err))
						continue
					}
				} else if helper.GetTimestamp() > fileMapping.ExpirationTime {
					err := CozeUploadFile(meta, uoloadFile, fileMapping)
					if err != nil {
						logger.SysError(fmt.Sprintf("update file failed: %v", err))
						continue
					}
					err = db_model.UpdateChannelFileMapping(fileMapping)
					if err != nil {
						logger.SysError(fmt.Sprintf("update file mapping failed: %v", err))
						continue
					}
				}

				var contentType string
				if strings.HasPrefix(uoloadFile.MimeType, "image/") {
					contentType = "image"
				} else {
					contentType = "file"
				}
				mergedContent = append(mergedContent, map[string]interface{}{
					"type":    contentType,
					"file_id": fileMapping.ChannelFileID,
				})
			}
			mergedJSON, _ := json.Marshal(mergedContent)
			cozeMessage := AdditionalMessage{
				Role:        message.Role,
				Content:     string(mergedJSON),
				ContentType: ContentTypeObjectString,
				Type:        typeStr,
			}
			cozeRequest.AdditionalMessages = append(cozeRequest.AdditionalMessages, cozeMessage)
		} else {
			cozeMessage := AdditionalMessage{
				Role:        message.Role,
				Content:     contentStr,
				ContentType: contentType,
				Type:        typeStr,
			}
			cozeRequest.AdditionalMessages = append(cozeRequest.AdditionalMessages, cozeMessage)
		}
	}
	return &cozeRequest
}

func StreamResponseCoze2OpenAI(cozeResponse *StreamResponseV3) (*openai.ChatCompletionsStreamResponse, *Response) {
	var response *Response
	var stopReason string
	var choice openai.ChatCompletionsStreamResponseChoice

	choice.Delta.Content = cozeResponse.Content
	choice.Delta.Role = "assistant"
	finishReason := stopReasonCoze2OpenAI(&stopReason)
	if finishReason != "null" {
		choice.FinishReason = &finishReason
	}
	var openaiResponse openai.ChatCompletionsStreamResponse
	openaiResponse.Object = "chat.completion.chunk"
	openaiResponse.Choices = []openai.ChatCompletionsStreamResponseChoice{choice}
	openaiResponse.Id = cozeResponse.ConversationId
	return &openaiResponse, response
}

func ResponseCoze2OpenAI(cozeResponse *Response) (*openai.TextResponse, string) {
	var responseText string
	for _, message := range cozeResponse.Messages {
		if message.Type == messagetype.Answer {
			responseText = message.Content
			break
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
		Id:      fmt.Sprintf("chatcmpl-%s", cozeResponse.ConversationId),
		Model:   "coze-bot",
		Object:  "chat.completion",
		Created: helper.GetTimestamp(),
		Choices: []openai.TextResponseChoice{choice},
	}
	return &fullTextResponse, cozeResponse.ConversationId
}

func StreamHandler(c *gin.Context, resp *http.Response) (*model.ErrorWithStatusCode, *string, string) {
	var responseText string
	createdTime := helper.GetTimestamp()
	scanner := bufio.NewScanner(resp.Body)
	// 设置更大的缓冲区以处理大型响应 (1MB)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanLines)

	common.SetEventStreamHeaders(c)
	conversationId := ""
	var modelName string

	eventStr := ""
	for scanner.Scan() {
		data := scanner.Text()
		if data == "" || data == "\n" {
			continue
		}
		// logger.SysLogf("coze stream : %s", data)
		if strings.HasPrefix(data, "event:") {
			eventStr = strings.TrimPrefix(data, "event:")
		}

		if eventStr != "conversation.message.delta" {
			continue
		}

		if len(data) < 5 || !strings.HasPrefix(data, "data:") {
			continue
		}
		data = strings.TrimPrefix(data, "data:")
		data = strings.TrimSuffix(data, "\r")

		var cozeResponse StreamResponseV3
		err := json.Unmarshal([]byte(data), &cozeResponse)
		if err != nil {
			logger.SysError("error unmarshalling stream response: " + err.Error())
			continue
		}

		if cozeResponse.Type != "answer" {
			continue
		}

		response, _ := StreamResponseCoze2OpenAI(&cozeResponse)

		if response == nil {
			continue
		}

		if response.Id != "" {
			conversationId = response.Id
		}

		for _, choice := range response.Choices {
			responseText += conv.AsString(choice.Delta.Content)
		}
		response.Model = modelName
		response.Created = createdTime

		err = render.ObjectData(c, response)
		if err != nil {
			logger.SysError(err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		logger.SysError("error reading stream: " + err.Error())
	}

	render.Done(c)

	err := resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil, conversationId
	}

	return nil, &responseText, conversationId
}

func Handler(c *gin.Context, resp *http.Response, promptTokens int, modelName string) (*model.ErrorWithStatusCode, *string, string) {
	responseBody, err := io.ReadAll(resp.Body)
	conversationId := ""
	if err != nil {
		return openai.ErrorWrapper(err, "read_response_body_failed", http.StatusInternalServerError), nil, conversationId
	}
	err = resp.Body.Close()
	if err != nil {
		return openai.ErrorWrapper(err, "close_response_body_failed", http.StatusInternalServerError), nil, conversationId
	}
	var cozeResponse Response
	err = json.Unmarshal(responseBody, &cozeResponse)
	if err != nil {
		return openai.ErrorWrapper(err, "unmarshal_response_body_failed", http.StatusInternalServerError), nil, conversationId
	}
	if cozeResponse.Code != 0 {
		return &model.ErrorWithStatusCode{
			Error: model.Error{
				Message: cozeResponse.Msg,
				Code:    cozeResponse.Code,
			},
			StatusCode: resp.StatusCode,
		}, nil, conversationId
	}
	fullTextResponse, conversationId := ResponseCoze2OpenAI(&cozeResponse)
	fullTextResponse.Model = modelName
	jsonResponse, err := json.Marshal(fullTextResponse)
	if err != nil {
		return openai.ErrorWrapper(err, "marshal_response_body_failed", http.StatusInternalServerError), nil, conversationId
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(resp.StatusCode)
	_, err = c.Writer.Write(jsonResponse)
	var responseText string
	if len(fullTextResponse.Choices) > 0 {
		responseText = fullTextResponse.Choices[0].Message.StringContent()
	}
	return nil, &responseText, conversationId
}

func CozeUploadFile(meta *meta.Meta, uploadFile *db_model.UploadFile, fileMapping *db_model.ChannelFileMapping) error {
	url := fmt.Sprintf("%s/v1/files/upload", meta.BaseURL)
	fileContent, err := storage.StorageInstance.Load(uploadFile.Key)
	if err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", uploadFile.FileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	if err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+meta.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}

	var result struct {
		Code int `json:"code"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return err
	}

	if result.Code != 0 {
		return fmt.Errorf("upload failed with code: %d", result.Code)
	}
	ExpirationTime := helper.GetTimestamp() + 3600*24*30

	fileMapping.Eid = uploadFile.Eid
	fileMapping.FileID = uploadFile.ID
	fileMapping.ChannelID = meta.ChannelId
	fileMapping.Model = "bot-" + strings.TrimPrefix(meta.ActualModelName, "bot-")
	fileMapping.ChannelFileID = result.Data.ID
	fileMapping.ExpirationTime = ExpirationTime
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}
	fileMapping.ApiResponse = string(jsonResult)

	return nil
}
