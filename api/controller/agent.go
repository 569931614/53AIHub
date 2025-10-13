package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/53AI/53AIHub/common"
	"github.com/53AI/53AIHub/common/utils"
	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/53AI/53AIHub/service"
	"github.com/gin-gonic/gin"
)

type AgentListRequest struct {
	Keyword      string `json:"keyword" form:"keyword" example:"Json"`
	GroupId      int64  `json:"group_id" form:"group_id" example:"0"`
	Offset       int    `json:"offset" form:"offset" example:"0"`
	Limit        int    `json:"limit" form:"limit" example:"10"`
	ChannelTypes string `json:"channel_types" form:"channel_types" example:"0,1,2"`
	AgentTypes   string `json:"agent_types" form:"agent_types" example:"0,1,2"`
}

type AgentsResponse struct {
	Count  int64          `json:"count"`
	Agents []*model.Agent `json:"agents"`
}

type AgentRequest struct {
	Name                 string  `json:"name" example:"OpenAI-ChatGPT"`
	Logo                 string  `json:"logo" example:"http://URL_ADDRESS.com/logo.png"`
	Sort                 int     `json:"sort" example:"0"`
	Description          string  `json:"description" example:"A ChatGPT based agent for general conversation"`
	Configs              string  `json:"configs" example:"{\"model\":\"gpt-3.5-turbo\",\"temperature\":0.7}"`
	Prompt               string  `json:"prompt" example:"You are a helpful AI assistant"`
	ChannelType          int     `json:"channel_type"`
	Model                string  `json:"model" example:"gpt-3.5-turbo"`
	GroupId              int64   `json:"group_id" example:"0"`
	UseCases             string  `json:"use_cases" example:"[]"`
	Tools                string  `json:"tools"  example:"[]"`
	CustomConfig         string  `json:"custom_config" example:"{}"`
	UserGroupIds         []int64 `json:"user_group_ids"`
	Enable               bool    `json:"enable" example:"true"`
	SubscriptionGroupIds []int64 `json:"subscription_group_ids"` // 订阅分组IDs
	Settings             string  `json:"settings" example:"{}"`
	AgentType            int     `json:"agent_type" example:"0"` // Agent type (0=App, 1=Workflow), default is 0
}

type UpdateAgentEnableRequest struct {
	Enable *bool `json:"enable" example:"true" binding:"required"` // Enable status (true=enabled, false=disabled)
}

