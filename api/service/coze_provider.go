package service

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/53AI/53AIHub/common/logger"
	"github.com/53AI/53AIHub/common/storage"
	"github.com/53AI/53AIHub/common/utils/coze"
	"github.com/53AI/53AIHub/model"
	db_model "github.com/53AI/53AIHub/model"
	"github.com/songquanpeng/one-api/relay/channeltype"
)

type CozeService struct {
	Provider db_model.Provider
}

func (ser *CozeService) GetCozeApiSdk() (*coze.CozeApi, error) {
	var baseUrl string
	if ser.Provider.ProviderType == model.ProviderTypeCozeCn {
		baseUrl = coze.CozeCnUrl
	} else if ser.Provider.ProviderType == model.ProviderTypeCozeCom {
		baseUrl = coze.CozeComUrl
	} else if ser.Provider.ProviderType == model.ProviderTypeCozeStudio {
		// coze-studio uses custom base_url from Provider.BaseURL
		if ser.Provider.BaseURL != nil && *ser.Provider.BaseURL != "" {
			baseUrl = *ser.Provider.BaseURL
		} else {
			return nil, errors.New("coze-studio requires custom base_url")
		}
	} else {
		baseUrl = coze.CozeComUrl
	}

	return &coze.CozeApi{
		BaseUrl: baseUrl,
	}, nil
}

func (ser *CozeService) HandlerAccessTokenByCode(coze string, callbackUrl string) error {
	if ser.Provider.ProviderType != model.ProviderTypeCozeCn && ser.Provider.ProviderType != model.ProviderTypeCozeCom {
		return errors.New("invalid provider type")
	}

	// coze-studio uses fixed AccessToken, skip OAuth flow
	if ser.Provider.ProviderType == model.ProviderTypeCozeStudio {
		return errors.New("coze-studio does not support OAuth flow")
	}

	api, err := ser.GetCozeApiSdk()
	if err != nil {
		return err
	}

	var config model.CozeConfig

	err = json.Unmarshal([]byte(ser.Provider.Configs), &config)
	if err != nil {
		return err
	}

	cozeApiToken, err := api.GetOAuthToken(config.ClientID, config.ClientSecret, coze, callbackUrl)
	if err != nil {
		return err
	}

	ser.Provider.AccessToken = cozeApiToken.AccessToken
	ser.Provider.RefreshToken = cozeApiToken.RefreshToken
	ser.Provider.ExpiresIn = cozeApiToken.ExpiresIn
	ser.Provider.IsAuthorized = true
	ser.Provider.AuthedTime = time.Now().UTC().UnixMilli()
	err = model.UpdateProvider(&ser.Provider)

	return err
}

func (ser *CozeService) HandlerAccessTokenByRefreshToken() error {
	if ser.Provider.ProviderType != model.ProviderTypeCozeCn && ser.Provider.ProviderType != model.ProviderTypeCozeCom {
		return errors.New("invalid provider type")
	}

	// coze-studio uses fixed AccessToken, skip refresh flow
	if ser.Provider.ProviderType == model.ProviderTypeCozeStudio {
		return errors.New("coze-studio does not support token refresh")
	}
	api, err := ser.GetCozeApiSdk()
	if err != nil {
		return err
	}

	var config model.CozeConfig
	err = json.Unmarshal([]byte(ser.Provider.Configs), &config)
	if err != nil {
		return err
	}
	cozeApiToken, err := api.RefreshOAuthToken(config.ClientID, config.ClientSecret, ser.Provider.RefreshToken)
	if err != nil {
		return err
	}
	ser.Provider.AccessToken = cozeApiToken.AccessToken
	ser.Provider.RefreshToken = cozeApiToken.RefreshToken
	ser.Provider.ExpiresIn = cozeApiToken.ExpiresIn
	ser.Provider.IsAuthorized = true
	err = model.UpdateProvider(&ser.Provider)
	if err != nil {
		return err
	}

	// update channel key
	existingChannel, err := model.GetFirstChannelByEidAndProviderId(ser.Provider.Eid, ser.Provider.ProviderID)
	if err != nil {
		return err
	}

	existingChannel.Key = ser.Provider.AccessToken
	err = model.UpdateChannel(existingChannel)
	return err
}

