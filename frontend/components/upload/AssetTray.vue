<template>
  <div class="asset-tray" :class="{ 'asset-tray--embedded': embedded }">
    <div class="tray-head">
      <div class="tray-title">
        <i class="fa-solid fa-shapes"></i>
        附图素材
        <span class="tray-count">{{ assets.length }}/{{ MAX_ASSETS }}</span>
      </div>
      <div class="tray-tools">
        <!-- 多图模式：新上传素材的归属页（每张主图有独立附图集） -->
        <select v-if="mode === 'multi'" v-model="ownerId" class="owner-select" title="新上传附图的归属页">
          <option v-for="(img, i) in images" :key="img.uploadId" :value="img.uploadId">
            归属 P{{ i + 1 }}{{ img.file?.name ? ' · ' + shortName(img.file.name) : '' }}
          </option>
        </select>
        <button class="tray-add" :disabled="!canAdd" @click="pickFiles">
          <i class="fas fa-plus"></i> 上传
        </button>
      </div>
    </div>
    <p class="tray-hint">
      仅上传需要作为真实素材使用的图形（Logo / 图标 / 插画 / 照片），标注时可在画布上框选引用；纯装饰元素无需上传
    </p>

    <div class="tray-body" @dragover.prevent @drop.prevent="onDrop">
      <div
        v-for="a in assets"
        :key="a.assetId"
        class="asset-card"
        :class="{ 'asset-card--err': !!a.uploadError, 'asset-card--linked': isLinked(a) }"
      >
        <div class="asset-thumb">
          <!-- 缩略图：同 TasksGrid 用 el-image + 失败占位 -->
          <el-image v-if="a.preview" :src="a.preview" fit="contain" class="asset-img" loading="lazy">
            <template #error><div class="asset-ph"><i class="fas fa-image"></i></div></template>
          </el-image>
          <div v-else class="asset-ph"><i class="fas fa-image"></i></div>
          <div v-if="a.uploading" class="asset-mask"><i class="fas fa-spinner fa-spin"></i></div>
          <div v-else-if="a.uploadError" class="asset-mask asset-mask--err">
            <button class="mini-btn" @click="store.retryAsset(a.assetId)">重试</button>
          </div>
          <!-- 已被标注区域引用：勾选标记，与标注弹层里的选中态一致 -->
          <span v-if="isLinked(a)" class="asset-linked" title="已被标注区域引用">
            <i class="fa-solid fa-check"></i>
          </span>
          <!-- 归属徽标：多图模式下点击切换所属主图；已有勾选时移到左上，避免重叠 -->
          <span
            v-if="mode === 'multi'"
            class="asset-owner"
            :class="{ 'asset-owner--none': !a.ownerId, 'asset-owner--shift': isLinked(a) }"
            :title="images.length > 1 ? '点击切换归属主图' : '未归属'"
            @click="cycleOwner(a)"
          >
            {{ ownerLabel(a) }}
          </span>
          <!-- hover 蒙层：未引用时可删除；已被标注引用的只能在标注弹层内解除/删除 -->
          <div class="asset-hover">
            <button v-if="!isLinked(a)" class="asset-del" title="删除素材" @click="confirmRemove(a)">
              <i class="fa-regular fa-trash-can"></i>
            </button>
          </div>
        </div>
        <input
          class="asset-label"
          :value="a.label"
          maxlength="30"
          :title="a.label"
          @change="store.updateAssetLabel(a.assetId, $event.target.value)"
        />
      </div>

      <div v-if="!assets.length" class="tray-empty" @click="pickFiles">
        <i class="fa-regular fa-images"></i>
        <span>暂无附图，点击或拖拽图片到此处上传</span>
      </div>
    </div>

    <input ref="fileInput" type="file" multiple accept="image/*" style="display: none" @change="onFileChosen" />
  </div>
</template>

<script setup>
/**
 * 附图素材托盘 —— 上传区子组件，单图/多图模式共用
 * 状态与上传动作统一挂在 code store（assets），不做独立存储。
 */
import { ref, computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCodeStore } from '~/stores/code'
import { useUiProjectStore } from '~/stores/uiProject'
import { MAX_ASSETS, isImageFile } from '~/stores/code/upload'

defineProps({
  // 嵌入单图右栏时去掉外层卡片壳（背景/边框由 hero-aside 提供）
  embedded: { type: Boolean, default: false },
})

