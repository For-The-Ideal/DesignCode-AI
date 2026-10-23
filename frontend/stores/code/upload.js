import { commonApi } from '~/api/common'
import { fileToBase64 } from '~/utils/index.js'
import { createCancelScope, cancelScope, cancelScopes } from '~/utils/cancel'
import { ElMessage } from 'element-plus'

export const createUploadActions = (ctx) => {
  const { state, isConcurrencyFull } = ctx

  // 取消作用域前缀：所有上传请求统一注册为 upload:<uploadId>
  const scopeKey = (uploadId) => `upload:${uploadId}`

  // 每次回调时重新定位，确保拿到的是响应式代理；图片已被移除则返回 undefined
  const findItem = (uploadId) => state.images.find(i => i.uploadId === uploadId)

  /** 取消单张在途上传 */
  const cancelUpload = (uploadId) => cancelScope(scopeKey(uploadId))

  /** 取消全部在途上传（清空时使用） */
  const cancelAllUploads = () => cancelScopes('upload:')

  /** File → { preview, width, height } */
  const fileToPreview = (file) => {
    return new Promise((resolve) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        const dataUrl = e.target.result
        const img = new Image()
        img.onload = () => {
          resolve({ preview: dataUrl, width: img.naturalWidth, height: img.naturalHeight })
          URL.revokeObjectURL(img.src)
        }
        img.src = dataUrl
      }
      reader.readAsDataURL(file)
    })
  }

  /** 添加图片到列表（含校验） */
  const addImages = async (files) => {
    // 并发满时不允许添加
    if (isConcurrencyFull.value) return

    const imageFiles = files.filter(f => f.type.startsWith('image/'))
    if (imageFiles.length === 0) return

    // 上限检查
    if (state.images.length >= state.maxImages) {
      ElMessage.warning(`最多上传 ${state.maxImages} 张图片`)
      return
    }

    const remaining = Math.min(state.maxImages - state.images.length, imageFiles.length)
    const toAdd = imageFiles.slice(0, remaining)

    for (const file of toAdd) {
      const { preview, width, height } = await fileToPreview(file)
      const uploadId = `u${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
      const scope = createCancelScope(scopeKey(uploadId))
      state.images.push({
        uploadId, // 内部上传标识，仅用于异步回调定位
        file,
        preview,
        naturalWidth: width,
        naturalHeight: height,
        cosUrl: '',
        uploading: true,
        uploadError: '',
        description: '',
        // 标注框列表: [{ id, x, y, w, h(0~1 相对比例), label, desc }]
        annotations: [],
      })
      // 传 uploadId 而非索引：上传期间删除/清空/重排不会导致错位
      uploadOne(uploadId, file, scope)
    }
  }

  /** 上传单张到 COS */
  const uploadOne = async (uploadId, file, scope) => {
    const isCancelled = () => scope.aborted

    try {
      const base64 = await fileToBase64(file)
      // 本地读文件期间已被取消，不再发起网络请求
      if (isCancelled()) return
      const base64Data = base64.split(',')[1]
      // signal 由统一取消作用域管理，取消时底层 fetch 立即 abort
      const res = await commonApi.uploadImage(base64Data, file.name, { signal: scope.signal })
      if (isCancelled()) return
      if (res && res.code === 200 && res.data?.url) {
        const item = findItem(uploadId)
        if (item) {
          item.cosUrl = res.data.url
          item.uploading = false
        }
      } else {
        throw new Error(res?.message || '上传失败')
      }
    } catch (e) {
      // 主动取消不视为失败，不弹错误也不写状态
      if (isCancelled()) return
      const item = findItem(uploadId)
      // 图片已被移除/清空时不再打扰用户
      if (item) {
        item.uploading = false
        item.uploadError = e.message || '上传失败'
        ElMessage.error(`"${file.name}" 上传失败: ${e.message}`)
      }
    } finally {
      // 请求结束（成功/失败/取消）后从作用域注册表移除
      scope.dispose()
    }
  }

  /** 重试上传（失败后手动触发） */
  const retryUpload = (uploadId) => {
    const item = findItem(uploadId)
    if (!item || item.uploading) return
    // createCancelScope 内部会自动废弃同 key 的旧作用域
    const scope = createCancelScope(scopeKey(uploadId))
    item.uploading = true
    item.uploadError = ''
    uploadOne(uploadId, item.file, scope)
  }

  /** 移除单张 */
  const removeImage = (idx) => {
    const item = state.images[idx]
    // 先取消在途上传，再移除
    if (item?.uploadId) cancelUpload(item.uploadId)
    state.images.splice(idx, 1)
  }

  return { addImages, removeImage, retryUpload, cancelUpload, cancelAllUploads }
}
