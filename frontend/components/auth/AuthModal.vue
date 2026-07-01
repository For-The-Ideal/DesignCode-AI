<!-- components/auth/AuthModal.vue -->
<template>
  <div v-border-gradient>
    <DialogModel ref="dialogRef" @close="emit('close')">
    <button v-if="isForgotPassword" class="back-btn" @click="goToLogin" type="button">
      <i class="fas fa-arrow-left"></i>
      <span>返回登录</span>
    </button>
    <button v-if="!isForgotPassword && !isCaptchaActive" class="close-btn" @click="close" type="button">
      <i class="fas fa-times"></i>
    </button>

    <div v-if="!isCaptchaActive" class="modal-icon">
      <div class="icon-ring">
        <i class="fas fa-microchip"></i>
      </div>
    </div>

    <h2 v-if="!isCaptchaActive" class="modal-title">
      {{ isForgotPassword ? '忘记密码' : isLogin ? '登录账户' : '注册账户' }}
    </h2>
    <p v-if="!isCaptchaActive" class="modal-subtitle">
      {{ isForgotPassword ? '输入邮箱地址，我们将发送重置链接' : isLogin ? '欢迎回来，登录继续你的创作' : '创建账户，开启 AI 设计转代码之旅' }}
    </p>

    <!-- 登录 / 注册 / 忘记密码 / 验证码 四态切换 -->
    <Login v-if="isLogin && !isCaptchaActive" @forgotPassword="goToForgotPassword" @requestCaptcha="onRequestCaptcha" />
    <CaptchaModal v-if="isCaptchaActive" ref="captchaModalRef" @confirm="onCaptchaConfirm" @back="onCaptchaBack" />
    <Register v-if="isRegister" @success="onRegisterSuccess" />
    <div v-if="isForgotPassword" class="forgot-form">
      <div class="forgot-input-group" :class="{ focused: forgotFocused }">
        <i class="fas fa-envelope input-icon"></i>
        <div class="email-wrap">
          <input
            ref="forgotInputRef"
            v-model="forgotEmail"
            type="email"
            placeholder="输入邮箱地址"
            maxlength="254"
            @focus="forgotFocused = true"
            @blur="forgotFocused = false"
          />
        </div>
        <span class="input-divider"></span>
        <select v-model="forgotProvider" class="provider-select" @change="changeProvider" @blur="forgotFocused = false">
          <option v-for="p in emailProviders" :key="p.value" :value="p.value">{{ p.label }}</option>
        </select>
        <svg class="border-svg" :class="{ animate: forgotFocused }">
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
        </svg>
      </div>
      <button type="button" class="submit-btn" :disabled="forgotLoading" @click="submitForgotPassword">
        <span v-if="!forgotLoading">发送重置链接</span>
        <span v-else class="loading-spinner"></span>
      </button>
      <p v-if="forgotMsg" :class="['forgot-msg', forgotMsgType]">{{ forgotMsg }}</p>
    </div>

    <div v-if="!isCaptchaActive && !isForgotPassword" class="switch-mode">
      <p class="switch-status">
        {{ isLogin ? '还没有账户？' : '已有账户？' }}
        <a href="#" @click.prevent="toggleMode">
          {{ isLogin ? '立即注册' : '立即登录' }}
        </a>
      </p>
    </div>
  </DialogModel>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import DialogModel from '@/components/dialog/DialogModel.vue'
import Login from './Login.vue'
import Register from './Register.vue'
import CaptchaModal from './CaptchaModal.vue'
import { loginApi } from '@/api/login'
import { useUserStore } from '@/stores/user'
import { validateEmail } from '@/utils/index'

const emit = defineEmits(['loginSuccess', 'registerSuccess', 'close'])

const dialogRef = ref(null)
const isLogin = ref(true)
const isRegister = ref(false)
const isForgotPassword = ref(false)
const isCaptchaActive = ref(false)
const captchaModalRef = ref(null)
const loginCreds = reactive({ email: '', password: '' })
const userStore = useUserStore()

// 忘记密码表单
const forgotEmail = ref('')
const forgotInputRef = ref(null)
const forgotFocused = ref(false)
const forgotLoading = ref(false)
const forgotProvider = ref('qq')
const forgotMsg = ref('')
const forgotMsgType = ref('')

