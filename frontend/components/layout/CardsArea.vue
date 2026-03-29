<template>
  <div class="hero-cards" ref="heroCardsRef">
    <!-- 浮动气泡卡片 -->
    <div
      v-for="(card, index) in cards"
      :key="index"
      class="floating-card"
      :class="`card-${index}`"
      :style="{
        top: card.currentTop,
        left: card.currentLeft,
        transition: card.isAnimating ? 'top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), left 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1)' : 'none',
        opacity: card.isAnimating ? 1 : 0,
        '--duration': card.duration,
        '--delay': card.delay
      }"
      @mouseenter="handleMouseEnter(index)"
      @mouseleave="handleMouseLeave(index)"
    >
      <div class="card-icon"><i :class="card.icon"></i></div>
      <div class="card-text">{{ card.text }}</div>
      <div class="card-badge" v-if="card.badge">{{ card.badge }}</div>
      <div class="card-glow"></div>
    </div>

    <!-- 代码预览卡片 -->
    <div
      class="preview-card"
      ref="previewCardRef"
      :style="{
        top: previewTop,
        opacity: previewOpacity,
        transition: 'top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), opacity 0.4s ease',
        animation: previewAnimationDone ? 'floatCard 5s ease-in-out infinite' : 'none'
      }"
    >
      <div class="preview-header">
        <div class="preview-dots">
          <span></span><span></span><span></span>
        </div>
        <span>generated.dart</span>
      </div>
      <div class="preview-code">
        <pre><code><span class="keyword">import</span> <span class="string">'package:flutter/material.dart'</span>;

<span class="keyword">class</span> <span class="class">GeneratedWidget</span> <span class="keyword">extends</span> <span class="class">StatelessWidget</span> {
  <span class="annotation">@override</span>
  <span class="class">Widget</span> <span class="function">build</span>(<span class="class">BuildContext</span> context) {
    <span class="keyword">return</span> <span class="class">Container</span>(
      <span class="property">child:</span> <span class="class">Text</span>(<span class="string">'DesignCode AI'</span>),
    );
  }
}</code></pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, watch } from 'vue'

const emit = defineEmits(['copy'])

// 卡片配置 - 添加方向属性
const cardsConfig = [
  { icon: "fas fa-image", text: "设计稿上传", badge: null, direction: "top" },
  { icon: "fas fa-brain", text: "AI 分析", badge: "AI", direction: "bottom" },
  { icon: "fas fa-code", text: "代码生成", badge: null, direction: "left" },
  { icon: "fas fa-stethoscope", text: "智能诊断", badge: "AI", direction: "right" },
  { icon: "fas fa-code-branch", text: "代码Diff", badge: "对比", direction: "random" },
]

const CARD_WIDTH = 140
const CARD_HEIGHT = 70
const TARGET_TOP = 100 // 代码卡片最终距离顶部的位置

const previewAnimationDone = ref(false)
const previewTop = ref('-100px')
const previewOpacity = ref(0)
const cards = ref([])
const hoveredIndex = ref(-1)

const previewCardRef = ref(null)
const heroCardsRef = ref(null)

// 获取代码卡片最终位置（相对于 hero-cards 容器）
const getFinalPreviewRect = () => {
  if (!previewCardRef.value || !heroCardsRef.value) return null
  const containerRect = heroCardsRef.value.getBoundingClientRect()
  const cardWidth = 350
  const cardHeight = previewCardRef.value?.offsetHeight || 320
  const left = (containerRect.width - cardWidth) / 2
  const top = TARGET_TOP
  return { left, top, width: cardWidth, height: cardHeight }
}

// 根据方向获取偏移量
const getOffsetByDirection = (direction, distance = 205) => {
  switch(direction) {
    case 'top':
      return { offsetX: 0, offsetY: -distance }
    case 'bottom':
      return { offsetX: 0, offsetY: distance }
    case 'left':
      return { offsetX: -distance, offsetY: 0 }
    case 'right':
      return { offsetX: distance, offsetY: 0 }
    case 'random':
    default:
      const angle = Math.random() * Math.PI * 2
      return {
        offsetX: Math.cos(angle) * distance,
        offsetY: Math.sin(angle) * distance
      }
  }
}

