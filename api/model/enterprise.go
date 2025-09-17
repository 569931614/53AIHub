package model

import (
	"errors"
	"time"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/common/utils/system"
	"github.com/53AI/53AIHub/config"
	"gorm.io/gorm"
)

type Enterprise struct {
	Eid              int64             `json:"id" gorm:"primaryKey;autoIncrement"`
	DisplayName      string            `json:"display_name" gorm:"not null" binding:"required" example:"Enterprise Name"`
	Logo             string            `json:"logo" gorm:"not null" binding:"required" example:"http://a.com/a.jpg"`
	Ico              string            `json:"ico" gorm:"type:varchar(100);default:'';not null" example:"http://a.com/favicon.ico"`
	Keywords         string            `json:"keywords" gorm:"type:text;not null" example:"AI,Hub,Agent"`
	Copyright        string            `json:"copyright" gorm:"type:varchar(255);default:'';not null" example:"© 2023 Company Name"`
	Type             string            `json:"type" gorm:"type:varchar(20);default:independent;not null;comment:'站点类型：independent、enterprise、industry'" example:"independent、enterprise、industry"`
	Banner           string            `json:"banner" gorm:"type:text;not null" example:"http://a.com/banner.jpg"`
	Language         string            `json:"language" gorm:"type:varchar(10);default:zh-cn;not null" binding:"required" example:"En"`
	Timezone         string            `json:"timezone" gorm:"type:varchar(20);default:UTC+8;not null" binding:"required" example:"UTC+8"`
	Domain           string            `json:"domain" gorm:"not null" binding:"required" example:"http://a.com"`
	Slogan           string            `json:"slogan" gorm:"not null" binding:"required" example:"Slogan Test"`
	Status           int               `json:"status" gorm:"type:int;default:1;not null" example:"1"`
	Description      string            `json:"description" gorm:"not null" example:"Description Test"`
	TemplateType     string            `json:"template_type" gorm:"type:text;not null" example:"default"`
	LayoutType       string            `json:"layout_type" gorm:"type:varchar(10);default:1;not null" example:"1"`
	WecomCorpID      string            `json:"wecom_corp_id" gorm:"type:varchar(100);default:'';not null" example:""`
	WecomInstallInfo *WecomInstallInfo `json:"wecom_install_info" gorm:"-"`
	BaseModel
}

type WecomInstallInfo struct {
	InstallWecomApp int           `json:"install_wecom_app" default:"0"`
	AuthCorpInfo    *AuthCorpInfo `json:"auth_corp_info"`
}

const (
	EnterpriseStatusNormal       = 1
	EnterpriseStatusDisabled     = 2
	EnterpriseStatusNotActivated = 0

	EnterpriseTypeIndependent = "independent"
	EnterpriseTypeEnterprise  = "enterprise"
	EnterpriseTypeIndustry    = "industry"
)

type EnterpriseTypeDescription struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var enterpriseTypeDescMap = map[string]string{
	EnterpriseTypeIndependent: "AI 独立站",
	EnterpriseTypeEnterprise:  "企业 AI 门户",
	EnterpriseTypeIndustry:    "行业 AI 门户",
}

func GetEnterpriseTypeDescription(key string) string {
	if desc, ok := enterpriseTypeDescMap[key]; ok {
		return desc
	}
	return ""
}

func GetAllEnterpriseTypeDescriptions() []EnterpriseTypeDescription {
	descriptions := make([]EnterpriseTypeDescription, 0, len(enterpriseTypeDescMap))
	for k, v := range enterpriseTypeDescMap {
		descriptions = append(descriptions, EnterpriseTypeDescription{Key: k, Value: v})
	}
	return descriptions
}

func GetEnterpriseModel(id int64) (*Enterprise, error) {
	var enterprise Enterprise
	err := DB.First(&enterprise, id).Error

	if err != nil {
		return nil, err
	}
	return &enterprise, nil
}

// GetEnterpriseByID Get enterprise information by ID
func GetEnterpriseByID(eid int64) (*Enterprise, error) {
	var enterprise Enterprise
	err := DB.Where("eid = ?", eid).First(&enterprise).Error
	if err != nil {
		return nil, err
	}
	return &enterprise, nil
}

func CreateEnterpriseModel(enterprise *Enterprise) error {
	return DB.Create(enterprise).Error
}

func (enterprise *Enterprise) Update() error {
	err := DB.Model(enterprise).
		Where("eid = ?", enterprise.Eid). // 添加主键条件
		Select("*").                      // 明确指定更新所有字段
		Updates(enterprise).Error
	return err
}

