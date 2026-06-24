import CryptoJS from "crypto-js";
import { ElMessage } from 'element-plus'
import {useEnv} from "@/composables/useEnv"

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
 * 加密数据
 */
export const setEncrypt = (data) => {
  if (!data) {
    console.warn("⚠️  加密数据为空")
    return null
  }
  try {
    const {cryptoKey,cryptoIv} = useEnv()
    const dataStr = typeof data === "string" ? data : JSON.stringify(data)
    const encrypted = CryptoJS.AES.encrypt(
      dataStr,
      CryptoJS.enc.Utf8.parse(cryptoKey),
      {
        iv: CryptoJS.enc.Utf8.parse(cryptoIv),
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
 * 复制文本到剪贴板
 * @param {string} text 要复制的文本
 */
export const handleCopy = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch {
    ElMessage.error('复制失败')
  }
}

/**
 * 格式化代码（缩进对齐）
 * @param {string} code 原始代码
 * @returns {string} 格式化后的代码
 */
export const handleFormat = (code) => {
  let indent = 0
  const formatted = code.split('\n').map(line => {
    const t = line.trim()
    if (!t) return ''
    if (/^[})]/.test(t)) indent = Math.max(0, indent - 1)
    const r = '  '.repeat(indent) + t
    if (/[({]$/.test(t)) indent++
    return r
  }).join('\n')
  ElMessage.success('代码已格式化')
  return formatted
}

/**
 * 下载代码文件
 * @param {string} text 代码文本
 * @param {string} language 语言标识（dart/typescript/html）
 */
export const handleDownload = (text, language) => {
  const ext = language === 'dart' ? 'dart' : language === 'typescript' ? 'tsx' : 'vue'
  const blob = new Blob([text], { type: 'text/plain' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `generated.${ext}`
  a.click()
  URL.revokeObjectURL(a.href)
  ElMessage.success('已下载')
}


// 读取文件为 base64 data URL
export function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}
