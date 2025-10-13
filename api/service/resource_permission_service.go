package service

import (
	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateResourcePermissions 更新资源权限的通用方法
// 该方法会先删除指定资源的所有现有权限，然后添加新的权限
func UpdateResourcePermissions(c *gin.Context, tx *gorm.DB, resourceID int64, resourceType string, groupIDs []int64) error {
	// 删除现有权限
	if err := tx.Where("resource_id = ? AND resource_type = ?", resourceID, resourceType).Delete(&model.ResourcePermission{}).Error; err != nil {
		return err
	}

	// 添加新权限
	if len(groupIDs) > 0 {
		// 使用map去重
		groupIDMap := make(map[int64]bool)
		uniqueGroupIDs := make([]int64, 0)

		for _, groupID := range groupIDs {
			if !groupIDMap[groupID] && groupID > 0 {
				groupIDMap[groupID] = true
				uniqueGroupIDs = append(uniqueGroupIDs, groupID)
			}
		}

		for _, groupID := range uniqueGroupIDs {
			permission := model.ResourcePermission{
				GroupID:      groupID,
				ResourceID:   resourceID,
				ResourceType: resourceType,
				Permission:   model.PermissionRead,
			}

			if err := tx.Create(&permission).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// UpdateAgentResourcePermissions 更新Agent资源权限的专用方法
// 该方法会先删除指定Agent的所有现有权限，然后根据企业类型过滤后添加新的权限
func UpdateAgentResourcePermissions(c *gin.Context, tx *gorm.DB, agentID int64, subscriptionGroupIds, userGroupIds []int64, enterprise *model.Enterprise) error {
	// 删除现有权限
	if err := tx.Where("resource_id = ? AND resource_type = ?", agentID, model.ResourceTypeAgent).Delete(&model.ResourcePermission{}).Error; err != nil {
		return err
	}

	// 统一合并两类分组ID（去重），默认不过滤企业类型（未知/空时等同于行业型）
	groupIDSet := make(map[int64]bool)
	for _, id := range subscriptionGroupIds {
		if id > 0 {
			groupIDSet[id] = true
		}
	}
	for _, id := range userGroupIds {
		if id > 0 {
			groupIDSet[id] = true
		}
	}
	// 转换为切片
	allGroupIds := make([]int64, 0, len(groupIDSet))
	for id := range groupIDSet {
		allGroupIds = append(allGroupIds, id)
	}
	// 仅当企业类型为独立站/企业内部时做类型过滤；行业型或未知类型不做过滤
	needFilter := enterprise != nil && (enterprise.Type == model.EnterpriseTypeIndependent || enterprise.Type == model.EnterpriseTypeEnterprise)
	if needFilter && len(allGroupIds) > 0 {
		var groups []model.Group
		if err := tx.Where("group_id IN ?", allGroupIds).Find(&groups).Error; err != nil {
			return err
		}
		filteredGroupIds := make([]int64, 0, len(groups))
		for _, group := range groups {
			if enterprise.Type == model.EnterpriseTypeIndependent {
				if group.GroupType == model.USER_GROUP_TYPE {
					filteredGroupIds = append(filteredGroupIds, group.GroupId)
				}
			} else if enterprise.Type == model.EnterpriseTypeEnterprise {
				if group.GroupType == model.INTERNAL_USER_GROUP_TYPE {
					filteredGroupIds = append(filteredGroupIds, group.GroupId)
				}
			}
		}
		allGroupIds = filteredGroupIds
	}

	// 添加新权限
	if len(allGroupIds) > 0 {
		for _, groupID := range allGroupIds {
			permission := model.ResourcePermission{
				GroupID:      groupID,
				ResourceID:   agentID,
				ResourceType: model.ResourceTypeAgent,
				Permission:   model.PermissionRead,
			}
			if err := tx.Create(&permission).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
