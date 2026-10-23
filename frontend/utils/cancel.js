/**
 * ═══════════════════════════════════════════════
 *  统一取消作用域管理器（Cancel Scope）
 *
 *  设计：
 *  - 请求层（utils/request.js）只认 options.signal，保持通用
 *  - 业务层按「作用域 key」创建 AbortController，并作为 signal 传入请求
 *  - 任何接口（上传/生成/轮询/搜索…）都可用同一套机制取消
 *
 *  用法：
 *    const scope = createCancelScope(`upload:${id}`)
 *    await httpRequest.post({ url, params, signal: scope.signal })
 *    scope.dispose()                    // 请求结束后从注册表移除
 *
 *    cancelScope(`upload:${id}`)        // 取消单个作用域
 *    cancelScopes('upload:')            // 按前缀批量取消
 *
 *  约定：
 *  - key 使用 `业务名:实例id` 格式（如 upload:u123、task:abc），
 *    前缀即业务标签，便于批量取消
 *  - 同 key 重复创建时自动取消旧作用域，防止泄漏
 *  - 取消后请求层 AbortError 会被 resolve 为 null，调用方需容忍 null
 * ═══════════════════════════════════════════════
 */

const scopes = new Map() // key → AbortController

/**
 * 创建（并注册）一个取消作用域
 * @param {string} key 作用域唯一标识，建议 `业务:实例id`
 * @returns {{ key: string, signal: AbortSignal, aborted: boolean } & { cancel: () => boolean, dispose: () => void }}
 */
export function createCancelScope(key) {
  // 同 key 已存在：先取消旧的，避免悬挂请求
  if (scopes.has(key)) {
    scopes.get(key).abort()
    scopes.delete(key)
  }
  const controller = new AbortController()
  scopes.set(key, controller)

  return {
    key,
    signal: controller.signal,
    get aborted() {
      return controller.signal.aborted
    },
    /** 取消自身并从注册表移除 */
    cancel() {
      return cancelScope(key)
    },
    /** 请求正常结束（成功/失败）后清理注册表，必须调用以免泄漏 */
    dispose() {
      if (scopes.get(key) === controller) scopes.delete(key)
    },
  }
}

/**
 * 取消指定作用域
 * @returns {boolean} 是否命中了在途请求
 */
export function cancelScope(key) {
  const controller = scopes.get(key)
  if (!controller) return false
  controller.abort()
  scopes.delete(key)
  return true
}

/**
 * 按 key 前缀批量取消（如 cancelScopes('upload:') 取消所有上传）
 * @returns {number} 取消的数量
 */
export function cancelScopes(prefix) {
  let count = 0
  for (const [key, controller] of scopes) {
    if (key.startsWith(prefix)) {
      controller.abort()
      scopes.delete(key)
      count++
    }
  }
  return count
}

/** 当前所有在途作用域 key（调试用） */
export function listScopes() {
  return [...scopes.keys()]
}
