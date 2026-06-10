<template>
  <section class="screen diagnostic-screen" ref="screenRef">
    <div class="screen-content">
      <div class="header-title">
        <h2>🧠 AI 智能诊断 · 优化工坊</h2>
        <p>深度扫描代码质量，动态生成个性化优化建议</p>
      </div>

      <!-- 紧凑型气泡布局 -->
      <div class="bubble-grid" ref="bubbleGridRef">
        <div class="center-caption" ref="captionRef">AI 诊断中心</div>

        <div
          v-for="(b, i) in cards"
          :key="i"
          class="bubble"
          :class="{ 'spread-done': animationDone }"
          :style="{
            top: b.currentTop,
            left: b.currentLeft,
            transition: `all 0.3s cubic-bezier(0.34, 1.2, 0.64, 1), top 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1), left 0.5s cubic-bezier(0.2, 0.9, 0.4, 1.1)`,
            opacity: b.isVisible ? 1 : 0,
          }"
        >
          <i :class="b.icon"></i>
          <span class="bubble-title">{{ b.title }}</span>
          <span class="bubble-desc">{{ b.desc }}</span>
        </div>
      </div>
    </div>
    <div class="scroll-hint"></div>
  </section>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useCardSpreadAnimation } from '@/composables/useCardSpreadAnimation'

const screenRef = ref(null)
const bubbleGridRef = ref(null)
const captionRef = ref(null)

const bubbleConfigs = [
  { icon: 'fas fa-search', title: '代码扫描', desc: '已分析 127 行代码' },
  { icon: 'fas fa-exclamation-triangle', title: '风险提示', desc: '检测到 2 处潜在缺陷' },
  { icon: 'fas fa-tachometer-alt', title: '性能诊断', desc: '渲染耗时可优化 15%' },
  { icon: 'fas fa-shield-alt', title: '安全检测', desc: '通过 OWASP Top 10' },
  { icon: 'fas fa-sync-alt', title: '重构建议', desc: '建议提取 1 个 composable' },
]

const { cards, animationDone, start } = useCardSpreadAnimation({
  cardWidth: 200,
  cardHeight: 44,
  radius: 240,
  anglePoints: [
    { angle: -130, rFactor: 1.1 },
    { angle: -30,  rFactor: 1.05 },
    { angle: 30,   rFactor: 1.1 },
    { angle: 120,  rFactor: 1.2 },
    { angle: 180,  rFactor: 0.9 },
  ],
})

let observer = null
onMounted(() => {
  observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting && !animationDone.value) {
        if (!captionRef.value) return
        const gridRect = bubbleGridRef.value.getBoundingClientRect()
        const capRect = captionRef.value.getBoundingClientRect()
        const cx = capRect.left + capRect.width / 2 - gridRect.left
        const cy = capRect.top + capRect.height / 2 - gridRect.top
        start(bubbleConfigs, cx, cy)
        observer.unobserve(entry.target)
      }
    })
  }, { threshold: 0.3 })
  if (screenRef.value) observer.observe(screenRef.value)
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})
</script>

<style scoped>
.screen {
  scroll-snap-align: start;
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: relative;
  padding: 80px 2rem 2rem;
  background: transparent;
}

.screen-content {
  max-width: 900px;
  width: 100%;
  margin: 0 auto;
}

.header-title {
  text-align: center;
  margin-bottom: 2.5rem;
}

.header-title h2 {
  font-family: 'Orbitron', monospace;
  font-size: 1.8rem;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  display: inline-block;
}

.header-title p {
  color: #acc9ff;
  margin-top: 0.25rem;
  font-size: 0.9rem;
}

.bubble-grid {
  position: relative;
  min-height: 480px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.center-caption {
  width: 110px;
  height: 110px;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  border: 1px solid #0ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  font-family: 'Orbitron', monospace;
  font-weight: bold;
  color: #0ff;
  font-size: 0.9rem;
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
  z-index: 2;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.bubble {
  position: absolute;
  background: rgba(10, 20, 35, 0.85);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 255, 255, 0.4);
  border-radius: 60px;
  padding: 8px 18px;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
  white-space: nowrap;
  z-index: 5;
  user-select: none;
  animation: bubbleBreath 3s ease-in-out infinite;
}
.bubble:nth-child(2) { animation-delay: 0s; }
.bubble:nth-child(3) { animation-delay: 0.6s; }
.bubble:nth-child(4) { animation-delay: 1.2s; }
.bubble:nth-child(5) { animation-delay: 1.8s; }
.bubble:nth-child(6) { animation-delay: 2.4s; }

.bubble:hover {
  border-color: #00ffff;
  box-shadow: 0 15px 30px rgba(0, 255, 255, 0.3);
  background: rgba(20, 35, 55, 0.95);
  z-index: 10;
  animation: none;
}

@keyframes bubbleBreath {
  0%, 100% { box-shadow: 0 5px 15px rgba(0,0,0,0.3); border-color: rgba(0,255,255,0.4); }
  50% { box-shadow: 0 8px 28px rgba(0,255,255,0.35); border-color: rgba(0,255,255,0.8); }
}

.bubble i {
  font-size: 20px;
  color: #00ffff;
  transition: text-shadow 0.2s;
}

.bubble:hover i {
  text-shadow: 0 0 6px #00ffff;
}

.bubble-title {
  font-size: 13px;
  font-weight: 600;
  color: white;
}

.bubble-desc {
  font-size: 10px;
  color: #aaffdd;
  margin-left: 5px;
}

.refresh-diagnostic {
  text-align: center;
  margin-top: 40px;
}

.btn-primary {
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  padding: 10px 24px;
  border-radius: 40px;
  color: white;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 255, 255, 0.4);
}

.scroll-hint {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  width: 30px;
  height: 30px;
}

@media (max-width: 768px) {
  .bubble-grid {
    min-height: auto;
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 20px;
  }
  .bubble {
    position: relative;
    top: auto !important;
    left: auto !important;
    right: auto !important;
    bottom: auto !important;
    transform: none !important;
    width: 100%;
    justify-content: center;
    margin: 0;
  }
  .center-caption {
    position: relative;
    top: auto;
    left: auto;
    transform: none;
    margin: 20px auto;
  }
  .bubble-5 {
    margin-top: 0;
  }
}
</style>