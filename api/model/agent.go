package model

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Agent struct {
	AgentID      int64  `json:"agent_id" gorm:"primaryKey;autoIncrement"`
	Eid          int64  `json:"eid" gorm:"not null;index"`
	Name         string `json:"name" gorm:"not null"`
	Logo         string `json:"logo" gorm:"not null"`
	Sort         int    `json:"sort" gorm:"default:0"`
	Description  string `json:"description" gorm:"not null"`
	ChannelType  int    `json:"channel_type" gorm:"default:0"`
	Model        string `json:"model" gorm:"not null"`
	Prompt       string `json:"prompt" gorm:"not null"`
	Configs      string `json:"configs" gorm:"not null;type:text"`
	Tools        string `json:"tools" gorm:"not null;type:text"`
	GroupID      int64  `json:"group_id" gorm:"type:int;default:0;not null"`
	UseCases     string `json:"use_cases" gorm:"not null;type:text"`
	CreatedBy    int64  `json:"created_by" gorm:"not null"`
	CustomConfig string `json:"custom_config" gorm:"not null;type:text"`
	Settings     string `json:"settings" gorm:"not null;type:text"`
	// UserGroupIds 表示内部用户组ID集合
	UserGroupIds []int64 `json:"user_group_ids" gorm:"-"`
	// SubscriptionGroupIds 表示注册用户订阅分组ID集合
	SubscriptionGroupIds []int64 `json:"subscription_group_ids" gorm:"-"`
	Enable               bool    `json:"enable" gorm:"default:false;comment:enable status"`
	ConversationCount    int64   `json:"conversation_count" gorm:"-"`
	AgentType            int     `json:"agent_type" gorm:"default:0"`
	BaseModel
}

const (
	AgentTypeApp      = 0
	AgentTypeWorkflow = 1
)

func (agent *Agent) Create() error {
	if agent.Eid == 0 {
		return errors.New("eid is empty")
	}
	// check if name exists
	var count int64
	DB.Model(&Agent{}).Where("eid = ? AND name = ?", agent.Eid, agent.Name).Count(&count)
	if count > 0 {
		return errors.New("name already exists")
	}

	result := DB.Create(agent)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (agent *Agent) Update() error {
	err := DB.Model(agent).Updates(agent).Error
	return err
}

func (agent *Agent) Delete() error {
	err := DB.Delete(agent).Error
	return err
}

func GetAgentByID(eid int64, agentID int64) (*Agent, error) {
	var agent Agent
	err := DB.Where("eid = ? AND agent_id = ?", eid, agentID).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func GetAgentListWithIDs(eid int64, keyword string, group_id int64, permittedAgentIDs []int64, channel_types []int, agent_types []int, offset int, limit int) (count int64, agents []*Agent, err error) {
	db := DB.Model(&Agent{}).Where("eid = ?", eid)
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	if group_id != 0 {
		db = db.Where("group_id =?", group_id)
	}

	if len(channel_types) > 0 {
		db = db.Where("channel_type IN?", channel_types)
	}

	if len(agent_types) > 0 {
		db = db.Where("agent_type IN?", agent_types)
	}

	if permittedAgentIDs != nil {
		if len(permittedAgentIDs) == 0 {
			return 0, []*Agent{}, nil
		}
		db = db.Where("agent_id IN ?", permittedAgentIDs)
	}

	db.Count(&count)

	err = db.Offset(offset).Limit(limit).Order("sort DESC").Order("agent_id DESC").Find(&agents).Error

	return count, agents, err
}

func GetAvailableAgentList(eid int64, agent_types []int, offset int, limit int) (count int64, agents []*Agent, err error) {
	db := DB.Model(&Agent{}).Where("eid = ?", eid).
		Where("Enable = ?", true)

	if len(agent_types) > 0 {
		db = db.Where("agent_type IN?", agent_types)
	}

	db.Count(&count)

	err = db.Offset(offset).Limit(limit).Order("sort DESC").Order("agent_id DESC").Find(&agents).Error

	return count, agents, err
}

func (a *Agent) GetUserGroupIds() ([]int64, error) {
	var permissions []ResourcePermission
	groupIds := make([]int64, 0)
	seen := make(map[int64]bool)

	err := DB.Where("resource_id = ? AND resource_type = ?", a.AgentID, ResourceTypeAgent).Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	for _, p := range permissions {
		if !seen[p.GroupID] {
			seen[p.GroupID] = true
			groupIds = append(groupIds, p.GroupID)
		}
	}

	return groupIds, nil
}

// LoadUserGroupIds
func (a *Agent) LoadUserGroupIds() error {
	ids, err := a.GetUserGroupIds()
	if err != nil {
		return err
	}
	a.UserGroupIds = ids
	return nil
}

func (a *Agent) LoadConversationCount() error {
	var count int64
	err := DB.Model(&Conversation{}).Where("agent_id =?", a.AgentID).Count(&count).Error
	if err != nil {
		return err
	}
	a.ConversationCount = count
	return nil
}

func UpdateAgentStatus(eid, agentID int64, enable *bool) error {
	return DB.Model(&Agent{}).
		Where("eid = ? AND agent_id = ?", eid, agentID).
		Update("enable", enable).Error
}

func GetAgentCountByEID(eid int64) (int64, error) {
	var count int64
	err := DB.Model(&Agent{}).Where("eid =?", eid).Count(&count).Error
	return count, err
}

func (a *Agent) GetProviderID() int64 {
	if a.CustomConfig == "" {
		return 0
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(a.CustomConfig), &config); err != nil {
		return 0
	}

	if providerID, exists := config["provider_id"]; exists {
		switch v := providerID.(type) {
		case float64:
			return int64(v)
		case int:
			return int64(v)
		case int64:
			return v
		case string:
			// Try to parse string as number
			if num, err := strconv.ParseInt(v, 10, 64); err == nil {
				return num
			}
		}
	}

	return 0
}

// LoadGroupIdsByType loads both subscription and internal user group IDs for the agent
func (a *Agent) LoadGroupIdsByType() error {
	// Query resource permissions joined with groups to get group types
	type row struct {
		GroupID   int64
		GroupType int64
	}
	var rows []row
	if err := DB.Table("resource_permissions").
		Select("resource_permissions.group_id as group_id, groups.group_type as group_type").
		Joins("JOIN groups ON groups.group_id = resource_permissions.group_id").
		Where("resource_permissions.resource_id = ? AND resource_permissions.resource_type = ?", a.AgentID, ResourceTypeAgent).
		Find(&rows).Error; err != nil {
		return err
	}

	// initialize as empty slices (not nil) so JSON encodes [] instead of null when empty
	subIDs := make([]int64, 0)
	userIDs := make([]int64, 0)
	seenSub := make(map[int64]bool)
	seenUser := make(map[int64]bool)
	for _, r := range rows {
		if r.GroupType == USER_GROUP_TYPE { // 订阅分组
			if !seenSub[r.GroupID] {
				seenSub[r.GroupID] = true
				subIDs = append(subIDs, r.GroupID)
			}
		} else if r.GroupType == INTERNAL_USER_GROUP_TYPE { // 内部用户组
			if !seenUser[r.GroupID] {
				seenUser[r.GroupID] = true
				userIDs = append(userIDs, r.GroupID)
			}
		}
	}
	// Assign back (empty slices when none found)
	a.SubscriptionGroupIds = subIDs
	a.UserGroupIds = userIDs
	return nil
}
