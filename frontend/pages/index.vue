<template>
  <div class="app">

    <div class="screens-container">
      <!-- 第一屏 -->
      <HeroSection @getStarted="scrollToSecond" @watchDemo="playDemo" />

      <!-- 第二屏 -->
      <!-- <SSEGenerator ref="secondScreenRef" /> -->

      <!-- 第三屏 -->
      <!-- <DiagnosticBubbles ref="thirdScreenRef" /> -->
    </div>
  </div>
</template>

<script setup>
import { ElMessage } from 'element-plus'
import HeroSection from "~/components/landing/HeroSection.vue";
import SSEGenerator from "~/components/landing/SSEGenerator.vue";
import DiagnosticBubbles from "~/components/landing/DiagnosticBubbles.vue";


const scrollToSecond = () => {
  const container = document.querySelector(".screens-container");
  if (container) container.scrollBy({ top: window.innerHeight, behavior: "smooth" });
};

const playDemo = () => {
  ElMessage.warning('演示视频即将上线')
};
</script>

<style scoped>
.app {
  position: relative;
  background: #03050b;
}
.screens-container {
  scroll-snap-type: y mandatory;
  overflow-y: scroll;
  height: 100vh;
  scroll-behavior: smooth;
}

html, body {
  overflow: hidden;
  height: 100%;
  margin: 0;
  padding: 0;
}
body {
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE/Edge */
}
body::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}
</style>

<style>
/* 全局隐藏滚动条（兼容所有浏览器） */
html, body {
  overflow: hidden;
  height: 100%;
  margin: 0;
  padding: 0;
}
body {
  overflow-y: auto;
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE/Edge */
}
body::-webkit-scrollbar {
  display: none; /* Chrome/Safari */
}
</style>
