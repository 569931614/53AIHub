package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
)

// autoAssignCozeStudioProvider automatically assigns a ProviderID for CozeStudio channels
// when ProviderID is 0 in the request
func autoAssignCozeStudioProvider(channel *model.Channel) error {
	// Check if this is a CozeStudio channel and ProviderID is 0
	if channel.Type == model.ChannelApiTypeCozeStudio && channel.ProviderID == 0 {
		// Get the first Provider with ProviderTypeCozeStudio for this enterprise
		provider, err := model.GetFirstProviderByEidAndProviderType(channel.Eid, model.ProviderTypeCozeStudio)
		if err != nil {
			return err
		}
		// Assign the ProviderID to the channel
		channel.ProviderID = provider.ProviderID
	}
	return nil
}

type ChannelRequest struct {
	Type         int     `json:"type" example:"1"`
	ModelType    *int    `json:"model_type" example:"1"`
	Key          string  `json:"key" example:"channel_key"`
	Name         string  `json:"name" example:"channel_name"`
	Models       string  `json:"models" example:"gpt-3.5-turbo"`
	Config       string  `json:"config" example:"{\"region\":\"us-east-1\"}"`
	ModelMapping *string `json:"model_mapping"`
	Weight       *uint   `json:"weight"`
	Priority     *int64  `json:"priority"`
	BaseURL      *string `json:"base_url"`
	Other        *string `json:"other"`
	ProviderID   *int64  `json:"provider_id" example:"181"`
}

// @Summary Create channel
// @Description Create new channel configuration
// @Tags Channel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel body ChannelRequest true "Channel data"
// @Success 200 {object} model.CommonResponse
// @Router /api/channels [post]
func CreateChannel(c *gin.Context) {
	var req ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	channel := model.Channel{
		Eid:          config.GetEID(c),
		Type:         req.Type,
		ModelType:    1,
		Key:          req.Key,
		Name:         req.Name,
		Models:       req.Models,
		Config:       req.Config,
		ModelMapping: req.ModelMapping,
		Weight:       req.Weight,
		Priority:     req.Priority,
		BaseURL:      req.BaseURL,
		Other:        req.Other,
		ProviderID:   0, // Default to 0 if not provided
	}

	// Set ProviderID if provided in request
	if req.ProviderID != nil {
		channel.ProviderID = *req.ProviderID
	}
	// Set ModelType: default 1; if provided and valid (1,2,3), use it
	if req.ModelType != nil {
		if model.IsValidModelType(*req.ModelType) {
			channel.ModelType = *req.ModelType
		}
	}

	channel.Models = model.ProcessModelNames(req.Models, channel.Type)
	if channel.Models == "" {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(strings.NewReader("models is required")))
		return
	}

	// Auto assign ProviderID for CozeStudio channels if needed
	if err := autoAssignCozeStudioProvider(&channel); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	if err := model.CreateChannel(&channel); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(channel))
}

// @Summary Get channel
// @Description Get channel configuration by ID
// @Tags Channel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel_id path int true "Channel ID"
// @Success 200 {object} model.CommonResponse
// @Router /api/channels/{channel_id} [get]
func GetChannel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("channel_id"), 10, 64)
	channel, err := model.GetChannelByID(id)

	if err != nil || channel.Eid != config.GetEID(c) {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(channel))
}

// @Summary Update channel
// @Description Update existing channel configuration
// @Tags Channel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel_id path int true "Channel ID"
// @Param channel body ChannelRequest true "Update data"
// @Success 200 {object} model.CommonResponse
// @Router /api/channels/{channel_id} [put]
func UpdateChannel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("channel_id"), 10, 64)
	channel, err := model.GetChannelByID(id)

	if err != nil || channel.Eid != config.GetEID(c) {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	var req ChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	channel.Models = model.ProcessModelNames(req.Models, channel.Type)

	if channel.Models == "" {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(strings.NewReader("models is required")))
		return
	}

	channel.Type = req.Type
	channel.Key = req.Key
	channel.Name = req.Name

	channel.Config = req.Config
	channel.ModelMapping = req.ModelMapping
	channel.Weight = req.Weight
	channel.Priority = req.Priority
	channel.BaseURL = req.BaseURL
	channel.Other = req.Other

	// Update ProviderID if provided in request
	if req.ProviderID != nil {
		channel.ProviderID = *req.ProviderID
	}
	// Update ModelType if provided and valid (1,2,3)
	if req.ModelType != nil {
		if model.IsValidModelType(*req.ModelType) {
			channel.ModelType = *req.ModelType
		}
	}

	// Auto assign ProviderID for CozeStudio channels if needed
	if err := autoAssignCozeStudioProvider(channel); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	if err := model.UpdateChannel(channel); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(channel))
}

// @Summary Delete channel
// @Description Delete channel by ID
// @Tags Channel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel_id path int true "Channel ID"
// @Success 200 {object} model.CommonResponse
// @Router /api/channels/{channel_id} [delete]
func DeleteChannel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("channel_id"), 10, 64)
	channel, err := model.GetChannelByID(id)

	if err == nil && channel.Eid == config.GetEID(c) {
		err = model.DeleteChannelByID(id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

// @Summary Get all channels
// @Description Get all channels for current enterprise
// @Tags Channel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_id query int false "Provider ID, 0 means platform-added keys, non-zero means get channels from other platforms" example:"0"
// @Param channel_types query string false "Channel type filters" example:"1,1001,1002"
// @Param model_type query string false "Model type filters: 1=LLM,2=Embedding,3=Rerank; comma-separated supported; 0 or empty means no filter" example:"1,3"
// @Success 200 {object} model.CommonResponse
// @Router /api/channels [get]
func GetChannels(c *gin.Context) {
	providerId, _ := strconv.ParseInt(c.Query("provider_id"), 10, 64)
	channelTypesStr := c.Query("channel_types")
	var channelTypes []int
	if channelTypesStr != "" {
		for _, s := range strings.Split(channelTypesStr, ",") {
			if t, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				channelTypes = append(channelTypes, t)
			}
		}
	}

	modelTypesStr := c.Query("model_type")
	var modelTypes []int
	if modelTypesStr != "" {
		for _, s := range strings.Split(modelTypesStr, ",") {
			if t, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				// Only accept defined model types; 0 or invalid values mean no filter
				if model.IsValidModelType(t) {
					modelTypes = append(modelTypes, t)
				}
			}
		}
	}

	channels, err := model.GetChannelsByEidAndParams(config.GetEID(c), providerId, channelTypes, modelTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(channels))
}
