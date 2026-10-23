import { useRequestHeaders, useNuxtApp } from 'nuxt/app';
import { setEncrypt } from "~/utils/index.js"

const DEFAULT_CONFIG = {
  SERVER_URL: '/api/ajaxData',
  TIMEOUT: 60000,
};

class HttpRequest {
  constructor() {
    this.pendingRequests = new Map(); // 管理多个请求
  }

  /**
   * 生成请求唯一标识
   */
  _getRequestKey(url, params) {
    return `${url}_${JSON.stringify(params)}`;
  }

  /**
   * 取消重复请求
   */
  _cancelDuplicate(key) {
    if (this.pendingRequests.has(key)) {
      this.pendingRequests.get(key).abort();
      this.pendingRequests.delete(key);
    }
  }

  async request(options) {
    const {
      url,
      method = 'POST',
      timeout = DEFAULT_CONFIG.TIMEOUT,
      params = {},
      isCookie = true,
      cancelDuplicate = true, // 是否取消重复请求
    } = options;

    // 处理重复请求
    const requestKey = this._getRequestKey(url, params);
    if (cancelDuplicate) {
      this._cancelDuplicate(requestKey);
    }

    const controller = new AbortController();
    this.pendingRequests.set(requestKey, controller);

    // 外部 signal（如取消作用域）联动内部 controller，任一方 abort 均生效
    let removeExternalAbortListener = null;
    if (options.signal) {
      if (options.signal.aborted) {
        controller.abort();
      } else {
        const onExternalAbort = () => controller.abort();
        options.signal.addEventListener('abort', onExternalAbort, { once: true });
        // 请求正常结束时手动移除，避免监听器残留在外部 signal 上
        removeExternalAbortListener = () =>
          options.signal.removeEventListener('abort', onExternalAbort);
      }
    }

    return new Promise(async (resolve) => {
      try {
        const req = useRequestHeaders();
        const { _route } = useNuxtApp();
        const isDesign = _route?.query?.design === 'show';

        const fetchOptions = {
          method: 'post',
          headers: {
            ...req,
            ...(isCookie && req.cookie ? { cookie: req.cookie } : {}),
            ...(isDesign ? { design: true } : {}),
          },
          body: isDesign ? { url, method, params } : { aes: setEncrypt({ url, method, params }) },
          timeout,
          signal: controller.signal,
        };

        // 可选配置
        if (options.lazy){
          fetchOptions.lazy = options.lazy;
        }

        if (options.server){
          fetchOptions.server = options.server;
        } 
        const startTime = Date.now();
        
        const res = await $fetch(DEFAULT_CONFIG.SERVER_URL, fetchOptions);

        // 请求成功，清除记录
        this.pendingRequests.delete(requestKey);
        removeExternalAbortListener?.();

        if (res.code !== 200) {
          console.error(`请求失败 [${url}]`, {
            耗时: `${Date.now() - startTime}ms`,
            参数: params,
            响应: res
          });
          resolve(res);
          return;
        }

        resolve(res);
      } catch (err) {
        // 清除记录
        this.pendingRequests.delete(requestKey);
        removeExternalAbortListener?.();
        
        // 取消导致的异常：ofetch 会把中断包装成 FetchError，真正的 AbortError 在 cause 里
        // 以内部 controller 的状态为准最可靠（外部 signal / 重复请求取消都会 abort 它）
        const isAborted =
          controller.signal.aborted ||
          err?.name === 'AbortError' ||
          err?.cause?.name === 'AbortError';
        if (isAborted) {
          console.log(`请求已取消: ${url}`);
          resolve(null);
          return;
        }
        
        console.error(`请求异常 [${url}]`, err);
        resolve({
          code: err.code || -1,
          status: err.status,
          message: err.message || '请求失败',
        });
      }
    });
  }

  // 取消所有请求
  cancelAll() {
    this.pendingRequests.forEach((controller) => {
      controller.abort();
    });
    this.pendingRequests.clear();
  }

  // 取消指定请求
  cancel(url, params) {
    const key = this._getRequestKey(url, params);
    if (this.pendingRequests.has(key)) {
      this.pendingRequests.get(key).abort();
      this.pendingRequests.delete(key);
    }
  }

  post(options) {
    return this.request({ method: 'POST', ...options });
  }

  get(options) {
    return this.request({ method: 'GET', ...options });
  }
}

const httpRequest = new HttpRequest();
export default httpRequest;