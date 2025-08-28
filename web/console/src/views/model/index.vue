<template>
  <Layout class="px-[60px] py-8">
    <Header title="模型设置" />
    <div
      v-loading="isLoading"
      class="flex-1 flex flex-col bg-white py-6 px-10 mt-3 box-border max-h-[calc(100vh-100px)] overflow-auto"
    >
      <div class="max-w-[600px]">
        <ElForm label-position="top">
          <ElFormItem label="逻辑推理">
            <ModelSelect v-model="logicValue" :type="MODEL_USE_TYPE.REASONING" />
            <!-- <ElSelect v-model="formData.data1" class="w-[600px] h-10">
              <ElOption></ElOption>
            </ElSelect> -->
          </ElFormItem>

          <ElFormItem label="向量嵌入">
            <ModelSelect v-model="vectorValue" :type="MODEL_USE_TYPE.EMBEDDING" />
          </ElFormItem>

          <ElFormItem label="检索设置">
            <div class="flex-1">
              <el-radio
                v-model="searchConfig.vector"
                :label="true"
                @change="handleSearchMethodChange('vector', $event)"
              >
                向量检索
              </el-radio>
              <el-radio
                v-model="searchConfig.fulltext"
                :label="true"
                @change="handleSearchMethodChange('fulltext', $event)"
              >
                全文检索
              </el-radio>
              <el-radio
                v-model="searchConfig.hybrid"
                :label="true"
                @change="handleSearchMethodChange('hybrid', $event)"
              >
                混合检索
              </el-radio>
              <!-- </el-radio-group> -->
              <div class="border rounded p-4 mt-4">
                <div v-if="searchConfig.hybrid">
                  <div class="text-[#182B50] text-sm">混合检索</div>
                  <div class="text-[#182b50] text-opacity-60 text-xs mt-1">
                    同时执行向量检索和全文检索，通过Rerank模型（必须启用）进行重新排序，匹配用户问题的最佳结果
                  </div>
                </div>
                <div v-else-if="searchConfig.vector">
                  <div class="text-[#182B50] text-sm">向量检索</div>
                  <div class="text-[#182b50] text-opacity-60 text-xs mt-1">通过向量检索查询内容最相似的语料切片</div>
                </div>
                <div v-else-if="searchConfig.fulltext">
                  <div class="text-[#182B50] text-sm">全文检索</div>
                  <div class="text-[#182b50] text-opacity-60 text-xs mt-1">通过全文检索查询包含词汇的语料切片</div>
                </div>
                <div class="border-b border-b-[#EBEEF5] mt-4" />

                <div v-if="searchConfig.hybrid" class="flex gap-4 mt-4">
                  <div
                    :class="[
                      searchConfig.rerank_model === RERANKING_MODE.WEIGHTED_SCORE
                        ? 'border-[#2563EB] bg-[#F6F9FE] text-[#2563EB]'
                        : 'text-[#182B50]',
                    ]"
                    class="flex-1 h-10 rounded border flex-center gap-1 cursor-pointer text-opacity-80 text-sm"
                    @click="handleRerankMode(RERANKING_MODE.WEIGHTED_SCORE)"
                  >
                    权重计算
                    <el-tooltip
                      content="通过调整分配的权重，重新排序策略确定是优先进行语义匹配还是关键字匹配。"
                      placement="top"
                    >
                      <svg-icon name="help" width="14px" color="A4AAB9" />
                    </el-tooltip>
                  </div>
                  <div
                    :class="[
                      searchConfig.rerank_model === RERANKING_MODE.RERANKING_MODEL
                        ? 'border-[#2563EB] bg-[#F6F9FE] text-[#2563EB]'
                        : 'text-[#182B50]',
                    ]"
                    class="flex-1 h-10 rounded border flex-center gap-1 cursor-pointer text-opacity-80 text-sm bg-[#F6F9FE]"
                    @click="handleRerankMode(RERANKING_MODE.RERANKING_MODEL)"
                  >
                    模型计算
                    <el-tooltip
                      content="重排序模型将根据候选文档列表与用户问题语义匹配度进行重新排序，从而改进语义排序的结果"
                      placement="top"
                    >
                      <svg-icon name="help" width="14px" color="A4AAB9" />
                    </el-tooltip>
                  </div>
                </div>
                <div
                  v-if="searchConfig.hybrid && searchConfig.rerank_model === RERANKING_MODE.WEIGHTED_SCORE"
                  class="mt-5 rounded border px-5 py-4"
                >
                  <el-slider
                    v-model="searchConfig.weights.vector_setting.vector_weight"
                    :min="0"
                    :step="0.1"
                    :max="1"
                    style="--el-slider-runway-bg-color: #7575ff"
                    @input="handleWeightChange"
                  />
                  <div class="flex items-center justify-between">
                    <div class="text-sm text-[#2590F5]">
                      语义 {{ searchConfig.weights.vector_setting.vector_weight }}
                    </div>
                    <div class="text-sm text-[#7575FF]">
                      关键词 {{ searchConfig.weights.keyword_setting.keyword_weight }}
                    </div>
                  </div>
                </div>

                <div
                  v-if="!searchConfig.hybrid || searchConfig.rerank_model === RERANKING_MODE.RERANKING_MODEL"
                  class="mt-4"
                >
                  <div class="flex items-center text-[#182B50] text-opacity-80 text-sm">
                    Rerank模型
                    <el-tooltip
                      effect="dark"
                      content="重排序模型将根据候选文档列表与用户问题语义匹配度进行重新排序，从而改进语义排序的结果"
                      placement="top"
                    >
                      <el-icon class="ml-2">
                        <QuestionFilled />
                      </el-icon>
                    </el-tooltip>

                    <el-switch
                      v-if="searchConfig.fulltext || searchConfig.vector"
                      v-model="searchConfig.reranking_enable"
                      class="ml-4"
                    />
                  </div>
                  <div v-if="searchConfig.reranking_enable || searchConfig.hybrid" class="mt-4">
                    <ModelSelect v-model="rerankValue" :type="MODEL_USE_TYPE.RERANKER" />
                  </div>
                </div>
                <div class="border-t border-dashed my-4" />
                <div class="flex items-center">
                  <div class="flex-none flex items-center w-[214px] text-sm text-[#182B50] text-opacity-80">
                    召回数量
                    <el-tooltip
                      effect="dark"
                      content="用于筛选与用户问题相似度最高的文本片段。系统同时会根据选用模型上下文窗口大小动态调整分段数量。"
                      placement="top"
                    >
                      <el-icon class="ml-2">
                        <QuestionFilled />
                      </el-icon>
                    </el-tooltip>
                  </div>
                  <div class="flex-1 flex items-center">
                    <el-slider v-model="searchConfig.top_k" :min="1" :step="1" :max="10" />
                    <span class="flex-none w-9 text-right text-[#182B50] text-opacity-80 text-sm">{{
                      searchConfig.top_k
                    }}</span>
                  </div>
                </div>
                <div class="flex items-center mt-5">
                  <div class="flex-none flex items-center w-[214px] text-sm text-[#182B50] text-opacity-80">
                    分数阈值
                    <el-tooltip effect="dark" content="用于设置文本片段筛选的相似度阈值。" placement="top">
                      <el-icon class="ml-2">
                        <QuestionFilled />
                      </el-icon>
                    </el-tooltip>

                    <el-switch v-model="searchConfig.score_threshold_enabled" class="ml-4" />
                  </div>
                  <div class="flex-1 flex items-center">
                    <el-slider
                      v-model="searchConfig.score_threshold"
                      :disabled="!searchConfig.score_threshold_enabled"
                      :min="0"
                      :step="0.01"
                      :format-tooltip="(value: number) => `${(value * 100).toFixed(0)}`"
                      :max="1"
                    />
                    <span class="flex-none w-9 text-right text-[#182B50] text-opacity-80 text-sm">{{
                      (searchConfig.score_threshold * 100).toFixed(0)
                    }}</span>
                  </div>
                </div>
              </div>
            </div>
          </ElFormItem>
        </ElForm>

        <ElButton v-debounce type="primary" class="mt-8 w-24 h-9" @click="handleSave">保存</ElButton>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { QuestionFilled } from '@element-plus/icons-vue'
