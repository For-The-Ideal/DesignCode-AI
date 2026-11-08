<template>
  <button class="generate-btn" :disabled="btnState.disabled" @click="onGenerate">
    <i :class="btnState.icon"></i>
    {{ btnState.text }}
  </button>
</template>

<script setup>
/**
 * 开始生成按钮：从 ConfigPanel 抽出。
 * 上传页把它放在配置项末尾；工作台右栏把它固定在底部，折叠配置区也能直接点。
 */
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCodeStore } from '~/stores/code'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'
import { debounce } from '~/utils/index'

const store = useCodeStore()
const { images, allUploaded, isSubmitting, isConcurrencyFull } = storeToRefs(store)
const userStore = useUserStore()
const commonStore = useCommonStore()

const canGenerate = computed(() => images.value.length > 0 && allUploaded.value && !isConcurrencyFull.value)

// 上传完成但未被任何标注区域引用的附图（所有图片的标注树都算）
const collectUnusedAssets = () => {
  const used = new Set()
  const walk = (nodes) => (nodes || []).forEach((n) => {
    if (n.assetId) used.add(n.assetId)
    walk(n.children)
  })
  images.value.forEach((img) => walk(img.annotationTree))
  return store.assets.filter((a) => !a.uploading && !a.uploadError && a.cosUrl && !used.has(a.assetId))
}

const onGenerate = debounce(250, async () => {
  if (!userStore.isLogin) {
    commonStore.setLoginModalVisible(true)
    return
  }
  // 已上传但未被任何区域使用的附图：提示一键清空；取消则什么都不动（不清理、不生成）
  const unused = collectUnusedAssets()
  if (unused.length) {
    try {
      await ElMessageBox.confirm(
        `有 ${unused.length} 张附图已上传但未被任何标注区域使用，继续将清空这些附图后开始生成，取消则不做任何操作。`,
        '未使用的附图',
        { confirmButtonText: '继续', cancelButtonText: '取消', type: 'warning' },
      )
    } catch {
      return
    }
    unused.forEach((a) => store.removeAsset(a.assetId))
    ElMessage.success(`已清空 ${unused.length} 张未使用的附图`)
  }
  await store.createTask()
})

const btnState = computed(() => {
  if (!userStore.isLogin) {
    return { text: '请先登录', icon: 'fa-solid fa-wand-magic-sparkles', disabled: true }
  }
  if (!images.value.length) {
    return { text: '请先上传设计稿', icon: 'fa-regular fa-image', disabled: true }
  }
  if (!canGenerate.value) {
    return { text: '开始生成', icon: 'fa-solid fa-wand-magic-sparkles', disabled: true }
  }
  if (isSubmitting.value) {
    return { text: '生成中...', icon: 'fa-solid fa-spinner fa-spin', disabled: true }
  }
  return { text: '开始生成', icon: 'fa-solid fa-wand-magic-sparkles', disabled: false }
})
</script>

<style scoped>
.generate-btn {
  width: 100%;
  padding: 12px 0;
  border-radius: 8px;
  border: none;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: linear-gradient(135deg, #60a5fa, #818cf8);
  box-shadow: 0 4px 20px rgba(96, 165, 250, 0.25);
  transition: all 0.3s;
}
.generate-btn:hover:not(:disabled) {
  box-shadow: 0 6px 30px rgba(96, 165, 250, 0.4);
  transform: translateY(-1px);
}
.generate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
