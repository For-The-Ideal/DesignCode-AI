import { commonApi } from '~/api/common'
import { fileToBase64 } from '~/utils/index.js'
import { createCancelScope, cancelScope } from '~/utils/cancel'
import { ElMessage } from 'element-plus'

// 附图素材约束：硬上限 10 张，单张 2MB
export const MAX_ASSETS = 10
export const ASSET_MAX_BYTES = 2 * 1024 * 1024

// 图片判定：Windows 下部分文件（截图/无扩展名/HEIC）的 MIME 为空，按扩展名兜底，避免多选时被静默丢弃
const IMAGE_EXT_RE = /\.(png|jpe?g|gif|webp|svg|bmp|avif|heic|heif|ico)$/i
export const isImageFile = (f) => !!f && (f.type?.startsWith('image/') || IMAGE_EXT_RE.test(f.name || ''))

export const createUploadActions = (ctx) => {
  const { state, images, assets, pageLinks, isConcurrencyFull } = ctx

  // 取消作用域前缀：所有上传请求统一注册为 upload:<uploadId>
  const scopeKey = (uploadId) => `upload:${uploadId}`

  // 跨桶定位（单图/多图数据独立，切换模式后回调仍能找到原条目）；图片已被移除则返回 undefined
  const findItem = (uploadId) => {
    for (const bucket of Object.values(state.modes)) {
      const it = bucket.images.find(i => i.uploadId === uploadId)
      if (it) return it
    }
    return undefined
  }

  /** 取消单张在途上传 */
  const cancelUpload = (uploadId) => cancelScope(scopeKey(uploadId))

  /** File → { preview, width, height }；读取/解码失败返回 null，绝不挂起（否则后续文件全部无法入列） */
  const fileToPreview = (file) => {
    return new Promise((resolve) => {
      const reader = new FileReader()
      reader.onerror = () => resolve(null)
      reader.onload = (e) => {
        const dataUrl = e.target.result
        const img = new Image()
        // 解码失败也保留 dataURL 预览（部分格式浏览器可显示但取不到尺寸）
        img.onload = () => resolve({ preview: dataUrl, width: img.naturalWidth, height: img.naturalHeight })
        img.onerror = () => resolve({ preview: dataUrl, width: 0, height: 0 })
        img.src = dataUrl
      }
      reader.readAsDataURL(file)
    })
  }

  /** 分配页面 ID：p1/p2...，按已有最大序号 +1（删除不回收，保证页间引用稳定） */
  const nextPageID = (list) => {
    const max = list.reduce((m, it) => {
      const n = parseInt(String(it.pageId || '').replace(/^p/, ''), 10)
      return Number.isFinite(n) && n > m ? n : m
    }, 0)
    return `p${max + 1}`
  }

  /** 添加图片到列表（含校验） */
  const addImages = async (files) => {
    // 并发满时不允许添加
    if (isConcurrencyFull.value) return

    const imageFiles = files.filter(isImageFile)
    if (imageFiles.length === 0) return

    // 上限检查
    if (images.value.length >= state.maxImages) {
      ElMessage.warning(`最多上传 ${state.maxImages} 张图片`)
      return
    }

    const remaining = Math.min(state.maxImages - images.value.length, imageFiles.length)
    const toAdd = imageFiles.slice(0, remaining)

    for (const file of toAdd) {
      const meta = await fileToPreview(file)
      if (!meta) {
        ElMessage.warning(`「${file.name}」读取失败，已跳过`)
        continue
      }
      const { preview, width, height } = meta
      const uploadId = `u${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
      const scope = createCancelScope(scopeKey(uploadId))
      images.value.push({
        uploadId, // 内部上传标识，仅用于异步回调定位；也是附图归属的 stable 外键
        file,
        preview,
        naturalWidth: width,
        naturalHeight: height,
        cosUrl: '',
        uploading: true,
        uploadError: '',
        // 页面说明（多图模式即 page_desc：用户填写的页面级 brief，随任务提交给生成模型）
        description: '',
        // 标注框列表: [{ id, x, y, w, h(0~1 相对比例), label, desc }]
        annotations: [],
        // ── 页面身份（多图模式下每张图是一页）──
        pageId: nextPageID(images.value),   // 稳定标识，页面关系引用它
        pageName: '',                        // 页面名（留空由后端按文件名/序号兜底）
        route: '',                           // 路由（留空由后端兜底：入口页 /，其余 /pN）
        // 入口页不落字段：恒为页面列表第 1 张，排序/增删后跟着列表走
        // ── 页面关系（UI 友好字段，提交时转成 page_links 边表）──
        parentPageId: '',                    // 从哪一页进入（空=独立页）
        linkTrigger: 'node_click',           // 进入方式：node_click | bottom_tab | auto
        // ── 多图画板坐标（仅前端排版；null = 尚未拖入画布）──
        board: null,
      })
      // 传 uploadId 而非索引：上传期间删除/清空/重排不会导致错位
      uploadOne(uploadId, file, scope)
    }

    // 主图变更后同步素材归属，避免悬空：
    // - 单图：唯一主图，全部素材重绑到新主图（更换/重传后素材跟随）
    // - 多图：原归属页已被删除的悬空素材自动挂到新上传的第一张
    if (images.value[0]?.uploadId) {
      const first = images.value[0].uploadId
      if (state.mode === 'single') {
        assets.value.forEach((a) => { a.ownerId = first })
      } else {
        assets.value.forEach((a) => {
          if (!images.value.some((i) => i.uploadId === a.ownerId)) a.ownerId = first
        })
      }
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

  /** 移除单张主图：其名下附图改绑到剩余第一张，避免产生孤儿素材 */
  const removeImage = (idx) => {
    const item = images.value[idx]
    if (!item) return
    // 先取消在途上传，再移除
    if (item.uploadId) cancelUpload(item.uploadId)
    images.value.splice(idx, 1)
    // 删除的是当前编辑页：当前页落到相邻页（优先同位置，其次上一页）
    if (state.activePageId === item.pageId) {
      state.activePageId = images.value[Math.min(idx, images.value.length - 1)]?.pageId || ''
    }
    const fallback = images.value[0]?.uploadId || ''
    assets.value.forEach((a) => { if (a.ownerId === item.uploadId) a.ownerId = fallback })
    // 页面关系清理：删除引用被删页的所有边，避免提交悬空边
    const links = pageLinks.value
    for (let i = links.length - 1; i >= 0; i--) {
      if (links[i].fromPageId === item.pageId || links[i].toPageId === item.pageId) links.splice(i, 1)
    }
    // 旧字段同步置空
    images.value.forEach((it) => { if (it.parentPageId === item.pageId) it.parentPageId = '' })
  }

  // ═══ 附图素材（与主图同一上传模块，作用域前缀 asset: 隔离） ═══

  const assetScopeKey = (assetId) => `asset:${assetId}`
  const findAsset = (assetId) => {
    for (const bucket of Object.values(state.modes)) {
      const it = bucket.assets.find(a => a.assetId === assetId)
      if (it) return it
    }
    return undefined
  }

  /** 取消单张在途附图上传 */
  const cancelAsset = (assetId) => cancelScope(assetScopeKey(assetId))

  /** 上传单张附图素材 */
  const uploadAsset = async (assetId, file, scope) => {
    const isCancelled = () => scope.aborted
    try {
      const base64 = await fileToBase64(file)
      if (isCancelled()) return
      const res = await commonApi.uploadImage(base64.split(',')[1], file.name, { signal: scope.signal })
      if (isCancelled()) return
      if (res && res.code === 200 && res.data?.url) {
        const asset = findAsset(assetId)
        if (asset) {
          asset.cosUrl = res.data.url
          asset.uploading = false
        }
      } else {
        throw new Error(res?.message || '上传失败')
      }
    } catch (e) {
      if (isCancelled()) return
      const asset = findAsset(assetId)
      if (asset) {
        asset.uploading = false
        asset.uploadError = e.message || '上传失败'
      }
    } finally {
      scope.dispose()
    }
  }

  /** 添加附图：过滤非图 / 超 2MB / 超出上限。ownerId = 所属主图的 uploadId */
  const addAssets = async (files, ownerId = '') => {
    const imageFiles = (files || []).filter(isImageFile)
    if (!imageFiles.length) return

    imageFiles.filter(f => f.size > ASSET_MAX_BYTES).forEach(f =>
      ElMessage.warning(`「${f.name}」超过 ${Math.round(ASSET_MAX_BYTES / 1024 / 1024)}MB，已跳过`),
    )
    const okFiles = imageFiles.filter(f => f.size <= ASSET_MAX_BYTES)

    const room = MAX_ASSETS - assets.value.length
    if (room <= 0) {
      ElMessage.warning(`附图最多 ${MAX_ASSETS} 张，请先删除不需要的素材`)
      return
    }
    if (okFiles.length > room) ElMessage.warning(`附图最多 ${MAX_ASSETS} 张，超出部分未添加`)

    for (const file of okFiles.slice(0, room)) {
      const meta = await fileToPreview(file)
      if (!meta) {
        ElMessage.warning(`「${file.name}」读取失败，已跳过`)
        continue
      }
      const { preview, width, height } = meta
      const assetId = `img${Date.now()}_${Math.random().toString(36).slice(2, 8)}`
      const scope = createCancelScope(assetScopeKey(assetId))
      assets.value.push({
        assetId, // 内部标识，标注框通过它引用素材
        ownerId: ownerId || images.value[0]?.uploadId || '', // 归属主图（多图模式必填，单图自动唯一）
        file,
        preview,
        naturalWidth: width,
        naturalHeight: height,
        label: file.name.replace(/\.[^.]+$/, '') || '未命名',
        size: file.size,
        mime: file.type,
        cosUrl: '',
        uploading: true,
        uploadError: '',
      })
      uploadAsset(assetId, file, scope)
    }
  }

  /** 重试单张附图上传 */
  const retryAsset = (assetId) => {
    const asset = findAsset(assetId)
    if (!asset || asset.uploading || !asset.file) return
    const scope = createCancelScope(assetScopeKey(assetId))
    asset.uploading = true
    asset.uploadError = ''
    uploadAsset(assetId, asset.file, scope)
  }

  /** 移除单张附图（在途上传先取消） */
  const removeAsset = (assetId) => {
    cancelScope(assetScopeKey(assetId))
    for (const bucket of Object.values(state.modes)) {
      const idx = bucket.assets.findIndex(a => a.assetId === assetId)
      if (idx > -1) {
        bucket.assets.splice(idx, 1)
        return
      }
    }
  }

  /** 清空指定主图名下的全部附图（取消在途上传；单图更换/删除主图时使用） */
  const clearAssetsOf = (ownerId) => {
    if (!ownerId) return
    for (const bucket of Object.values(state.modes)) {
      bucket.assets.filter(a => a.ownerId === ownerId).forEach(a => cancelScope(assetScopeKey(a.assetId)))
      bucket.assets = bucket.assets.filter(a => a.ownerId !== ownerId)
    }
  }

  /** 改素材标签（画布引用与 Prompt 清单都用它） */
  const updateAssetLabel = (assetId, label) => {
    const asset = findAsset(assetId)
    if (asset && String(label || '').trim()) asset.label = String(label).trim().slice(0, 30)
  }

  /** 改素材归属主图（多图模式下切页/拖拽重排后仍以 uploadId 为准） */
  const setAssetOwner = (assetId, ownerId) => {
    const asset = findAsset(assetId)
    if (asset) asset.ownerId = ownerId || ''
  }

  return {
    addImages, removeImage, retryUpload, cancelUpload,
    addAssets, retryAsset, removeAsset, clearAssetsOf, cancelAsset, updateAssetLabel, setAssetOwner,
  }
}
