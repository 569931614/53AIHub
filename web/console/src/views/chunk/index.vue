<template>
  <Layout class="px-[60px] py-8">
    <Header :title="$t('module.chunk_setting')" />

    <div v-loading="isLoading" class="flex-1 flex flex-col bg-white p-6 mt-3 overflow-y-auto">
      <div class="flex-1 flex flex-col overflow-y-auto">
        <div class="flex items-center gap-2.5 mb-4">
          <div class="w-1 h-4 bg-[#2563EB]"></div>
          <span class="text-sm text-[#1D1E1F] mr-2">将知识标题添加至</span>
          <el-checkbox v-model="setting.chunking_config.knowledge_chunking.include_title" class="!mr-0" :value="true"
            >知识点</el-checkbox
          >
          <el-checkbox v-model="setting.chunking_config.index_chunking.include_title" :value="true">索引块</el-checkbox>
        </div>

        <!-- 知识点配置 -->
        <div class="border rounded">
          <div class="flex items-center gap-2 p-5 border-b">
            <img class="size-10" :src="$getRealPath({ url: '/images/chunk/knowledge.webp' })" />
            <div>
              <h4 class="text-sm text-[#1D1E1F]">知识点</h4>
              <p class="text-xs text-[#999999]">将知识原文按拆分规则分拆成若干个知识点</p>
            </div>
          </div>
          <div class="p-5 flex flex-col gap-4">
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">拆分规则</div>
              <el-radio-group v-model="setting.knowledge_chunking_type">
                <el-radio :value="SPLIT_TYPE.NONE">不拆分</el-radio>
                <el-radio :value="SPLIT_TYPE.HEADING">
                  <div class="flex items-center">
                    <el-dropdown @command="handleChangeHeading('knowledge', $event)">
                      <div class="flex items-center gap-2">
                        {{ getHeadingLabel('knowledge') }}
                        <el-icon><ArrowDown /></el-icon>
                      </div>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item v-for="item in headerList" :key="item.type" :command="item.type">
                            {{ item.label }}
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </el-radio>
                <el-radio :value="SPLIT_TYPE.CUSTOM">
                  <div class="flex items-center gap-2">
                    <span>指定</span>
                    <div>
                      <el-input v-model="setting.knowledge_chunking_input" />
                    </div>
                    <el-dropdown @command="handleChangeCustom('knowledge', $event)">
                      <div class="flex items-center gap-2">
                        常用
                        <el-icon><ArrowDown /></el-icon>
                      </div>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item v-for="item in commonList" :key="item.value" :command="item.value">
                            {{ item.label }}
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </el-radio>
              </el-radio-group>
            </div>
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">最大长度</div>
              <div>
                <el-input-number
                  v-model="setting.chunking_config.knowledge_chunking.max_length"
                  class="el-input-number--left"
                  size="large"
                  :min="maxLength.min"
                  :max="maxLength.max"
                  :controls="false"
                  @blur="handleBlurMaxLength('knowledge')"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 索引块配置 -->
        <div class="border rounded mt-4">
          <div class="flex items-center gap-2 p-5 border-b">
            <img class="size-10" :src="$getRealPath({ url: '/images/chunk/index.webp' })" />
            <div>
              <h4 class="text-sm text-[#1D1E1F]">索引块</h4>
              <p class="text-xs text-[#999999]">将知识点按拆分规则分拆成若干个索引块</p>
            </div>
          </div>
          <div class="p-5 flex flex-col gap-4">
            <div class="text-sm text-[#1D1E1F] font-semibold">默认索引</div>

            <!-- 拆分规则 -->
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">拆分规则</div>
              <el-radio-group v-model="setting.index_chunking_type">
                <el-radio :value="SPLIT_TYPE.NONE">不拆分</el-radio>
                <el-radio :value="SPLIT_TYPE.HEADING">
                  <div class="flex items-center">
                    <el-dropdown @command="handleChangeHeading('index', $event)">
                      <div class="flex items-center gap-2">
                        {{ getHeadingLabel('index') }}
                        <el-icon><ArrowDown /></el-icon>
                      </div>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item v-for="item in headerList" :key="item.type" :command="item.type">
                            {{ item.label }}
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </el-radio>
                <el-radio :value="SPLIT_TYPE.CUSTOM">
                  <div class="flex items-center gap-2">
                    <span>指定</span>
                    <div>
                      <el-input v-model="setting.index_chunking_input" />
                    </div>
                    <el-dropdown @command="handleChangeCustom('index', $event)">
                      <div class="flex items-center gap-2">
                        常用
                        <el-icon><ArrowDown /></el-icon>
                      </div>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item v-for="item in commonList" :key="item.value" :command="item.value">
                            {{ item.label }}
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </el-radio>
              </el-radio-group>
            </div>

            <!-- 最大长度 -->
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">最大长度</div>
              <div>
                <el-input-number
                  v-model="setting.chunking_config.index_chunking.max_length"
                  class="el-input-number--left"
                  size="large"
                  :min="0"
                  :max="1000000"
                  :controls="false"
                  @blur="handleBlurMaxLength('index')"
                />
              </div>
            </div>

            <!-- 内容概要 -->
            <div class="text-sm text-[#1D1E1F] font-semibold">内容概要</div>
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">生成规则</div>
              <el-radio-group v-model="setting.chunking_config.content_summary.generation_method">
                <el-radio :value="SUMMARY_GENERATION.MANUAL">人工撰写</el-radio>
                <el-radio :value="SUMMARY_GENERATION.AI">AI生成</el-radio>
              </el-radio-group>
            </div>

            <!-- 常见问法 -->
            <div class="text-sm text-[#1D1E1F] font-semibold">常见问法</div>
            <div class="flex items-center">
              <div class="flex-none w-24 text-sm text-[#4F5052]">生成规则</div>
              <el-radio-group v-model="setting.chunking_config.common_questions.generation_method">
                <el-radio :value="QUESTION_GENERATION.MANUAL">人工撰写</el-radio>
                <el-radio :value="QUESTION_GENERATION.AI">AI生成</el-radio>
              </el-radio-group>
            </div>
          </div>
        </div>
      </div>
      <div class="flex-none">
        <el-button v-debounce type="primary" size="large" class="mt-4" @click="handleConfirm">保存</el-button>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ArrowDown } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'