func (ser *CozeService) CheckAndRefreshToken() (ok bool, err error) {
	// coze-studio uses fixed AccessToken, no need to refresh
	if ser.Provider.ProviderType == model.ProviderTypeCozeStudio {
		return false, nil
	}

	if ser.Provider.ExpiresIn <= time.Now().Unix() {
		logger.SysLogf("Coze RefreshToken: eid = %d", ser.Provider.Eid)
		err := ser.HandlerAccessTokenByRefreshToken()
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (ser *CozeService) GetAllWorkspace() ([]*coze.Workspace, error) {
	logger.SysLogf("CozeService.GetAllWorkspace: Using provider - ID: %d, Name: %s, Type: %d, AccessToken: %s",
		ser.Provider.ProviderID, ser.Provider.Name, ser.Provider.ProviderType,
		ser.Provider.AccessToken[:min(len(ser.Provider.AccessToken), 20)]+"...")

	_, err := ser.CheckAndRefreshToken()
	if err != nil {
		return nil, err
	}

	api, err := ser.GetCozeApiSdk()
	if err != nil {
		return nil, err
	}

	var config model.CozeConfig
	err = json.Unmarshal([]byte(ser.Provider.Configs), &config)
	if err != nil {
		return nil, err
	}

	var allWorkspaces []*coze.Workspace
	page := 1
	pageSize := 50

	for {
		if page > 20 {
			break
		}
		logger.SysLogf("CozeService.GetAllWorkspace: Fetching page %d for provider %d", page, ser.Provider.ProviderID)
		workspacesResp, err := api.GetWorkspaces(&ser.Provider, page, pageSize)
		if err != nil {
			logger.SysLogf("CozeService.GetAllWorkspace: Error fetching workspaces for provider %d on page %d: %v", ser.Provider.ProviderID, page, err)
			return nil, err
		}
		logger.SysLogf("CozeService.GetAllWorkspace: Got %d workspaces on page %d for provider %d", len(workspacesResp.Workspaces), page, ser.Provider.ProviderID)
		if len(workspacesResp.Workspaces) == 0 {
			break
		}
		for _, workspace := range workspacesResp.Workspaces {
			allWorkspaces = append(allWorkspaces, &workspace)
		}
		page++
	}

	logger.SysLogf("CozeService.GetAllWorkspace: Total workspaces found for provider %d: %d", ser.Provider.ProviderID, len(allWorkspaces))
	return allWorkspaces, nil
}

func (ser *CozeService) GetAllBot(workspaceId string) ([]*coze.Bot, error) {
	// 使用API工具类中的认证检查方式
	api, err := ser.GetCozeApiSdk()
	if err != nil {
		logger.SysErrorf("GetCozeApiSdk failed: %v", err)
		return nil, err
	}

	// 通过API工具类检查和刷新token
	if err := api.RefreshTokenIfNeeded(&ser.Provider); err != nil {
		logger.SysErrorf("RefreshTokenIfNeeded failed: %v", err)
		return nil, err
	}

	var config db_model.CozeConfig
	err = json.Unmarshal([]byte(ser.Provider.Configs), &config)
	if err != nil {
		logger.SysErrorf("Failed to unmarshal config: %v", err)
		return nil, err
	}

	logger.SysLogf("Getting bots for workspace: %s, page_size: %d", workspaceId, 50)

	var allBots []*coze.Bot
	page := 1
	pageSize := 50

	for {
		if page > 20 {
			logger.SysLogf("Reached max page limit: %d", page)
			break
		}

		logger.SysLogf("Fetching page %d for workspace: %s", page, workspaceId)
		// 调整参数顺序，确保与API定义一致
		botsResp, err := api.GetPublishedBots(&ser.Provider, workspaceId, page, pageSize)
		if err != nil {
			logger.SysErrorf("GetPublishedBots failed for workspace %s, page %d: %v", workspaceId, page, err)
			return nil, err
		}

		logger.SysLogf("Fetched %d bots for workspace %s on page %d", len(botsResp.SpaceBots), workspaceId, page)

		if len(botsResp.SpaceBots) == 0 {
			logger.SysLogf("No more bots found for workspace %s", workspaceId)
			break
		}

		for _, bot := range botsResp.SpaceBots {
			logger.SysLogf("Name: %s, icon: %s", bot.BotName, bot.IconURL)
			allBots = append(allBots, &bot)
		}

		// 如果返回的记录数小于分页大小，说明已经到达最后一页
		if len(botsResp.SpaceBots) < pageSize {
			logger.SysLogf("Reached last page for workspace %s", workspaceId)
			break
		}
		page++
	}

	logger.SysLogf("Total bots fetched for workspace %s: %d", workspaceId, len(allBots))
	return allBots, nil
}

// CacheBotIconWithUploadFile 使用UploadFile表缓存bot图标
func (ser *CozeService) CacheBotIconWithUploadFile(botID string, iconURL string, eid int64) (string, error) {
	// 清理URL，只保留?之前的部分（去除查询参数）
	cleanedIconURL := iconURL
	if idx := strings.Index(iconURL, "?"); idx != -1 {
		cleanedIconURL = iconURL[:idx]
	}

	// 计算清理后URL的哈希值
	urlHash := fmt.Sprintf("%x", md5.Sum([]byte(cleanedIconURL)))

	// 检查UploadFile表中是否已存在该URL的记录
	existingFile := ser.findUploadFileByHash(urlHash)
	if existingFile != nil {
		// 已存在，直接返回预览URL
		return existingFile.GetPreviewFullUrl(), nil
	}

	// 不存在，需要下载并保存图标
	var iconData []byte
	var err error
	for i := 0; i < 3; i++ {
		iconData, err = ser.downloadIcon(iconURL)
		if err == nil {
			break
		}
		logger.SysLogf("下载图标失败，第 %d 次重试: %v", i+1, err)
		time.Sleep(200 * time.Millisecond)
	}
	if err != nil {
		return "", fmt.Errorf("下载图标失败: %w", err)
	}

	// 获取文件扩展名
	ext := path.Ext(cleanedIconURL)
	if ext == "" {
		// 尝试从Content-Type获取扩展名
		// 这里简化处理，直接默认使用png
		ext = ".png"
	}

	// 生成previewKey
	previewKey, err := db_model.GetPreviewKey(urlHash, ext)
	if err != nil {
		return "", fmt.Errorf("生成预览键失败: %w", err)
	}

	// 保存新文件到存储系统
	fileKey := db_model.GetFileKey(fmt.Sprintf("coze_bot_%s_icon%s", botID, ext), 0, 0)
	if err := storage.StorageInstance.Save(iconData, fileKey); err != nil {
		return "", fmt.Errorf("保存图标文件失败: %w", err)
	}

	// 创建UploadFile记录
	uploadFile := &db_model.UploadFile{
		FileName:   fmt.Sprintf("coze_bot_%s_icon%s", botID, ext),
		Key:        fileKey,
		Eid:        eid,
		UserID:     0, // 系统文件
		Size:       int64(len(iconData)),
		Extension:  ext,
		MimeType:   http.DetectContentType(iconData),
		Hash:       urlHash, // 使用URL哈希
		PreviewKey: previewKey,
	}

	if err := uploadFile.Save(); err != nil {
		// 如果保存记录失败，删除已上传的文件
		storage.StorageInstance.Delete(fileKey)
		return "", fmt.Errorf("保存上传文件记录失败: %w", err)
	}

	return uploadFile.GetPreviewFullUrl(), nil
}

// findUploadFileByHash 根据哈希查找UploadFile记录
func (ser *CozeService) findUploadFileByHash(hash string) *db_model.UploadFile {
	var uploadFile db_model.UploadFile
	if err := db_model.DB.Where("hash = ?", hash).First(&uploadFile).Error; err != nil {
		if err.Error() == "record not found" {
			return nil
		}
		logger.SysErrorf("查询UploadFile失败: %v", err)
		return nil
	}
	return &uploadFile
}

// downloadIcon 下载图标文件
func (ser *CozeService) downloadIcon(iconURL string) ([]byte, error) {
	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", iconURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置User-Agent避免被拦截
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	if ser.Provider.AccessToken != "" {
		logger.SysLogf("添加认证头到图标请求")
		req.Header.Set("Authorization", "Bearer "+ser.Provider.AccessToken)
	} else {
		logger.SysLogf("警告：没有AccessToken可用于图标请求")
	}

	logger.SysLogf("下载Coze图标，URL: %s", iconURL)

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载图标失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载图标失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}

	return data, nil
}

// UpdateCozeChannel asynchronously updates the Coze channel's model list with new bot IDs.
// It ensures no duplicate bot IDs are added to the channel's models.
// The update process runs in a separate goroutine to avoid blocking the main execution.
// Parameters:
//   - botIds: slice of bot IDs to be added to the channel's model list
//
// Returns:
//   - error: returns nil as the update process runs asynchronously
func (ser *CozeService) UpdateCozeChannel(botIds []string, provider *model.Provider) error {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.SysErrorf("Panic in UpdateCozeChannel: %v", r)
			}
		}()

		// Get the existing channel for the current enterprise ID with type Coze (34)
		existingChannel, err := model.GetFirstChannelByEidAndProviderType(ser.Provider.Eid, int64(channeltype.Coze), provider.ProviderID)
		if err != nil {
			logger.SysErrorf("Failed to get Channel: %v", err)
			return
		}

		// Parse existing models and store them in a map for deduplication
		existingModels := make(map[string]bool)
		if existingChannel.Models != "" {
			for _, model := range strings.Split(existingChannel.Models, ",") {
				model = strings.TrimSpace(model)
				if model != "" {
					existingModels[model] = true
				}
			}
		}

		// Process new bot IDs and add them if they don't exist
		var updatedBotIds []string
		for _, botId := range botIds {
			if botId == "" {
				continue
			}
			formattedBotId := "bot-" + botId
			if !existingModels[formattedBotId] {
				updatedBotIds = append(updatedBotIds, formattedBotId)
				existingModels[formattedBotId] = true // Mark as added to prevent duplicates
			}
		}

		// If no new bot IDs to add, return early
		if len(updatedBotIds) == 0 {
			logger.SysLogf("No new bot IDs to add")
			return
		}

		// Rebuild the complete bot IDs list from the map to ensure uniqueness
		var allBotIds []string
		for model := range existingModels {
			allBotIds = append(allBotIds, model)
		}

		// Prepare channel update with default configuration
		configStr := `{"region":"","sk":"","ak":"","user_id":"53AIHub","vertex_ai_project_id":"","vertex_ai_adc":""}`
		baseURL := ser.Provider.GetBaseURLByProviderType()

		// Create channel object with updated information
		channel := &model.Channel{
			ChannelID:  existingChannel.ChannelID,
			Eid:        ser.Provider.Eid,
			Name:       ser.Provider.Name,
			Key:        ser.Provider.AccessToken,
			Type:       channeltype.Coze,
			ProviderID: ser.Provider.ProviderID,
			BaseURL:    &baseURL,
			Models:     strings.Join(allBotIds, ","),
			Status:     model.ChannelStatusEnabled,
			Config:     configStr,
		}

		// Update the channel in the database
		if err := model.UpdateChannel(channel); err != nil {
			logger.SysErrorf("Failed to update Channel: %v", err)
			return
		}

		logger.SysLogf("Successfully updated Channel %d, added %d new bot IDs", channel.ChannelID, len(updatedBotIds))
	}()

	return nil
}
