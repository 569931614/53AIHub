package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/53AI/53AIHub/common/utils"
	"github.com/53AI/53AIHub/config"
	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
)

// @Summary 获取导航列表
// @Description 获取当前企业的所有导航项
// @Tags Navigation
// @Produce json
// @Success 200 {object} model.CommonResponse{data=[]model.Navigation} "成功响应"
// @Router /api/navigations [get]
func GetNavigations(c *gin.Context) {
	eid := config.GetEID(c)
	user, err := model.GetLoginUser(c)

	if err == nil {
		eid = user.Eid
	}

	navigations, err := model.GetNavigationsByEid(eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	for _, nav := range navigations {
		if nav.Type == 3 { // 添加类型判断
			nav.LoadContent()
		}
		// 处理空图标URL，设置默认图标
		if nav.Icon == "" {
			nav.Icon = "/static/images/navigation/default.svg"
		}

	}
	c.JSON(http.StatusOK, model.Success.ToResponse(navigations))
}

type NavigationRequest struct {
	Name     string `json:"name" binding:"required"`      // 导航名称
	JumpPath string `json:"jump_path" binding:"required"` // 跳转路径
	Config   string `json:"config" binding:"required"`    // 配置信息 json字符串
	Sort     int    `json:"sort"`                         // 排序
	Icon     string `json:"icon"`                         // 图标
	Type     int    `json:"type"`                         // 导航类型 2: 外部链接 3: 自定义页面
}

// @Summary 创建导航
// @Description 创建新的导航项
// @Tags Navigation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param navigation body NavigationRequest true "导航数据"
// @Success 200 {object} model.CommonResponse{data=model.Navigation} "成功响应"
// @Router /api/navigations [post]
func CreateNavigation(c *gin.Context) {
	var req NavigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	if req.Type == model.NavigationTypeSystemBuiltIn {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(errors.New("系统内置导航不能创建")))
		return
	}

	// 处理空图标URL，设置默认图标
	icon := req.Icon
	if icon == "" {
		icon = "/static/images/navigation/default.svg"
	}

	//从DTO映射到持久化模型
	nav := model.Navigation{
		Name:     req.Name,
		JumpPath: req.JumpPath,
		Config:   req.Config,
		Sort:     req.Sort,
		Icon:     icon,
		Type:     req.Type,
		Eid:      config.GetEID(c), // 服务端填充租户/企业ID，避免被客户端篡改
	}

	//持久化
	if err := nav.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(err))
		return
	}

	log := model.SystemLog{
		Eid:      nav.Eid,
		UserID:   config.GetUserId(c),
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleNavigation,
		Action:   model.SystemLogActionCreate,
		Content:  fmt.Sprintf("新建导航【%s】", nav.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(nav))
}

// @Summary 获取导航详情
// @Description 根据ID获取导航项详情
// @Tags Navigation
// @Produce json
// @Security BearerAuth
// @Param nav_id path int true "导航ID"
// @Success 200 {object} model.CommonResponse{data=model.Navigation} "成功响应"
// @Router /api/navigations/{nav_id} [get]
func GetNavigation(c *gin.Context) {
	navID, _ := strconv.Atoi(c.Param("nav_id"))
	eid := config.GetEID(c)
	nav, err := model.GetNavigationByID(eid, int64(navID))
	if err != nil {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return
	}
	if nav.Type == 3 { // 添加类型判断
		nav.LoadContent()
	}

	// 处理空图标URL，设置默认图标
	if nav.Icon == "" {
		nav.Icon = "/static/images/navigation/default.svg"
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(nav))
}