import ModelSelect from '@/components/Model/select.vue'
import { chunkSettingApi, ModelSetting } from '@/api/modules/chunk-setting'
import { deepCopy } from '@/utils'
import { MODEL_USE_TYPE } from '@/constants/platform/config'

const RERANKING_MODE = {
  WEIGHTED_SCORE: 'weighted_score',
  RERANKING_MODEL: 'reranking_model',
}

type RerankingMode = (typeof RERANKING_MODE)[keyof typeof RERANKING_MODE]

const defaultSetting: ModelSetting = {
  created_time: 1672502400,
  eid: 1,
  file_id: 1,
  id: 1,
  library_id: 1,
  model_config: {
    logic_reasoning: {
      channel_id: 0,
      model_name: 'string',
    },
    search_config: {
      fulltext: false,
      hybrid: false,
      rerank_model: RERANKING_MODE.RERANKING_MODEL,
      score_threshold: 0,
      top_k: 0,
      vector: true,
      rerank_channel_id: 0,
      rerank_model_name: '',
      reranking_enable: false,
      score_threshold_enabled: false,
      weights: {
        keyword_setting: {
          keyword_weight: 0,
        },
        vector_setting: {
          vector_weight: 0,
        },
      },
    },
    vector_embedding: {
      channel_id: 0,
      model_name: 'string',
    },
    version: 'string',
  },
  updated_time: 1672502400,
}