import { chunkSettingApi } from '@/api/modules/chunk-setting'
import type { ChunkSetting } from '@/api/modules/chunk-setting'
import { CHUNK_SETTING_DEFAULT, SUMMARY_GENERATION, QUESTION_GENERATION, SPLIT_TYPE } from '@/constants/chunk'
import { deepCopy } from '@/utils'

interface Setting extends ChunkSetting {
  knowledge_chunking_type: string
  knowledge_chunking_head: string
  knowledge_chunking_input: string
  index_chunking_type: string
  index_chunking_head: string
  index_chunking_input: string
}

// 常量配置
const CONFIG = {
  maxLength: { min: 50, max: 50000 },
  headerList: [
    { type: 'h1', label: '一级标题（H1）' },
    { type: 'h2', label: '二级标题（H2）' },
    { type: 'h3', label: '三级标题（H3）' },
    { type: 'h4', label: '四级标题（H4）' },
    { type: 'h5', label: '五级标题（H5）' },
  ],
  commonList: [
    { label: '1 个换行符（\\n）', value: '\\n' },
    { label: '2 个换行符（\\n\\n）', value: '\\n\\n' },
    { label: '句号（。）', value: '。' },
    { label: '感叹号（！）', value: '！' },
    { label: '问号（？）', value: '？' },
    { label: '分号（；）', value: '；' },
    { label: '分割线（---）', value: '---' },
  ],
}

const { maxLength, headerList, commonList } = CONFIG

const isLoading = ref(false)
const setting = ref<Setting>(
  deepCopy({
    ...CHUNK_SETTING_DEFAULT,
    knowledge_chunking_type: SPLIT_TYPE.HEADING,
    knowledge_chunking_head: headerList[0].type,
    knowledge_chunking_input: commonList[0].value,
    index_chunking_type: SPLIT_TYPE.HEADING,
    index_chunking_head: headerList[0].type,
    index_chunking_input: commonList[0].value,
  })
)

