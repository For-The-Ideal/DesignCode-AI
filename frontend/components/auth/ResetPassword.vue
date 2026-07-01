<template>
  <DialogModel ref="dialogRef" width="440px" v-border-gradient>
    <button class="close-btn" @click="handleClose" type="button">
      <i class="fas fa-times"></i>
    </button>

    <div class="modal-icon">
      <div class="icon-ring">
        <i class="fas fa-key"></i>
      </div>
    </div>

    <h2 class="modal-title">设置新密码</h2>
    <p class="modal-subtitle">请为您的账号设置一个新密码</p>

    <form @submit.prevent="handleSubmit" class="auth-form">
      <!-- 新密码 -->
      <div class="input-field" :class="{ focused: focusedInput === 'new', error: errors.password }">
        <i class="fas fa-lock-open input-icon"></i>
        <input
          :type="showPassword ? 'text' : 'password'"
          v-model="form.password"
          placeholder="新密码（至少6位）"
          maxlength="20"
          minlength="6"
          @focus="focusedInput = 'new'"
          @blur="focusedInput = null"
          @input="clearError('password')"
        />
        <button type="button" class="password-eye" @click="showPassword = !showPassword">
          <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
        </button>
        <svg class="border-svg" :class="{ animate: focusedInput === 'new' }">
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
        </svg>
      </div>
      <p v-if="errors.password" class="field-error">{{ errors.password }}</p>

      <!-- 确认密码 -->
      <div class="input-field" :class="{ focused: focusedInput === 'confirm', error: errors.confirm }">
        <i class="fas fa-check-circle input-icon"></i>
        <input
          :type="showConfirm ? 'text' : 'password'"
          v-model="form.confirm"
          placeholder="确认新密码"
          maxlength="20"
          minlength="6"
          @focus="focusedInput = 'confirm'"
          @blur="focusedInput = null"
          @input="clearError('confirm')"
        />
        <button type="button" class="password-eye" @click="showConfirm = !showConfirm">
          <i :class="showConfirm ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
        </button>
        <svg class="border-svg" :class="{ animate: focusedInput === 'confirm' }">
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
          <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
        </svg>
      </div>
      <p v-if="errors.confirm" class="field-error">{{ errors.confirm }}</p>

      <!-- 服务端错误 -->
      <p v-if="serverError" class="server-error">{{ serverError }}</p>

      <button type="submit" class="submit-btn" :disabled="loading">
        <span v-if="!loading">重置密码</span>
        <span v-else class="loading-spinner"></span>
      </button>
    </form>
  </DialogModel>
</template>

<script setup>
import { ref, reactive } from 'vue'
import DialogModel from '@/components/dialog/DialogModel.vue'
import { loginApi } from '@/api/login'
import { ElMessage } from 'element-plus'

const props = defineProps({
  token: { type: String, default: '' }
})

const emit = defineEmits(['close', 'success'])

const dialogRef = ref(null)
const loading = ref(false)
const focusedInput = ref(null)
const showPassword = ref(false)
const showConfirm = ref(false)
const serverError = ref('')

const form = reactive({
  password: '',
  confirm: ''
})

const errors = reactive({
  password: '',
  confirm: ''
})

const clearError = (field) => { errors[field] = '' }

const open = () => {
  dialogRef.value?.open()
  resetForm()
}

const close = () => {
  dialogRef.value?.close()
  emit('close')
}

const handleClose = () => {
  dialogRef.value?.close()
  emit('close')
}

const resetForm = () => {
  form.password = ''
  form.confirm = ''
  errors.password = ''
  errors.confirm = ''
  serverError.value = ''
  loading.value = false
  focusedInput.value = null
  showPassword.value = false
  showConfirm.value = false
}