const store = useCodeStore()
const uiStore = useUiProjectStore()
const { assets, images, mode, isSubmitting } = storeToRefs(store)

// 已被标注节点引用的素材：显示勾选且禁止在托盘删除（解除/删除只能在标注弹层内完成）
const linkedIds = computed(() => {
  const set = new Set()
  uiStore.nodeList.forEach((n) => { if (n.assetId) set.add(n.assetId) })
  return set
})
const isLinked = (a) => linkedIds.value.has(a.assetId)

const fileInput = ref(null)
const canAdd = computed(() => !isSubmitting.value && assets.value.length < MAX_ASSETS)

// 多图模式：新上传素材的归属页（跟随主图增删自动纠正，默认第一页）
const ownerId = ref('')
watch(
  () => images.value.map((i) => i.uploadId),
  (ids) => {
    if (!ids.length) { ownerId.value = ''; return }
    if (!ids.includes(ownerId.value)) ownerId.value = ids[0]
  },
  { immediate: true },
)

// 素材归属页序（P1/P2…），找不到主图时为未归属
const ownerIndex = (a) => images.value.findIndex((i) => i.uploadId === a.ownerId)
const ownerLabel = (a) => {
  const idx = ownerIndex(a)
  return idx > -1 ? `P${idx + 1}` : '未归属'
}

// 点击徽标：按主图顺序循环切换归属
function cycleOwner(a) {
  if (mode.value !== 'multi' || images.value.length < 2) return
  const idx = ownerIndex(a)
  const next = images.value[(idx + 1) % images.value.length]
  store.setAssetOwner(a.assetId, next.uploadId)
}

const shortName = (name) => name.replace(/\.[^.]+$/, '').slice(0, 10)

// 新增素材挂到当前归属页（单图不传，store 自动挂唯一主图）
const newOwner = () => (mode.value === 'multi' ? ownerId.value : '')

function pickFiles() {
  if (isSubmitting.value) return
  if (!canAdd.value) {
    ElMessage.warning(`附图最多 ${MAX_ASSETS} 张`)
    return
  }
  fileInput.value?.click()
}

function onFileChosen(e) {
  const files = Array.from(e.target.files || [])
  if (files.length) store.addAssets(files, newOwner())
  if (fileInput.value) fileInput.value.value = ''
}

function onDrop(e) {
  if (isSubmitting.value) return
  const files = Array.from(e.dataTransfer.files || []).filter(isImageFile)
  if (files.length) store.addAssets(files, newOwner())
}