// 设置初始位置
const setInitialPositions = async () => {
  await nextTick()
  const container = heroCardsRef.value
  if (!container) return

  const finalRect = getFinalPreviewRect()
  if (!finalRect) return

  // 代码卡片初始位置：顶部外 + 透明
  previewTop.value = '-100px'
  previewOpacity.value = 0

  const finalCenterX = finalRect.left + finalRect.width / 2
  const finalCenterY = finalRect.top + finalRect.height / 2

  // 生成每个气泡的目标位置（根据方向，距离 170-240px）
  const targetPositions = []
  for (let i = 0; i < cardsConfig.length; i++) {
    const config = cardsConfig[i]
    // 距离：170-240px 随机
    const distance = 170 + Math.random() * 70
    const { offsetX, offsetY } = getOffsetByDirection(config.direction, distance)
    targetPositions.push({ offsetX, offsetY })
  }

  const newCards = []
  for (let i = 0; i < cardsConfig.length; i++) {
    const config = cardsConfig[i]
    const target = targetPositions[i]
    const duration = 3 + Math.random() * 3
    const delay = Math.random() * 2
    newCards.push({
      ...config,
      currentTop: `${finalCenterY - CARD_HEIGHT / 2}px`,
      currentLeft: `${finalCenterX - CARD_WIDTH / 2}px`,
      targetTop: `${finalCenterY + target.offsetY - CARD_HEIGHT / 2}px`,
      targetLeft: `${finalCenterX + target.offsetX - CARD_WIDTH / 2}px`,
      duration: `${duration}s`,
      delay: `${delay}s`,
      isAnimating: false,
      isVisible: false,
    })
  }
  cards.value = newCards
}

// 鼠标移入效果
const handleMouseEnter = (index) => {
  hoveredIndex.value = index
}

const handleMouseLeave = (index) => {
  hoveredIndex.value = -1
}

// 触发动画：代码滑入 -> 气泡依次扩散
const startAnimation = () => {
  const finalRect = getFinalPreviewRect()
  if (!finalRect) return

  // 第一步：代码卡片滑入 + 淡入
  previewTop.value = `${finalRect.top}px`
  previewOpacity.value = 1

  const onPreviewTransitionEnd = () => {
    if (!previewCardRef.value) return
    previewCardRef.value.removeEventListener('transitionend', onPreviewTransitionEnd)
    // 第二步：气泡依次扩散
    startBubblesAnimation()
  }

  if (previewCardRef.value) {
    previewCardRef.value.addEventListener('transitionend', onPreviewTransitionEnd)
  } else {
    setTimeout(() => {
      startBubblesAnimation()
    }, 500)
  }
}

// 依次启动每个气泡的动画
const startBubblesAnimation = () => {
  const total = cards.value.length
  let currentIndex = 0

  const animateNext = () => {
    if (currentIndex >= total) {
      // 所有气泡动画完成，开启浮动动画
      previewAnimationDone.value = true
      return
    }

    const card = cards.value[currentIndex]
    card.isAnimating = true
    card.isVisible = true
    card.currentTop = card.targetTop
    card.currentLeft = card.targetLeft

    // 监听动画结束
    const onTransitionEnd = (event) => {
      const target = event.target
      if (target && target.classList && target.classList.contains('floating-card')) {
        target.removeEventListener('transitionend', onTransitionEnd)
        currentIndex++
        animateNext()
      }
    }

    // 找到对应的 DOM 元素
    setTimeout(() => {
      const elements = document.querySelectorAll('.floating-card')
      if (elements[currentIndex]) {
        elements[currentIndex].addEventListener('transitionend', onTransitionEnd)
      } else {
        setTimeout(() => {
          currentIndex++
          animateNext()
        }, 100)
      }
    }, 50)
  }

  // 先显示所有卡片（但位置在中心，透明）
  cards.value.forEach(card => {
    card.isVisible = true
  })
  
  // 开始第一个动画
  animateNext()
}

onMounted(() => {
  setInitialPositions().then(() => {
    startAnimation()
  })
})

const copyCode = async () => {
  const code = `import 'package:flutter/material.dart';

class GeneratedWidget extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      child: Text('DesignCode AI'),
    );
  }
}`
  await navigator.clipboard.writeText(code)
  emit('copy')
  const btn = document.querySelector(".preview-copy")
  if (btn) {
    const original = btn.innerHTML
    btn.innerHTML = '<i class="fas fa-check"></i>'
    setTimeout(() => {
      btn.innerHTML = original
    }, 1500)
  }
}
</script>

<style scoped>
/* 卡片容器 */
.hero-cards {
  flex: 1;
  position: relative;
  min-height: 400px;
  width: 100%;
  overflow: visible;
}

/* 浮动卡片基础样式 */
.floating-card {
  position: absolute;
  background: rgba(10, 20, 30, 0.8);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.25);
  border-radius: 20px;
  padding: 12px 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  z-index: 5;
  cursor: pointer;
  white-space: nowrap;
  opacity: 0;
  /* 优化过渡属性：增加 transform, box-shadow 和 border-color 的平滑度 */
  transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1), 
              box-shadow 0.3s ease,
              border-color 0.3s ease,
              opacity 0.4s ease;
  will-change: transform;
}

/* 悬浮效果：向上移动并放大 */
.floating-card:hover {
  transform: translateY(-10px) scale(1.08);
  border-color: #00ffff;
  box-shadow: 0 20px 40px rgba(0, 255, 255, 0.35);
}

