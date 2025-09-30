package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/53AI/53AIHub/model"
	adaptor53AI "github.com/53AI/53AIHub/service/hub_adaptor/53AI"
	"github.com/53AI/53AIHub/service/hub_adaptor/appbuilder"
	"github.com/53AI/53AIHub/service/hub_adaptor/bailian"
	"github.com/53AI/53AIHub/service/hub_adaptor/coze"
	"github.com/53AI/53AIHub/service/hub_adaptor/custom"
	"github.com/53AI/53AIHub/service/hub_adaptor/dify"
	"github.com/53AI/53AIHub/service/hub_adaptor/n8n"
	Hub_openai "github.com/53AI/53AIHub/service/hub_adaptor/openai"
	Hub_tencent "github.com/53AI/53AIHub/service/hub_adaptor/tencent"
	"github.com/53AI/53AIHub/service/hub_adaptor/yuanqi"
	"github.com/songquanpeng/one-api/relay/adaptor"
	"github.com/songquanpeng/one-api/relay/adaptor/aiproxy"
	"github.com/songquanpeng/one-api/relay/adaptor/ali"
	"github.com/songquanpeng/one-api/relay/adaptor/anthropic"
	"github.com/songquanpeng/one-api/relay/adaptor/aws"
	"github.com/songquanpeng/one-api/relay/adaptor/baidu"
	"github.com/songquanpeng/one-api/relay/adaptor/cloudflare"
	"github.com/songquanpeng/one-api/relay/adaptor/cohere"
	"github.com/songquanpeng/one-api/relay/adaptor/deepl"
	"github.com/songquanpeng/one-api/relay/adaptor/gemini"
	"github.com/songquanpeng/one-api/relay/adaptor/ollama"
	"github.com/songquanpeng/one-api/relay/adaptor/palm"
	"github.com/songquanpeng/one-api/relay/adaptor/proxy"
	"github.com/songquanpeng/one-api/relay/adaptor/replicate"
	"github.com/songquanpeng/one-api/relay/adaptor/tencent"
	"github.com/songquanpeng/one-api/relay/adaptor/vertexai"
	"github.com/songquanpeng/one-api/relay/adaptor/xunfei"
	"github.com/songquanpeng/one-api/relay/adaptor/zhipu"
	"github.com/songquanpeng/one-api/relay/apitype"
)

func GetAdaptor(apiType int) adaptor.Adaptor {
	switch apiType {
	case apitype.AIProxyLibrary:
		return &aiproxy.Adaptor{}
	case apitype.Ali:
		return &ali.Adaptor{}
	case apitype.Anthropic:
		return &anthropic.Adaptor{}
	case apitype.AwsClaude:
		return &aws.Adaptor{}
	case apitype.Baidu:
		return &baidu.Adaptor{}
	case apitype.Gemini:
		return &gemini.Adaptor{}
	case apitype.OpenAI:
		return &Hub_openai.Adaptor{}
	case apitype.PaLM:
		return &palm.Adaptor{}
	case apitype.Tencent:
		return &tencent.Adaptor{}
	case apitype.Xunfei:
		return &xunfei.Adaptor{}
	case apitype.Zhipu:
		return &zhipu.Adaptor{}
	case apitype.Ollama:
		return &ollama.Adaptor{}
	case apitype.Coze:
		return &coze.Adaptor{}
	case apitype.Cohere:
		return &cohere.Adaptor{}
	case apitype.Cloudflare:
		return &cloudflare.Adaptor{}
	case apitype.DeepL:
		return &deepl.Adaptor{}
	case apitype.VertexAI:
		return &vertexai.Adaptor{}
	case apitype.Proxy:
		return &proxy.Adaptor{}
	case apitype.Replicate:
		return &replicate.Adaptor{}
	case model.ChannelApiDify:
		return &dify.Adaptor{}
	case model.ChannelApi53AI:
		return &adaptor53AI.Adaptor{}
	case model.ChannelApiVolcengine:
		return &Hub_openai.Adaptor{}
	case model.ChannelApiAppBuilder:
		return &appbuilder.Adaptor{}
	case model.ChannelApiBailian:
		return &bailian.Adaptor{}
	case model.ChannelApiYuanqi:
		return &yuanqi.Adaptor{}
	case model.ChannelApiTypeFastGpt:
		return &Hub_openai.Adaptor{}
	case model.ChannelApiTypeMaxKB:
		return &Hub_openai.Adaptor{}
	case model.ChannelApiTypeN8n:
		return &n8n.Adaptor{}
	case model.ChannelApiTypeCozeStudio:
		return &coze.Adaptor{}
	case model.ChannelApiTypeTencent:
		return &Hub_tencent.Adaptor{}
	}

	return nil
}

func SetCustomConfig(a *adaptor.Adaptor, customConfig *custom.CustomConfig) error {
	switch v := (*a).(type) {
	case *dify.Adaptor:
		v.CustomConfig = customConfig
	case *adaptor53AI.Adaptor:
		v.CustomConfig = customConfig
	case *coze.Adaptor:
		v.CustomConfig = customConfig
	case *Hub_openai.Adaptor:
		v.CustomConfig = customConfig
	case *appbuilder.Adaptor:
		// Appbuilder 7 day Expired
		if customConfig.ConversationId == "" || customConfig.ConversationExpirationTime <= time.Now().Unix() {
			appBuilderConversationId, err := v.GetConversationId()
			if err != nil {
				return err
			}
			customConfig.ConversationId = appBuilderConversationId
			if customConfig.ConversationId == "" {
				return errors.New("appbuilder create conversation error: conversation id is empty")
			}
			customConfig.ConversationExpirationTime = time.Now().Unix() + 60*60*24*7
		}
		v.CustomConfig = customConfig
	case *bailian.Adaptor:
		// The SessionId in Bailian is customized for multi-turn conversations.
		v.CustomConfig = customConfig
		v.CustomConfig.ConversationId = fmt.Sprintf("53AIHub_%d", customConfig.AIHubConversationId)
	case *yuanqi.Adaptor:
		v.CustomConfig = customConfig
	case *n8n.Adaptor:
		v.CustomConfig = customConfig
	case *Hub_tencent.Adaptor:
		v.CustomConfig = customConfig
	}
	return nil
}

func GetCustomConfig(a *adaptor.Adaptor) *custom.CustomConfig {
	switch v := (*a).(type) {
	case *dify.Adaptor:
		return v.CustomConfig
	case *adaptor53AI.Adaptor:
		return v.CustomConfig
	case *coze.Adaptor:
		return v.CustomConfig
	case *Hub_openai.Adaptor:
		return v.CustomConfig
	case *appbuilder.Adaptor:
		return v.CustomConfig
	case *bailian.Adaptor:
		return v.CustomConfig
	case *yuanqi.Adaptor:
		return v.CustomConfig
	case *n8n.Adaptor:
		return v.CustomConfig
	case *Hub_tencent.Adaptor:
		return v.CustomConfig
	}
	return nil
}
