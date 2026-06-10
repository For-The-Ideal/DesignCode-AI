<template>
  <div class="hero-cards" ref="heroCardsRef" :class="{ 'preview-animation-done': previewAnimationDone }">
    <!-- 浮动气泡卡片 -->
    <div
      v-for="(card, index) in cards"
      :key="index"
      class="floating-card"
      :class="`card-${index}`"
      :style="{
        top: card.currentTop,
        left: card.currentLeft,
        transition: card.isAnimating ? 'top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), left 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), opacity 0.5s ease' : '',
        opacity: previewAnimationDone ? '' : (card.isVisible ? 1 : 0),
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
        transition: previewAnimationDone 
          ? 'top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), opacity 0.4s ease, transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.4s ease, border-color 0.4s ease, background 0.4s ease'
          : 'top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), opacity 0.4s ease',
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
import { ref, onMounted, watch } from 'vue'
import { useCardSpreadAnimation } from '@/composables/useCardSpreadAnimation'

const emit = defineEmits(['copy'])

const cardsConfig = [
  { icon: "fas fa-image", text: "设计稿上传", badge: null },
  { icon: "fas fa-brain", text: "AI 分析", badge: "AI" },
  { icon: "fas fa-code", text: "代码生成", badge: null },
  { icon: "fas fa-stethoscope", text: "智能诊断", badge: "AI" },
  { icon: "fas fa-code-branch", text: "代码Diff", badge: "对比" },
]

const CARD_WIDTH = 140
const CARD_HEIGHT = 70
const TARGET_TOP = 120
const RADIUS = 280

const previewAnimationDone = ref(false)
const previewTop = ref('-100px')
const previewOpacity = ref(0)
const hoveredIndex = ref(-1)

const previewCardRef = ref(null)
const heroCardsRef = ref(null)

// ── 卡片扩散动画 composable ──
const anglePoints = [
  { angle: -150, rFactor: 1.15 },
  { angle: -15,  rFactor: 1.1 },
  { angle: 15,   rFactor: 1.25 },
  { angle: 160,  rFactor: 1.15 },
  { angle: 90,   rFactor: 0.85 },
]

// 获取代码卡片最终位置
const getFinalPreviewRect = () => {
  if (!previewCardRef.value || !heroCardsRef.value) return null
  const containerRect = heroCardsRef.value.getBoundingClientRect()
  const cardWidth = 350
  const cardHeight = previewCardRef.value?.offsetHeight || 320
  const left = (containerRect.width - cardWidth) / 2
  const top = TARGET_TOP
  return { left, top, width: cardWidth, height: cardHeight, containerWidth: containerRect.width }
}

// 使用扩散动画 composable
const { cards, animationDone: spreadDone, start: startSpread } = useCardSpreadAnimation({
  cardWidth: CARD_WIDTH,
  cardHeight: CARD_HEIGHT,
  radius: RADIUS,
  anglePoints,
})

// 触发动画：代码滑入 → 气泡扩散
const startAnimation = () => {
  const finalRect = getFinalPreviewRect()
  if (!finalRect) return

  // 设定 composable 需要的中心点（preview 卡中心）
  const cx = finalRect.left + finalRect.width / 2
  const cy = finalRect.top + finalRect.height / 2

  // 第一步：代码卡片滑入 + 淡入
  previewTop.value = `${finalRect.top}px`
  previewOpacity.value = 1

  const onPreviewTransitionEnd = () => {
    if (!previewCardRef.value) return
    previewCardRef.value.removeEventListener('transitionend', onPreviewTransitionEnd)
    // 第二步：气泡扩散
    startSpread(cardsConfig, cx, cy)
  }

  if (previewCardRef.value) {
    previewCardRef.value.addEventListener('transitionend', onPreviewTransitionEnd)
  } else {
    setTimeout(() => startSpread(cardsConfig, cx, cy), 500)
  }
}

// 鼠标移入效果
const handleMouseEnter = (index) => {
  hoveredIndex.value = index
}

const handleMouseLeave = (index) => {
  hoveredIndex.value = -1
}

onMounted(() => {
  watch(spreadDone, (val) => {
    if (val) previewAnimationDone.value = true
  })
  startAnimation()
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
  min-height: 500px;
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
  transition: transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1), 
              box-shadow 0.4s ease,
              border-color 0.4s ease,
              opacity 0.4s ease,
              background 0.4s ease;
  will-change: transform;
}

/* 悬浮效果：向上移动并放大 */
/* 悬浮时停止基础动画，应用 Hover 样式 */
.floating-card:hover {
  animation-play-state: paused;
  transform: translateY(-15px) scale(1.12) !important;
  background: rgba(15, 30, 45, 0.9);
  border-color: #00ffff !important;
  box-shadow: 0 25px 50px rgba(0, 255, 255, 0.4),
              0 0 15px rgba(0, 255, 255, 0.2) !important;
  z-index: 10;
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

/* 所有扩散完成后才添加浮动和呼吸动画 */
.preview-animation-done .floating-card {
  opacity: 1;
  animation: subtleFloat var(--duration, 4s) ease-in-out infinite,
             glowBreathing var(--duration, 4s) ease-in-out infinite;
  animation-delay: var(--delay, 0s);
}

@keyframes subtleFloat {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}

@keyframes glowBreathing {
  0%, 100% { 
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.1);
    opacity: 0.85;
  }
  50% { 
    box-shadow: 0 15px 30px rgba(0, 255, 255, 0.25);
    border-color: rgba(0, 255, 255, 0.6);
    opacity: 1;
  }
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
  z-index: 4;
}

@keyframes floatCard {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(-5px); }
}

.preview-card:hover {
  border-color: #00ffff;
  transform: translateX(-50%) translateY(-12px) scale(1.03);
  box-shadow: 0 25px 50px rgba(0, 255, 255, 0.35),
              0 0 15px rgba(0, 255, 255, 0.1);
  background: rgba(15, 30, 45, 0.95);
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