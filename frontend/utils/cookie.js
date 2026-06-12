/**
 * Cookie 工具（token）
 * 前后端统一 cookie 名称：token
 * 基于 Nuxt3 useCookie 实现，支持 SSR
 */
const MAX_AGE = 60 * 60 * 24 // 1 天

const token = {
  /** 获取 token */
  get() {
    try {
      return useCookie('token').value ?? null
    } catch {
      return null
    }
  },

  /** 写入 token */
  set(value, maxAge = MAX_AGE) {
    const c = useCookie('token', { maxAge, path: '/' })
    c.value = value
  },

  /** 清除 token */
  remove() {
    const c = useCookie('token', { path: '/' })
    c.value = null
  },
}

export default token
