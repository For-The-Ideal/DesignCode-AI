<template>
  <footer class="app-footer">
    <div class="footer-inner">
      <span class="copyright">© {{ year }} DesignCode AI. All Rights Reserved.</span>
      <a
        class="icp"
        href="https://beian.miit.gov.cn/"
        target="_blank"
        rel="noopener noreferrer"
      >
        粤ICP备XXXXXXXX号
      </a>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
.app-footer {
  width: 100%;
  padding: 24px 2rem;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.3);
}
.footer-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  font-size: 0.8rem;
  color: #666;
}
.copyright {
  color: #555;
}
.icp {
  color: #555;
  text-decoration: none;
  transition: color 0.2s;
}
.icp:hover {
  color: #888;
}
</style>