func (enterprise *Enterprise) Delete() error {
	return DB.Delete(enterprise).Error
}

func (enterprise *Enterprise) PartialUpdateEnterprise(updateData map[string]interface{}) error {
	// 定义允许更新的字段白名单
	allowedFields := map[string]bool{
		"display_name":  true,
		"logo":          true,
		"ico":           true,
		"keywords":      true,
		"copyright":     true,
		"type":          true,
		"banner":        true,
		"language":      true,
		"timezone":      true,
		"domain":        true,
		"slogan":        true,
		"status":        true,
		"description":   true,
		"template_type": true,
		"layout_type":   true,
	}

	// 过滤非法字段
	filteredUpdate := make(map[string]interface{})
	for field, value := range updateData {
		if allowedFields[field] {
			filteredUpdate[field] = value
		}
	}

	// 自动添加更新时间
	filteredUpdate["updated_time"] = time.Now().Unix()

	return DB.Model(enterprise).
		Where("eid = ?", enterprise.Eid).
		Updates(filteredUpdate).Error
}

func GetEnterpriseName(eid int64) (string, error) {
	var displayName string
	err := DB.Model(&Enterprise{}).
		Select("display_name").
		Where("eid = ?", eid).
		Limit(1).
		Scan(&displayName).Error

	if err != nil {
		return "", err
	}

	return displayName, nil
}

// InitializeSystem checks if the system is in SAAS mode
// If not and no enterprise exists, it creates a default enterprise, user group, and admin user
func InitializeSystem() error {
	// Check if in SAAS mode
	isSaas := config.IS_SAAS
	if isSaas {
		logger.SysLogf("System is running in SAAS mode, skipping default enterprise and user initialization")
		return nil
	}

	// Check if any enterprise exists
	var count int64
	if err := DB.Model(&Enterprise{}).Count(&count).Error; err != nil {
		logger.SysLogf("Failed to check enterprise records: %s", err.Error())
		return err
	}

	// If enterprises already exist, no need to initialize
	if count > 0 {
		logger.SysLogf("Enterprise records already exist, skipping default enterprise and user initialization")
		return nil
	}

	logger.SysLogf("System is running in non-SAAS mode with no enterprise records, starting initialization")

	// Start a transaction
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		logger.SysLogf("Failed to begin transaction: %s", tx.Error.Error())
		return tx.Error
	}

	// 1. Create default enterprise
	enterprise := Enterprise{
		DisplayName:  "53AI Hub",
		Logo:         "https://img.ibos.cn/common/agenthub/agent/53ai.png",
		Language:     "zh-cn",
		Description:  "Default enterprise created during system initialization",
		LayoutType:   "1",
		TemplateType: "",
		Status:       EnterpriseStatusNormal,
	}

	if err := tx.Create(&enterprise).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create default enterprise: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created default enterprise, ID: %d", enterprise.Eid)

	// 2. Create default user group
	defaultGroup := Group{
		Eid:       enterprise.Eid,
		CreatedBy: 0, // System created
		GroupName: "免费版",
		GroupType: USER_GROUP_TYPE,
		Sort:      0,
	}

	if err := tx.Create(&defaultGroup).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create default user group: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created default user group, ID: %d", defaultGroup.GroupId)

	initAILinkData(tx, enterprise.Eid)
	logger.SysLog("Successfully created default ai_link group")

	agentGroup := Group{
		Eid:       enterprise.Eid,
		CreatedBy: 0, // System created
		GroupName: "默认",
		GroupType: AGENT_TYPE,
		Sort:      0,
	}

	if err := tx.Create(&agentGroup).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create default agent group: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created default agent group, ID: %d", agentGroup.GroupId)

	// 创建订阅设置
	subscriptionSetting := &SubscriptionSetting{
		GroupId:   defaultGroup.GroupId,
		LogoUrl:   "//img.ibos.cn/common/agenthub/subscription/vip-1.png", // 默认为空
		AiEnabled: false,
		IsDefault: true,
	}

	if err := tx.Create(subscriptionSetting).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create subscription setting: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created subscription setting, ID: %d", subscriptionSetting.SettingId)

	yearRelation := &SubscriptionRelation{
		SettingId: subscriptionSetting.SettingId,
		Amount:    0,
		Currency:  "CNY",
		TimeUnit:  "year",
		Type:      1,
	}

	if err := tx.Create(yearRelation).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create yearly subscription relation: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created yearly subscription relation")

	monthRelation := &SubscriptionRelation{
		SettingId: subscriptionSetting.SettingId,
		Amount:    0,
		Currency:  "CNY",
		TimeUnit:  "month",
		Type:      1,
	}

	if err := tx.Create(monthRelation).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create monthly subscription relation: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created monthly subscription relation")

	pointsRelation := &SubscriptionRelation{
		SettingId: subscriptionSetting.SettingId,
		Amount:    0,
		Currency:  "",
		TimeUnit:  "month",
		Type:      2,
	}

	if err := tx.Create(pointsRelation).Error; err != nil {
		tx.Rollback()
		logger.SysLogf("Failed to create points subscription relation: %s", err.Error())
		return err
	}
	logger.SysLogf("Successfully created points subscription relation")

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.SysLogf("Failed to commit transaction: %s", err.Error())
		return err
	}

	logger.SysLogf("System initialization completed successfully")

	// 执行版本检查（异步执行，不阻塞初始化流程）
	go func() {
		if resp, statisticScript, err := system.CheckVersionAndReturn(); err != nil {
			logger.SysLogf("Async version check failed: %v", err)
		} else {
			if resp.Action == "install" && statisticScript != "" {
				// 保存到数据库
				setting := &Setting{
					Eid:   enterprise.Eid,
					Key:   string(ThirdPartyStatisticHeader),
					Value: statisticScript,
				}

				if err := CreateSetting(setting); err != nil {
					logger.SysLogf("Failed to save statistic setting: %v", err)
				}
				logger.SysLogf("Successfully generated and saved statistic script")
			}

		}
	}()

	logger.SysLogf("\033[34m" + `
                    @                
                   ###                
           /###################\      
          |#####################|     
        ##|#####################|##  
       ###|####    #####    ####|###  
        ##|#####################|##  
          |#####################|     
           \###################/      
                  ######              
                  ###                 
                  @
	` + "\033[0m")
	logger.SysLogf("\033[32m\n" +
		"#################################\n" +
		"#  Email: admin@53ai.com        #\n" +
		"#  Password: admin888           #\n" +
		"#################################\n" +
		"\033[0m")
	return nil
}

