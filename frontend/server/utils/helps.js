
import CryptoJS from "crypto-js";
import { cryptoConfig } from "~/config/index.js";

/**
 * 解密数据
 */
export const getDecrypt = (encryptedStr) => {
  if (!encryptedStr) {
    console.warn("⚠️  解密数据为空")
    return null
  }
  try {
    const decrypted = CryptoJS.AES.decrypt(
      encryptedStr,
      CryptoJS.enc.Utf8.parse(cryptoConfig.key),
      {
        iv: CryptoJS.enc.Utf8.parse(cryptoConfig.iv),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      }
    )
    const resultUtf8 = decrypted.toString(CryptoJS.enc.Utf8)
    try {
      const parsed = JSON.parse(resultUtf8)
      console.log("✅ 解密成功，解析为JSON:", parsed)
      return parsed
    } catch (jsonError) {
      return resultUtf8
    }
  } catch (error) {
    return null
  }
}