// @Summary 更新导航
// @Description 更新现有导航项
// @Tags Navigation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param nav_id path int true "导航ID"
// @Param request body NavigationRequest true "更新数据"
// @Success 200 {object} model.CommonResponse "成功响应"
// @Router /api/navigations/{nav_id} [put]
func UpdateNavigation(c *gin.Context) {
	navID, _ := strconv.Atoi(c.Param("nav_id"))
	eid := config.GetEID(c)
	var req NavigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	nav, err := model.GetNavigationByID(eid, int64(navID))
	if handleNotFound(c, nav, err) {
		return
	}
	oldNav := *nav

	// 处理空图标URL，设置默认图标
	icon := req.Icon
	if icon == "" {
		icon = "/static/images/navigation/default.svg"
	}

	if err := model.UpdateNavigation(navID, map[string]interface{}{
		"name":      req.Name,
		"icon":      req.Icon,
		"jump_path": req.JumpPath,
		"sort":      req.Sort,
		"config":    req.Config,
		"type":      req.Type,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	nav, err = model.GetNavigationByID(eid, int64(navID))
	fieldMap := map[string]string{
		"Name": "名称",
	}
	model.LogEntityChange(
		fmt.Sprintf("导航【%s】", oldNav.Name),
		model.SystemLogActionUpdate,
		eid,
		config.GetUserId(c),
		config.GetUserNickname(c),
		model.SystemLogModuleNavigation,
		oldNav,
		nav,
		utils.GetClientIP(c),
		fieldMap,
	)

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

// @Summary 删除导航
// @Description 根据ID删除导航项
// @Tags Navigation
// @Produce json
// @Security BearerAuth
// @Param nav_id path int true "导航ID"
// @Success 200 {object} model.CommonResponse "成功响应"
// @Router /api/navigations/{nav_id} [delete]
func DeleteNavigation(c *gin.Context) {
	navID, _ := strconv.Atoi(c.Param("nav_id"))
	eid := config.GetEID(c)

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var nav model.Navigation
	err := tx.Where("eid = ? AND navigation_id = ?", eid, navID).First(&nav).Error
	if handleNotFound(c, &nav, err) {
		tx.Rollback()
		return
	}

	if nav.Type == model.NavigationTypeCustomPage {
		if err := tx.Where("navigation_id = ?", navID).Delete(&model.NavigationContent{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
	}

	if err := tx.Where("navigation_id = ?", navID).Delete(&model.Navigation{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	log := model.SystemLog{
		Eid:      nav.Eid,
		UserID:   config.GetUserId(c),
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleNavigation,
		Action:   model.SystemLogActionDelete,
		Content:  fmt.Sprintf("删除导航【%s】", nav.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

type UpdateNavigationStatusRequest struct {
	Status int `json:"status"`
}

// @Summary 更新导航状态
// @Description 启用/禁用导航项
// @Tags Navigation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param nav_id path int true "导航ID"
// @Param request body UpdateNavigationStatusRequest true "状态数据"
// @Success 200 {object} model.CommonResponse "成功响应"
// @Router /api/navigations/{nav_id}/status [patch]
func UpdateNavigationStatus(c *gin.Context) {
	navID, _ := strconv.Atoi(c.Param("nav_id"))

	req := &UpdateNavigationStatusRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	eid := config.GetEID(c)
	nav, err := model.GetNavigationByID(eid, int64(navID))
	if handleNotFound(c, nav, err) {
		return
	}

	if err := model.UpdateNavigationStatus(navID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}

	statusText := "显示"
	if req.Status == model.NavigationStatusDisabled {
		statusText = "不显示"
	}
	log := model.SystemLog{
		Eid:      eid,
		UserID:   config.GetUserId(c),
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleNavigation,
		Action:   model.SystemLogActionToggle,
		Content:  fmt.Sprintf("%s导航【%s】", statusText, nav.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

// 在文件顶部定义全局类型
type NavigationSortItem struct {
	ID   int `json:"id" binding:"required"`
	Sort int `json:"sort"`
}

// 修改SortNavigations方法
// @Summary 排序导航项
// @Description 批量更新导航项排序
// @Tags Navigation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sortList body []NavigationSortItem true "排序列表"
// @Success 200 {object} model.CommonResponse "成功响应"
// @Router /api/navigations/sort [post]
func SortNavigations(c *gin.Context) {
	var sortList []NavigationSortItem // 使用新定义的类型

	if err := c.ShouldBindJSON(&sortList); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	paramList := make([]struct {
		ID   int
		Sort int
	}, len(sortList))
	for i, item := range sortList {
		paramList[i].ID = item.ID
		paramList[i].Sort = item.Sort
	}

	if err := model.BatchUpdateNavigationSort(paramList); err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	c.JSON(http.StatusOK, model.Success.ToResponse(nil))
}

func handleNotFound(c *gin.Context, nav *model.Navigation, err error) bool {
	if err != nil || nav == nil || nav.Eid != config.GetEID(c) {
		c.JSON(http.StatusNotFound, model.NotFound.ToResponse(nil))
		return true
	}
	return false
}

// @Summary 初始化系统导航
// @Description 初始化企业系统默认导航
// @Tags Navigation
// @Accept json
// @Produce json
// @Param request body []NavigationRequest true "导航列表"
// @Success 201 {object} model.CommonResponse{data=[]model.Navigation} "创建成功"
// @Router /api/navigations/init [post]
func InitSystemNavigation(c *gin.Context) {
	eid := config.GetEID(c)
	user, err := model.GetLoginUser(c)
	if err == nil {
		eid = user.Eid
	}

	navs, err := model.GetNavigationsByEid(eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
		return
	}
	if len(navs) > 0 {
		c.JSON(http.StatusOK, model.Success.ToResponse(navs))
		return
	}

	var req []NavigationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	tx := model.DB.Begin()
	for _, nav := range req {
		createdTime := time.Now().UTC().UnixMilli()
		nav := model.Navigation{
			Name:     nav.Name,
			JumpPath: nav.JumpPath,
			Sort:     nav.Sort,
			Config:   nav.Config,
			Eid:      eid,
			Type:     model.NavigationTypeSystemBuiltIn,
			BaseModel: model.BaseModel{
				CreatedTime: createdTime,
				UpdatedTime: createdTime,
			},
		}
		if err := tx.Create(&nav).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
		navs = append(navs, &nav)
	}
	tx.Commit()

	c.JSON(http.StatusCreated, model.Success.ToResponse(nil))
}

type CreateNavigationContentRequest struct {
	HtmlContent string `json:"html_content" binding:"required"`
}

// @Summary 创建/更新导航内容
// @Description 当导航类型为自定义页时，自动创建或更新内容
// @Tags Navigation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param nav_id path int true "导航ID"
// @Param request body CreateNavigationContentRequest true "内容数据"
// @Success 200 {object} model.CommonResponse{data=model.NavigationContent} "成功响应"
// @Router /api/navigations/{nav_id}/content [post]
func CreateNavigationContent(c *gin.Context) {
	navID, _ := strconv.Atoi(c.Param("nav_id"))
	eid := config.GetEID(c)

	nav, err := model.GetNavigationByID(eid, int64(navID))
	if handleNotFound(c, nav, err) {
		return
	}

	if nav.Type != model.NavigationTypeCustomPage {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse("非自定义页导航类型"))
		return
	}

	var req CreateNavigationContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse(err))
		return
	}

	tx := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var content *model.NavigationContent
	if nav.Type == 3 { // 添加类型判断
		nav.LoadContent()
	}
	if nav.Content != nil {
		content = nav.Content
		content.HtmlContent = req.HtmlContent
		if err := content.Update(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
	} else {
		content = &model.NavigationContent{
			NavigationID: int64(navID),
			HtmlContent:  req.HtmlContent,
		}
		if err := content.Create(); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.DBError.ToResponse(nil))
			return
		}
	}

	tx.Commit()

	log := model.SystemLog{
		Eid:      nav.Eid,
		UserID:   config.GetUserId(c),
		Nickname: config.GetUserNickname(c),
		Module:   model.SystemLogModuleNavigation,
		Action:   model.SystemLogActionUpdate,
		Content:  fmt.Sprintf("编辑导航【%s】-页面编辑 ", nav.Name),
		IP:       utils.GetClientIP(c),
	}
	model.CreateSystemLog(&log)

	c.JSON(http.StatusOK, model.Success.ToResponse(content))
}

func GetNavigationContent(c *gin.Context) {
	navID, err := strconv.ParseInt(c.Param("nav_id"), 10, 64)
	if err != nil || navID <= 0 {
		c.JSON(http.StatusBadRequest, model.ParamError.ToResponse("无效的导航ID"))
		return
	}

	content, err := model.GetNavigationContentByID(navID)
	if err != nil {
		c.JSON(http.StatusOK, model.NotFound.ToResponse(nil))
		return
	}

	c.JSON(http.StatusOK, model.Success.ToResponse(content))
}