const emailProviders = [
  { value: 'qq', domain: 'qq.com', label: 'QQ邮箱' },
  { value: '163', domain: '163.com', label: '163邮箱' },
  { value: 'gmail', domain: 'gmail.com', label: 'Gmail' },
]

const open = () => {
  dialogRef.value?.open()
  goToLogin()
}
const close = () => { dialogRef.value?.close() }

const goToLogin = () => {
  isLogin.value = true
  isRegister.value = false
  isForgotPassword.value = false
}
const goToRegister = () => {
  isLogin.value = false
  isRegister.value = true
  isForgotPassword.value = false
}
const goToForgotPassword = () => {
  isLogin.value = false
  isRegister.value = false
  isForgotPassword.value = true
}
const toggleMode = () => {
  if (isLogin.value) goToRegister()
  else goToLogin()
}

const onRequestCaptcha = ({ email, password }) => {
  loginCreds.email = email
  loginCreds.password = password
  isCaptchaActive.value = true
}

const onCaptchaBack = () => {
  isCaptchaActive.value = false
}

const changeProvider = () => {
  forgotEmail.value = ''
  forgotFocused.value = true
  forgotInputRef.value?.focus()
}

const onCaptchaConfirm = async ({ captchaId, captchaCode }) => {
  try {
    const res = await loginApi.login({
      email: loginCreds.email.toLowerCase(),
      password: loginCreds.password,
      captcha_id: captchaId,
      captcha_code: captchaCode,
    })
    if(res.code !=200){
      captchaModalRef.value?.refresh()
      captchaModalRef.value?.showError(res.message || '验证失败，请重试')
      return
    }
    isCaptchaActive.value = false
    emit('loginSuccess', res.data)
    close()
  } catch (e) {
    captchaModalRef.value?.showError('网络错误，请重试')
    captchaModalRef.value?.refresh()
  } finally {
    captchaModalRef.value?.resetLoading()
  }
}

const onRegisterSuccess = (data) => {
  goToLogin()
}
const submitForgotPassword = async () => {
  const email = forgotEmail.value.trim()
  if (!email) return

  const check = validateEmail(email.toLowerCase())
  if (!check.valid) {
    forgotMsg.value = check.message
    forgotMsgType.value = 'error'
    return
  }

  // 校验邮箱域名是否匹配所选服务商
  const domain = email.split('@')[1]?.toLowerCase()
  const matched = emailProviders.find(p => domain?.includes(p.domain))
  if (matched && matched.value !== forgotProvider.value) {
    forgotMsg.value = `邮箱与所选服务商不匹配，请选择${matched.label}或修改邮箱`
    forgotMsgType.value = 'error'
    return
  }

  forgotLoading.value = true
  forgotMsg.value = ''
  try {
    const res = await loginApi.forgotPassword({ email: email.toLowerCase() })
    if (res.code != 200) {
      forgotMsg.value = res.message || '发送失败，请重试'
      forgotMsgType.value = 'error'
      return
    }
    forgotMsg.value = '重置链接已发送，请检查您的邮箱'
    forgotMsgType.value = 'success'
  } catch (e) {
    forgotMsg.value = '网络错误，请重试'
    forgotMsgType.value = 'error'
  } finally {
    forgotLoading.value = false
  }
}

defineExpose({ open, close })
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
  z-index: 2;
}
.back-btn:hover { color: #00ffff; }

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 50%;
  color: #888;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
}
.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  transform: rotate(90deg);
}

.modal-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
}
.icon-ring {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.15));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0, 255, 255, 0.3);
  animation: pulse 2s ease-in-out infinite;
}
.icon-ring i {
  font-size: 28px;
  color: #00ffff;
}
@keyframes pulse {
  0%, 100% { box-shadow: 0 0 0 0 rgba(0, 255, 255, 0.2); }
  50% { box-shadow: 0 0 0 8px rgba(0, 255, 255, 0); }
}

