<template>
  <div class="reset-page">
    <div v-if="errorMsg" class="error-card">
      <div class="error-icon">
        <i class="fas fa-exclamation-triangle"></i>
      </div>
      <h2 class="error-title">重置链接无效</h2>
      <p class="error-desc">{{ errorMsg }}</p>
      <NuxtLink to="/" class="back-link">
        <i class="fas fa-arrow-left"></i>
        返回首页
      </NuxtLink>
    </div>

    <ResetPassword
      v-if="resetToken"
      ref="resetRef"
      :token="resetToken"
      @success="onResetSuccess"
      @close="onResetClose"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ResetPassword from '@/components/auth/ResetPassword.vue'

const route = useRoute()
const router = useRouter()

const resetToken = ref('')
const resetEmail = ref('')
const resetRef = ref(null)
const errorMsg = ref('')

onMounted(() => {
  const token = route.query.token
  const email = route.query.email

  if (!token) {
    errorMsg.value = '缺少重置令牌，请重新申请密码重置'
    return
  }

  resetToken.value = token
  resetEmail.value = email || ''

  // 自动打开重置弹窗
  setTimeout(() => {
    resetRef.value?.open()
  }, 300)
})

const onResetSuccess = () => {
  // 重置成功，跳转首页，AuthModal 会自动检测需要登录
  router.push('/')
}

const onResetClose = () => {
  // 用户关闭弹窗，跳转首页
  router.push('/')
}
</script>

<style scoped>
.reset-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 200px);
  padding: 40px 20px;
}
.error-card {
  text-align: center;
  max-width: 400px;
  padding: 48px 32px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 71, 87, 0.2);
  border-radius: 20px;
}
.error-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 71, 87, 0.1);
  border-radius: 50%;
  font-size: 24px;
  color: #ff4757;
}
.error-title {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
}
.error-desc {
  font-size: 14px;
  color: #888;
  margin-bottom: 24px;
  line-height: 1.6;
}
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #00ffff;
  text-decoration: none;
  font-size: 14px;
  transition: opacity 0.2s;
}
.back-link:hover { opacity: 0.8; }
</style>
