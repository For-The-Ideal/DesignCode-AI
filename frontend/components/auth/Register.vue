<template>
  <form @submit.prevent="handleRegister" class="auth-form" v-border-gradient>
    <!-- 用户名 -->
    <div class="input-field" :class="{ focused: focusedInput === 'username', error: errors.username }">
      <i class="fas fa-user input-icon"></i>
      <input type="text" v-model="form.username" placeholder="用户名" maxlength="20"
        @focus="onFocus('username')" @blur="onBlur('username')" @input="clearError('username')">
      <svg class="border-svg" :class="{ animate: animatedInput === 'username' }">
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
      </svg>
    </div>
    <p v-if="errors.username" class="field-error">{{ errors.username }}</p>

    <!-- 邮箱 -->
    <div class="input-field" :class="{ focused: focusedInput === 'email', error: errors.email }">
      <i class="fas fa-envelope input-icon"></i>
      <input type="email" v-model="form.email" placeholder="邮箱" maxlength="20"
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
        @focus="onFocus('password')" @blur="onBlur('password')" @input="clearError('password')">
      <button type="button" class="password-eye" @click="showPassword = !showPassword">
        <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
      </button>
      <svg class="border-svg" :class="{ animate: animatedInput === 'password' }">
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
      </svg>
    </div>
    <p v-if="errors.password" class="field-error">{{ errors.password }}</p>

    <!-- 确认密码 -->
    <div class="input-field" :class="{ focused: focusedInput === 'confirm', error: errors.confirm }">
      <i class="fas fa-check-circle input-icon"></i>
      <input :type="showConfirm ? 'text' : 'password'" v-model="form.confirm" placeholder="确认密码" maxlength="20"
        @focus="onFocus('confirm')" @blur="onBlur('confirm')" @input="clearError('confirm')">
      <button type="button" class="password-eye" @click="showConfirm = !showConfirm">
        <i :class="showConfirm ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
      </button>
      <svg class="border-svg" :class="{ animate: animatedInput === 'confirm' }">
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
        <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
      </svg>
    </div>
    <p v-if="errors.confirm" class="field-error">{{ errors.confirm }}</p>

    <button type="submit" class="submit-btn" :disabled="loading">
      <span v-if="!loading">注册</span>
      <span v-else class="loading-spinner"></span>
    </button>
  </form>
</template>

<script setup>
import { ref, reactive } from 'vue'

const emit = defineEmits(['success'])

const loading = ref(false)
const showPassword = ref(false)
const showConfirm = ref(false)
const focusedInput = ref(null)
const animatedInput = ref(null)

const form = reactive({ username: '', email: '', password: '', confirm: '' })
const errors = reactive({ username: '', email: '', password: '', confirm: '' })

const clearError = (field) => { errors[field] = '' }
const onFocus = (id) => { focusedInput.value = id; animatedInput.value = id }
const onBlur = (id) => {
  if (focusedInput.value === id) { focusedInput.value = null; animatedInput.value = null }
}

const handleRegister = async () => {
  let valid = true
  if (!form.username)  { errors.username  = '请输入用户名'; valid = false }
  if (!form.email)     { errors.email     = '请输入邮箱'; valid = false }
  if (!form.password)  { errors.password  = '请输入密码'; valid = false }
  if (!form.confirm)   { errors.confirm   = '请确认密码'; valid = false }
  if (!valid) return
  if (form.password !== form.confirm) {
    errors.confirm = '两次输入的密码不一致'; return
  }
  if (form.password.length < 6) {
    errors.password = '密码长度至少6位'; return
  }

  loading.value = true
  // TODO: 对接真实注册 API
  await new Promise(resolve => setTimeout(resolve, 800))
  alert('注册成功！请登录')
  emit('success')
  loading.value = false
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