// 计算属性
const getHeadingLabel = (type: 'knowledge' | 'index') => {
  const chunkHead = type === 'knowledge' ? 'knowledge_chunking_head' : 'index_chunking_head'
  const label = headerList.find(item => item.type === setting.value[chunkHead])?.label
  return label || headerList[0].label
}

const handleChangeHeading = (type: 'knowledge' | 'index', value: string) => {
  const chunkHead = type === 'knowledge' ? 'knowledge_chunking_head' : 'index_chunking_head'
  setting.value[chunkHead] = value
}

const handleChangeCustom = (type: 'knowledge' | 'index', value: string) => {
  const chunkInput = type === 'knowledge' ? 'knowledge_chunking_input' : 'index_chunking_input'
  setting.value[chunkInput] = value
}

const handleBlurMaxLength = (type: 'knowledge' | 'index') => {
  const config =
    type === 'knowledge'
      ? setting.value.chunking_config.knowledge_chunking
      : setting.value.chunking_config.index_chunking

  if (type === 'knowledge') {
    config.max_length = Math.max(Math.min(config.max_length, maxLength.max), maxLength.min)
    // 同步更新索引块最大长度
    const indexConfig = setting.value.chunking_config.index_chunking
    indexConfig.max_length = Math.max(Math.min(indexConfig.max_length, config.max_length), maxLength.min)
  } else {
    const knowledgeConfig = setting.value.chunking_config.knowledge_chunking
    config.max_length = Math.max(Math.min(config.max_length, knowledgeConfig.max_length), maxLength.min)
  }
}

const handleConfirm = async () => {
  const data = { chunking_config: { ...setting.value.chunking_config } }

  // 统一处理拆分规则
  const processSplitRule = (type: 'knowledge' | 'index') => {
    const chunkType = type === 'knowledge' ? 'knowledge_chunking_type' : 'index_chunking_type'
    const chunkInput = type === 'knowledge' ? 'knowledge_chunking_input' : 'index_chunking_input'
    const chunkHead = type === 'knowledge' ? 'knowledge_chunking_head' : 'index_chunking_head'
    const config = type === 'knowledge' ? data.chunking_config.knowledge_chunking : data.chunking_config.index_chunking

    if (setting.value[chunkType] === SPLIT_TYPE.NONE) {
      config.split_rule = SPLIT_TYPE.NONE
    } else if (setting.value[chunkType] === SPLIT_TYPE.HEADING) {
      config.split_rule = setting.value[chunkHead]
    } else {
      config.split_rule = setting.value[chunkInput].trim() || 'none'
    }
  }

  processSplitRule('knowledge')
  processSplitRule('index')

  await chunkSettingApi.chunkingConfig.update(data)
  ElMessage.success($t('message_status.save_success'))
}

// 设置拆分规则
const setSplitRule = (config: Setting, prefix: 'knowledge_chunking' | 'index_chunking') => {
  const splitRule = config.chunking_config[prefix].split_rule
  const isHeading = headerList.some(item => item.type === splitRule)
  const isNone = splitRule === SPLIT_TYPE.NONE
  if (isHeading) {
    config[`${prefix}_type`] = SPLIT_TYPE.HEADING
    config[`${prefix}_head`] = splitRule
    config[`${prefix}_input`] = commonList[1].value
  } else if (isNone) {
    config[`${prefix}_type`] = SPLIT_TYPE.NONE
    config[`${prefix}_head`] = headerList[1].type
    config[`${prefix}_input`] = commonList[1].value
  } else {
    config[`${prefix}_type`] = SPLIT_TYPE.CUSTOM
    config[`${prefix}_head`] = headerList[1].type
    config[`${prefix}_input`] = splitRule
  }
}

onMounted(async () => {
  isLoading.value = true
  const data = await chunkSettingApi.chunkingConfig.get()
  const config = {
    ...data,
    knowledge_chunking_type: SPLIT_TYPE.HEADING,
    knowledge_chunking_head: headerList[0].type,
    knowledge_chunking_input: commonList[1].value,
    index_chunking_type: SPLIT_TYPE.HEADING,
    index_chunking_head: headerList[0].type,
    index_chunking_input: commonList[1].value,
  }

  setSplitRule(config, 'knowledge_chunking')
  setSplitRule(config, 'index_chunking')
  setting.value = config
  isLoading.value = false
})
</script>

<style scoped></style>
