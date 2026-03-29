import CryptoJS from "crypto-js";
import { error } from "node:console";

/**
 * 加密方法
 * @param data 需要加密的数据 (字符串或对象)
 * @returns 加密后的 Base64 字符串
 */
export const encryptParam = (data: any): string => {
  const config = useRuntimeConfig();
  const secretKey = config.public.cryptoSecret;
  if (!secretKey) {
    throw error("未知参数:secretKey=" + secretKey);
  }
  // 如果是对象，先转换为字符串
  const text = typeof data === "object" ? JSON.stringify(data) : String(data);

  // 使用 AES 加密
  return CryptoJS.AES.encrypt(text, secretKey).toString();
};

/**
 * 解密方法
 * @param ciphertext 加密后的 Base64 字符串
 * @returns 解密后的原始数据 (如果原先是 JSON 字符串，会自动解析为对象)
 */
export const decryptParam = (ciphertext: string): any => {
  const config = useRuntimeConfig();
  const secretKey = config.public.cryptoSecret;
  if (!secretKey) {
    throw error("未知参数:secretKey=" + secretKey);
  }
  try {
    // 使用 AES 解密
    const bytes = CryptoJS.AES.decrypt(ciphertext, secretKey);
    const originalText = bytes.toString(CryptoJS.enc.Utf8);

    // 尝试解析为 JSON 对象
    try {
      return JSON.parse(originalText);
    } catch (e) {
      return originalText;
    }
  } catch (error) {
    console.error("解密失败:", error);
    return null;
  }
};
