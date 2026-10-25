<template>
  <!-- 上传拖拽区：仅在未达上限时展示大块；单图传满后自动收起，避免空壳 -->
  <div
    v-if="canAddMore"
    class="upload-zone"
    :class="{ 'upload-zone--disabled': disabled }"
    @click="onZoneClick"
    @dragover.prevent
    @drop.prevent="onZoneDrop"
  >
    <i class="fa-solid fa-cloud-arrow-up upload-icon"></i>
    <div class="upload-text">
      <span class="upload-text-bold">{{ uploadHint.title }}</span>
    </div>
    <div class="upload-hint">
      {{ uploadHint.desc }}
    </div>
  </div>

  <!-- 隐藏文件选择 -->
  <input type="file" ref="fileInput" multiple accept="image/*" style="display:none" @change="onFileChange" />

  <!-- 预览区：单图=左图右信息两栏；多图=网格 + 拖拽 -->
  <div v-if="images.length" class="preview-area" :class="{ 'preview-area--locked': isSubmitting }">
    <!-- ── 单图 hero：左侧手机壳包裹设计稿 + 右侧信息/操作 ── -->
    <div v-if="images.length === 1" class="hero-panel">
      <!-- 左：复用 mobile.vue 手机壳（裸机模式），屏内用插槽传入设计稿原图 -->
      <div class="hero-figure" @click="openAnno(0)">
        <MobilePreview bare :screen-height="'auto'">
          <div class="img-frame" :style="{ aspectRatio: frameRatio(images[0]) }">
            <img :src="images[0].preview" alt="设计稿" draggable="false" />
            <span
              v-for="anno in images[0].annotations || []"
              :key="anno.id"
              class="anno-box"
              :class="{ 'anno-box--child': anno.top === false }"
              :style="annoStyle(anno)"
            >
              <span class="anno-num">{{ anno.label }}</span>
            </span>
          </div>

          <div v-if="images[0].uploading" class="upload-status uploading">
            <i class="fas fa-spinner fa-spin"></i>
            <span class="uploading-text">上传中…</span>
          </div>
          <div v-else-if="images[0].uploadError" class="upload-status failed">
            <i class="fas fa-circle-exclamation"></i>
            <span class="failed-text">上传失败</span>
          </div>
        </MobilePreview>
      </div>

      <!-- 右：文件信息卡 + 标注面板卡 + 底部操作 -->
      <aside class="hero-aside">
        <!-- 文件信息 -->
        <div class="meta-card">
          <div class="meta-icon"><i class="fa-regular fa-image"></i></div>
          <div class="meta-main">
            <div class="meta-name" :title="images[0].file?.name">{{ images[0].file?.name || '设计稿' }}</div>
            <div class="meta-dim">{{ images[0].naturalWidth || '--' }} × {{ images[0].naturalHeight || '--' }}</div>
          </div>
          <span class="status-pill" :class="imgState(images[0])">
            <i :class="imgStateIcon[imgState(images[0])]"></i>{{ imgStateText[imgState(images[0])] }}
          </span>
        </div>

        <!-- 区域标注 -->
        <div class="anno-card">
          <div class="anno-card-head">
            <span class="anno-card-title">区域标注</span>
            <span v-if="images[0].annotations?.length" class="anno-card-count"><b>{{ annoTopCount(images[0]) }}</b> 个</span>
          </div>

          <div v-if="images[0].uploading" class="anno-empty">
            <i class="fas fa-spinner fa-spin"></i><span>上传中，完成后可标注</span>
          </div>
          <div v-else-if="images[0].uploadError" class="anno-empty error">
            <i class="fas fa-circle-exclamation"></i><span>上传失败，暂不可标注</span>
          </div>
          <div v-else-if="images[0].annotations?.length" class="anno-list">
            <div class="anno-chips">
              <span v-for="a in topAnnos(images[0])" :key="a.id" class="anno-chip"><em>{{ a.label }}</em></span>
            </div>
            <button class="anno-edit" @click="openAnno(0)"><i class="fa-regular fa-pen-to-square"></i> 编辑标注区域</button>
          </div>
          <button v-else class="anno-add" @click="openAnno(0)">
            <i class="fa-regular fa-square-dashed"></i>
            <span>点击框选设计稿区域</span>
            <small>描述越细，AI 理解越准</small>
          </button>
        </div>

        <!-- 底部操作 -->
        <div class="hero-foot">
          <button class="foot-btn" @click="onReplace"><i class="fas fa-rotate"></i> 更换</button>
          <button class="foot-btn danger" @click="store.clearAll()"><i class="fa-regular fa-trash-can"></i> 删除</button>
        </div>
      </aside>
    </div>

    <!-- ── 多图：网格 + 拖拽 ── -->
    <template v-else>
      <div class="preview-grid" v-sortable="{ list: images, disabled: dragDisabled }">
        <div
          v-for="(img, i) in images"
          :key="i"
          class="preview-card"
          :class="{ 'preview-card--disabled': isNotReady(img) }"
          data-sortable-item
          :data-idx="i"
          @click="openAnno(i)"
        >
          <!-- ── 图区 ── -->
          <div class="card-media">
            <div class="media-scroll">
              <div class="img-frame" :style="{ aspectRatio: frameRatio(img) }">
                <img :src="img.preview" alt="设计稿" draggable="false" />
                <span
                  v-for="anno in img.annotations || []"
                  :key="anno.id"
                  class="anno-box"
                  :class="{ 'anno-box--child': anno.top === false }"
                  :style="annoStyle(anno)"
                >
                  <span class="anno-num">{{ anno.label }}</span>
                </span>
              </div>
            </div>

            <!-- 上传状态 -->
            <div v-if="img.uploading" class="upload-status uploading">
              <i class="fas fa-spinner fa-spin"></i>
              <span class="uploading-text">上传中…</span>
            </div>
            <div v-else-if="img.uploadError" class="upload-status failed">
              <i class="fas fa-circle-exclamation"></i>
              <span class="failed-text">上传失败</span>
              <button class="retry-btn" @click.stop="onRetry(img)">
                <i class="fas fa-rotate-right"></i> 重试
              </button>
            </div>

            <!-- 悬浮工具条 -->
            <div class="card-toolbar">
              <span class="card-index">{{ String(i + 1).padStart(2, '0') }}</span>
              <span class="toolbar-spacer"></span>
              <button
                class="tool-icon"
                :disabled="isNotReady(img)"
                :title="isNotReady(img) ? '图片未上传完成，暂不可标注' : '标注'"
                @click.stop="openAnno(i)"
              >
                <i class="fa-regular fa-pen-to-square"></i>
              </button>
              <button class="tool-icon tool-icon--danger" title="删除" @click.stop="store.removeImage(i)">
                <i class="fa-regular fa-trash-can"></i>
              </button>
            </div>
          </div>

          <!-- ── 信息区 ── -->
          <div class="card-info">
            <div class="info-line">
              <span class="file-name">{{ img.file?.name || '设计稿' }}</span>
              <span class="file-size">{{ img.naturalWidth || '--' }} x {{ img.naturalHeight || '--' }}</span>
            </div>
            <div v-if="img.uploading" class="anno-summary empty">
              <i class="fas fa-spinner fa-spin"></i>
              上传中，完成后可标注
            </div>
            <div v-else-if="img.uploadError" class="anno-summary error">
              <i class="fas fa-exclamation-circle"></i>
              上传失败，暂不可标注
            </div>
            <div v-else-if="img.annotations?.length" class="anno-summary">
              <i class="fa-solid fa-location-target"></i>
              {{ annoTopCount(img) }} 个顶层区域 · {{ annoLabels(img) }}
            </div>
            <div v-else class="anno-summary empty">
              <i class="fa-regular fa-square-dashed"></i>
              未标注，点击进入标注编辑器
            </div>
          </div>
        </div>
      </div>

      <div class="upload-actions">
        <span class="upload-count">{{ images.length }} 张已选择</span>
        <div class="action-btns">
          <button v-if="!canAddMore" class="replace-btn" @click="onReplace">
            <i class="fas fa-rotate"></i> 更换图片
          </button>
          <button class="clear-btn" @click="store.clearAll()">清空全部</button>
        </div>
      </div>

      <div class="upload-tip">
        <i class="fa-regular fa-lightbulb"></i>
        <span class="tip-bold">提示：</span>
        拖拽卡片可调整页面顺序，框选区域并填写描述能帮 AI 更准地理解设计稿
      </div>
    </template>
  </div>

  <!-- 结构编辑器（新：中间画布框选 + 右侧思维导图树 + 描述表单） -->
  <ProjectEditor ref="projectEditorRef" :analyzing="analyzing" @confirm="onEditorConfirm" @close="onEditorClose" />
