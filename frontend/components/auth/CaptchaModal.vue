<!-- components/auth/CaptchaModal.vue -->
<template>
  <div class="captcha-form">
    <button class="back-btn" @click="$emit('back')" type="button">
      <i class="fas fa-arrow-left"></i>
      <span>返回登录</span>
    </button>

    <div class="modal-icon">
      <div class="icon-ring">
        <i class="fas fa-shield-alt"></i>
      </div>
    </div>

    <h3 class="captcha-title">安全验证</h3>
    <p class="captcha-subtitle">请输入图片中的验证码</p>

    <!-- 验证码图片 -->
    <div class="captcha-img-wrap" @click="refresh" title="点击刷新验证码">
      <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
      <span v-else class="captcha-placeholder">
        <i class="fas fa-spinner fa-spin"></i>
      </span>
    </div>

    <!-- 输入区 -->
    <div class="captcha-input-row">
      <input
        ref="codeInputRef"
        v-model="code"
        type="text"
        placeholder="请输入验证码"
        maxlength="4"
        @keyup.enter="handleConfirm"
      />
      <button type="button" class="refresh-btn" @click="refresh" title="刷新验证码">
        <i class="fas fa-sync-alt"></i>
      </button>
    </div>

    <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>

    <button class="confirm-btn" :disabled="!code || loading" @click="handleConfirm">
      <span v-if="!loading">确认验证</span>
      <span v-else class="loading-spinner"></span>
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { loginApi } from '@/api/login'

const emit = defineEmits(['confirm', 'back'])

const codeInputRef = ref(null)
const code = ref('')
const errorMsg = ref('')
const loading = ref(false)

const captchaId = ref('')
const captchaImage = ref('')

const fetchCaptcha = async () => {
  try {
    const res = await loginApi.captcha()
    if (res.code === 200) {
      captchaId.value = res.data.captcha_id
      captchaImage.value = res.data.captcha_image
    }
  } catch (e) {
    console.error('验证码获取失败', e)
    errorMsg.value = '验证码加载失败，请重试'
  }
}

const refresh = () => {
  code.value = ''
  errorMsg.value = ''
  captchaImage.value = ''
  fetchCaptcha()
}

const showError = (msg) => {
  errorMsg.value = msg
}

const resetLoading = () => {
  loading.value = false
}

const init = () => {
  code.value = ''
  errorMsg.value = ''
  loading.value = false
  fetchCaptcha()
  nextTick(() => codeInputRef.value?.focus())
}

const handleConfirm = () => {
  if (!code.value || loading.value) return
  loading.value = true
  errorMsg.value = ''
  emit('confirm', {
    captchaId: captchaId.value,
    captchaCode: code.value,
  })
}

onMounted(() => { init() })
defineExpose({ refresh, showError, resetLoading, init })
</script>

<style scoped>
.back-btn {
  position: absolute;
  top: 20px;
  left: 20px;
  background: transparent;
  border: none;
  color: #888;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: color 0.3s;
  padding: 0;
}
.back-btn:hover { color: #00ffff; }

.modal-icon {
  display: flex;
  justify-content: center;
  margin-top: 16px;
  margin-bottom: 20px;
}
.icon-ring {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.15));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #00ffff;
  animation: iconPulse 2s ease-in-out infinite;
}
@keyframes iconPulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(0, 255, 255, 0.3); }
  50% { box-shadow: 0 0 0 8px rgba(0, 255, 255, 0); }
}

.captcha-title {
  font-size: 24px;
  font-weight: 700;
  text-align: center;
  background: linear-gradient(135deg, #fff, #00ffff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 6px;
}
.captcha-subtitle {
  font-size: 13px;
  color: #888;
  text-align: center;
  margin-bottom: 24px;
}

.captcha-img-wrap {
  width: 100%;
  height: 60px;
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid rgba(0, 255, 255, 0.2);
  transition: border-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  margin-bottom: 20px;
}
.captcha-img-wrap:hover { border-color: rgba(0, 255, 255, 0.5); }
.captcha-img-wrap img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.captcha-placeholder {
  color: #666;
  font-size: 22px;
}
.captcha-placeholder .fa-spin { color: #00ffff; }

.captcha-input-row {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}
.captcha-input-row input {
  flex: 1;
  padding: 10px 16px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 28px;
  color: #fff;
  font-size: 16px;
  letter-spacing: 6px;
  text-align: center;
  outline: none;
  transition: all 0.3s;
}
.captcha-input-row input:focus {
  border-color: rgba(0, 255, 255, 0.5);
  box-shadow: 0 0 0 1px rgba(0, 255, 255, 0.3);
}
.captcha-input-row input::placeholder {
  color: #4b5563;
  letter-spacing: 0;
  font-size: 14px;
}
.refresh-btn {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 50%;
  color: #888;
  cursor: pointer;
  font-size: 16px;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.refresh-btn:hover {
  color: #00ffff;
  border-color: rgba(0, 255, 255, 0.5);
  transform: rotate(180deg);
}

.error-msg {
  color: #ff4757;
  font-size: 12px;
  text-align: center;
  margin: 0 0 12px;
  animation: errorFadeIn 0.25s ease;
}
@keyframes errorFadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}

.confirm-btn {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 40px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}
.confirm-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 255, 255, 0.3);
}
.confirm-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