.card-glow {
  position: absolute;
  inset: 0;
  border-radius: 20px;
  background: radial-gradient(circle at 30% 20%, rgba(0, 255, 255, 0.15), transparent);
  opacity: 0;
  /* 调快过渡速度，更具响应感 */
  transition: opacity 0.4s ease;
  pointer-events: none;
}

.floating-card:hover .card-glow {
  opacity: 1;
}

.card-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 20px;
  padding: 2px 8px;
  font-size: 9px;
  font-weight: 600;
  color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  /* 优化徽章动画 */
  transition: transform 0.4s cubic-bezier(0.34, 1.2, 0.64, 1);
}

.floating-card:hover .card-badge {
  transform: scale(1.15) translateY(-2px);
}

.card-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s cubic-bezier(0.34, 1.2, 0.64, 1);
}

.floating-card:hover .card-icon {
  transform: scale(1.1);
}

.card-icon i {
  font-size: 20px;
  color: #00ffff;
  transition: color 0.2s ease, text-shadow 0.2s ease;
}

.floating-card:hover .card-icon i {
  color: #ffffff;
  text-shadow: 0 0 8px #00ffff;
}

.card-text {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  transition: color 0.2s ease;
}

.floating-card:hover .card-text {
  color: #00ffff;
}

/* 不同卡片样式 */
.floating-card:nth-child(1) { border-color: rgba(0, 255, 255, 0.4); }
.floating-card:nth-child(1) .card-icon i { color: #00ffff; }

.floating-card:nth-child(2) { border-color: rgba(255, 0, 255, 0.4); }
.floating-card:nth-child(2) .card-icon i { color: #ff00ff; }

.floating-card:nth-child(3) { border-color: rgba(0, 255, 200, 0.4); }
.floating-card:nth-child(3) .card-icon i { color: #00ffaa; }

.floating-card:nth-child(4) { border-color: rgba(255, 100, 0, 0.4); }
.floating-card:nth-child(4) .card-icon i { color: #ffaa00; }

.floating-card:nth-child(5) { border-color: rgba(100, 255, 100, 0.4); }
.floating-card:nth-child(5) .card-icon i { color: #88ff88; }

/* 所有扩散完成后才添加浮动动画 */
.preview-animation-done .floating-card {
  animation: floatRandom var(--duration, 4s) ease-in-out infinite;
  animation-delay: var(--delay, 0s);
}

@keyframes floatRandom {
  0%, 100% { transform: translateY(0) translateX(0); }
  25% { transform: translateY(-8px) translateX(3px); }
  75% { transform: translateY(5px) translateX(-2px); }
}

/* 代码预览卡片 */
.preview-card {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 350px;
  background: rgba(10, 20, 30, 0.9);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 20px;
  overflow: hidden;
  transition: all 0.3s;
  z-index: 4;
}

@keyframes floatCard {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(-5px); }
}

.preview-card:hover {
  border-color: #00ffff;
  transform: translateX(-50%) translateY(-5px);
  box-shadow: 0 15px 35px rgba(0, 255, 255, 0.2);
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.4);
  border-bottom: 1px solid rgba(0, 255, 255, 0.1);
}

.preview-dots {
  display: flex;
  gap: 6px;
}

.preview-dots span {
  width: 10px;
  height: 10px;
  
  border-radius: 50%;
}
.preview-dots span:nth-child(1) { background: #ff5f56; }
.preview-dots span:nth-child(2) { background: #ffbd2e; }
.preview-dots span:nth-child(3) { background: #27c93f; }

.preview-header span {
  flex: 1;
  font-family: monospace;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  text-align: center;
}


.preview-code {
  padding: 16px;
  font-family: "Fira Code", monospace;
  font-size: 12px;
  line-height: 1.5;
  overflow-x: auto;
}

.preview-code pre {
  margin: 0;
  color: #ccc;
  text-align: left;
  white-space: pre-wrap;
}

.preview-code .keyword { color: #ff79c6; }
.preview-code .string { color: #f1fa8c; }
.preview-code .class { color: #50fa7b; }
.preview-code .annotation { color: #8be9fd; }
.preview-code .function { color: #8be9fd; }
.preview-code .property { color: #ffb86c; }

/* 响应式 */
@media (max-width: 768px) {
  .floating-card {
    position: relative;
    top: auto !important;
    left: auto !important;
    display: inline-flex;
    margin: 8px;
    opacity: 1 !important;
  }

  .hero-cards {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    align-items: center;
    min-height: auto;
  }

  .preview-card {
    position: relative;
    left: auto;
    transform: none;
    width: 100%;
    margin-top: 20px;
    top: auto !important;
  }
  
  @keyframes floatCard {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-5px); }
  }
}
</style>