</template>

<script setup>
import { ref, computed, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import ProjectEditor from '~/components/uiProject/ProjectEditor.client.vue'
import MobilePreview from '~/components/previewTempLate/mobile.vue'
import { useCodeStore } from '~/stores/code'
import { useUiProjectStore } from '~/stores/uiProject'
import { uiProjectApi } from '~/api/uiProject'
import { createCancelScope, cancelScopes } from '~/utils/cancel'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'

const userStore = useUserStore()
const commonStore = useCommonStore()
const { isLogin } = storeToRefs(userStore)
const openLoginModal = () => commonStore.setLoginModalVisible(true)

const store = useCodeStore()
const uiStore = useUiProjectStore()
const { images, isSubmitting, isMaxReached, isConcurrencyFull, uploadHint } = storeToRefs(store)

// ── 从 store 派生 ──
const disabled = computed(() => {
  return isMaxReached.value || isSubmitting.value || isConcurrencyFull.value
})

// 是否还能继续添加图片（未达上限）。单图传满后收起大块拖拽区，避免空白
const canAddMore = computed(() => images.value.length < store.maxImages)

// ── 文件选择 ──
const fileInput = ref(null)

const onZoneClick = () => {
  if (!isLogin.value) return openLoginModal()
  if (disabled.value) return
  fileInput.value?.click()
}

const onFileChange = (e) => {
  const files = Array.from(e.target.files)
  if (files.length) store.addImages(files)
  if (fileInput.value) fileInput.value.value = ''
}

const onZoneDrop = (e) => {
  if (!isLogin.value) return openLoginModal()
  if (disabled.value) return
  const files = Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/'))
  if (files.length) store.addImages(files)
}

// 单图更换：先清空当前（会取消在途上传），再唤起选择框
const onReplace = () => {
  if (isSubmitting.value) return
  store.clearAll()
  fileInput.value?.click()
}

// ── 结构编辑（新树编辑器入口）──
const projectEditorRef = ref(null)
const editingIndex = ref(0)
// 顶层识别中：编辑器打开后后台调 AI 回填 topRegions，用该状态驱动工具条「识别中」提示
const analyzing = ref(false)

// 图片未就绪（上传中 / 失败 / 无 COS 地址）时禁止标注
const isNotReady = (img) => !!img.uploading || !!img.uploadError || !img.cosUrl

// 真实宽高比 -> CSS aspect-ratio：尺寸未就绪时回退 auto（不影响完整显示）
const frameRatio = (img) => {
  const w = img.naturalWidth || img.treeWidth || 0
  const h = img.naturalHeight || img.treeHeight || 0
  return w && h ? `${w} / ${h}` : 'auto'
}

const openAnno = (idx) => {
  // 仅“提交/生成中”锁定编辑；不能用 disabled（含 isMaxReached，单图传满即为 true 会误拦编辑）
  if (isSubmitting.value) return
  const img = images.value[idx]
  if (!img) return
  if (isNotReady(img)) {
    if (img.uploadError) ElMessage.warning('图片上传失败，请先点击重试')
    return
  }
  editingIndex.value = idx
  const url = img.cosUrl
  // 初始化项目 store：优先用 naturalWidth，回退到 img.treeWidth/treeHeight（刷新后 File 丢失时的持久化尺寸）
  const W = img.naturalWidth || img.treeWidth || 0
  const H = img.naturalHeight || img.treeHeight || 0
  if (uiStore.image.url !== url) {
    // 换图：清旧识别结果，并载入该图此前已确认的树（无则清空），避免串图
    uiStore.setTopRegions([])
    uiStore.loadTree(img.annotationTree || [])
  } else if (!uiStore.hasNodes && img.annotationTree?.length) {
    // 同图但 store 空（如刷新后）：恢复已确认的树继续编辑
    uiStore.loadTree(img.annotationTree)
  }
  uiStore.setImage({ url, width: W, height: H })
  projectEditorRef.value?.open()
  // 进入编辑器：顶层描述目录为空时按需调 AI（不阻塞框选，回填后选到同名枚举自动填充）
  ensureTopRegions(url)
}

// 顶层区域识别：幂等（已有目录则跳过），失败/未就绪时降级为手动框选填写
async function ensureTopRegions(imageUrl) {
  if (!imageUrl) return
  if (uiStore.topRegions && uiStore.topRegions.length) return
  // 统一取消作用域：关闭/确认编辑器或组件卸载时可 cancelScopes('region:') 中止未响应的识别请求
  const scope = createCancelScope(`region:${imageUrl}`)
  analyzing.value = true
  try {
    const res = await uiProjectApi.analyzeRegions(imageUrl, { signal: scope.signal })
    // 被取消时 request.js 返回 null：静默返回，不误报“识别暂不可用”
    if (!res || scope.aborted) return
    if (res.code === 200) {
      const regions = Array.isArray(res.data) ? res.data : []
      uiStore.setTopRegions(regions)
      if (regions.length) {
        ElMessage.success(`已识别 ${regions.length} 类顶层区域，框选时自动填充描述`)
      }
    } else {
      ElMessage.warning(res.message || '顶层识别暂不可用，可手动框选并填写描述')
    }
  } finally {
    scope.dispose()
    analyzing.value = false
  }
}

// 结构编辑器「确认标注」：把框写回图片（相对坐标供缩略展示）+ 存完整树供下单；随后关闭编辑器
function onEditorConfirm(tree) {
  const img = images.value[editingIndex.value]
  if (!img) return
  // 尺寸：优先 naturalWidth（本会话加载过），回退到上次确认持久化的 treeWidth
  const W = img.naturalWidth || uiStore.image.width || img.treeWidth || 0
  const H = img.naturalHeight || uiStore.image.height || img.treeHeight || 0
  const list = []
  const walk = (nodes) => (nodes || []).forEach((n) => {
    list.push({
      id: n.nodeKey,
      label: n.name || '未命名',
      top: !n.parentId,
      status: n.status,
      x: W ? n.bbox.x / W : 0,
      y: H ? n.bbox.y / H : 0,
      w: W ? n.bbox.w / W : 0,
      h: H ? n.bbox.h / H : 0,
    })
    walk(n.children)
  })
  walk(tree)
  img.annotations = list
  img.annotationTree = tree
  // 持久化原图尺寸（naturalWidth 刷新后会丢），供重开编辑器与下单时坐标换算兜底
  if (W && H) { img.treeWidth = W; img.treeHeight = H }
  projectEditorRef.value?.close()
  ElMessage.success('标注已确认，设置生成参数后点「开始生成」')
}

// 编辑器关闭（工具条「关闭」/确认后 close）：中止仍在途的顶层识别请求
// （确认后 close() 也会 emit close，故 onEditorConfirm 无需重复取消）
function onEditorClose() {
  cancelScopes('region:')
}

// 组件卸载（如编辑器开着就离开页面）：一并取消识别请求，避免悬挂浪费带宽
onBeforeUnmount(() => {
  cancelScopes('region:')
})

// 上传失败后重试
const onRetry = (img) => {
  if (!img.uploadId) return
  store.retryUpload(img.uploadId)
}

// 卡片摘要：只统计/展示顶层区域（子级不重复列出）
const annoStyle = (anno) => ({
  left: anno.x * 100 + '%',
  top: anno.y * 100 + '%',
  width: anno.w * 100 + '%',
  height: anno.h * 100 + '%',
})
const topAnnos = (img) => (img.annotations || []).filter((a) => a.top !== false)
const annoLabels = (img) => {
  const labels = topAnnos(img).map((a) => a.label || '未命名').filter(Boolean)
  return labels.join('、')
}
const annoTopCount = (img) => topAnnos(img).length

// 单图上传状态 -> 用于右栏状态药丸的样式类 / 文案 / 图标
const imgState = (img) => (img.uploading ? 'uploading' : img.uploadError ? 'error' : 'ready')
const imgStateText = { uploading: '上传中', error: '上传失败', ready: '已就绪' }
const imgStateIcon = { uploading: 'fas fa-spinner fa-spin', error: 'fas fa-circle-exclamation', ready: 'fa-solid fa-circle-check' }

// ── 拖拽排序 ──
// 具体拖拽 / 落点重排逻辑已抽象为全局自定义指令 v-sortable（directives/sortable.js）
// 单图场景禁用拖拽，代码保留以备后期多图复用
const dragDisabled = computed(() => images.value.length <= 1)
</script>

<style scoped>
/* ── 上传区 ── */
.upload-zone {
  border: 2px dashed rgba(0, 255, 255, 0.1);
  border-radius: 16px;
  padding: 40px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.upload-zone:hover {
  border-color: rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.02);
}

.upload-zone--disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-color: rgba(255, 255, 255, 0.08) !important;
}
.upload-zone--disabled:hover {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.08) !important;
}

