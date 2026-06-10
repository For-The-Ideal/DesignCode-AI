<template>
  <div class="particles-fixed">
    <div v-for="i in 100" :key="i" class="particle" :style="getParticleStyle()"></div>
  </div>
</template>

<script setup>
const getParticleStyle = () => {
  const size = Math.random() * 3 + 1
  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${Math.random() * 100}%`,
    animationDuration: `${Math.random() * 20 + 10}s`,
    animationDelay: `${Math.random() * 10}s`,
  }
}
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
}

.particle {
  position: absolute;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 50%;
  opacity: 0;
  animation: floatParticle linear infinite;
}

@keyframes floatParticle {
  0%   { transform: translateY(100vh) rotate(0deg);  opacity: 0; }
  10%  { opacity: 0.4; }
  90%  { opacity: 0.4; }
  100% { transform: translateY(-100vh) rotate(360deg); opacity: 0; }
}
</style>
