import { readBody, getRequestHeaders } from "h3"
import { getData } from "../utils/request.js"
import { getDecrypt } from "../utils/helps.js"
import { serverLog } from "../utils/logger.js"

export default defineEventHandler(async (event) => {
  const startTime = Date.now()
  const body = await readBody(event)
  if (!body) {
    return { code: 400, message: "参数错误!", data: {} }
  }
  const req = getRequestHeaders(event)
  const isBrowser = !!req.client

  // client 直发：body 就是 { url, method, params }
  // SSR 请求：body = { aes: 加密串 }，需解密
  const { url = "", method = "post", params = {} } = isBrowser
    ? body
    : await getDecrypt(body.aes)

  // 记录解密后的请求参数
  serverLog('REQ', { method, url, params, isBrowser })

  const result = await getData({ method, url, params, req })
  const duration = Date.now() - startTime

  // 记录响应——正常 / 错误 区分
  if (result.code === 200) {
    serverLog('RES', { method, url, status: 200, duration })
  } else {
    serverLog('ERR', { method, url, status: result.code || 500, message: result.message || '未知错误', duration })
  }

  return result
})
