package controller

import (
	"net/http"

	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
)

// NavigationIcon 导航图标信息
type NavigationIcon struct {
	Key   string `json:"key"`   // 图标标识
	Name  string `json:"name"`  // 图标名称
	URL   string `json:"url"`   // 图标URL
	Style string `json:"style"` // 图标样式类
}

// @Summary 获取导航默认图标列表
// @Description 获取系统预置的导航图标列表
// @Tags Navigation
// @Produce json
// @Success 200 {object} model.CommonResponse{data=[]NavigationIcon} "成功响应"
// @Router /api/navigations/icons [get]
func GetNavigationIcons(c *gin.Context) {
	// 预置的默认图标列表
	icons := []NavigationIcon{
		{
			Key:   "home",
			Name:  "首页",
			URL:   "/static/images/navigation/home.svg",
			Style: "icon-home",
		},
		{
			Key:   "dashboard",
			Name:  "仪表盘",
			URL:   "/static/images/navigation/dashboard.svg",
			Style: "icon-dashboard",
		},
		{
			Key:   "chat",
			Name:  "聊天",
			URL:   "/static/images/navigation/chat.svg",
			Style: "icon-chat",
		},
		{
			Key:   "ai",
			Name:  "AI助手",
			URL:   "/static/images/navigation/ai.svg",
			Style: "icon-ai",
		},
		{
			Key:   "document",
			Name:  "文档",
			URL:   "/static/images/navigation/document.svg",
			Style: "icon-document",
		},
		{
			Key:   "setting",
			Name:  "设置",
			URL:   "/static/images/navigation/setting.svg",
			Style: "icon-setting",
		},
		{
			Key:   "user",
			Name:  "用户",
			URL:   "/static/images/navigation/user.svg",
			Style: "icon-user",
		},
		{
			Key:   "group",
			Name:  "群组",
			URL:   "/static/images/navigation/group.svg",
			Style: "icon-group",
		},
		{
			Key:   "analysis",
			Name:  "分析",
			URL:   "/static/images/navigation/analysis.svg",
			Style: "icon-analysis",
		},
		{
			Key:   "report",
			Name:  "报告",
			URL:   "/static/images/navigation/report.svg",
			Style: "icon-report",
		},
		{
			Key:   "calendar",
			Name:  "日历",
			URL:   "/static/images/navigation/calendar.svg",
			Style: "icon-calendar",
		},
		{
			Key:   "notification",
			Name:  "通知",
			URL:   "/static/images/navigation/notification.svg",
			Style: "icon-notification",
		},
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(icons))
}