// @Summary Create a new agent
// @Description Create agent with configurable parameters. agent_type: 0=App (default), 1=Workflow
// @Tags Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Agent body AgentRequest true "Agent Data"
// @Success 200 {object} model.CommonResponse{data=model.Agent} "Success"
// @Router /api/agents [post]
func CreateAgent(c *gin.Context) {
	var agentReq AgentRequest
	if err := c.ShouldBindJSON(&agentReq); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}
	// Check if user is admin
	if !common.IsAdmin(c) {
		c.JSON(http.StatusForbidden, model.AuthFailed.ToResponse(nil))
		return
	}

	agentReq.Model = model.ProcessModelNames(agentReq.Model, agentReq.ChannelType)
	if agentReq.Model == "" {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(errors.New("model is required")))
		return
	}

	params := map[string]interface{}{
		"from": "agent",
	}
	_, err := service.IsFeatureAvailable(c, "agent", params)
	if err != nil {
		c.JSON(http.StatusForbidden, model.FeatureNotAvailableError.ToResponse(err))
		return
	}

	// Start transaction
	tx := model.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	agent := model.Agent{
		Eid:          config.GetEID(c),
		Name:         agentReq.Name,
		Logo:         agentReq.Logo,
		ChannelType:  agentReq.ChannelType,
		Sort:         agentReq.Sort,
		Description:  agentReq.Description,
		Model:        agentReq.Model,
		Prompt:       agentReq.Prompt,
		Configs:      agentReq.Configs,
		Tools:        agentReq.Tools,
		CustomConfig: agentReq.CustomConfig,
		GroupID:      agentReq.GroupId,
		UseCases:     agentReq.UseCases,
		CreatedBy:    config.GetUserId(c),
		Enable:       agentReq.Enable,
		Settings:     agentReq.Settings,
		AgentType:    agentReq.AgentType, // 添加 AgentType 字段，默认为 0
	}

	if err := tx.Create(&agent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	allGroupIds := make([]int64, 0)

	if len(agentReq.SubscriptionGroupIds) > 0 {
		allGroupIds = append(allGroupIds, agentReq.SubscriptionGroupIds...)
	}

	if len(agentReq.UserGroupIds) > 0 {
		allGroupIds = append(allGroupIds, agentReq.UserGroupIds...)
	}

	// Add permissions for user groups
	if len(allGroupIds) > 0 {
		for _, groupID := range allGroupIds {
			permission := model.ResourcePermission{
				GroupID:      groupID,
				ResourceID:   agent.AgentID,
				ResourceType: model.ResourceTypeAgent,
				Permission:   model.PermissionRead,
			}
			if err := tx.Create(&permission).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Parse CustomConfig to get agent_type
	var customConfig map[string]interface{}
	if err := json.Unmarshal([]byte(agentReq.CustomConfig), &customConfig); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToErrorResponse(err))
		return
	}

	agentType, ok := customConfig["agent_type"].(string)
	if !ok {
		agentType = "unknown"
	}

	log := model.SystemLog{
		Eid:      agent.Eid,
		UserID:   agent.CreatedBy,
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleAgent,
		Action:   model.SystemLogActionCreate,
		Content:  fmt.Sprintf("新建智能体【】名称：【%s】；类型：%s", agent.Name, model.GetChannelDescription(agentType)),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	// Ensure both subscription_group_ids and user_group_ids are populated
	if err := agent.LoadGroupIdsByType(); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	c.JSON(http.StatusOK, model.Success.ToResponse(agent))
}

// @Summary Get agent details
// @Description Get agent by AgentID
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Param agent_id path int true "Agent ID"
// @Success 200 {object} model.CommonResponse{data=model.Agent} "Success"
// @Router /api/agents/{agent_id} [get]
func GetAgent(c *gin.Context) {
	agent_id, err := strconv.ParseInt(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}
	eid := config.GetEID(c)
	agent, err := model.GetAgentByID(eid, agent_id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	// Check if user is admin or has permission to access this agent
	if !common.IsAdmin(c) {
		hasPermission, err := model.CheckPermission(config.GetUserGroupID(c), agent_id, model.ResourceTypeAgent, model.PermissionRead)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, model.AuthFailed.ToResponse(nil))
			return
		}
	}

	if err := agent.LoadGroupIdsByType(); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	c.JSON(http.StatusOK, model.Success.ToResponse(agent))
}

// @Summary Update agent
// @Description Update existing agent details. agent_type: 0=App (default), 1=Workflow
// @Tags Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param agent_id path int true "Agent ID"
// @Param agent body AgentRequest true "Agent data" example:{"name":"OpenAI-ChatGPT","description":"ChatGPT","configs":"{\"model\":\"gpt-3.5-turbo\",\"temperature\":0.7}","prompt":"你好","model":"gpt-3.5-turbo","group_id":0,"use_cases":"[]","tools":"[]","user_group_ids":[1,2,3],"agent_type":0}
// @Success 200 {object} model.CommonResponse{data=model.Agent} "Success"
// @Router /api/agents/{agent_id} [put]
func UpdateAgent(c *gin.Context) {
	agent_id, err := strconv.ParseInt(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}
	eid := config.GetEID(c)
	agent, err := model.GetAgentByID(eid, agent_id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	// Check if user is admin
	if !common.IsAdmin(c) {
		c.JSON(http.StatusForbidden, model.AuthFailed.ToResponse(nil))
		return
	}

	var agentReq AgentRequest
	if err := c.ShouldBindJSON(&agentReq); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}

	agentReq.Model = model.ProcessModelNames(agentReq.Model, agentReq.ChannelType)
	if agentReq.Model == "" {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(errors.New("model is required")))
		return
	}

	// Start transaction
	tx := model.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	oldAgent := *agent

	// Update agent fields
	agent.Name = agentReq.Name
	agent.Description = agentReq.Description
	agent.Model = agentReq.Model
	agent.Prompt = agentReq.Prompt
	agent.Configs = agentReq.Configs
	agent.Tools = agentReq.Tools
	agent.GroupID = agentReq.GroupId
	agent.UseCases = agentReq.UseCases
	agent.ChannelType = agentReq.ChannelType
	agent.Sort = agentReq.Sort
	agent.Logo = agentReq.Logo
	agent.CustomConfig = agentReq.CustomConfig
	agent.Enable = agentReq.Enable
	agent.Settings = agentReq.Settings
	agent.AgentType = agentReq.AgentType // 添加 AgentType 字段更新

	if err := tx.Save(agent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// 获取企业信息
	enterprise, err := model.GetEnterpriseByID(eid)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	// 确定分组ID
	var allGroupIds []int64
	groupIDSet := make(map[int64]bool)

	switch enterprise.Type {
	case model.EnterpriseTypeIndustry:
		// 取所有分组ID
		if len(agentReq.SubscriptionGroupIds) > 0 {
			for _, id := range agentReq.SubscriptionGroupIds {
				groupIDSet[id] = true
			}
		}
		if len(agentReq.UserGroupIds) > 0 {
			for _, id := range agentReq.UserGroupIds {
				groupIDSet[id] = true
			}
		}
	case model.EnterpriseTypeIndependent:
		// 只取订阅分组ID
		if len(agentReq.SubscriptionGroupIds) > 0 {
			for _, id := range agentReq.SubscriptionGroupIds {
				groupIDSet[id] = true
			}
		}
	case model.EnterpriseTypeEnterprise:
		// 只取用户分组ID
		if len(agentReq.UserGroupIds) > 0 {
			for _, id := range agentReq.UserGroupIds {
				groupIDSet[id] = true
			}
		}
	}

	// 转换为切片
	allGroupIds = make([]int64, 0, len(groupIDSet))
	for id := range groupIDSet {
		allGroupIds = append(allGroupIds, id)
	}

	// 使用通用方法更新资源权限（带过滤）
	if err := service.UpdateAgentResourcePermissions(c, tx, agent.AgentID, agentReq.SubscriptionGroupIds, agentReq.UserGroupIds, enterprise); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Prepare for logging
	fieldMap := map[string]string{
		"Name":        "名称",
		"Description": "描述",
	}
	model.LogEntityChange(
		fmt.Sprintf("智能体【%s】", oldAgent.Name),
		model.SystemLogActionUpdate,
		eid,
		config.GetUserId(c),
		config.GetUserNickname(c),
		model.SystemLogModuleAgent,
		oldAgent,
		agent,
		utils.GetClientIP(c),
		fieldMap,
	)

	if err := agent.LoadGroupIdsByType(); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	c.JSON(http.StatusOK, model.Success.ToResponse(agent))
}

