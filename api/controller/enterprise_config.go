package controller

import (
	"errors"
	"net/http"

	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 只返回类型和启用状态
type ConfigTypeStatus struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

// GetEnterpriseConfigTypes 获取所有企业配置的类型及启用状态
// @Summary 获取所有企业配置的类型及启用状态
// @Description 获取所有企业配置的列表，只返回对应类型 type 的 enable 状态
// @Tags EnterpriseConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.CommonResponse{data=[]ConfigTypeStatus}
// @Router /api/enterprise-configs [get]
func GetEnterpriseConfigTypes(c *gin.Context) {
	eid := config.GetEID(c)

	configs, err := service.GetEnterpriseConfigs(eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	result := make([]ConfigTypeStatus, 0)
	for _, config := range configs {
		result = append(result, ConfigTypeStatus{
			Type:    config.Type,
			Enabled: config.Enabled,
		})
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(result))
}

// GetEnterpriseConfig 获取单个企业配置详情
// @Summary 获取单个企业配置详情
// @Description 通过 Type 获取单个企业配置详情
// @Tags EnterpriseConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "配置类型"
// @Success 200 {object} model.CommonResponse{data=model.EnterpriseConfig}
// @Router /api/enterprise-configs/{type} [get]
func GetEnterpriseConfig(c *gin.Context) {
	eid := config.GetEID(c)
	configType := c.Param("type")

	config, err := service.GetEnterpriseConfigByType(eid, configType)
	if err != nil {
		// 如果没有找到配置，返回一个默认配置
		defaultcontent, err := model.GetEnterpriseConfigDefaultContent(configType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		}
		config = &model.EnterpriseConfig{
			EID:     eid,
			Type:    configType,
			Content: defaultcontent,
			Enabled: false,
		}
		c.JSON(http.StatusOK, model.Success.ToResponse(config))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(config))
}

// IsEnterpriseConfigEnabled 检查企业配置是否启用
// @Summary 检查企业配置是否启用
// @Description 检查指定类型的企业配置是否启用（content中enable字段为true表示启用）
// @Tags EnterpriseConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "配置类型"
// @Success 200 {object} model.CommonResponse{data=bool}
// @Router /api/enterprise-configs/{type}/enabled [get]
func IsEnterpriseConfigEnabled(c *gin.Context) {
	eid := config.GetEID(c)
	configType := c.Param("type")

	enabled, err := service.IsEnterpriseConfigEnabled(eid, configType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, model.Success.ToResponse(false))
			return
		}
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(enabled))
}

type SaveEnterpriseConfigRequest struct {
	// smtp {\"smtp_host\":\"smtp_host.com\",\"smtp_username\":\"smtp_username@xx.com\",\"smtp_port\":\"465\",\"smtp_password\":\"xxxxxx\",\"smtp_from\":\"smtp_username@xx.com\",\"smtp_is_ssl\":true,\"smtp_to\":\"smtp_to\"}
	Content string `json:"content" binding:"required"`
	Enabled bool   `json:"enabled"`
}

// SaveEnterpriseConfig 保存单个企业配置详情
// @Summary 保存单个企业配置详情
// @Description 通过 Type 保存单个企业配置详情，如果不存在则创建，否则更新
// @Tags EnterpriseConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "配置类型"
// @Param config body SaveEnterpriseConfigRequest true "企业配置"
// @Success 200 {object} model.CommonResponse{data=model.EnterpriseConfig}
// @Router /api/enterprise-configs/{type} [post]
func SaveEnterpriseConfig(c *gin.Context) {
	eid := config.GetEID(c)
	configType := c.Param("type")

	var req SaveEnterpriseConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	config, err := service.SaveEnterpriseConfig(eid, configType, req.Content, req.Enabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(config))
}

// ToggleEnterpriseConfig 开关某个企业配置
// @Summary 开关某个企业配置
// @Description 开关某个 EnterpriseConfig，通过 Type
// @Tags EnterpriseConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "配置类型"
// @Success 200 {object} model.CommonResponse{data=bool}
// @Router /api/enterprise-configs/{type}/toggle [put]
func ToggleEnterpriseConfig(c *gin.Context) {
	eid := config.GetEID(c)
	configType := c.Param("type")

	enabled, err := service.ToggleEnterpriseConfig(eid, configType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, model.NotFound.ToResponse("请先保存设置"))
			return
		}
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(enabled))
}
