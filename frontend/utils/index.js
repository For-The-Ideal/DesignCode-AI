import CryptoJS from "crypto-js";
import { cryptoConfig } from "~/config/index.js";

// 创建原生通知的函数
export const showElectronNotification = (options) => {
  // 创建图标
  let icon = null;
  if (options.iconPath) {
    // icon = nativeImage.createFromPath(path.join(__dirname, options.iconPath))
    icon = "";
  }

  // 创建通知对象
  const notification = new Notification({
    title: options.title || "通知",
    body: options.body || "",
    icon: icon,
    silent: options.silent || false,
    timeoutType: options.persistent ? "never" : "default",
    urgency: options.urgent ? "critical" : "normal",
    actions: options.actions || [],
    closeButtonText: options.closeButtonText || "关闭",
    ...options
  });

  // 显示通知
  notification.show();

  // 事件监听
  notification.on("click", (event) => {
    console.log("通知被点击");
    if (options.onClick) options.onClick(event);

    // 聚焦所有窗口
    const windows = BrowserWindow.getAllWindows();
    windows.forEach((win) => {
      if (win.isMinimized()) win.restore();
      win.focus();
    });
  });

  notification.on("close", (event) => {
    console.log("通知已关闭");
    if (options.onClose) options.onClose(event);
  });

  notification.on("action", (event, index) => {
    console.log(`操作按钮被点击: ${index}`);
    if (options.actions && options.actions[index]) {
      const action = options.actions[index];
      if (action.click) action.click();
    }
  });

  return notification;
};

/**
 * 加密数据（同步版本）
 */
export const setEncrypt = (data) => {
  if (!data) {
    console.warn("⚠️  加密数据为空")
    return null
  }
  try {
    const dataStr = typeof data === "string" ? data : JSON.stringify(data)
    const encrypted = CryptoJS.AES.encrypt(
      dataStr,
      CryptoJS.enc.Utf8.parse(cryptoConfig.key),
      {
        iv: CryptoJS.enc.Utf8.parse(cryptoConfig.iv),
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      }
    )
    const result = encrypted.toString()
    return result
  } catch (error) {
    console.error("❌ 加密错误:", error)
    return null
  }
}

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