async function confirmRemove(a) {
  // 已被标注引用的素材：托盘内禁止删除，需先在标注弹层内解除
  if (isLinked(a)) {
    ElMessage.warning('该素材已被标注区域引用，请在标注弹层内解除或删除')
    return
  }
  try {
    await ElMessageBox.confirm(`确定删除素材「${a.label}」？`, '删除附图', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  store.removeAsset(a.assetId)
}
</script>

<style scoped>
.asset-tray {
  margin-top: 16px;
  padding: 14px 16px;
  border-radius: 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 嵌入单图右栏：无壳、紧凑网格 */
.asset-tray--embedded { margin-top: 0; padding: 0; background: transparent; border: none; }
.asset-tray--embedded .tray-body { grid-template-columns: repeat(auto-fill, minmax(72px, 1fr)); min-height: 72px; }
.asset-tray--embedded .tray-empty { min-height: 72px; }

.tray-head { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.tray-tools { display: flex; align-items: center; gap: 8px; }
.owner-select {
  max-width: 150px;
  padding: 5px 8px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: #0d1117;
  color: rgba(255, 255, 255, 0.75);
  font-size: 11px;
  outline: none;
  cursor: pointer;
}
.owner-select:focus { border-color: rgba(0, 255, 255, 0.45); }
.tray-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.85);
}
.tray-title i { color: #00ffff; font-size: 13px; }
.tray-count { font-size: 11px; color: rgba(255, 255, 255, 0.4); font-variant-numeric: tabular-nums; }

.tray-add {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.35);
  background: transparent;
  color: #00ffff;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.tray-add:hover:not(:disabled) { background: rgba(0, 255, 255, 0.1); }
.tray-add:disabled { opacity: 0.4; cursor: not-allowed; }

.tray-hint { margin: 0; font-size: 11px; line-height: 1.6; color: rgba(255, 255, 255, 0.4); }

.tray-body {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(92px, 1fr));
  gap: 10px;
  min-height: 100px;
}

.tray-empty {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 2px dashed rgba(0, 255, 255, 0.16);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.4);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.tray-empty:hover { border-color: rgba(0, 255, 255, 0.45); color: #00ffff; }
.tray-empty i { font-size: 20px; }

.asset-card {
  position: relative;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.03);
  overflow: hidden;
  transition: border-color 0.2s;
}
.asset-card:hover { border-color: rgba(0, 255, 255, 0.5); }
.asset-card--err { border-color: rgba(255, 71, 87, 0.5); }
/* 已关联标注：与标注弹层选中态一致（青色描边 + 光晕） */
.asset-card--linked { border-color: #00ffff; box-shadow: 0 0 10px rgba(0, 255, 255, 0.22); }
.asset-card--linked:hover { border-color: #00ffff; }

/* 已关联标注：右上角青色勾选圆点，样式与弹层内 .asset-opt-check 完全一致 */
.asset-linked {
  position: absolute;
  top: 2px;
  right: 2px;
  z-index: 4;
  width: 15px;
  height: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #00ffff;
  color: #04110f;
  font-size: 8px;
}

/* 已关联标注：hover 锁定提示，替代删除按钮 */
.asset-locked {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(0, 20, 20, 0.8);
  border: 1px solid rgba(0, 255, 255, 0.4);
  color: #00ffff;
  font-size: 11px;
}

.asset-thumb { position: relative; aspect-ratio: 1; background: #0a0a0f; }
.asset-img { width: 100%; height: 100%; display: block; }
.asset-thumb :deep(.el-image__inner) {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
  /* 棋盘格底：透明 PNG 图标也能看清边界 */
  background:
    conic-gradient(rgba(255, 255, 255, 0.05) 0 25%, transparent 0 50%, rgba(255, 255, 255, 0.05) 0 75%, transparent 0)
    0 0 / 16px 16px;
}
.asset-ph {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.asset-ph i { font-size: 22px; color: rgba(255, 255, 255, 0.12); }
.asset-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
}
.asset-mask--err { color: #ff4757; }
.mini-btn {
  padding: 3px 10px;
  border-radius: 6px;
  border: 1px solid #00ffff;
  background: rgba(0, 0, 0, 0.5);
  color: #00ffff;
  font-size: 11px;
  cursor: pointer;
}

/* hover 蒙层：覆盖缩略图，透明不可点，避免误触删除 */
.asset-hover {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 8, 12, 0.55);
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.18s ease;
}
.asset-card:hover .asset-hover { opacity: 1; }

/* 中央删除按钮：仅按钮本身可点 */
.asset-del {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  border: 1px solid rgba(255, 71, 87, 0.55);
  background: rgba(20, 10, 14, 0.75);
  color: #ff6b7a;
  font-size: 18px;
  cursor: pointer;
  pointer-events: auto;
  transition: all 0.18s ease;
}
.asset-del:hover {
  border-color: #ff4757;
  background: rgba(255, 71, 87, 0.22);
  color: #fff;
  transform: scale(1.08);
}

/* 归属徽标：P1/P2…，未归属红色告警（置顶于 hover 蒙层之上，悬停时仍可切归属） */
.asset-owner {
  position: absolute;
  top: 4px;
  right: 4px;
  z-index: 3;
  padding: 1px 6px;
  border-radius: 6px;
  background: rgba(0, 0, 0, 0.65);
  border: 1px solid rgba(0, 255, 255, 0.4);
  color: #00ffff;
  font-size: 10px;
  line-height: 1.5;
  cursor: pointer;
  user-select: none;
}
.asset-owner--none { border-color: rgba(255, 71, 87, 0.55); color: #ff4757; }
/* 右上角让位给勾选：归属徽标移到左上 */
.asset-owner--shift { left: 4px; right: auto; }

.asset-label {
  width: 100%;
  padding: 5px 8px;
  border: none;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: transparent;
  color: rgba(255, 255, 255, 0.75);
  font-size: 11px;
  text-align: center;
  outline: none;
}
.asset-label:focus { background: rgba(0, 255, 255, 0.06); color: #fff; }
</style>
