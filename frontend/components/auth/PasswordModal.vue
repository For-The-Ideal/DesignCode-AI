<!-- components/auth/PasswordModal.vue -->
<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div v-if="visible" class="modal-overlay" @click.self="close">
        <div class="modal-container">
          <button class="close-btn" @click="close">
            <i class="fas fa-times"></i>
          </button>

          <div class="modal-icon">
            <div class="icon-ring">
              <i class="fas fa-key"></i>
            </div>
          </div>

          <h2 class="modal-title">修改密码</h2>
          <p class="modal-subtitle">请输入当前密码并设置新密码</p>

          <form @submit.prevent="handleSubmit" class="auth-form">
            <!-- 当前密码 -->
            <div class="input-field" :class="{ focused: focusedInput === 'old' }">
              <i class="fas fa-lock input-icon"></i>
              <input
                :type="showOld ? 'text' : 'password'"
                v-model="form.oldPassword"
                placeholder="当前密码"
                maxlength="20"
                @focus="focusedInput = 'old'"
                @blur="focusedInput = null"
              />
              <button type="button" class="password-eye" @click="showOld = !showOld">
                <i :class="showOld ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: focusedInput === 'old' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>

            <!-- 新密码 -->
            <div class="input-field" :class="{ focused: focusedInput === 'new1' }">
              <i class="fas fa-lock-open input-icon"></i>
              <input
                :type="showNew ? 'text' : 'password'"
                v-model="form.newPassword"
                placeholder="新密码"
                maxlength="20"
                @focus="focusedInput = 'new1'"
                @blur="focusedInput = null"
              />
              <button type="button" class="password-eye" @click="showNew = !showNew">
                <i :class="showNew ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: focusedInput === 'new1' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>

            <!-- 确认新密码 -->
            <div class="input-field" :class="{ focused: focusedInput === 'new2' }">
              <i class="fas fa-check-circle input-icon"></i>
              <input
                :type="showConfirm ? 'text' : 'password'"
                v-model="form.confirmPassword"
                placeholder="确认新密码"
                maxlength="20"
                @focus="focusedInput = 'new2'"
                @blur="focusedInput = null"
              />
              <button type="button" class="password-eye" @click="showConfirm = !showConfirm">
                <i :class="showConfirm ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: focusedInput === 'new2' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>

            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="!loading">确认修改</span>
              <span v-else class="loading-spinner"></span>
            </button>
          </form>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { ref, reactive } from 'vue'

const emit = defineEmits(['close'])

const visible = ref(false)
const loading = ref(false)
const focusedInput = ref(null)
const showOld = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const open = () => {
  visible.value = true
  resetForm()
}

const close = () => {
  visible.value = false
  emit('close')
}

const resetForm = () => {
  form.oldPassword = ''
  form.newPassword = ''
  form.confirmPassword = ''
  loading.value = false
  focusedInput.value = null
  showOld.value = false
  showNew.value = false
  showConfirm.value = false
}

const handleSubmit = async () => {
  if (!form.oldPassword) { alert('请输入当前密码'); return }
  if (!form.newPassword) { alert('请输入新密码'); return }
  if (form.newPassword.length < 6) { alert('新密码长度至少6位'); return }
  if (form.newPassword !== form.confirmPassword) { alert('两次输入的新密码不一致'); return }

  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 800))

  if (form.oldPassword !== '123456') {
    alert('当前密码错误')
    loading.value = false
    return
  }

  alert('密码修改成功！')
  close()
}

defineExpose({ open, close })
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0; left: 0;
  width: 100%; height: 100%;
  background: rgba(0,0,0,0.85);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  position: relative;
  width: 440px;
  max-width: 90%;
  background: rgba(12,20,28,0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0,255,255,0.3);
  border-radius: 32px;
  padding: 40px 32px;
  box-shadow: 0 25px 50px -12px rgba(0,0,0,0.5);
  animation: modalSlideUp 0.3s ease;
}

@keyframes modalSlideUp {
  from { opacity: 0; transform: translateY(30px) scale(0.96); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

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
  background: linear-gradient(135deg, rgba(255, 0, 255, 0.15), rgba(0, 255, 255, 0.15));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0, 255, 255, 0.3);
  box-shadow: 0 0 20px rgba(0, 255, 255, 0.08);
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

/* 表单 */
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
.input-field.focused .input-icon {
  color: #00ffff;
}

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

/* SVG 光带边框 */
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

/* 过渡动画 */
.modal-fade-enter-active { transition: all 0.3s ease; }
.modal-fade-leave-active { transition: all 0.25s ease; }
.modal-fade-enter-from,
.modal-fade-leave-to { opacity: 0; }
</style>
