import { defineStore } from 'pinia'
import { channelApi } from '@/api/modules/channel'
import { deepCopy } from '@/utils'

export const useChannelStore = defineStore('channel-store', {
  state: () => ({
    channel_list: [],
    channel_models: [],
  }),
  actions: {
    async loadListData({ reset = true } = {}) {
      if (reset || !this.channel_list.length) this.channel_list = await channelApi.list()
      return JSON.parse(JSON.stringify(this.channel_list))
    },
    async save({ data = {} } = {}) {
      return channelApi.save({ data })
    },
    async delete({ data: { channel_id } }: { data: { channel_id: number } }) {
      return channelApi.delete({ data })
    },
    async loadChannelModels(owner = 'all') {
      if (!this.channel_models.length) {
        const list = await channelApi.modelList({ params: { owner } })
        this.channel_models = list
      }
      if (owner === 'all') return deepCopy(this.channel_models)
      const list = this.channel_models.filter((item: any) => item.owned_by === owner)
      return deepCopy(list)
    },
    async loadCozeWorkspaceList() {
      const list = await channelApi.cozeWorkspaceList()
      return list
    },
  },
})
