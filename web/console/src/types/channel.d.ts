export interface ChannelConfig {
  channel_id: number
  key: string
  base_url: string
  models: string[]
  model?: string
  config: {
    agent_type: 'chat' | 'completion' | 'workflow'
  }
  name: string
}

export interface ChannelInfo {
  value: ChannelConfig
}

export interface ChannelForm {
  key: string
  base_url: string
  models: string[]
  model: string
  config: {
    agent_type: 'chat' | 'completion' | 'workflow'
  }
}
