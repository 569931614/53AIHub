package router

import (
	"github.com/53AI/53AIHub/controller"
	"github.com/53AI/53AIHub/middleware"
	"github.com/53AI/53AIHub/model"
	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	apiRouter := router.Group("/api")
	// apiRouter.Use(middleware.CORS())
	apiRouter.Use(middleware.Logger())

	enterpriseRoute := apiRouter.Group("/enterprises")
	{
		enterpriseRoute.GET("/is_saas", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetIsSaas)
		enterpriseRoute.GET("/homepage", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetHomePage)

		enterpriseRoute.GET("/current", controller.GetCurrentEnterprise)

		enterpriseRoute.GET("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetEnterprise)
		enterpriseRoute.PUT("/:id", controller.UpdateEnterprise)
		enterpriseRoute.PATCH("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdateEnterpriseAttribute)
		enterpriseRoute.DELETE("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.DeleteEnterprise)
		enterpriseRoute.POST("", middleware.UserTokenAuth(model.RoleAdminUser), controller.CreateEnterprise)
		enterpriseRoute.GET("/banner", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetEnterpriseBanner)
		enterpriseRoute.PUT("/banner", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdateEnterpriseBanner)
		enterpriseRoute.GET("/template_type", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetEnterpriseTemplateType)
		enterpriseRoute.PUT("/template_type", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdateEnterpriseTemplateType)
	}

	enterpriseConfigRoute := apiRouter.Group("/enterprise-configs")
	{
		enterpriseConfigRoute.GET("", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetEnterpriseConfigTypes)
		enterpriseConfigRoute.GET("/:type", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetEnterpriseConfig)
		enterpriseConfigRoute.GET("/:type/enabled", controller.IsEnterpriseConfigEnabled)
		enterpriseConfigRoute.POST("/:type", middleware.UserTokenAuth(model.RoleAdminUser), controller.SaveEnterpriseConfig)
		enterpriseConfigRoute.PUT("/:type/toggle", middleware.UserTokenAuth(model.RoleAdminUser), controller.ToggleEnterpriseConfig)
	}

	commonRoute := apiRouter.Group("")
	{
		commonRoute.POST("/register", controller.PasswordRegister)
		commonRoute.POST("/login", controller.Login)
		commonRoute.POST("/logout", middleware.UserTokenAuth(model.RoleGuestUser), controller.Logout)
		commonRoute.POST("/sms_login", controller.SmsLogin)
		commonRoute.POST("/check_account", controller.CheckAccountExists)
		commonRoute.POST("/upload", controller.Upload)
		commonRoute.GET("/is_init", controller.IsInit)
		commonRoute.GET("/preview/*key", controller.PreviewFile)
		commonRoute.GET("/response_codes", controller.GetAllResponseCodes)
		commonRoute.POST("/reset_password", controller.ResetPassword)
	}

	emailRoute := apiRouter.Group("/email")
	{
		emailRoute.POST("/send_verification", controller.SendVerificationEmail)
		emailRoute.POST("/send_test", middleware.UserTokenAuth(model.RoleAdminUser), controller.SendTestEmail)
	}

	userRoute := apiRouter.Group("/users")
	userRoute.GET("/me", middleware.UserTokenAuth(model.RoleCommonUser), controller.GetCurrentUser)
	userRoute.PUT("/password", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdateUserPassword)
	userRoute.PATCH("/:id/mobile", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdateUserMobile)
	userRoute.PATCH("/:id/email", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdateUserEmail)
	userRoute.PUT("/me", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdateCurrentUser)
	userRoute.POST("/system_log", middleware.UserTokenAuth(model.RoleCommonUser), controller.CreateSystemLogs)
	userRoute.PUT("/:id/default_subscription", middleware.UserTokenAuth(model.RoleCommonUser), controller.SetUserToDefaultSubscription)
	userRoute.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		userRoute.POST("", controller.EnterpriseAddUser)
		userRoute.GET("", controller.EnterpriseUsers)
		userRoute.DELETE("/:id", controller.DeleteEnterpriseUser)
		userRoute.PUT("/:id", controller.UpdateEnterpriseUser)
		userRoute.GET("/:user_id/agents/:agent_id/messages", controller.GetUserMessages)
		userRoute.GET("/:user_id/conversations", controller.GetUserConversations)
		userRoute.PUT("/batch/admin", controller.SetUserAsAdmin)
		userRoute.DELETE("/batch/admin", controller.UnsetUserAsAdmin)
		userRoute.POST("/internal/batch", controller.BatchAddInternalUsers)
		userRoute.PUT("/register/to/internal", controller.RegisterUserToInternal)
		userRoute.GET("/internal", controller.GetInternalUsers)
		userRoute.PATCH("/:id/status", controller.UpdateUserStatus)
		userRoute.PUT("/internal/:id", controller.UpdateInternalUser)
		userRoute.GET("/admin", controller.EnterpriseUsers)
		userRoute.GET("/organization", controller.GetOrganizationUserList)
	}

	groupRoute := apiRouter.Group("/groups")
	groupRoute.GET("type/current/:group_type", controller.GetGroups)
	groupRoute.POST("/prompt", middleware.UserTokenAuth(model.RoleCommonUser), controller.CreateGroup)
	groupRoute.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		groupRoute.POST("", controller.CreateGroup)
		groupRoute.GET("/:id", controller.GetGroup)
		groupRoute.PUT("/:id", controller.UpdateGroup)
		groupRoute.DELETE("/:id", controller.DeleteGroup)
		groupRoute.POST("type/:group_type", controller.BatchSubmitGroups)
		groupRoute.GET("type/:group_type", controller.GetGroups)
		groupRoute.POST("/:id/agents", controller.AddAgentsToGroup)
		groupRoute.DELETE("/:id/agents", controller.RemoveAgentsFromGroup)
		groupRoute.GET("/:id/agents", controller.GetGroupAgents)
		groupRoute.DELETE("/:id/users", controller.RemoveUsersFromGroup)
		groupRoute.GET("/:id/users", controller.GetGroupUsers)
		groupRoute.POST("/:id/users/batch", controller.BatchAddUsersToGroup)
	}

	aiLinkRoute := apiRouter.Group("/ai_links")
	aiLinkRoute.GET("/current", controller.GetCurrentSiteAILinks)
	aiLinkRoute.GET("/default", controller.GetDefaultAILinks)
	aiLinkRoute.GET("/:id", middleware.UserTokenAuth(model.RoleCommonUser), controller.GetAILink)
	aiLinkRoute.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		aiLinkRoute.POST("", controller.CreateAILink)
		aiLinkRoute.GET("", controller.GetAILinks)
		aiLinkRoute.PUT("/:id", controller.UpdateAILink)
		aiLinkRoute.DELETE("/:id", controller.DeleteAILink)
		aiLinkRoute.POST("/batch/sort", controller.BatchSortAILinks)
	}

	settingRoute := apiRouter.Group("/settings")
	{
		settingRoute.POST("", middleware.UserTokenAuth(model.RoleGuestUser), controller.CreateSetting)
		settingRoute.GET("/:id", controller.GetSetting)
		settingRoute.PUT("/:id", middleware.UserTokenAuth(model.RoleGuestUser), controller.UpdateSetting)
		settingRoute.DELETE("/:id", middleware.UserTokenAuth(model.RoleGuestUser), controller.DeleteSetting)
		settingRoute.GET("", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetSettings)
		settingRoute.GET("/group/:group_name", controller.GetSettingsByGroup)
		settingRoute.GET("/key/:key", controller.GetSettingByKey)
		settingRoute.POST("/default_links", middleware.UserTokenAuth(model.RoleGuestUser), controller.BatchUpdateDefaultPromptLinks) // 批量更新默认提示词链接
		settingRoute.GET("/default_links", middleware.UserTokenAuth(model.RoleGuestUser), controller.GetDefaultPromptLinks)          // 获取默认提示词链接
	}

	channelGroup := apiRouter.Group("/channels")
	channelGroup.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		channelGroup.POST("", controller.CreateChannel)
		channelGroup.GET("", controller.GetChannels)
		channelGroup.GET("/:channel_id", controller.GetChannel)
		channelGroup.PUT("/:channel_id", controller.UpdateChannel)
		channelGroup.DELETE("/:channel_id", controller.DeleteChannel)
		channelGroup.GET("/test/:channel_id", controller.TestChannel)
		channelGroup.GET("/models", controller.ListAllModels)
	}

	agentGroup := apiRouter.Group("/agents")
	agentGroup.GET("/current", controller.GetCurrentAgents)
	agentGroup.GET("/available", controller.GetAvailableAgents)
	agentGroup.Use(middleware.UserTokenAuth(model.RoleGuestUser))
	{
		agentGroup.POST("", controller.CreateAgent)
		agentGroup.GET("", controller.GetAgents)
		agentGroup.GET("/group", controller.GetAgentsByGroup)
		agentGroup.GET("/:agent_id", controller.GetAgent)
		agentGroup.PUT("/:agent_id", controller.UpdateAgent)
		agentGroup.DELETE("/:agent_id", controller.DeleteAgent)
		agentGroup.GET("/:agent_id/messages", controller.GetMessagesByUserAndAgent)
		agentGroup.PATCH("/:agent_id/status", controller.UpdateAgentStatus)
		agentGroup.GET("/internal_users", controller.GetInternalUserAgents)
		agentGroup.GET("/:agent_id/conversations", controller.GetAgentConversations)
	}

	conversationGroup := apiRouter.Group("/conversations")
	conversationGroup.Use(middleware.UserTokenAuth(model.RoleGuestUser))
	{
		conversationGroup.POST("", controller.CreateConversation)
		conversationGroup.GET("", controller.GetConversations)
		conversationGroup.GET("/:conversation_id", controller.GetConversation)
		conversationGroup.PUT("/:conversation_id", controller.UpdateConversation)
		conversationGroup.DELETE("/:conversation_id", controller.DeleteConversation)
		//conversationGroup.POST("/:conversation_id/messages", controller.CreateMessage)
		conversationGroup.GET("/:conversation_id/messages", controller.GetMessagesByConversation)
	}

	subscription := apiRouter.Group("/subscriptions")
	{
		subscription.GET("/settings", controller.GetSubscriptionList)
		subscription.
			POST("/batch", middleware.UserTokenAuth(model.RoleAdminUser), controller.BatchSubscriptionOperation)
	}

	providerRouter := apiRouter.Group("/providers")
	providerRouter.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		providerRouter.POST("", controller.CreateProvider)
		providerRouter.GET("", controller.GetProviders)
		providerRouter.PUT("/:id", controller.UpdateProvider)
		providerRouter.DELETE("/:id", controller.DeleteProvider)
	}

	callbackRouter := apiRouter.Group("/callback")
	{
		callbackRouter.GET("/cozecn/auth/:eid", controller.CozeCallBack)
		callbackRouter.GET("/cozecom/auth/:eid", controller.CozeCallBack)
	}

	cozeRouter := apiRouter.Group("/coze")
	cozeRouter.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		cozeRouter.GET("/workspaces", controller.GetCozeAllWorkspaces)
		cozeRouter.GET("/workspaces/:workspace_id/bots", controller.GetCozeAllBots)
	}

	AppBuilderRouter := apiRouter.Group("/appbuilder")
	AppBuilderRouter.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		AppBuilderRouter.GET("/bots", controller.GetAppBuilderAllBots)
	}

	ai53Router := apiRouter.Group("/53ai")
	ai53Router.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		ai53Router.GET("/bots", controller.Get53AIAllBots)
		ai53Router.GET("/workflows", controller.Get53AIAllWorkflows)
		ai53Router.GET("/parameters/:botId", controller.Get53AIAppParameters)
	}

	apiV1Router := router.Group("/v1")
	apiV1Router.Use(middleware.CORS())
	apiV1Router.Use(middleware.Logger())
	apiV1Router.Use(middleware.RelayTokenAuth())
	{
		apiV1Router.POST("/chat/completions", controller.Relay)
		apiV1Router.POST("/workflow/run", controller.WorkflowRun)
		apiV1Router.POST("/rerank", controller.Rerank)
	}

	paySettingRouter := apiRouter.Group("/pay_settings")
	paySettingRouter.GET("/type/:type", controller.GetPaySettingByType)
	{
		paySettingRouter.GET("", middleware.UserTokenAuth(model.RoleGuestUser), controller.GetPaySettings)
		paySettingRouter.GET("/:id", middleware.UserTokenAuth(model.RoleGuestUser), controller.GetPaySetting)
		paySettingRouter.POST("", middleware.UserTokenAuth(model.RoleAdminUser), controller.CreatePaySetting)
		// paySettingRouter.PUT("/:id", controller.UpdatePaySetting)
		paySettingRouter.DELETE("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.DeletePaySetting)
		paySettingRouter.PATCH("/:id/config", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdatePayConfig)
		paySettingRouter.PATCH("/:id/status", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdatePayStatus)
	}

	orderRouter := apiRouter.Group("/orders")
	{
		orderRouter.POST("", middleware.UserTokenAuth(model.RoleCommonUser), controller.CreateOrder)
		orderRouter.PUT("/:id/manual", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdateManualTransferOrder)
		orderRouter.GET("", middleware.UserTokenAuth(model.RoleCommonUser), controller.GetOrders)
		orderRouter.GET("/me", middleware.UserTokenAuth(model.RoleCommonUser), controller.GetOrders)
		orderRouter.GET("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetOrder)
		orderRouter.PATCH("/:id/status", middleware.UserTokenAuth(model.RoleAdminUser), controller.UpdateOrderStatus) // Only manual transfers can be marked as paid
		orderRouter.DELETE("/:id", middleware.UserTokenAuth(model.RoleAdminUser), controller.DeleteOrder)             // Only manual transfers can be deleted, but paid ones cannot be deleted
		orderRouter.GET("/status/:order_id", middleware.UserTokenAuth(model.RoleCommonUser), controller.QueryOrderStatus)
		orderRouter.POST("/:id/confirm", middleware.UserTokenAuth(model.RoleCommonUser), controller.ConfirmManualPayment)
		orderRouter.GET("/user", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetUserOrders)
		orderRouter.POST("/:id/close", middleware.UserTokenAuth(model.RoleCommonUser), controller.CloseOrder)
		orderRouter.GET("/trade/:order_id", middleware.UserTokenAuth(model.RoleAdminUser), controller.QueryTradeOrder)
		orderRouter.POST("/trade/:order_id/refund", middleware.UserTokenAuth(model.RoleAdminUser), controller.RefunTradeOrder)
	}

	paymentRouter := apiRouter.Group("/payment")
	{
		paymentRouter.GET("/available", controller.GetAvailablePayTypes)
		// Payment notification routes
		paymentRouter.POST("/wechat/notify/:id", controller.WechatPayNotify)
		paymentRouter.POST("/alipay/notify/:id", controller.AlipayNotify)
	}

	// Department routes
	departmentGroup := apiRouter.Group("/departments")
	departmentGroup.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		departmentGroup.POST("", controller.CreateDepartment)
		departmentGroup.GET("", controller.GetDepartments)
		departmentGroup.GET("/:did", controller.GetDepartment)
		departmentGroup.PUT("/:did", controller.UpdateDepartment)
		departmentGroup.DELETE("/:did", controller.DeleteDepartment)
		departmentGroup.GET("/children/:pdid", controller.GetChildDepartments)
		departmentGroup.GET("/tree", controller.GetDepartmentTree)
		departmentGroup.POST("/sync/:from", controller.SyncOrganization)
		departmentGroup.POST("/bind-member", controller.DepartmentBindMember)
		departmentGroup.DELETE("/bind-member", controller.DepartmentUnbindMember)
	}

	promptGroup := apiRouter.Group("/prompts")
	{
		promptGroup.GET("", controller.GetPrompts)
		promptGroup.GET("/admin", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetPrompts)
		promptGroup.POST("/system", middleware.UserTokenAuth(model.RoleAdminUser), controller.CreatePrompt)
		promptGroup.POST("/personal", middleware.UserTokenAuth(model.RoleCommonUser), controller.CreatePrompt)
		promptGroup.GET("/:pid", controller.GetPrompt)
		promptGroup.PUT("/:pid", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdatePrompt)
		promptGroup.DELETE("/:pid", middleware.UserTokenAuth(model.RoleCommonUser), controller.DeletePrompt)
		promptGroup.PATCH("/:pid/like", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdatePromptLike)
		promptGroup.GET("/:pid/groups", middleware.UserTokenAuth(model.RoleCommonUser), controller.GetPromptGroups)
		promptGroup.PATCH("/:pid/status", middleware.UserTokenAuth(model.RoleCommonUser), controller.UpdatePromptStatus)
	}

	navigationRoute := apiRouter.Group("/navigations")
	navigationRoute.GET("", controller.GetNavigations)
	navigationRoute.POST("/init", controller.InitSystemNavigation)
	navigationRoute.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		navigationRoute.GET("/:nav_id", controller.GetNavigation)
		navigationRoute.POST("", controller.CreateNavigation)
		navigationRoute.PUT("/:nav_id", controller.UpdateNavigation)
		navigationRoute.DELETE("/:nav_id", controller.DeleteNavigation)
		navigationRoute.PATCH("/:nav_id/status", controller.UpdateNavigationStatus)
		navigationRoute.POST("/sort", controller.SortNavigations)
		navigationRoute.POST("/:nav_id/content", controller.CreateNavigationContent)
		navigationRoute.GET("/:nav_id/content", controller.GetNavigationContent)
	}

	systemLogRouter := apiRouter.Group("/system_logs")
	systemLogRouter.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		systemLogRouter.GET("/modules", controller.GetModules)
		systemLogRouter.GET("/actions", controller.GetActions)
		systemLogRouter.GET("", controller.GetSystemLogs)
	}

	maxKB := apiRouter.Group("/maxkb")
	{
		maxKB.GET("/application/profile", middleware.UserTokenAuth(model.RoleAdminUser), controller.GetMaxKBApplicationProfile)
	}

	difyRouter := apiRouter.Group("/dify")
	difyRouter.Use(middleware.UserTokenAuth(model.RoleAdminUser))
	{
		difyRouter.GET("/info/:channelId", controller.GetDifyAppInfo)
		difyRouter.GET("/parameters/:channelId", controller.GetDifyAppParameters)
	}

	sharesAuth := apiRouter.Group("/shares")
	sharesAuth.Use(middleware.UserTokenAuth(model.RoleGuestUser))
	{
		sharesAuth.POST("", controller.CreateShare)
	}

	sharesPublic := apiRouter.Group("/shares")
	{
		sharesPublic.GET("/:share_id", controller.GetShare)
	}
}
