package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/common/utils/appbuilder"
	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service"
	"github.com/gin-gonic/gin"
)

// GetAppBuilderAllBots Get all AppBuilder bots
// @Summary Get all AppBuilder bots list
// @Description Get all bots list from AppBuilder
// @Tags AppBuilder
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider_id query int false "Provider ID (optional, for backward compatibility)"
// @Success 200 {object} model.CommonResponse{data=[]appbuilder.AppInfo}
// @Router /api/appbuilder/bots [get]
func GetAppBuilderAllBots(c *gin.Context) {
	eid := config.GetEID(c)
	providerID, _ := strconv.ParseInt(c.DefaultQuery("provider_id", "0"), 10, 64)
	provider, err := model.GetProviderByEidAndProviderTypeWithOptionalID(eid, int64(model.ProviderTypeAppBuilder), providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ProviderNoFoundError.ToResponse(err))
		return
	}

	ser := service.AppBuilderService{
		Provider: &provider,
	}

	apps, err := ser.GetAllDescribeApps()
	if err != nil {
		logger.SysLogf("GetAppBuilderAllBots: %s", err.Error())
		if len(apps) == 0 {
			c.JSON(http.StatusInternalServerError, model.ProviderNoFoundError.ToResponse(err))
			return
		}
	}
	// update or channel
	UpdateAppBuilderChannel(provider, apps)
	c.JSON(http.StatusOK, model.Success.ToResponse(apps))
}

func UpdateAppBuilderChannel(provider model.Provider, apps []appbuilder.AppInfo) error {
	var botIds []string
	for _, bot := range apps {
		if !bot.IsPublished {
			continue
		}
		botIds = append(botIds, "bot-"+bot.ID)
	}
	// todo Call interface user ID all use the same
	configStr := `{"region":"","sk":"","ak":"","user_id":"53AIHub","vertex_ai_project_id":"","vertex_ai_adc":""}`
	// Create or update channel record
	baseURL := provider.GetBaseURLByProviderType()
	channel := &model.Channel{
		Eid:        provider.Eid,
		Name:       provider.Name,
		Key:        provider.AccessToken,
		Type:       model.ChannelApiAppBuilder,
		ProviderID: provider.ProviderID,
		BaseURL:    &baseURL, // Fix: assign address of baseURL since Channel.BaseURL is *string
		Models:     strings.Join(botIds, ","),
		Status:     model.ChannelStatusEnabled,
		Config:     configStr,
	}

	existingChannel, err := model.GetFirstChannelByEidAndProviderId(channel.Eid, channel.ProviderID)
	if err != nil {
		// Create new record
		return model.CreateChannel(channel) // Fix: return error instead of assigning i
	} else {
		// Update existing record
		channel.ChannelID = existingChannel.ChannelID
		channel.CreatedTime = existingChannel.CreatedTime
		return model.UpdateChannel(channel)
	}
}