// @Summary Delete agent
// @Description Delete agent by ID
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Param agent_id path int true "Agent ID"
// @Success 200 {object} model.CommonResponse "Success"
// @Router /api/agents/{agent_id} [delete]
func DeleteAgent(c *gin.Context) {
	agent_id, err := strconv.ParseInt(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}

	eid := config.GetEID(c)
	agent, err := model.GetAgentByID(eid, agent_id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	// Check if user is admin
	if !common.IsAdmin(c) {
		c.JSON(http.StatusForbidden, model.AuthFailed.ToResponse(nil))
		return
	}

	// Start transaction
	tx := model.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Delete agent
	if err := tx.Delete(agent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Delete associated permissions
	if err := tx.Where("resource_id = ? AND resource_type = ?", agent_id, model.ResourceTypeAgent).Delete(&model.ResourcePermission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	log := model.SystemLog{
		Eid:      agent.Eid,
		UserID:   agent.CreatedBy,
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleAgent,
		Action:   model.SystemLogActionDelete,
		Content:  fmt.Sprintf("删除智能体【%s】", agent.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

// @Summary Get agents list
// @Description Retrieve paginated list of agents
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Param keyword query string false "Keyword"
// @Param group_id query int false "Group ID"
// @Param offset query int false "Offset"
// @Param limit query int false "Limit" default(10)
// @Param channel_types query string false "Channel types , split by comma"
// @Param agent_types query string false "Agent types , split by comma"
// @Success 200 {object} model.CommonResponse{data=AgentsResponse} "Success"
// @Router /api/agents [get]
func GetAgents(c *gin.Context) {
	var agentListRequest AgentListRequest
	if err := c.ShouldBindQuery(&agentListRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	// Set default limit to 10 if not specified
	if agentListRequest.Limit == 0 {
		agentListRequest.Limit = 10
	}

	// If user is not admin, get only agents they have permission to access
	var total int64
	var agents []*model.Agent
	var err error
	channelTypes := splitChannelTypesString(agentListRequest.ChannelTypes)
	agentTypes := splitAgentTypesString(agentListRequest.AgentTypes)
	if common.IsAdmin(c) {
		total, agents, err = model.GetAgentListWithIDs(
			config.GetEID(c), agentListRequest.Keyword, agentListRequest.GroupId,
			nil, channelTypes, agentTypes, agentListRequest.Offset, agentListRequest.Limit)
	} else {
		// Get list of agent IDs the user has permission to access
		permittedAgentIDs, getErr := model.GetResourcesByGroupAndType(config.GetUserGroupID(c), model.ResourceTypeAgent)
		if getErr != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}

		total, agents, err = model.GetAgentListWithIDs(
			config.GetEID(c), agentListRequest.Keyword, agentListRequest.GroupId,
			permittedAgentIDs, channelTypes, agentTypes,
			agentListRequest.Offset, agentListRequest.Limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	for _, agent := range agents {
		if err := agent.LoadGroupIdsByType(); err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
		Count:  total,
		Agents: agents,
	}))
}

// @Summary Get agents by group
// @Description Retrieve paginated list of agents by specified group
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Param keyword    query string false "Search keyword"
// @Param group_id   query int    false  "Target group ID"
// @Param offset     query int    false "Pagination offset" default(0)
// @Param limit      query int    false "Pagination limit"  default(10)
// @Param channel_types query string false "Channel types , split by comma"
// @Param agent_types query string false "Agent types , split by comma"
// @Success 200 {object} model.CommonResponse{data=AgentsResponse} "Success response with agent list"
// @Router /api/agents/group [get]
func GetAgentsByGroup(c *gin.Context) {
	var agentListRequest AgentListRequest
	if err := c.ShouldBindQuery(&agentListRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	// Set default limit to 10 if not specified
	if agentListRequest.Limit == 0 {
		agentListRequest.Limit = 10
	}

	channelTypes := splitChannelTypesString(agentListRequest.ChannelTypes)
	agentTypes := splitAgentTypesString(agentListRequest.AgentTypes)
	var total, agents, err = model.GetAgentListWithIDs(
		config.GetEID(c), agentListRequest.Keyword, agentListRequest.GroupId,
		nil, channelTypes, agentTypes, agentListRequest.Offset, agentListRequest.Limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	for _, agent := range agents {
		if err := agent.LoadGroupIdsByType(); err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
		Count:  total,
		Agents: agents,
	}))
}

// @Summary Get available agents
// @Description Retrieve paginated list of available agents
// @Tags Agent
// @Produce json
// @Param offset query int false "Pagination offset" default(0)
// @Param limit query int false "Pagination limit" default(10)
// @Param agent_types query string false "Agent types , split by comma"
// @Success 200 {object} model.CommonResponse{data=AgentsResponse} "Success response with available agent list"
// @Router /api/agents/available [get]
func GetAvailableAgents(c *gin.Context) {
	var agentListRequest AgentListRequest
	if err := c.ShouldBindQuery(&agentListRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if agentListRequest.Limit == 0 {
		agentListRequest.Limit = 10
	}

	agentTypes := splitAgentTypesString(agentListRequest.AgentTypes)
	var total, agents, err = model.GetAvailableAgentList(
		config.GetEID(c),
		agentTypes,
		agentListRequest.Offset,
		agentListRequest.Limit,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	for _, agent := range agents {
		if err = agent.LoadGroupIdsByType(); err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
		if err = agent.LoadConversationCount(); err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
			return
		}
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
		Count:  total,
		Agents: agents,
	}))
}

// Get Current Agents
// @Summary Get current agent list
// @Description Get agents list under the first agent-type group of current enterprise (no pagination required)
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.CommonResponse{data=AgentsResponse} "Success response containing agent list"
// @Router /api/agents/current [get]
func GetCurrentAgents(c *gin.Context) {
	eid := config.GetEID(c)

	authHeader := c.GetHeader("Authorization")
	authHeader = strings.Replace(authHeader, "Bearer ", "", 1)

	var theGroup *model.Group
	var err error

	if authHeader == "" {
		group, err := model.GetFirstGroupByEid(eid, model.USER_GROUP_TYPE)
		if err != nil {
			c.JSON(http.StatusNotFound, model.DBError.ToResponse(err))
			return
		}
		theGroup = &group
	} else {
		user := model.ValidateAccessToken(authHeader)
		if user == nil {
			c.JSON(http.StatusUnauthorized, model.UnauthorizedError.ToResponse(nil))
			return
		}

		if user.GroupId > 0 {
			theGroup, err = model.GetGroupByID(user.GroupId)
			if err != nil {
				c.JSON(http.StatusNotFound, model.DBError.ToResponse(err))
				return
			}
		}
	}

	group, err := model.GetGroupWithAgents(theGroup.GroupId, true)
	if err != nil {
		c.JSON(http.StatusNotFound, model.DBError.ToResponse(err))
		return
	}

	for i := range group.Agents {
		var count int64
		model.DB.Model(&model.Conversation{}).
			Where("eid = ? AND agent_id = ?", eid, group.Agents[i].AgentID).
			Count(&count)

		group.Agents[i].ConversationCount = count
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
		Count:  int64(len(group.Agents)),
		Agents: group.Agents,
	}))
}

// @Summary Update agent status
// @Description Update agent enable/disable status
// @Tags Agent
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param agent_id path int true "Agent ID"
// @Param request body UpdateAgentEnableRequest true "Enable status"
// @Success 200 {object} model.CommonResponse "Success"
// @Router /api/agents/{agent_id}/status [patch]
func UpdateAgentStatus(c *gin.Context) {
	agentID, err := strconv.ParseInt(c.Param("agent_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}

	var updateAgentEnableReq UpdateAgentEnableRequest
	if err = c.ShouldBindJSON(&updateAgentEnableReq); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(nil))
		return
	}

	eid := config.GetEID(c)
	_, err = model.GetAgentByID(eid, agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}

	err = model.UpdateAgentStatus(eid, agentID, updateAgentEnableReq.Enable)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	agent, err := model.GetAgentByID(eid, agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}
	statusText := "启用"
	if !agent.Enable {
		statusText = "禁用"
	}

	log := model.SystemLog{
		Eid:      eid,
		UserID:   config.GetUserId(c),
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleAgent,
		Action:   model.SystemLogActionToggle,
		Content:  fmt.Sprintf("%s智能体【%s】", statusText, agent.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

func splitChannelTypesString(channelTypesStr string) []int {
	var channelTypes []int
	if channelTypesStr != "" {
		strSlice := strings.Split(channelTypesStr, ",")
		for _, s := range strSlice {
			if i, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				channelTypes = append(channelTypes, i)
			}
		}
	}
	return channelTypes
}

func splitAgentTypesString(agentTypesStr string) []int {
	var agentTypes []int
	if agentTypesStr != "" {
		strSlice := strings.Split(agentTypesStr, ",")
		for _, s := range strSlice {
			if i, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
				agentTypes = append(agentTypes, i)
			}
		}
	}
	return agentTypes
}

// GetInternalUserAgents retrieves available agents for a specific internal user
// @Summary Get available agents for a specific internal user
// @Description Get all available agents for a specific internal user (including agents from user's groups and department groups, with duplicates removed)
// @Tags Agent
// @Produce json
// @Security BearerAuth
// @Success 200 {object} model.CommonResponse{data=AgentsResponse} "Success response containing agent list"
// @Router /api/agents/internal_users [get]
func GetInternalUserAgents(c *gin.Context) {
	eid := config.GetEID(c)

	userID := config.GetUserId(c)
	user, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, model.DBError.ToResponse(err))
		return
	}

	if user.Eid != eid || user.Type != model.UserTypeInternal {
		c.JSON(http.StatusForbidden, model.AuthFailed.ToResponse("User is not an internal member of the current enterprise"))
		return
	}

	var groupIDs []int64
	groupIDs, err = model.GetGroupsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	var dids []int64
	dids, err = model.GetMemberDidsByBID(eid, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	var departmentGroupIDs []int64
	if len(dids) == 0 {
		dids = []int64{0}
	}

	for _, did := range dids {
		deptGroupIDs, err := model.GetGroupsByDepartmentID(did)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
			return
		}
		departmentGroupIDs = append(departmentGroupIDs, deptGroupIDs...)
	}

	if len(departmentGroupIDs) > 0 {
		groupIDMap := make(map[int64]bool)
		for _, id := range groupIDs {
			groupIDMap[id] = true
		}

		for _, id := range departmentGroupIDs {
			if !groupIDMap[id] {
				groupIDs = append(groupIDs, id)
				groupIDMap[id] = true
			}
		}
	}

	if len(groupIDs) == 0 {
		c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
			Count:  0,
			Agents: []*model.Agent{},
		}))
		return
	}

	var allAgents []*model.Agent
	agentMap := make(map[int64]*model.Agent)

	for _, groupID := range groupIDs {
		group, err := model.GetGroupWithAgents(groupID, true)
		if err != nil {
			continue
		}

		for i := range group.Agents {
			if group.Agents[i].Enable {
				agentMap[group.Agents[i].AgentID] = group.Agents[i]
			}
		}
	}

	for _, agent := range agentMap {
		var count int64
		model.DB.Model(&model.Conversation{}).
			Where("eid = ? AND agent_id = ?", eid, agent.AgentID).
			Count(&count)

		agent.ConversationCount = count
		allAgents = append(allAgents, agent)
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(AgentsResponse{
		Count:  int64(len(allAgents)),
		Agents: allAgents,
	}))
}
