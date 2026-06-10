<template>
  <form @submit.prevent="handleLogin" class="auth-form" v-border-gradient>
    <!-- 邮箱 -->
    <div class="input-field" :class="{ focused: focusedInput === 'email', error: errors.email }">
      <i class="fas fa-envelope input-icon"></i>
      <input type="email" v-model="form.email" placeholder="邮箱地址" maxlength="50"
        @focus="onFocus('email')" @blur="onBlur('email')" @input="clearError('email')">
      <svg class="border-svg" :class="{ animate: animatedInput === 'email' }">
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
      </svg>
    </div>
    <p v-if="errors.email" class="field-error">{{ errors.email }}</p>

    <!-- 密码 -->
    <div class="input-field" :class="{ focused: focusedInput === 'password', error: errors.password }">
      <i class="fas fa-lock input-icon"></i>
      <input :type="showPassword ? 'text' : 'password'" v-model="form.password" placeholder="密码" maxlength="20"
        @focus="onFocus('password')" @blur="onBlur('password')" @input="clearError('password')" @keyup.enter="handleLogin">
      <button type="button" class="password-eye" @click="showPassword = !showPassword">
        <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
      </button>
      <svg class="border-svg" :class="{ animate: animatedInput === 'password' }">
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
      </svg>
    </div>
    <p v-if="errors.password" class="field-error">{{ errors.password }}</p>

    <!-- 验证码 -->
    <div class="captcha-row">
      <div class="input-field captcha-input" :class="{ focused: focusedInput === 'captchaCode', error: errors.captchaCode }">
        <i class="fas fa-shield-alt input-icon"></i>
        <input type="text" v-model="form.captchaCode" placeholder="验证码" maxlength="4"
          @focus="onFocus('captchaCode')" @blur="onBlur('captchaCode')" @input="clearError('captchaCode')" @keyup.enter="handleLogin">
        <svg class="border-svg" :class="{ animate: animatedInput === 'captchaCode' }">
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
        </svg>
      </div>
      <div class="captcha-img-wrap" @click="fetchCaptcha" title="点击刷新验证码">
        <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-img" />
        <span v-else class="captcha-placeholder">
          <i v-if="captchaLoading" class="fas fa-spinner fa-spin"></i>
          <i v-else class="fas fa-sync-alt"></i>
        </span>
      </div>
    </div>
    <p v-if="errors.captchaCode" class="field-error">{{ errors.captchaCode }}</p>

    <!-- 记住我 / 忘记密码 -->
    <div class="form-options">
      <label class="checkbox">
        <input type="checkbox" v-model="rememberMe">
        <span class="checkmark"></span>
        <span>记住我</span>
      </label>
      <a href="#" class="forgot-link" @click.prevent="$emit('forgotPassword')">忘记密码？</a>
    </div>

    <button type="submit" class="submit-btn" :disabled="loading">
      <span v-if="!loading">登录</span>
      <span v-else class="loading-spinner"></span>
    </button>
  </form>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { loginApi } from '@/api/login'

const emit = defineEmits(['success', 'forgotPassword'])

const loading = ref(false)
const rememberMe = ref(false)
const showPassword = ref(false)
const focusedInput = ref(null)
const animatedInput = ref(null)

const form = reactive({ email: '', password: '', captchaCode: '' })
const errors = reactive({ email: '', password: '', captchaCode: '' })

const captchaId = ref('')
const captchaImage = ref('')
const captchaLoading = ref(false)

const clearError = (field) => { errors[field] = '' }
const onFocus = (id) => { focusedInput.value = id; animatedInput.value = id }
const onBlur = (id) => {
  if (focusedInput.value === id) { focusedInput.value = null; animatedInput.value = null }
}

const fetchCaptcha = async () => {
  captchaLoading.value = true
  try {
    const res = await loginApi.captcha()
    if (res.code === 200) {
      captchaId.value = res.data.captcha_id
      captchaImage.value = res.data.captcha_image
    }
  } catch (e) {
    console.error('验证码获取失败', e)
  } finally {
    captchaLoading.value = false
  }
}

