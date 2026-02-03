/**
 * 格式化工具函数
 */

/**
 * 格式化金额
 * @param {number} amount - 金额
 * @param {number} decimals - 小数位数
 * @param {string} currency - 货币符号
 */
export function formatMoney(amount, decimals = 2, currency = '¥') {
  if (amount === null || amount === undefined || isNaN(amount)) {
    return `${currency}0.00`
  }

  const num = Number(amount)
  return num.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }) + (currency ? ` ${currency}` : '')
}

/**
 * 格式化数字（带千分位）
 */
export function formatNumber(num, decimals = 0) {
  if (num === null || num === undefined || isNaN(num)) {
    return '0'
  }

  return Number(num).toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

/**
 * 格式化百分比
 */
export function formatPercent(value, decimals = 2) {
  if (value === null || value === undefined || isNaN(value)) {
    return '0%'
  }

  return `${(Number(value) * 100).toFixed(decimals)}%`
}

/**
 * 格式化日期
 * @param {string|Date} date - 日期
 * @param {string} format - 格式
 */
export function formatDate(date, format = 'YYYY-MM-DD') {
  if (!date) return ''

  const d = new Date(date)
  if (isNaN(d.getTime())) return ''

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 格式化日期时间（简短格式）
 */
export function formatDateTimeShort(date) {
  return formatDate(date, 'MM-DD HH:mm')
}

/**
 * 格式化日期时间（完整格式）
 */
export function formatDateTimeFull(date) {
  return formatDate(date, 'YYYY-MM-DD HH:mm:ss')
}

/**
 * 格式化相对时间
 */
export function formatRelativeTime(date) {
  if (!date) return ''

  const d = new Date(date)
  const now = new Date()
  const diff = now.getTime() - d.getTime()

  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)

  if (seconds < 60) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  if (months < 12) return `${months}个月前`
  return `${years}年前`
}

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 B'

  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 格式化手机号（隐藏中间4位）
 */
export function formatPhone(phone) {
  if (!phone) return ''
  const str = String(phone)
  if (str.length !== 11) return phone
  return str.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

/**
 * 格式化身份证号（隐藏中间部分）
 */
export function formatIdCard(idCard) {
  if (!idCard) return ''
  const str = String(idCard)
  if (str.length !== 18) return idCard
  return str.replace(/(\d{6})\d{8}(\d{4})/, '$1********$2')
}

/**
 * 格式化银行卡号（隐藏中间部分）
 */
export function formatBankCard(cardNo) {
  if (!cardNo) return ''
  const str = String(cardNo).replace(/\s/g, '')
  if (str.length < 10) return cardNo

  // 每4位一组
  const grouped = str.replace(/(\d{4})(?=\d)/g, '$1 ')
  // 显示前4后4，中间用****代替
  const first4 = str.slice(0, 4)
  const last4 = str.slice(-4)
  return `${first4} **** **** ${last4}`
}

/**
 * 截断文本
 */
export function truncate(text, maxLength = 50, suffix = '...') {
  if (!text) return ''
  const str = String(text)
  if (str.length <= maxLength) return str
  return str.slice(0, maxLength) + suffix
}

/**
 * 高亮关键词
 */
export function highlightKeywords(text, keywords, className = 'highlight') {
  if (!text || !keywords) return text

  const regex = new RegExp(`(${keywords})`, 'gi')
  return text.replace(regex, `<span class="${className}">$1</span>`)
}

/**
 * 格式化数量（大数字简化）
 */
export function formatQuantity(num) {
  if (num === null || num === undefined || isNaN(num)) return '0'

  const n = Number(num)
  if (n >= 100000000) {
    return (n / 100000000).toFixed(2) + '亿'
  }
  if (n >= 10000) {
    return (n / 10000).toFixed(2) + '万'
  }
  return n.toString()
}

/**
 * 获取状态颜色类型（用于Vant Tag）
 */
export function getStatusColor(status) {
  const colorMap = {
    // 成功状态
    success: 'success',
    approved: 'success',
    completed: 'success',
    active: 'success',

    // 警告状态
    warning: 'warning',
    pending: 'warning',
    processing: 'warning',

    // 危险状态
    danger: 'danger',
    rejected: 'danger',
    failed: 'danger',
    cancelled: 'danger',
    expired: 'danger',

    // 默认状态
    default: 'default',
    draft: 'default',
    unknown: 'default',

    // 主要状态
    primary: 'primary',
  }

  return colorMap[status] || 'default'
}