.upload-icon {
  font-size: 40px;
  color: rgba(0, 255, 255, 0.25);
  margin-bottom: 12px;
}

.upload-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 6px;
}
.upload-text-bold {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.upload-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* ── 预览卡片区 ── */
.preview-area {
  margin-top: 16px;
}
.preview-area--locked {
  pointer-events: none;
  opacity: 0.45;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

/* ── 单图 hero：左侧手机壳包裹设计稿 + 右侧信息/操作 ── */
/* 右栏顶部对齐：机身高度随图变化时，右侧信息栏不被拉伸填充出空白 */
.hero-panel {
  display: flex;
  align-items: flex-start;
  gap: 22px;
  padding: 20px;
  border-radius: 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  transition: border-color 0.3s;
}
.hero-panel:hover { border-color: rgba(0, 255, 255, 0.18); }

/* 左：手机壳居中，高度撑满整栏 */
.hero-figure {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
/* 屏内设计稿图与标注框由默认插槽传入，定位相对 MobilePreview 的 .flutter-screen（position:relative） */
/* 同宽原型：机身定宽、高度随图比例展开（整图完整/无黑边/不变形，标注框按 frame 百分比对齐） */
.hero-figure .img-frame {
  width: 100%;
  height: auto;
  max-width: 100%;
}
.hero-figure .img-frame img {
  width: 100%;
  height: auto;
  object-fit: fill;
}

/* 右：信息 / 标注 / 操作，纵向撑满 */
.hero-aside {
  flex: 1 1 0;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 4px;
}

/* 文件信息卡 */
.meta-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.07);
}
.meta-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 11px;
  font-size: 16px;
  color: #00ffff;
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.15);
}
.meta-main { flex: 1 1 0; min-width: 0; display: flex; flex-direction: column; gap: 4px; }
.meta-name {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.meta-dim { font-size: 12px; color: rgba(255, 255, 255, 0.4); font-variant-numeric: tabular-nums; letter-spacing: .5px; }
.status-pill {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}
.status-pill i { font-size: 10px; }
.status-pill.ready { color: #37e0a0; background: rgba(55, 224, 160, 0.1); border: 1px solid rgba(55, 224, 160, 0.25); }
.status-pill.uploading { color: #00ffff; background: rgba(0, 255, 255, 0.08); border: 1px solid rgba(0, 255, 255, 0.2); }
.status-pill.error { color: #ff6b7a; background: rgba(255, 71, 87, 0.1); border: 1px solid rgba(255, 71, 87, 0.28); }

/* 标注面板卡 */
.anno-card {
  flex: 1 1 0;
  min-height: 240px;   /* 不依赖机身高度：右栏顶部对齐时仍自成一块舒适区域，避免 flex:1 塔缩 */
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-radius: 14px;
  background: rgba(0, 255, 255, 0.025);
  border: 1px solid rgba(0, 255, 255, 0.1);
}
.anno-card-head { display: flex; align-items: center; justify-content: space-between; }
.anno-card-title { font-size: 12px; font-weight: 700; letter-spacing: 1px; color: rgba(0, 255, 255, 0.7); }
.anno-card-count { font-size: 12px; color: rgba(255, 255, 255, 0.5); }
.anno-card-count b { color: #00ffff; font-size: 16px; margin-right: 2px; }
.anno-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  text-align: center;
}
.anno-empty i { font-size: 20px; }
.anno-empty.error { color: rgba(255, 107, 122, 0.8); }
.anno-empty.error i { font-size: 22px; }
.anno-add {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  border-radius: 10px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.55);
  background: rgba(255, 255, 255, 0.02);
  border: 1px dashed rgba(0, 255, 255, 0.28);
  transition: all 0.2s;
}
.anno-add:hover { color: #00ffff; border-color: rgba(0, 255, 255, 0.5); background: rgba(0, 255, 255, 0.06); }
.anno-add i { font-size: 22px; color: rgba(0, 255, 255, 0.6); }
.anno-add span { font-size: 13px; font-weight: 600; }
.anno-add small { font-size: 11px; color: rgba(255, 255, 255, 0.3); }
.anno-list { flex: 1; display: flex; flex-direction: column; gap: 14px; min-height: 0; }
.anno-chips { flex: 1; display: flex; flex-wrap: wrap; gap: 8px; align-content: flex-start; overflow-y: auto; }
.anno-chip {
  display: inline-flex;
  align-items: center;
  padding: 5px 11px;
  border-radius: 8px;
  font-size: 12px;
  color: rgba(0, 255, 255, 0.9);
  background: rgba(0, 255, 255, 0.08);
  border: 1px solid rgba(0, 255, 255, 0.18);
}
.anno-chip em { font-style: normal; font-weight: 600; }
.anno-edit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  color: #04110f;
  background: linear-gradient(135deg, #00ffff, #24d3c4);
  border: none;
  transition: filter 0.2s;
}
.anno-edit:hover { filter: brightness(1.08); }

/* 底部操作 */
.hero-foot { display: flex; gap: 10px; }
.foot-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 9px;
  border-radius: 9px;
  font-size: 12px;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.2s;
}
.foot-btn:hover { color: #fff; background: rgba(255, 255, 255, 0.08); border-color: rgba(255, 255, 255, 0.2); }
.foot-btn.danger:hover { color: #ff6b7a; background: rgba(255, 71, 87, 0.1); border-color: rgba(255, 71, 87, 0.3); }

.preview-card {
  position: relative;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  overflow: visible;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.preview-card:hover { border-color: rgba(0, 255, 255, 0.25); }
.preview-card--disabled { cursor: not-allowed; }
/* 拖拽态视觉由 v-sortable 指令添加以下两个 class */
.preview-card.sortable-dragging {
  opacity: 0;
  transform: scale(0.8);
}
.preview-card.sortable-over {
  border-color: transparent;
  transform: none;
}
.preview-card.sortable-over::after {
  content: '';
  position: absolute;
  inset: -3px;
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 14px;
  pointer-events: none;
  animation: dropPulse 1s ease-in-out infinite;
}

@keyframes dropPulse {
  0%, 100% { border-color: rgba(0, 255, 255, 0.3); box-shadow: none; }
  50% { border-color: rgba(0, 255, 255, 0.7); box-shadow: 0 0 12px rgba(0, 255, 255, 0.15); }
}

/* ── 图区 ── */
.card-media {
  position: relative;
  height: 220px;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 12px 12px 0 0;
  overflow: hidden;
}
/* 居中容器：整图按比例完整显示，两侧/上下可留空，不滚动 */
.media-scroll {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
/* 收缩包裹 img：定高 + aspect-ratio 让宽度按比例推导，整图等比完整显示（两侧留空不裁切不滚动） */
.img-frame {
  position: relative;
  height: 100%;
  max-width: 100%;
}
.img-frame img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
  user-select: none;
}

/* 标注框（缩略展示，不可交互） */
.anno-box {
  position: absolute;
  border: 1.5px solid rgba(0, 255, 255, 0.9);
  background: rgba(0, 255, 255, 0.08);
  pointer-events: none;
}
.anno-num {
  position: absolute;
  top: -8px;
  left: -2px;
  min-width: 14px;
  height: 14px;
  padding: 0 4px;
  border-radius: 7px;
  background: #00ffff;
  color: #0a0a0f;
  font-size: 9px;
  font-weight: 700;
  line-height: 14px;
  text-align: center;
  white-space: nowrap;
}
/* 子级区域：虚线黄框，与顶层青实框区分 */
.anno-box--child {
  border-color: rgba(255, 207, 77, 0.75);
  border-style: dashed;
  background: rgba(255, 207, 77, 0.06);
}
.anno-box--child .anno-num { background: #ffcf4d; }

/* 上传状态 */
.upload-status {
  position: absolute;
  z-index: 3;
  display: flex;
  align-items: center;
  justify-content: center;
}
/* 上传中：整图遮罩 + 居中 loading */
.upload-status.uploading {
  inset: 0;
  flex-direction: column;
  gap: 8px;
  background: rgba(0, 0, 0, 0.45);
  color: #00ffff;
  font-size: 26px;
}
.upload-status.uploading .uploading-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 1px;
}
/* 失败：整图遮罩 + 重试入口 */
.upload-status.failed {
  inset: 0;
  flex-direction: column;
  gap: 8px;
  background: rgba(0, 0, 0, 0.55);
  color: #ff4757;
  font-size: 22px;
}
.upload-status.failed .failed-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.75);
  letter-spacing: 1px;
}
.retry-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
  padding: 5px 14px;
  border-radius: 16px;
  border: 1px solid rgba(255, 71, 87, 0.5);
  background: rgba(255, 71, 87, 0.12);
  color: #ff6b7a;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.retry-btn:hover {
  background: rgba(255, 71, 87, 0.25);
  color: #fff;
}

/* 悬浮工具条 */
.card-toolbar {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.65));
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 3;
}
.preview-card:hover .card-toolbar { opacity: 1; }
.card-index {
  font-size: 13px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.85);
}
.toolbar-spacer { flex: 1; }
.tool-icon {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.2);
  color: rgba(0, 255, 255, 0.7);
  font-size: 12px;
  width: 26px;
  height: 26px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.tool-icon:hover {
  color: #00ffff;
  background: rgba(0, 255, 255, 0.12);
  border-color: rgba(0, 255, 255, 0.5);
}
.tool-icon--danger { color: rgba(255, 255, 255, 0.5); border-color: rgba(255, 255, 255, 0.15); }
.tool-icon:disabled {
  opacity: 0.35;
  cursor: not-allowed;
  color: rgba(255, 255, 255, 0.4);
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.3);
}
.tool-icon--danger:hover {
  color: #ff4757;
  background: rgba(255, 71, 87, 0.12);
  border-color: rgba(255, 71, 87, 0.4);
}

