declare namespace Agent {
	interface UseCase {
		type: 'case' | 'scene';
		id: string;
		input_text?: string;
		output_text?: string;
		image?: string;
		scene?: string;
		desc?: string;
	}

	interface CompletionParams {
		temperature: number;
		top_p: number;
		presence_penalty: number;
		frequency_penalty: number;
	}

	interface CustomConfig {
		agent_type: string;
		coze_workspace_id: string;
		coze_bot_id: string;
		app_builder_bot_id: string;
		file_parse: {
			enable: boolean;
		};
		image_parse: {
			enable: boolean;
		};
	}

	interface State {
		agent_id: number;
		eid: number;
		name: string;
		logo: string;
		sort: number;
		description: string;
		channel_type: number;
		model: string;
		prompt: string;
		configs: string;  // 或者可以使用 CompletionParams 类型
		tools: string;    // 可以根据实际内容定义具体的工具类型
		group_id: number;
		use_cases: string;  // 或者可以使用 UseCase[] 类型
		created_by: number;
		custom_config: string;  // 或者可以使用 CustomConfig 类型
		user_group_ids: number[];
		user_group_names: string[];
    settings: {
      opening_statement: string
      suggested_questions: string[]
      file_parse: {
        enable: boolean
      }
      image_parse: {
        enable: boolean
      }
      input_fields: Field[]
      output_fields: Field[]
    }
		enable: boolean;
		created_time: number;
		updated_time: number;
	}

  interface RelateAgent {
    agent_id: number
    name: string
    logo: string
    description: string
    input_fields: Field[]
    field_mapping: Record<string, string>
    execution_rule: 'auto' | 'manual'
  }

	interface FormData  {
		logo: string
		name: string
		group_id: number
		description: string
		channel_type: number
		model: string
		sort: number
		prompt: string
		user_group_ids: number[]
		subscription_group_ids: number[]
		tools: any[]
		use_cases: any[]
		configs: Record<string, any>
		enable: boolean
		custom_config: {
			agent_type: string
      agent_mode: string
      provider_id: number
			coze_workspace_id: string
			coze_bot_id: string
      coze_bot_url: string
			app_builder_bot_id: string
			chat53ai_agent_id: string
      channel_config: Record<string, any>
		}
    settings: {
      opening_statement: string
      suggested_questions: string[]
      file_parse: {
        enable: boolean
      }
      image_parse: {
        enable: boolean
      }
      relate_agents: RelateAgent[]
      input_fields: Field[]
      output_fields: Field[]
    }
	}


  interface Field {
    id: string
    // 变量名
    variable: string
    // 类型
    type: string
    // 显示名
    label: string
    // 描述
    desc: string
    // 是否必填
    required: boolean
    // 是否多选
    multiple: boolean
    // 选项
    options: FieldOption[]
    // 最大长度
    max_length: number
    // 是否显示字数统计
    show_word_limit: boolean
    // 时间格式
    date_format: string
    // 是否系统变量
    is_system: boolean
    // 文件类型
    file_type: 'all' | 'custom'
    // 支持文件格式
    file_accept: string[]
    // 单个文件上限
    file_size: number
    // 上传最大数量
    file_limit: number
  }

  interface FieldOption {
    id: string
    label: string
    value: string
  }
}