const isLoading = ref(false)
const setting = ref<ModelSetting>(deepCopy(defaultSetting))

const logicValue = computed({
  get: () => {
    const { channel_id, model_name } = setting.value.model_config.logic_reasoning
    return channel_id && model_name ? `${channel_id}_${model_name}` : ''
  },
  set: value => {
    const [channel_id, model_name] = value.split('_')
    setting.value.model_config.logic_reasoning.channel_id = Number(channel_id)
    setting.value.model_config.logic_reasoning.model_name = model_name
  },
})

const vectorValue = computed({
  get: () => {
    const { channel_id, model_name } = setting.value.model_config.vector_embedding
    return channel_id && model_name ? `${channel_id}_${model_name}` : ''
  },
  set: value => {
    const [channel_id, model_name] = value.split('_')
    setting.value.model_config.vector_embedding.channel_id = Number(channel_id)
    setting.value.model_config.vector_embedding.model_name = model_name
  },
})

const rerankValue = computed({
  get: () => {
    const { rerank_channel_id, rerank_model_name } = setting.value.model_config.search_config
    return rerank_channel_id && rerank_model_name ? `${rerank_channel_id}_${rerank_model_name}` : ''
  },
  set: value => {
    const [channel_id, model_name] = value.split('_')
    setting.value.model_config.search_config.rerank_channel_id = Number(channel_id)
    setting.value.model_config.search_config.rerank_model_name = model_name
  },
})

const searchConfig = computed(() => {
  return setting.value.model_config.search_config
})

const handleSearchMethodChange = (method: keyof ModelSetting['model_config']['search_config'], event: boolean) => {
  setting.value.model_config.search_config.hybrid = false
  setting.value.model_config.search_config.vector = false
  setting.value.model_config.search_config.fulltext = false
  setting.value.model_config.search_config[method] = event
}

const handleRerankMode = (mode: RerankingMode) => {
  setting.value.model_config.search_config.rerank_model = mode
}

const handleWeightChange = (value: number) => {
  setting.value.model_config.search_config.score_threshold = (100 - value * 100) / 100
}

const handleSave = async () => {
  if (!logicValue.value) {
    ElMessage.error(window.$t('message_status.logic_reasoning_required'))
    return
  }
  if (!vectorValue.value) {
    ElMessage.error(window.$t('message_status.vector_embedding_required'))
    return
  }

  const model_config = setting.value.model_config
  const showRerank =
    (model_config.search_config.hybrid && model_config.search_config.rerank_model === RERANKING_MODE.RERANKING_MODEL) ||
    (!model_config.search_config.hybrid && model_config.search_config.reranking_enable)
  if (showRerank && !rerankValue.value) {
    ElMessage.error(window.$t('message_status.rerank_model_required'))
    return
  }
  await chunkSettingApi.modelConfig.update({
    model_config: setting.value.model_config,
  })
  ElMessage.success(window.$t('message_status.save_success'))
}

const loadConfig = async () => {
  isLoading.value = true
  const data = await chunkSettingApi.modelConfig.get()
  const { search_config } = data.model_config
  if (search_config.rerank_model === '') {
    search_config.rerank_model = RERANKING_MODE.RERANKING_MODEL
  }

  // 判断是否只有一个搜索方法，如果不只一个，则默认设置为vector为true
  const searchMethods = [search_config.fulltext, search_config.vector, search_config.hybrid]
  const notOnlyOne = searchMethods.filter(Boolean).length !== 1
  if (notOnlyOne) {
    search_config.vector = false
  }

  setting.value = data
  isLoading.value = false
}

onMounted(async () => {
  loadConfig()
})
</script>

<style scoped>
:deep(.el-form-item__label) {
  padding-right: 0 !important;
}
</style>