.modal-title {
  font-size: 28px;
  font-weight: 700;
  text-align: center;
  background: linear-gradient(135deg, #fff, #00ffff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  margin-bottom: 8px;
}
.modal-subtitle {
  font-size: 13px;
  color: #888;
  text-align: center;
  margin-bottom: 32px;
}

.divider {
  text-align: center;
  margin: 24px 0 20px;
  position: relative;
}
.divider::before,
.divider::after {
  content: '';
  position: absolute;
  top: 50%;
  width: calc(50% - 30px);
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(0, 255, 255, 0.3), transparent);
}
.divider::before { left: 0; }
.divider::after { right: 0; }
.divider span {
  background: rgba(12, 20, 28, 0.8);
  padding: 0 16px;
  font-size: 12px;
  color: #666;
  position: relative;
  z-index: 1;
}

.social-btn {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  color: #aaa;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 20px;
}
.social-btn:hover {
  transform: translateY(-2px);
  background: rgba(0, 255, 255, 0.1);
  color: #00ffff;
  border-color: rgba(0, 255, 255, 0.3);
}

.switch-mode {
  text-align: center;
  font-size: 13px;
  color: #888;
}
.switch-status{
  margin-top: 20px;
}
.switch-mode a {
  color: #00ffff;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

/* 忘记密码表单 */
.forgot-form {
  margin-bottom: 24px;
}
/* 忘记密码 — 账号+服务商输入组 */
.forgot-input-group {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 28px;
  padding: 0 18px;
  height: 52px;
  transition: all 0.3s;
  margin-bottom: 18px;
}
.forgot-input-group.focused {
  background: rgba(255, 255, 255, 0.08);
}
.forgot-input-group .email-wrap {
  flex: 1;
  min-width: 0;
}
.forgot-input-group .email-wrap input {
  width: 100%;
  background: transparent;
  border: none;
  outline: none;
  color: #e5e7eb;
  font-size: 15px;
  padding: 0;
}
.forgot-input-group .email-wrap input::placeholder {
  color: #4b5563;
}
.forgot-input-group .input-divider {
  width: 1px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  margin: 0 10px;
  flex-shrink: 0;
}
.forgot-input-group .provider-select {
  background: transparent;
  border: none;
  outline: none;
  color: #00ffff;
  font-size: 13px;
  cursor: pointer;
  flex-shrink: 0;
  padding: 4px 0;
}
.forgot-input-group .provider-select option {
  background: #1a1a2e;
  color: #e5e7eb;
}
.forgot-input-group .input-icon {
  color: #6b7280;
  font-size: 16px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
  margin-right: 8px;
}
.forgot-input-group.focused .input-icon {
  color: #00ffff;
}

/* 忘记密码 — 提示消息 */
.forgot-msg {
  text-align: center;
  font-size: 13px;
  margin-top: 12px;
  padding: 8px 12px;
  border-radius: 8px;
}
.forgot-msg.error {
  color: #ff4757;
  background: rgba(255, 71, 87, 0.08);
}
.forgot-msg.success {
  color: #00ff88;
  background: rgba(0, 255, 136, 0.08);
}

.forgot-form .submit-btn {
  width: 100%;
  height: 50px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.2), rgba(255, 0, 255, 0.2));
  border: 1px solid rgba(0, 255, 255, 0.4);
  border-radius: 28px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
}
.forgot-form .submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.3), rgba(255, 0, 255, 0.3));
  box-shadow: 0 0 24px rgba(0, 255, 255, 0.25);
  transform: translateY(-1px);
}
.forgot-form .submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 忘记密码 SVG 光带边框 */
.forgot-form .border-svg {
  position: absolute;
  top: -2px;
  left: -2px;
  width: calc(100% + 4px);
  height: calc(100% + 4px);
  pointer-events: none;
  z-index: 0;
  opacity: 0;
  transition: opacity 0.2s;
  overflow: visible;
}
.forgot-form .border-svg.animate { opacity: 1; }
.forgot-form .border-bg {
  fill: none;
  stroke: rgba(0, 255, 255, 0.1);
  stroke-width: 1.5;
}
.forgot-form .moving-stroke {
  fill: none;
  stroke: url(#borderGradient);
  stroke-width: 2;
  stroke-linecap: round;
  filter: drop-shadow(0 0 4px rgba(0, 255, 255, 0.6));
  stroke-dasharray: 30 70;
  animation: forgotStrokeMove 3s linear infinite;
}
@keyframes forgotStrokeMove {
  from { stroke-dashoffset: 100; }
  to   { stroke-dashoffset: 0; }
}

.loading-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #00ffff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>