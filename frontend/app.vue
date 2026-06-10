<template>
  <div>
    <ParticlesBackground />
    <AppHeader />
    <NuxtPage />
    <AppFooter />
  </div>
</template>

<script setup>
import AppFooter from "~/components/layout/AppFooter.vue";
import AppHeader from "~/components/layout/AppHeader.vue";
import ParticlesBackground from "~/components/layout/ParticlesBackground.vue";
</script>

<style lang="scss" scoped>
/* 确保全局样式生效 */
.app-root {
  min-height: 100vh;
  background: #0a0a0f;
}
</style>
