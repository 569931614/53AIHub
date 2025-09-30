package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
)

type ProviderRequest struct {
	Name         string  `json:"name" example:"Coze.cn"`                                                                // Platform name
	ProviderType int64   `json:"provider_type" example:"1"`                                                             // Platform type: 1 for coze.cn, 2 for coze.com
	Configs      string  `json:"configs" example:"{\"client_id\":\"client_id\", \"client_secret\": \"client_secret\"}"` // Platform configuration
	AccessToken  string  `json:"access_token" example:"access_token"`
	BaseURL      *string `json:"base_url" example:""`
}

func validateCozeConfig(providerType int64, configStr string) error {

	switch providerType {
	case model.ProviderTypeCozeCn, model.ProviderTypeCozeCom:
		var config model.CozeConfig
		if err := json.Unmarshal([]byte(configStr), &config); err != nil {
			return fmt.Errorf("invalid coze.cn config: %v", err)
		}
		if config.ClientID == "" || config.ClientSecret == "" {
			return errors.New("client_id and client_secret are required for coze.cn config")
		}
	}

	return nil
}

func checkSaveAccessToken(ProviderType int64, req ProviderRequest) (bool, error) {
	saveAccessToken := false
	switch ProviderType {
	case model.ProviderTypeAppBuilder, model.ProviderType53AI, model.ProviderTypeCozeStudio:
		if req.AccessToken == "" {
			return saveAccessToken, errors.New("access_token is required for provider")
		}
		saveAccessToken = true
	case model.ProviderTypeCozeCn, model.ProviderTypeCozeCom:
		if req.AccessToken != "" {
			return saveAccessToken, errors.New("access_token is not required for provider")
		}
	}
	return saveAccessToken, nil
}

// @Summary Create provider
// @Description Create new provider configuration
// @Tags Provider
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param provider body ProviderRequest true "Provider data"
// @Success 200 {object} model.CommonResponse
// @Router /api/providers [post]
func CreateProvider(c *gin.Context) {
	var req ProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if err := validateCozeConfig(req.ProviderType, req.Configs); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	provider := model.Provider{
		Eid:          config.GetEID(c),
		Name:         req.Name,
		ProviderType: req.ProviderType,
		Configs:      req.Configs,
		BaseURL:      req.BaseURL,
	}
	saveAccessToken, err := checkSaveAccessToken(req.ProviderType, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if saveAccessToken {
		provider.AccessToken = req.AccessToken
		provider.IsAuthorized = true
	}

	if err := model.CreateProvider(&provider); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(provider))
}

// @Summary Delete provider
// @Description Delete provider by ID
// @Tags Provider
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Provider ID"
// @Success 200 {object} model.CommonResponse
// @Router /api/providers/{id} [delete]
func DeleteProvider(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	eid := config.GetEID(c)

	// Start transaction to ensure atomicity
	tx := model.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(tx.Error))
		return
	}

	// First delete all channels associated with this provider
	if err := tx.Where("provider_id = ? AND eid = ?", id, eid).Delete(&model.Channel{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	// Then delete the provider itself
	if err := tx.Where("provider_id = ? AND eid = ?", id, eid).Delete(&model.Provider{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

// @Summary Get providers
// @Description Get all providers for current enterprise
// @Tags Provider
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param providerType query int false "Provider type filter (0 for all, 1 for coze.cn, 2 for coze.com, 3 for bailian, 4 for volcengine)"
// @Param name query string false "Provider name filter (supports like search)"
// @Success 200 {object} model.CommonResponse
// @Router /api/providers [get]
// GetProviders Get all provider configurations for current enterprise
func GetProviders(c *gin.Context) {
	eid := config.GetEID(c)
	providerType, _ := strconv.ParseInt(c.DefaultQuery("providerType", "0"), 10, 64)
	name := c.Query("name")
	providers, err := model.GetProvidersByEidWithFilters(eid, providerType, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(providers))
}

// @Summary Update provider
// @Description Update provider configuration
// @Tags Provider
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Provider ID"
// @Param provider body ProviderRequest true "Provider data"
// @Success 200 {object} model.CommonResponse
// @Router /api/providers/{id} [put]
func UpdateProvider(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	eid := config.GetEID(c)
	existingProvider, err := model.GetProviderByID(id, eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	var req ProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if err := validateCozeConfig(req.ProviderType, req.Configs); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	// Reset authorization if configs changed
	if existingProvider.Configs != req.Configs {
		existingProvider.IsAuthorized = false
	}

	existingProvider.Name = req.Name
	existingProvider.ProviderType = req.ProviderType
	existingProvider.Configs = req.Configs
	if req.BaseURL != nil {
		existingProvider.BaseURL = req.BaseURL
	}

	saveAccessToken, err := checkSaveAccessToken(req.ProviderType, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if saveAccessToken {
		existingProvider.AccessToken = req.AccessToken
		existingProvider.IsAuthorized = true
	}

	if err := model.UpdateProvider(existingProvider); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(existingProvider))
}
