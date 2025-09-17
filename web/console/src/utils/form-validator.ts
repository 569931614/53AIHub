import { Ref } from 'vue'
import { ElForm } from 'element-plus'

/**
 * 表单字段验证工具函数
 * @param formRef 表单实例引用
 * @param field 要验证的字段名
 * @returns 验证结果
 */
export const validateFormField = async (
  formRef: Ref<InstanceType<typeof ElForm> | undefined>,
  field: string
): Promise<boolean> => {
  if (!formRef.value) {
    return false
  }

  try {
    await formRef.value.validateField(field)
    return true
  } catch (error) {
    return false
  }
}