const handleSubmit = async () => {
  let valid = true
  if (!form.password)          { errors.password = '请输入新密码'; valid = false }
  if (!form.confirm)           { errors.confirm  = '请确认新密码'; valid = false }
  if (!valid) return
  if (form.password.length < 6) { errors.password = '密码长度至少6位'; return }
  if (form.password !== form.confirm) { errors.confirm = '两次输入的密码不一致'; return }

  loading.value = true
  serverError.value = ''
  try {
    const res = await loginApi.resetPassword({
      token: props.token,
      password: form.password,
    })
    if(res.code != 200){
      ElMessage.error(res.message || '重置失败')
      serverError.value = res.message || '重置失败'
      return
    }
      ElMessage.success('密码重置成功')
      close()
      emit('success')
  } catch (e) {
    serverError.value = '网络错误，请重试'
  } finally {
    loading.value = false
  }
}

defineExpose({ open, close })
</script>

<style scoped>
.close-btn {
  position: absolute;
  top: 20px; right: 20px;
  width: 32px; height: 32px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.15);
  border-radius: 50%;
  color: #888;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
}
.close-btn:hover {
  background: rgba(255,255,255,0.1);
  color: #fff;
  transform: rotate(90deg);
}
.modal-icon {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}
.icon-ring {
  width: 56px; height: 56px;
  background: linear-gradient(135deg, rgba(255,0,255,0.15), rgba(0,255,255,0.15));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0,255,255,0.3);
  box-shadow: 0 0 20px rgba(0,255,255,0.08);
}
.icon-ring i {
  font-size: 24px;
  color: #ff00ff;
}
.modal-title {
  font-size: 26px;
  font-weight: 700;
  text-align: center;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}
.modal-subtitle {
  text-align: center;
  color: #6b7280;
  font-size: 14px;
  margin-bottom: 28px;
}
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.input-field {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255,255,255,0.04);
  border-radius: 28px;
  padding: 0 18px;
  height: 52px;
  transition: all 0.3s;
}
.input-field.focused {
  background: rgba(255,255,255,0.08);
}
.input-field input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e5e7eb;
  font-size: 15px;
  padding: 0 8px;
}
.input-field input::placeholder { color: #4b5563; }
.input-icon {
  color: #6b7280;
  font-size: 16px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}
.input-field.focused .input-icon { color: #00ffff; }
.password-eye {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
  font-size: 14px;
  transition: color 0.2s;
  flex-shrink: 0;
}
.password-eye:hover { color: #00ffff; }
.field-error {
  color: #ff4757;
  font-size: 12px;
  margin: -14px 0 0 18px;
  animation: errorFadeIn 0.25s ease;
}
.server-error {
  text-align: center;
  color: #ff4757;
  font-size: 13px;
  padding: 8px;
  background: rgba(255,71,87,0.08);
  border-radius: 8px;
}
@keyframes errorFadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to   { opacity: 1; transform: translateY(0); }
}
.border-svg {
  position: absolute;
  top: 0; left: 0;
  width: 100%; height: 100%;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s;
}
.input-field.focused .border-svg,
.border-svg.animate { opacity: 1; }
.border-bg {
  fill: none;
  stroke: rgba(255,255,255,0.08);
  stroke-width: 1;
}
.moving-stroke {
  fill: none;
  stroke: url(#borderGradient);
  stroke-width: 1.5;
  stroke-dasharray: 30 70;
  stroke-dashoffset: 0;
  animation: moveBorder 2.5s linear infinite;
}
@keyframes moveBorder {
  to { stroke-dashoffset: -100; }
}
.submit-btn {
  height: 50px;
  background: linear-gradient(135deg, rgba(0,255,255,0.2), rgba(255,0,255,0.2));
  border: 1px solid rgba(0,255,255,0.4);
  border-radius: 28px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 6px;
}
.submit-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, rgba(0,255,255,0.3), rgba(255,0,255,0.3));
  box-shadow: 0 0 24px rgba(0,255,255,0.25);
  transform: translateY(-1px);
}
.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.loading-spinner {
  display: inline-block;
  width: 20px; height: 20px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #00ffff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>