const handleLogin = async () => {
  let valid = true
  if (!form.email)    { errors.email    = '请输入邮箱地址'; valid = false }
  if (!form.password)    { errors.password    = '请输入密码'; valid = false }
  if (!form.captchaCode) { errors.captchaCode = '请输入验证码'; valid = false }
  if (!valid) return

  loading.value = true
  try {
    const res = await loginApi.login({
      email: form.email,
      password: form.password,
      captcha_id: captchaId.value,
      captcha_code: form.captchaCode,
    })
    if (res.code === 200) {
      if (rememberMe.value) localStorage.setItem('user', JSON.stringify(res.data))
      emit('success', res.data)
    } else {
      fetchCaptcha()
      form.captchaCode = ''
    }
  } catch (e) {
    fetchCaptcha()
    form.captchaCode = ''
  } finally {
    loading.value = false
  }
}


</script>

<style scoped>
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-field {
  position: relative;
  width: 100%;
}
.input-icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
  font-size: 16px;
  transition: color 0.3s;
  z-index: 2;
}
.input-field input {
  width: 100%;
  padding: 14px 16px 14px 48px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.2);
  border-radius: 40px;
  color: #fff;
  font-size: 14px;
  outline: none;
  transition: all 0.3s;
  position: relative;
  z-index: 1;
}
.input-field.focused input {
  box-shadow: 0 0 0 1px rgba(0, 255, 255, 0.3);
}
.input-field.error input {
  border-color: rgba(255, 71, 87, 0.6);
  box-shadow: 0 0 0 1px rgba(255, 71, 87, 0.3);
}
.input-field.error .input-icon {
  color: #ff4757;
}
.input-field input:focus ~ .input-icon {
  color: #00ffff;
}

.field-error {
  color: #ff4757;
  font-size: 12px;
  margin: -14px 0 0 16px;
  animation: errorFadeIn 0.25s ease;
}
@keyframes errorFadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* 验证码 */
.captcha-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
}
.captcha-input {
  flex: 1;
}
.captcha-img-wrap {
  width: 130px;
  height: 48px;
  border-radius: 40px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid rgba(0, 255, 255, 0.2);
  transition: border-color 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.4);
  flex-shrink: 0;
}
.captcha-img-wrap:hover {
  border-color: rgba(0, 255, 255, 0.5);
}
.captcha-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.captcha-placeholder {
  color: #666;
  font-size: 20px;
}
.captcha-placeholder .fa-spin {
  color: #00ffff;
}

/* SVG 光带边框 */
.border-svg {
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
.border-svg.animate { opacity: 1; }
.border-bg {
  fill: none;
  stroke: rgba(0, 255, 255, 0.1);
  stroke-width: 1.5;
}
.moving-stroke {
  fill: none;
  stroke: url(#borderGradient);
  stroke-width: 2;
  stroke-linecap: round;
  filter: drop-shadow(0 0 4px rgba(0, 255, 255, 0.6));
  stroke-dasharray: 30 70;
  animation: moveStroke 3s linear infinite;
}
@keyframes moveStroke {
  from { stroke-dashoffset: 100; }
  to   { stroke-dashoffset: 0; }
}

.password-eye {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  background: transparent;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 16px;
  transition: color 0.3s;
  z-index: 2;
}
.password-eye:hover { color: #00ffff; }
.input-field.focused .password-eye { color: #00ffff; }

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  margin: 8px 0 4px;
}
.checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: #888;
}
.checkbox input { display: none; }
.checkmark {
  width: 16px;
  height: 16px;
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 4px;
  position: relative;
  transition: all 0.3s;
}
.checkbox input:checked + .checkmark {
  background: #00ffff;
  border-color: #00ffff;
}
.checkbox input:checked + .checkmark::after {
  content: '';
  position: absolute;
  left: 4px;
  top: 1px;
  width: 5px;
  height: 9px;
  border: solid #0a0a0f;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}
.forgot-link {
  color: #00ffff;
  text-decoration: none;
  font-size: 12px;
}
.forgot-link:hover { opacity: 0.8; }

.submit-btn {
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
  margin-top: 16px;
}
.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 255, 255, 0.3);
}
.submit-btn:disabled {
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
