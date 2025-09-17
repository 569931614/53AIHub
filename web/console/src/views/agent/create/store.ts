import { defineStore } from 'pinia'
import { AGENT_TYPE, agentApi } from '@/api/modules/agent'
import { CHANNEL_TYPE_VALUE_MAP, channelApi } from '@/api/modules/channel'
import { groupApi } from '@/api/modules/group'
import { useEnterpriseStore } from '@/stores'
import {
  AGENT_TYPES,
  getAgentByAgentType,
  getModelChannelTypes,
  BACKEND_AGENT_TYPE,
  AGENT_MODES,
} from '@/constants/platform/config'
import type { AgentType } from '@/constants/platform/config'
import { GROUP_TYPE } from '@/constants/group'

const enterprise_store = useEnterpriseStore()

const DEFAULT_COMPLETION_PARAMS = {
  temperature: 0.2,
  top_p: 0.75,
  presence_penalty: 0.5,
  frequency_penalty: 0.5,
}

export const useAgentFormStore = defineStore('agent-form-store', {
  state: (): {
    saving: boolean
    loading: boolean
    agent_id: number
    agent_type: string
    form_data: Agent.FormData
    agent_data: Record<string, any>
    group_options: any[]
    model_options: any[]
    coze_workspace_options: any[]
    coze_bot_options: any[]
    app_builder_bot_options: any[]
    chat53ai_app_options: any[]
    is_new: boolean
  } => ({
    saving: false,
    loading: false,
    agent_id: 0,
    agent_type: 'prompt',
    form_data: {
      logo: '',
      name: '',
      group_id: 0,
      description: '',
      channel_type: 0,
      model: '',
      sort: 0,
      prompt: '',
      user_group_ids: [],
      subscription_group_ids: [],
      tools: [],
      use_cases: [],
      configs: {},
      custom_config: {
        agent_type: 'prompt',
        agent_mode: 'chat',
        coze_workspace_id: '',
        coze_bot_id: '',
        coze_bot_url: '',
        app_builder_bot_id: '',
        chat53ai_agent_id: '',
        channel_config: {},
      },
      settings: {
        opening_statement: '',
        suggested_questions: [],
        file_parse: {
          enable: false,
        },
        image_parse: {
          enable: false,
        },
        relate_agents: [],
        input_fields: [],
        output_fields: [],
      },
    },
    agent_data: {},
    group_options: [],
    model_options: [],
    coze_workspace_options: [],
    coze_bot_options: [],
    app_builder_bot_options: [],
    chat53ai_app_options: [],
    is_new: false,
  }),
  getters: {
    agent_option_data: state => {
      return getAgentByAgentType(state.agent_type as AgentType)
    },
    support_file: state => {
      return state.agent_type !== AGENT_TYPE.PROMPT
    },
    support_image: state => {
      if (state.agent_type !== AGENT_TYPE.PROMPT) return true

      return (
        state.model_options
          .find((item: { model_options: Array<{ value: string }> }) =>
            item.model_options.some((row: { value: string }) => row.value === state.form_data.model)
          )
          ?.model_options.find(
            (item: { value: string; vision?: boolean }) => item.value === state.form_data.model
          )?.vision || false
      )
    },
    is_independent: () => enterprise_store.info.is_independent,
  },

  actions: {
    async loadDetailData() {
      if (!this.agent_id) return Promise.resolve()

      this.loading = true
      try {
        this.agent_data = await agentApi.detail({ data: { agent_id: this.agent_id } })

        this.updateFormData()
      } finally {
        this.loading = false
      }
    },

    updateFormData() {
      this.agent_type = this.agent_data.agent_type || this.agent_type || ''
      Object.assign(this.form_data, {
        logo: this.agent_data.logo || this.agent_option_data?.icon || '',
        name: this.agent_data.name || '',
        group_id: +this.agent_data.group_id || 0,
        description: this.agent_data.description || '',
        channel_type: +this.agent_data.channel_type || 0,
        model: this.agent_data.model || '',
        sort: +this.agent_data.sort || 0,
        prompt: this.agent_data.prompt || '',
        user_group_ids: this.agent_data.user_group_ids || [],
        subscription_group_ids: this.agent_data.user_group_ids || [],
        tools: this.agent_data.tools || [],
        use_cases: this.agent_data.use_cases || [],
        configs:
          this.agent_data.configs && Object.keys(this.agent_data.configs).length > 0
            ? this.agent_data.configs
            : { completion_params: DEFAULT_COMPLETION_PARAMS },
        enable: !!+this.agent_data.enable || false,
        custom_config: {
          agent_type: this.agent_type,
          channel_id: 0,
          coze_workspace_id: '',
          coze_bot_id: '',
          coze_bot_url: '',
          app_builder_bot_id: '',
          chat53ai_agent_id: '',
          channel_config: {},
          ...(this.agent_data.custom_config || {}),
        },
        settings: {
          opening_statement: '',
          suggested_questions: [],
          // 文件解析
          file_parse: {
            enable: false,
          },
          // 图片解析
          image_parse: {
            enable: false,
          },
          // 关联场景
          relate_agents: [],
          input_fields: [],
          output_fields: [],
          ...(this.agent_data.settings || {}),
        },
      })
      const custom_config = this.agent_data.custom_config || {}
      if (this.agent_type === AGENT_TYPE.PROMPT) {
        this.form_data.model =
          `${custom_config.channel_id}_${this.agent_data.channel_type}_${this.agent_data.model}` ||
          ''
      }
      // 文件解析 迁移到 setting来
      if (custom_config.file_parse) {
        this.form_data.settings.file_parse = custom_config.file_parse
        this.form_data.settings.image_parse = custom_config.image_parse
        delete custom_config.file_parse
        delete custom_config.image_parse
      }
    },

    async loadGroupOptions() {
      const list = await groupApi.list({ params: { group_type: GROUP_TYPE.AGENT } })
      this.group_options = (list || []).map(item => {
        item.value = +item.group_id || 0
        item.label = item.group_name || ''
        return item
      })
      if (this.group_options.length && !this.form_data.group_id)
        this.form_data.group_id = this.group_options[0].value
      if (!this.group_options.find(item => item.value === this.form_data.group_id))
        this.form_data.group_id = ''
    },

    async loadChannelOptions() {
      const list = await channelApi.list()
      const model_List = (list || []).filter((item = {}) =>
        getModelChannelTypes().includes(item.channel_type)
      )

      const all_model_options = []
      this.model_options = model_List.reduce((acc, item) => {
        ;((item.value = item.channel_type),
          (item.icon = window.$getRealPath({ url: `/images/platform/${item.icon}.png` })))

        const options = (item.model_options || []).map((option = {}) => {
          option.value = `${item.channel_id}_${item.channel_type}_${option.value}`
          option.icon = window.$getRealPath({ url: `/images/platform/${option.icon}.png` })
          return option
        })
        const model = acc.find(res => res.value === item.channel_type)
        if (model) {
          model.options.push(...options)
        } else {
          item.options = options
          acc.push(item)
        }
        all_model_options.push(...options)

        return acc
      }, [])
      if (all_model_options.length && !this.form_data.model)
        this.form_data.model = all_model_options[0].value
      if (
        [AGENT_TYPE.PROMPT].includes(this.agent_type) &&
        !all_model_options.find(item => item.value === this.form_data.model)
      )
        this.form_data.model = ''
    },

    async loadCozeWorkspaceOptions() {
      const list = await channelApi.cozeWorkspaceList()
      this.coze_workspace_options = (list || []).map(item => {
        item.value = item.id || 0
        item.label = item.name || ''
        item.icon = item.icon_url || ''
        return item
      })
      if (this.coze_workspace_options.length && !this.form_data.custom_config.coze_workspace_id) {
        this.form_data.custom_config.coze_workspace_id = this.coze_workspace_options[0].value
      }
      if (
        !this.coze_workspace_options.find(
          item => item.value === this.form_data.custom_config.coze_workspace_id
        )
      )
        this.form_data.custom_config.coze_workspace_id = ''
    },

    async loadCozeBotOptions(workspace_id: string) {
      if (!workspace_id) return

      const list = await channelApi.cozeBotList({ params: { workspace_id } })
      this.coze_bot_options = (list || []).map(item => {
        item.value = item.bot_id || 0
        item.label = item.bot_name || ''
        item.icon = item.icon_url || ''
        return item
      })
      if (this.coze_bot_options.length && !this.form_data.custom_config.coze_bot_id) {
        const bot = this.coze_bot_options[0]
        if (!this.agent_id) {
          this.form_data.logo = bot.icon
          this.form_data.name = bot.label
          this.form_data.description = bot.description || ''
        }
        this.form_data.custom_config.coze_bot_id = bot.value
      }
      if (
        !this.coze_bot_options.find(item => item.value === this.form_data.custom_config.coze_bot_id)
      )
        this.form_data.custom_config.coze_bot_id = ''
    },

    async loadAppBuilderBotOptions() {
      const list = await channelApi.appBuilderBotList()
      this.app_builder_bot_options = (list || []).map(item => {
        item.value = item.id || 0
        if (item.value) item.value = `bot-${item.value}`
        item.label = item.name || ''
        item.icon = item.icon || getAgentByAgentType(AGENT_TYPES.APP_BUILDER).icon
        return item
      })
      if (this.app_builder_bot_options.length && !this.form_data.custom_config.app_builder_bot_id) {
        const bot = this.app_builder_bot_options[0]
        if (!this.agent_id) {
          this.form_data.name = bot.label
          this.form_data.description = bot.description || ''
        }
        this.form_data.custom_config.app_builder_bot_id = bot.value
      }
      if (
        !this.app_builder_bot_options.find(
          item => item.value === this.form_data.custom_config.app_builder_bot_id
        )
      )
        this.form_data.custom_config.app_builder_bot_id = ''
    },

    async load53aiAppOptions() {
      let list = []
      if (this.agent_type === AGENT_TYPES['53AI_AGENT']) {
        list = await channelApi.chat53aiAgentList()
      } else {
        list = await channelApi.chat53aiWorkflowList()
      }

      this.chat53ai_app_options = (list || []).map(item => {
        item.value = item.bot_id || 0
        item.label = item.name || ''
        item.icon = item.logo
        return item
      })
      if (this.chat53ai_app_options.length && !this.form_data.custom_config.chat53ai_agent_id) {
        const agent = this.chat53ai_app_options[0]
        if (!this.agent_id) {
          this.form_data.logo = agent.logo
          this.form_data.name = agent.name
          this.form_data.description = agent.description || ''
        }
        this.form_data.custom_config.chat53ai_agent_id = agent.value
        this.form_data.settings.opening_statement = agent.opening_statement
        this.form_data.settings.suggested_questions = agent.suggested_questions.map(item => {
          return {
            id: Math.random().toString(36).substring(2, 15),
            content: item,
          }
        })
      }
      if (
        !this.chat53ai_app_options.find(
          item => item.value === this.form_data.custom_config.chat53ai_agent_id
        )
      )
        this.form_data.custom_config.chat53ai_agent_id = ''
    },
    resetState() {
      this.saving = false
      this.loading = false
      this.agent_id = 0
      this.agent_type = 'prompt'
      this.is_new = false
      this.form_data = {
        logo: '',
        name: '',
        group_id: 0,
        description: '',
        channel_type: 0,
        model: '',
        sort: 0,
        prompt: '',
        user_group_ids: [],
        subscription_group_ids: [],
        tools: [],
        use_cases: [],
        configs: { completion_params: DEFAULT_COMPLETION_PARAMS },
        custom_config: {
          agent_type: 'prompt',
          agent_mode: 'chat',
          channel_id: 0,
          coze_workspace_id: '',
          coze_bot_id: '',
          app_builder_bot_id: '',
          channel_config: {},
          file_parse: {
            enable: false,
          },
          image_parse: {
            enable: false,
          },
        },
        settings: {
          opening_statement: '',
          suggested_questions: [],
          file_parse: {
            enable: false,
          },
          image_parse: {
            enable: false,
          },
          relate_agents: [],
          input_fields: [],
          output_fields: [],
        },
      }
      this.agent_data = {}
      this.group_options = []
      this.model_options = []
      this.coze_workspace_options = []
      this.coze_bot_options = []
      this.app_builder_bot_options = []
    },

    async saveAgentData({ hideToast = false } = {}) {
      const {
        logo = '',
        name = '',
        group_id = 0,
        description = '',
        model = '',
        channel_type = 0,
        prompt = '',
        user_group_ids = [],
        subscription_group_ids = [],
        use_cases = [],
        tools = [],
        sort = 0,
        configs = {},
        enable,
        custom_config = {},
        settings = {},
      } = this.form_data
      const data = {
        agent_id: this.agent_id || 0,
        agent_type: BACKEND_AGENT_TYPE.AGENT,
        channel_type,
        model,
        logo,
        name,
        group_id: +group_id || 0,
        description,
        sort,
        prompt,
        user_group_ids,
        subscription_group_ids,
        use_cases,
        tools,
        configs,
        enable,
        custom_config,
        settings,
      }
      const agentConfig = getAgentByAgentType(this.agent_type as AgentType)

      if (!channel_type) data.channel_type = CHANNEL_TYPE_VALUE_MAP.get(this.agent_type) || 0

      switch (this.agent_type) {
        case AGENT_TYPE.PROMPT:
          data.custom_config.channel_id = +model.split('_')[0] || 0
          data.channel_type = +model.split('_')[1] || 0
          data.model = model.split('_')[2] || ''
          break
        case AGENT_TYPE.COZE_AGENT_CN:
          data.model = custom_config.coze_bot_id || ''
          break
        // case AGENT_TYPE.COZE_AGENT_OSV:

        // case AGENT_TYPE.COZE_WORKFLOW_OSV:
        //   const params = custom_config.coze_bot_url.split('/')
        //   data.model =

        case AGENT_TYPE.COZE_WORKFLOW_CN:
          const params = new URLSearchParams(custom_config.coze_bot_url.split('?')[1])
          data.model = `workflow-${params.get('workflow_id')}` || ''
          break

        case AGENT_TYPE.APP_BUILDER:
          data.model = custom_config.app_builder_bot_id || ''
          break
        case AGENT_TYPE['53AI_AGENT']:
          data.model = custom_config.chat53ai_agent_id || ''
          break
        case AGENT_TYPE['53AI_WORKFLOW']:
          data.model = `workflow-${custom_config.chat53ai_agent_id}` || ''
          break
      }
      if (agentConfig && agentConfig.mode === AGENT_MODES.COMPLETION) {
        data.agent_type = BACKEND_AGENT_TYPE.WORKFLOW
      }
      data.custom_config.agent_type = this.agent_type
      data.custom_config.agent_mode = getAgentByAgentType(this.agent_type).mode || 'chat'
      this.saving = true
      const result_data = await agentApi.save({ data }).finally(() => {
        this.saving = false
      })
      if (!hideToast) ElMessage.success(window.$t('action_save_success'))
      this.agent_data = result_data
      this.agent_id = result_data.agent_id
      return this.agent_data
    },
  },
})
