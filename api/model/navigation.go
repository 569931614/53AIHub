package model

import "errors"

type Navigation struct {
	NavigationID int64              `json:"navigation_id" gorm:"primaryKey;autoIncrement;comment:自增id"`
	Eid          int64              `json:"eid" gorm:"not null;comment:企业ID"`
	Name         string             `json:"name" gorm:"size:50;not null;comment:名称"`
	Icon         string             `json:"icon" gorm:"size:255;comment:图标URL或标识"`
	Type         int                `json:"type" gorm:"not null;comment:类型1系统内置,2外部链接,3自定义页"` // 类型(1=系统内置,2=外部链接,3=自定义页)
	JumpPath     string             `json:"jump_path" gorm:"size:255;not null;comment:跳转路径"`
	Sort         int                `json:"sort" gorm:"default:0;comment:排序"`
	Config       string             `json:"config" gorm:"type:text;not null;comment:配置"`
	Status       int                `json:"status" gorm:"not null;default:1;comment:状态1启用,2禁用"` // 状态(1=启用,2=禁用)
	Content      *NavigationContent `json:"content" gorm:"-"`
	BaseModel
}

const (
	NavigationTypeSystemBuiltIn = 1
	NavigationTypeExternalLink  = 2
	NavigationTypeCustomPage    = 3

	NavigationStatusEnabled  = 1
	NavigationStatusDisabled = 2
)

func (nav *Navigation) Create() error {
	if nav.Eid == 0 {
		return errors.New("eid不能为空")
	}
	if nav.JumpPath == "" {
		return errors.New("跳转路径不能为空")
	}

	// var count int64
	// DB.Model(&Navigation{}).Where("eid = ? AND jump_path = ?", nav.Eid, nav.JumpPath).Count(&count)
	// if count > 0 {
	// 	return errors.New("当前企业下该跳转路径已存在")
	// }

	result := DB.Create(nav)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (nav *Navigation) Update() error {
	return DB.Model(nav).Updates(nav).Error
}

func (nav *Navigation) Delete() error {
	return DB.Delete(nav).Error
}

func GetNavigationByID(eid int64, navigationID int64) (*Navigation, error) {
	var nav Navigation
	err := DB.Where("eid = ? AND navigation_id = ?", eid, navigationID).First(&nav).Error
	if err != nil {
		return nil, err
	}

	return &nav, nil
}

func (nav *Navigation) LoadContent() error {
	var content NavigationContent
	err := DB.Where("navigation_id = ?", nav.NavigationID).First(&content).Error
	if err != nil {
		return err
	}
	nav.Content = &content
	return nil
}

func UpdateNavigation(id int, updates map[string]interface{}) error {
	return DB.Model(&Navigation{}).Where("navigation_id = ?", id).Updates(updates).Error
}

func DeleteNavigation(id int) error {
	return DB.Where("navigation_id = ?", id).Delete(&Navigation{}).Error
}

func UpdateNavigationStatus(id int, status int) error {
	return DB.Model(&Navigation{}).Where("navigation_id = ?", id).Update("status", status).Error
}

func BatchUpdateNavigationSort(sortList []struct {
	ID   int
	Sort int
}) error {
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range sortList {
		if err := tx.Model(&Navigation{}).Where("navigation_id = ?", item.ID).Update("sort", item.Sort).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func GetNavigationsByEid(eid int64) ([]*Navigation, error) {
	var navs []*Navigation
	err := DB.Where("eid =?", eid).Order("sort DESC").Find(&navs).Error
	if err != nil {
		return nil, err
	}
	return navs, nil
}
