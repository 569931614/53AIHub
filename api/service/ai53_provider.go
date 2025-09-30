package service

import (
	"github.com/53AI/53AIHub/common/utils/ai53"
	"github.com/53AI/53AIHub/model"
)

type AI53Service struct {
	Provider *model.Provider
}

func (a *AI53Service) GetApiSdk() *ai53.AI53Api {
	baseUrl := a.Provider.GetBaseURLByProviderType()
	return &ai53.AI53Api{
		BaseUrl:   baseUrl,
		AuthToken: a.Provider.AccessToken,
	}
}

func (a *AI53Service) GetAllBots() ([]ai53.AppResponse, error) {
	api := a.GetApiSdk()
	var allBots []ai53.AppResponse
	offset := 0
	limit := 100

	for {
		resp, err := api.GetApps(offset, limit)
		if err != nil {
			return nil, err
		}

		allBots = append(allBots, resp...)

		if len(resp) < limit {
			break
		}
		offset += limit
	}

	return allBots, nil
}

func (a *AI53Service) GetAllWorkflows() ([]ai53.AppResponse, error) {
	api := a.GetApiSdk()
	allWorkflows := make([]ai53.AppResponse, 0)
	offset := 0
	limit := 100

	for {
		resp, err := api.GetWorkflows(offset, limit)
		if err != nil {
			return nil, err
		}

		allWorkflows = append(allWorkflows, resp...)

		if len(resp) < limit {
			break
		}
		offset += limit
	}

	return allWorkflows, nil
}

func (a *AI53Service) GetAppParameters(botId string) (interface{}, error) {
	// 获取 API SDK 实例 (复用现有模式)
	api := a.GetApiSdk()

	// 调用 API 方法
	return api.GetAppParameters(botId)
}
