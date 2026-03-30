/**
 * 统一日志工具
 * 开发环境输出所有日志，生产环境只输出错误日志
 */

const isDev = import.meta.env.DEV

export const logger = {
  /**
   * 普通日志 - 仅开发环境
   */
  log: (...args) => {
    if (isDev) {
      console.log('[LOG]', ...args)
    }
  },

  /**
   * 信息日志 - 仅开发环境
   */
  info: (...args) => {
    if (isDev) {
      console.info('[INFO]', ...args)
    }
  },

  /**
   * 警告日志 - 仅开发环境
   */
  warn: (...args) => {
    if (isDev) {
      console.warn('[WARN]', ...args)
    }
  },

  /**
   * 错误日志 - 始终输出
   */
  error: (...args) => {
    console.error('[ERROR]', ...args)
  },

  /**
   * 调试日志 - 仅开发环境
   */
  debug: (...args) => {
    if (isDev) {
      console.log('[DEBUG]', ...args)
    }
  },

  /**
   * 分组开始 - 仅开发环境
   */
  group: (label) => {
    if (isDev) {
      console.group(label)
    }
  },

  /**
   * 分组结束 - 仅开发环境
   */
  groupEnd: () => {
    if (isDev) {
      console.groupEnd()
    }
  },

  /**
   * 表格输出 - 仅开发环境
   */
  table: (data) => {
    if (isDev) {
      console.table(data)
    }
  },

  /**
   * 时间追踪开始 - 仅开发环境
   */
  time: (label) => {
    if (isDev) {
      console.time(label)
    }
  },

  /**
   * 时间追踪结束 - 仅开发环境
   */
  timeEnd: (label) => {
    if (isDev) {
      console.timeEnd(label)
    }
  }
}

export default logger
