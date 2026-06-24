
import CryptoJS from "crypto-js";

/**
 * 解密数据
 */
export const getDecrypt = (encryptedStr) => {
  if (!encryptedStr) {
    console.warn("⚠️  解密数据为空")
    return null
  }
  try {
    const runConfig = useRuntimeConfig()
    const decrypted = CryptoJS.AES.decrypt(
      encryptedStr,
      CryptoJS.enc.Utf8.parse(runConfig.public.cryptoKey),
      {
        iv: CryptoJS.enc.Utf8.parse(runConfig.public.cryptoIv),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      }
    )
    const resultUtf8 = decrypted.toString(CryptoJS.enc.Utf8)
    try {
      const parsed = JSON.parse(resultUtf8)
      return parsed
    } catch (jsonError) {
      console.log("解密失败，解析为JSON:", jsonError)
      return resultUtf8
    }
  } catch (error) {
    console.log("解密失败，解析为JSON:", error)
    return null
  }
}
