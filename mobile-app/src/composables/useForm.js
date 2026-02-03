import { ref, computed, watch } from 'vue'
import { showDialog } from 'vant'

/**
 * useForm - 表单处理Composable
 *
 * @param {Object} formData - 表单初始数据
 * @param {Object} rules - 验证规则
 * @param {Object} options - 配置选项
 * @returns {Object} 表单处理相关方法
 */
export function useForm(formData, rules = {}, options = {}) {
  const {
    onSubmit = null,
    validateOnChange = false,
  } = options

  const form = ref({ ...formData })
  const errors = ref({})
  const touched = ref({})
  const submitting = ref(false)

  // 获取字段值
  const getValue = (field) => form.value[field]

  // 设置字段值
  const setValue = (field, value) => {
    form.value[field] = value
    if (validateOnChange) {
      validateField(field)
    }
  }

  // 批量设置值
  const setValues = (values) => {
    Object.assign(form.value, values)
    if (validateOnChange) {
      validate()
    }
  }

  // 重置表单
  const reset = () => {
    form.value = { ...formData }
    errors.value = {}
    touched.value = {}
  }

  // 验证单个字段
  const validateField = (field) => {
    const value = form.value[field]
    const fieldRules = rules[field]

    if (!fieldRules) {
      return true
    }

    // 清除之前的错误
    delete errors.value[field]

    // 必填验证
    if (fieldRules.required && (value === undefined || value === null || value === '')) {
      errors.value[field] = fieldRules.message || `${field} 是必填项`
      return false
    }

    // 最小长度验证
    if (fieldRules.minLength && String(value).length < fieldRules.minLength) {
      errors.value[field] = fieldRules.message || `${field} 长度不能少于 ${fieldRules.minLength}`
      return false
    }

    // 最大长度验证
    if (fieldRules.maxLength && String(value).length > fieldRules.maxLength) {
      errors.value[field] = fieldRules.message || `${field} 长度不能超过 ${fieldRules.maxLength}`
      return false
    }

    // 自定义验证
    if (fieldRules.validator && typeof fieldRules.validator === 'function') {
      const result = fieldRules.validator(value, form.value)
      if (result !== true) {
        errors.value[field] = result || fieldRules.message || `${field} 验证失败`
        return false
      }
    }

    // 正则验证
    if (fieldRules.pattern && !fieldRules.pattern.test(value)) {
      errors.value[field] = fieldRules.message || `${field} 格式不正确`
      return false
    }

    return true
  }

  // 验证所有字段
  const validate = () => {
    errors.value = {}
    let isValid = true

    for (const field in rules) {
      if (!validateField(field)) {
        isValid = false
      }
    }

    return isValid
  }

  // 获取错误信息
  const getError = (field) => errors.value[field] || ''

  // 是否有错误
  const hasError = computed(() => Object.keys(errors.value).length > 0)

  // 是否有效
  const isValid = computed(() => !hasError.value)

  // 标记字段为已触摸
  const touch = (field) => {
    touched.value[field] = true
  }

  // 提交表单
  const submit = async () => {
    // 先验证
    if (!validate()) {
      const firstError = Object.values(errors.value)[0]
      showDialog({
        message: firstError,
        confirmButtonText: '确定',
      })
      return false
    }

    submitting.value = true
    try {
      if (onSubmit) {
        const result = await onSubmit(form.value)
        return result
      }
      return form.value
    } catch (error) {
      console.error('表单提交失败:', error)
      throw error
    } finally {
      submitting.value = false
    }
  }

  // 监听表单变化
  if (validateOnChange) {
    watch(
      form,
      (newVal) => {
        Object.keys(newVal).forEach(field => {
          if (touched.value[field]) {
            validateField(field)
          }
        })
      },
      { deep: true }
    )
  }

  return {
    form,
    errors,
    touched,
    submitting,
    hasError,
    isValid,
    getValue,
    setValue,
    setValues,
    reset,
    validate,
    validateField,
    getError,
    touch,
    submit,
  }
}

/**
 * 常用验证规则
 */
export const ValidationRules = {
  // 必填
  required: (message = '此项为必填项') => ({
    required: true,
    message,
  }),

  // 手机号
  phone: (message = '请输入正确的手机号') => ({
    pattern: /^1[3-9]\d{9}$/,
    message,
  }),

  // 邮箱
  email: (message = '请输入正确的邮箱地址') => ({
    pattern: /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/,
    message,
  }),

  // 身份证号
  idCard: (message = '请输入正确的身份证号') => ({
    pattern: /^[1-9]\d{5}(18|19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$/,
    message,
  }),

  // 密码（6-20位，包含字母和数字）
  password: (message = '密码需6-20位，包含字母和数字') => ({
    pattern: /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,20}$/,
    message,
  }),

  // 最小长度
  minLength: (length, message) => ({
    minLength: length,
    message: message || `最少需要 ${length} 个字符`,
  }),

  // 最大长度
  maxLength: (length, message) => ({
    maxLength: length,
    message: message || `最多只能 ${length} 个字符`,
  }),

  // 自定义验证器
  custom: (validator, message = '验证失败') => ({
    validator,
    message,
  }),
}