var GroupData = []struct {
	GroupName string
	GroupType int64
	Sort      int64
}{
	{"AI搜索", 2, 6},
	{"智能对话", 2, 5},
	{"办公提效", 2, 4},
	{"图片处理", 2, 3},
	{"视频制作", 2, 2},
	{"AI学习", 2, 1},
}

var AILinkData = []struct {
	GroupName   string
	Name        string
	Logo        string
	URL         string
	Description string
	Sort        int64
}{
	// AI搜索分组
	{"AI搜索", "百度AI+", "https://hubapi.53ai.com/api/preview/b5970a3697479df6b00d73ab827dabb2.png", "https://chat.baidu.com", "百度官方ai搜索", 0},
	{"AI搜索", "天工AI", "https://hubapi.53ai.com/api/preview/432dfdbb2ade2e941a331fdc25ee29f5.png", "https://www.tiangong.cn/", "国内首个对标 ChatGPT 的双千亿级大语言模型，也是一个对话式AI助手", 0},
	{"AI搜索", "同花顺问财", "https://hubapi.53ai.com/api/preview/c65e9d65c42a1bdabfd7e09635dec05a.png", "https://www.iwencai.com", "同花顺旗下专业的智能选股平台", 0},
	{"AI搜索", "秘塔搜索", "https://hubapi.53ai.com/api/preview/710cd2a90fc7a38d8e78798af1fc597a.png", "https://metaso.cn", "没有广告，直达结果", 0},
	{"AI搜索", "Perplexity AI", "https://hubapi.53ai.com/api/preview/b2d85e0aa413297b2dccd0837fba6f28.png", "https://perplexity.ai", "知识的起点", 0},
	{"AI搜索", "知乎直答", "https://hubapi.53ai.com/api/preview/8698388b9dfc34d995a6238b120365d8.png", "https://zhida.zhihu.com", "用提问发现世界", 0},

	// 智能对话分组
	{"智能对话", "360智脑", "https://hubapi.53ai.com/api/preview/4a83fd5e7a31d0dd816d4f57237f13c5.png", "https://i.360.com/", "360搜索最新推出的AI对话聊天大模型", 0},
	{"智能对话", "百度AI伙伴", "https://hubapi.53ai.com/api/preview/afbc2525ffca738ba39989d486e97223.png", "https://chat.baidu.com/", "百度最新上线的AI搜索对话工具", 0},
	{"智能对话", "智谱清言", "https://hubapi.53ai.com/api/preview/b0072ad41d46626043cf1b2e3b2ce374.png", "https://chatglm.cn/", "Chatglm,千亿参数对话模型,支持多轮对话", 0},
	{"智能对话", "豆包", "https://hubapi.53ai.com/api/preview/d98b75d99fba38975312841a3c85aa72.png", "https://www.doubao.com/", "抖音旗下AI工具，你的智能助手", 0},
	{"智能对话", "ChatGPT", "https://hubapi.53ai.com/api/preview/bcade7d1cebca9273da445ffc8671711.png", "https://chat.openai.com", "Chatgpt.com", 0},
	{"智能对话", "通义千问", "https://hubapi.53ai.com/api/preview/ea1ad076efc73a30c8eaf1e86fc193cc.png", "https://tongyi.aliyun.com", "阿里巴巴旗下的一款智能体机器人，它利用自然语言处理技术，为用户提供智能化的语音交互服务", 0},
	{"智能对话", "零一万知", "https://hubapi.53ai.com/api/preview/f03bced2dfe845dec2d897cffcb3ce1b.png", "https://www.wanzhi.com/", "集AI对话聊天、文档阅读和PPT创作于一体的一站式AI工作平台", 0},
	{"智能对话", "讯飞星火", "https://hubapi.53ai.com/api/preview/4417ab5f7607452ccd8a3174616d7f56.png", "https://xinghuo.xfyun.cn", "懂你的AI助手", 0},
	{"智能对话", "文心一言", "https://hubapi.53ai.com/api/preview/eee853619f4fcbd7f15622198101630c.png", "https://yiyan.baidu.com/", "文心一言是百度研发的知识增强大语言模型，能够与人对话互动，回答问题，协助创作", 0},
	{"智能对话", "腾讯元宝", "https://hubapi.53ai.com/api/preview/433b8834406d66420558b6f093f0fed1.png", "https://yuanbao.tencent.com", "腾讯元宝是一款基于腾讯混元大模型的AI产品，为用户提供多元化的AI能力", 0},

	// 办公提效分组
	{"办公提效", "秒出PPT", "https://hubapi.53ai.com/api/preview/e3d748b2fc4a7f108090552e0b0dfc18.png", "https://10sppt.com/", "10S快速生成PPT", 0},
	{"办公提效", "AIPPT", "https://hubapi.53ai.com/api/preview/872850cdbb1fec8bc54581982572d4aa.png", "https://www.aippt.cn/", "AI一键生成PPT", 0},
	{"办公提效", "笔尖写作", "https://hubapi.53ai.com/api/preview/23addd994bc064fd2d48d8b0adbad6bd.png", "https://www.bijianxiezuo.com/", "高质量Ai写作利器", 0},
	{"办公提效", "ChatPPT", "https://hubapi.53ai.com/api/preview/f95c1d7469c53aff1c2d896677ce504b.png", "https://chat-ppt.com/", "对话式创作演示文稿，1400+类指令支持", 0},
	{"办公提效", "百度橙篇", "https://hubapi.53ai.com/api/preview/d6121e6ed2e190ad67ef05fc2897fc84.png", "https://cp.baidu.com", "写长文神器", 0},
	{"办公提效", "歌者PPT", "https://hubapi.53ai.com/api/preview/5379914644c44119e771865e00a1a565.png", "https://gezhe.com/", "永久免费的 PPT 智能生成工具", 0},
	{"办公提效", "万彩AI", "https://hubapi.53ai.com/api/preview/dd71cc93f0324a7d985e51ea931f8396.png", "https://ai.kezhan365.com/", "万彩AI，让创意轻松落地", 0},
	{"办公提效", "标智客", "https://hubapi.53ai.com/api/preview/b55f51c783bd473d3f9a1b3d1ebcf147.png", "https://www.logomaker.com.cn/", "智能LOGO设计生成", 0},
	{"办公提效", "Wegic", "https://hubapi.53ai.com/api/preview/de7b13cb2c80afc5cc010b2a6615d69a.png", "https://wegic.ai/", "即时设计团队推出的 AI 网页生成工具", 0},
	{"办公提效", "有道写作", "https://hubapi.53ai.com/api/preview/a5a84f10db33b8d49f5c242ba52b3a47.png", "https://write.youdao.com", "网易有道出品的智能英文写作修改和润色工具", 0},

	// 图片处理分组
	{"图片处理", "美图抠图", "https://hubapi.53ai.com/api/preview/b9f3a740af7c87e18af40d8ed8e50a8c.png", "https://cutout.designkit.com/", "美图秀秀推出的AI智能抠图工具，一键移除背景", 0},
	{"图片处理", "美图设计室", "https://hubapi.53ai.com/api/preview/179d3e00e9ebddf92da34330ad6e2097.png", "https://www.designkit.com/", "一款功能强大、易于使用的图像处理和照片编辑软件，提供了丰富功能", 0},
	{"图片处理", "一键抠图", "https://hubapi.53ai.com/api/preview/17a1c172c241b5ac79a5e9eb9ab58561.png", "https://www.yijiankoutu.com/", "在线一键抠图换背景", 0},

	// 视频制作分组
	{"视频制作", "百度度加", "https://hubapi.53ai.com/api/preview/919eb97d2b02114f475046c68fe3e70b.png", "https://aigc.baidu.com/", "度加剪辑是百度官方出品的口播自媒体必备剪辑工具，简洁好用", 0},
	{"视频制作", "鬼手剪辑", "https://hubapi.53ai.com/api/preview/8244602231503e7ebb429ffd393ccab7.png", "https://cn.jollytoday.com", "视频AI翻译、硬字幕翻译和视频去字幕的专业视频剪辑工具", 0},
	{"视频制作", "快手可灵", "https://hubapi.53ai.com/api/preview/52193bfb4d03a28fccfd827bbb450e04.png", "https://app.klingai.com/", "快手旗下图片生成和视频生成大模型工具", 0},
	{"视频制作", "抖音即创", "https://hubapi.53ai.com/api/preview/707620ea6742afae608c9c109b51a33d.png", "https://aic.oceanengine.com", "专注于智能创意生产与管理分析", 0},
	{"视频制作", "pika", "https://hubapi.53ai.com/api/preview/7f2201d2f291251a64c626dffa5d9d2d.png", "https://pika.art", "文本生成电影工具", 0},
	{"视频制作", "腾讯智影", "https://hubapi.53ai.com/api/preview/2bf37ee21997235dea73084e21795987.png", "https://zenvideo.qq.com", "腾讯智影AI绘画，只需简单的描述就可为您生成独一无二的创意画作", 0},

	// AI学习分组
	{"AI学习", "LangGPT", "https://hubapi.53ai.com/api/preview/41193ab845ca040c9ea34b0a7fa1bb80.png", "https://langgptai.feishu.cn/", "人人都能写出高质量提示词", 0},
	{"AI学习", "通往AGI之路", "https://hubapi.53ai.com/api/preview/30f58543638b8e9eab9a242f2c1594ed.png", "https://waytoagi.feishu.cn/", "一个全面系统的AI学习路径", 0},
}

