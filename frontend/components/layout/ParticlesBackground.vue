<template>
  <div class="particles-fixed" :style="{ opacity: isReducedMotion ? 0.3 : 1 }">
    <div
      v-for="(style, idx) in particleStyles"
      :key="idx"
      class="particle"
      :style="style"
    ></div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'

// 可配置的 props
const props = defineProps({
  count: {
    type: Number,
    default: 100, // 降低默认数量，提升性能
  },
  speedRange: {
    type: Array,
    default: () => [10, 30], // [最小秒数, 最大秒数]
  },
  sizeRange: {
    type: Array,
    default: () => [1, 10], // [最小px, 最大px]
  },
  colors: {
    type: Array,
    default: () => [
      'linear-gradient(135deg, #00ffff, #ff00ff)',
      'linear-gradient(135deg, #00ffaa, #ffaa00)',
      'linear-gradient(135deg, #ff66cc, #66ffcc)',
    ],
  },
})

// 检测用户是否偏好减少动画
const prefersReducedMotion = ref(false)

onMounted(() => {
  if (window.matchMedia) {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
    prefersReducedMotion.value = mediaQuery.matches
    const handler = (e) => (prefersReducedMotion.value = e.matches)
    mediaQuery.addEventListener('change', handler)
    // 清理监听器（组件卸载时）
    // 注意：由于粒子组件通常全局存在，不必须，但为了规范可保留
  }
})

const isReducedMotion = computed(() => prefersReducedMotion.value)

// 生成随机数函数
const random = (min, max) => min + Math.random() * (max - min)

// 生成单个粒子的样式（静态，只在初始化时执行一次）
const generateParticleStyle = () => {
  const size = random(props.sizeRange[0], props.sizeRange[1])
  const duration = random(props.speedRange[0], props.speedRange[1])
  const delay = random(0, duration * 0.8) // 延迟不超过动画周期的 80%
  // 随机起始位置偏移（避免全部从底部开始）
  const startOffsetY = random(0, 100) // vh 单位
  const colorIndex = Math.floor(Math.random() * props.colors.length)
  const opacityVal = random(0.2, 0.6) // 透明度变化范围

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${random(0, 100)}%`,
    '--duration': `${duration}s`,
    '--delay': `${delay}s`,
    '--start-offset': `${startOffsetY}vh`,
    '--particle-color': props.colors[colorIndex],
    '--opacity': opacityVal,
  }
}

// 存储粒子样式（普通数组，不需要响应式深度追踪）
const particleStyles = ref([])

onMounted(() => {
  // 生成所有粒子的样式
  const styles = Array.from({ length: props.count }, () => generateParticleStyle())
  particleStyles.value = styles
})
</script>

<style scoped>
.particles-fixed {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
  transition: opacity 0.3s ease;
}

.particle {
  position: absolute;
  bottom: 0; /* 配合 transform 从下往上 */
  background: var(--particle-color, linear-gradient(135deg, #00ffff, #ff00ff));
  border-radius: 50%;
  opacity: 0;
  animation: floatParticle var(--duration, 20s) linear infinite;
  animation-delay: var(--delay, 0s);
  will-change: transform; /* 性能优化提示 */
}

/* 减少动画模式下，只保留极简效果或完全禁用 */
@media (prefers-reduced-motion: reduce) {
  .particle {
    animation: none;
    opacity: 0.1;
  }
}

@keyframes floatParticle {
  0% {
    transform: translateY(100vh) translateX(0) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: var(--opacity, 0.4);
  }
  90% {
    opacity: var(--opacity, 0.4);
  }
  100% {
    transform: translateY(calc(-1 * var(--start-offset, 0vh))) translateX(20px) rotate(360deg);
    opacity: 0;
  }
}
</style>