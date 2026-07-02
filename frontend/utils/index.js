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


/**
 * 验证邮箱格式
 * @param {string} email 邮箱地址
 * @returns {{ valid: boolean, message?: string }} 验证结果
 */
export const validateEmail = (email) => {
  if (!email) return { valid: false, message: '请输入邮箱地址' }
  if (email.length > 254) return { valid: false, message: '邮箱地址过长' }
  
  // RFC 5322 简化版邮箱正则
  const reg = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$/
  if (!reg.test(email)) return { valid: false, message: '邮箱格式不正确' }
  return { valid: true }
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



// ------------------------------------------ 节流 和 防抖 ----------------------------------------
/**
 * 节流。适用于速率限制
 * 对 resize 和 scroll 等事件处理。
 *
 * @param {number} delay -                   延迟零或更大，单位毫秒。一般在100或250左右。
 * @param {Function} callback -              回调，每个 delay 后执行 callback。
 * @param {object} [options] -               配置选项。
 * @param {boolean} [options.noTrailing] -   不尾随执行 callback，可选，默认 false。
 *                                           noTrailing 为 true，最后一次 delay 内，不执行 callback。
 *                                           noTrailing 为 false，最后一次 delay 内，执行 callback。
 *                                           （在 delay 内未调用 throttled-function 后，内部计数器将重置）。
 * @param {boolean} [options.noLeading] -    首次不执行 callback，可选，默认为 false。
 *                                           noLeading 为 false，在开始时，首次执行 callback。
 *                                           noLeading 为 true，在开始时，不执行 callback。
 *                                           注意: noLeading = true 和 noTrailing = true，则 callback 永不执行。
 * @param {boolean} [options.debounceMode] - debounceMode 为 true，最后一次 delay 后执行 clear，不执行 callback。
 *                                           debounceMode 为 false，最后一次 delay 后执行 callback。
 *
 * @returns {Function}   返回 throttle 函数本身。
 */
export function throttle(
  delay, callback, options
) {
  const {
    noTrailing = false,
    noLeading = false,
    debounceMode = undefined
  } = options || {};

  let timeoutID;
  let cancelled = false;
  let lastExec = 0;

  function clearExistingTimeout() {
    if (timeoutID) {
      clearTimeout(timeoutID);
    }
  }

  function cancel(options) {
    const { upcomingOnly = false } = options || {};
    clearExistingTimeout();
    cancelled = !upcomingOnly;
  }

  function wrapper(...arguments_) {
    let self = this;
    let elapsed = Date.now() - lastExec;

    if (cancelled) {
      return;
    }

    function exec() {
      lastExec = Date.now();
      callback.apply(self, arguments_);
    }

    function clear() {
      timeoutID = undefined;
    }

    if (!noLeading && debounceMode && !timeoutID) {
      exec();
    }

    clearExistingTimeout();

    if (debounceMode === undefined && elapsed > delay) {
      if (noLeading) {
        lastExec = Date.now();
        if (!noTrailing) {
          timeoutID = setTimeout(debounceMode ? clear : exec, delay);
        }
      } else {
        exec();
      }
    } else if (noTrailing !== true) {
      timeoutID = setTimeout(
        debounceMode ? clear : exec,
        debounceMode === undefined ? delay - elapsed : delay
      );
    }
  }
  wrapper.cancel = cancel;
  return wrapper;
}

/**
 * 防抖
 * 保证 callback 在 delay 内只执行一次，在逻辑开头调用，或者结束时调用。
 *
 * @param {number} delay -               延迟零或更大，单位毫秒。一般在100或250左右。
 * @param {Function} callback -          回调，在每个 delay 后执行 callback。
 * @param {object} [options] -           配置选项。
 * @param {boolean} [options.atBegin] -  首次执行 callback，可选，默认 false。
 *                                       atBegin 为 false，在开始时，不执行 callback。
 *                                       atBegin 为 true，在开始时，首次执行 callback。
 *                                       （在 delay 内未调用 throttled-function，内部计数器将重置）。
 *
 * @returns {Function}   返回 debounced 函数本身。
 */
export function debounce(
  delay, callback, options
) {
  const { atBegin = false } = options || {};
  return throttle(delay, callback, { debounceMode: atBegin !== false });
}