func initAILinkData(tx *gorm.DB, eid int64) error {
	groups := make([]Group, 0, len(GroupData))
	currentTime := time.Now().UnixMilli()

	for _, g := range GroupData {
		groups = append(groups, Group{
			Eid:       eid,
			GroupName: g.GroupName,
			GroupType: g.GroupType,
			Sort:      g.Sort,
			BaseModel: BaseModel{
				CreatedTime: currentTime,
				UpdatedTime: currentTime,
			},
		})
	}

	if err := tx.CreateInBatches(groups, 100).Error; err != nil {
		tx.Rollback()
		return err
	}

	groupMap := make(map[string]int64)
	var insertedGroups []Group
	if err := tx.Where("eid = ?", eid).Find(&insertedGroups).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, g := range insertedGroups {
		groupMap[g.GroupName] = g.GroupId
	}

	aiLinks := make([]AILink, 0, len(AILinkData))
	for _, link := range AILinkData {
		groupID, exists := groupMap[link.GroupName]
		if !exists {
			continue
		}

		aiLinks = append(aiLinks, AILink{
			Eid:         eid,
			GroupID:     groupID,
			Name:        link.Name,
			Logo:        link.Logo,
			URL:         link.URL,
			Description: link.Description,
			Sort:        link.Sort,
			BaseModel: BaseModel{
				CreatedTime: currentTime,
				UpdatedTime: currentTime,
			},
		})
	}

	if err := tx.CreateInBatches(aiLinks, 200).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func GetEnterpriseByWecomCorpID(wecomCorpID string) (*Enterprise, error) {
	var enterprise Enterprise
	err := DB.Where("wecom_corp_id = ?", wecomCorpID).First(&enterprise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &enterprise, nil
}

func (e *Enterprise) LoadWecomCorpInfo(suiteID string, loadType int) error {
	e.WecomInstallInfo = &WecomInstallInfo{
		InstallWecomApp: 0,
	}
	if e.WecomCorpID == "" || suiteID == "" {
		return nil
	}

	wc, err := GetWecomCorp(suiteID, e.WecomCorpID)
	if wc == nil || err != nil {
		return nil
	}

	e.WecomInstallInfo.InstallWecomApp = 1
	if loadType == 1 {
		e.WecomInstallInfo.AuthCorpInfo = wc.GetAuthCorpInfo()
	}

	return nil
}