/* ── 信息区 ── */
.card-info {
  padding: 10px 12px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.info-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}
.file-name {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.75);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-size {
  flex-shrink: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
}

.anno-summary {
  font-size: 12px;
  color: rgba(0, 255, 255, 0.55);
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.anno-summary i { font-size: 11px; flex-shrink: 0; }
.anno-summary.empty { color: rgba(255, 255, 255, 0.25); }
.anno-summary.error { color: rgba(255, 71, 87, 0.75); }

/* ── 底部操作与提示 ── */
.upload-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0 4px;
}
.action-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}
.replace-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 5px 12px;
  border-radius: 16px;
  border: 1px solid rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.08);
  color: rgba(0, 255, 255, 0.8);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.replace-btn:hover { background: rgba(0, 255, 255, 0.16); color: #00ffff; }
.upload-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}
.clear-btn {
  background: none;
  border: none;
  color: #ff6b6b;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
}
.clear-btn:hover { background: rgba(255, 0, 0, 0.1); }
.clear-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  pointer-events: none;
}

.upload-tip {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
}
.upload-tip i {
  color: rgba(0, 255, 255, 0.3);
  font-size: 13px;
}
.tip-bold {
  color: rgba(255, 255, 255, 0.5);
  font-weight: 600;
}
</style>
