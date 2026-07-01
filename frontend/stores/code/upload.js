import { commonApi } from '~/api/common'
import { fileToBase64 } from '~/utils/index.js'
import { ElMessage } from 'element-plus'

export const createUploadActions = (ctx) => {
  const { state, isConcurrencyFull } = ctx

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
      const idx = state.images.length
      state.images.push({
        file,
        preview,
        naturalWidth: width,
        naturalHeight: height,
        cosUrl: '',
        uploading: true,
        uploadError: '',
        description: '',
      })
      uploadOne(idx, file)
    }
  }

  /** 上传单张到 COS */
  const uploadOne = async (idx, file) => {
    try {
      const base64 = await fileToBase64(file)
      const base64Data = base64.split(',')[1]
      const res = await commonApi.uploadImage(base64Data, file.name)
      if (res.code === 200 && res.data?.url) {
        state.images[idx].cosUrl = res.data.url
        state.images[idx].uploading = false
      } else {
        throw new Error(res.message || '上传失败')
      }
    } catch (e) {
      state.images[idx].uploading = false
      state.images[idx].uploadError = e.message || '上传失败'
      ElMessage.error(`"${file.name}" 上传失败: ${e.message}`)
    }
  }

  /** 移除单张 */
  const removeImage = (idx) => {
    state.images.splice(idx, 1)
  }

  return { addImages, removeImage }
}
