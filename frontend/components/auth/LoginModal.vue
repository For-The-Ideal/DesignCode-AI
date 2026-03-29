<!-- components/auth/LoginModal.vue -->
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
              <i class="fas fa-microchip"></i>
            </div>
          </div>
          
          <h2 class="modal-title">{{ isLogin ? '登录账户' : '注册账户' }}</h2>
          <p class="modal-subtitle">
            {{ isLogin ? '欢迎回来，登录继续你的创作' : '创建账户，开启 AI 设计转代码之旅' }}
          </p>
          
          <!-- 登录表单 -->
          <form v-if="isLogin" @submit.prevent="handleLogin" class="auth-form">
            <div class="input-field" :class="{ focused: focusedInput === 'login-username' }">
              <i class="fas fa-user input-icon"></i>
              <input 
                type="text" 
                v-model="loginForm.username" 
                placeholder="用户名"
                maxlength="20"
                @focus="handleFocus('login-username')"
                @blur="handleBlur('login-username')"
              >
              <!-- SVG 光带边框层 -->
              <svg class="border-svg" :class="{ animate: animatedInput === 'login-username' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <div class="input-field" :class="{ focused: focusedInput === 'login-password' }">
              <i class="fas fa-lock input-icon"></i>
              <input 
                :type="showLoginPassword ? 'text' : 'password'" 
                v-model="loginForm.password" 
                placeholder="密码"
                maxlength="20"
                @focus="handleFocus('login-password')"
                @blur="handleBlur('login-password')"
                @keyup.enter="handleLogin"
              >
              <button type="button" class="password-eye" @click="showLoginPassword = !showLoginPassword">
                <i :class="showLoginPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: animatedInput === 'login-password' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <div class="form-options">
              <label class="checkbox">
                <input type="checkbox" v-model="rememberMe">
                <span class="checkmark"></span>
                <span>记住我</span>
              </label>
              <a href="#" class="forgot-link">忘记密码？</a>
            </div>
            
            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="!loading">登录</span>
              <span v-else class="loading-spinner"></span>
            </button>
          </form>
          
          <!-- 注册表单 -->
          <form v-else @submit.prevent="handleRegister" class="auth-form">
            <div class="input-field" :class="{ focused: focusedInput === 'reg-username' }">
              <i class="fas fa-user input-icon"></i>
              <input 
                type="text" 
                v-model="registerForm.username" 
                placeholder="用户名"
                maxlength="20"
                @focus="handleFocus('reg-username')"
                @blur="handleBlur('reg-username')"
              >
              <svg class="border-svg" :class="{ animate: animatedInput === 'reg-username' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <div class="input-field" :class="{ focused: focusedInput === 'reg-email' }">
              <i class="fas fa-envelope input-icon"></i>
              <input 
                type="email" 
                v-model="registerForm.email" 
                placeholder="邮箱"
                maxlength="20"
                @focus="handleFocus('reg-email')"
                @blur="handleBlur('reg-email')"
              >
              <svg class="border-svg" :class="{ animate: animatedInput === 'reg-email' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <div class="input-field" :class="{ focused: focusedInput === 'reg-password' }">
              <i class="fas fa-lock input-icon"></i>
              <input 
                :type="showRegisterPassword ? 'text' : 'password'" 
                v-model="registerForm.password" 
                placeholder="密码"
                maxlength="20"
                @focus="handleFocus('reg-password')"
                @blur="handleBlur('reg-password')"
              >
              <button type="button" class="password-eye" @click="showRegisterPassword = !showRegisterPassword">
                <i :class="showRegisterPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: animatedInput === 'reg-password' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <div class="input-field" :class="{ focused: focusedInput === 'reg-confirm' }">
              <i class="fas fa-check-circle input-icon"></i>
              <input 
                :type="showRegisterConfirm ? 'text' : 'password'" 
                v-model="registerForm.confirm" 
                placeholder="确认密码"
                maxlength="20"
                @focus="handleFocus('reg-confirm')"
                @blur="handleBlur('reg-confirm')"
              >
              <button type="button" class="password-eye" @click="showRegisterConfirm = !showRegisterConfirm">
                <i :class="showRegisterConfirm ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
              </button>
              <svg class="border-svg" :class="{ animate: animatedInput === 'reg-confirm' }">
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="border-bg" />
                <rect x="0" y="0" rx="26" ry="26" width="100%" height="100%" pathLength="100" class="moving-stroke" />
              </svg>
            </div>
            
            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="!loading">注册</span>
              <span v-else class="loading-spinner"></span>
            </button>
          </form>
          
          <div class="divider"><span>或</span></div>
          
          <div class="social-login">
            <button class="social-btn" @click="socialLogin('github')"><i class="fab fa-github"></i></button>
            <button class="social-btn" @click="socialLogin('google')"><i class="fab fa-google"></i></button>
            <button class="social-btn" @click="socialLogin('discord')"><i class="fab fa-discord"></i></button>
          </div>
          
          <div class="switch-mode">
            <p>
              {{ isLogin ? '还没有账户？' : '已有账户？' }}
              <a href="#" @click.prevent="toggleMode">
                {{ isLogin ? '立即注册' : '立即登录' }}
              </a>
            </p>
          </div>
          
          <div class="demo-hint">
            <i class="fas fa-info-circle"></i>
            <span>演示账号: demo / 123456</span>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const emit = defineEmits(['login-success', 'register-success', 'close'])

// 状态
const visible = ref(false)
const isLogin = ref(true)
const loading = ref(false)
const rememberMe = ref(false)
const focusedInput = ref(null)
const animatedInput = ref(null)

const showLoginPassword = ref(false)
const showRegisterPassword = ref(false)
const showRegisterConfirm = ref(false)

// 表单数据
const loginForm = reactive({ username: '', password: '' })
const registerForm = reactive({ username: '', email: '', password: '', confirm: '' })

// 聚焦处理：记录聚焦输入框，并启动边框动画
const handleFocus = (inputId) => {
  focusedInput.value = inputId
  animatedInput.value = inputId
}

// 失焦处理：清除聚焦和动画状态
const handleBlur = (inputId) => {
  if (focusedInput.value === inputId) {
    focusedInput.value = null
    animatedInput.value = null
  }
}

// 弹窗控制
const open = () => {
  visible.value = true
  resetForms()
}
const close = () => {
  visible.value = false
  emit('close')
}
const resetForms = () => {
  loginForm.username = ''
  loginForm.password = ''
  registerForm.username = ''
  registerForm.email = ''
  registerForm.password = ''
  registerForm.confirm = ''
  loading.value = false
  focusedInput.value = null
  animatedInput.value = null
  showLoginPassword.value = false
  showRegisterPassword.value = false
  showRegisterConfirm.value = false
}
const toggleMode = () => {
  isLogin.value = !isLogin.value
  resetForms()
}

// 登录/注册逻辑
const handleLogin = async () => {
  if (!loginForm.username || !loginForm.password) {
    alert('请输入用户名和密码')
    return
  }
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 800))
  if (loginForm.username === 'demo' && loginForm.password === '123456') {
    const user = { username: loginForm.username, email: 'demo@example.com' }
    if (rememberMe.value) localStorage.setItem('user', JSON.stringify(user))
    emit('login-success', user)
    close()
  } else {
    alert('用户名或密码错误 (演示账号: demo / 123456)')
  }
  loading.value = false
}

const handleRegister = async () => {
  if (!registerForm.username || !registerForm.email || !registerForm.password) {
    alert('请填写完整信息'); return
  }
  if (registerForm.password !== registerForm.confirm) {
    alert('两次输入的密码不一致'); return
  }
  if (registerForm.password.length < 6) {
    alert('密码长度至少6位'); return
  }
  loading.value = true
  await new Promise(resolve => setTimeout(resolve, 800))
  alert('注册成功！请登录')
  isLogin.value = true
  resetForms()
  loading.value = false
}

const socialLogin = (provider) => {
  alert(`${provider.toUpperCase()} 登录功能开发中...`)
}

// 注入全局 SVG 渐变定义（确保所有边框动画能使用渐变）
onMounted(() => {
  if (!document.getElementById('global-gradient-def')) {
    const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg')
    svg.setAttribute('style', 'position:absolute; width:0; height:0')
    svg.id = 'global-gradient-def'
    const defs = document.createElementNS('http://www.w3.org/2000/svg', 'defs')
    const grad = document.createElementNS('http://www.w3.org/2000/svg', 'linearGradient')
    grad.setAttribute('id', 'borderGradient')
    grad.setAttribute('x1', '0%')
    grad.setAttribute('y1', '0%')
    grad.setAttribute('x2', '100%')
    grad.setAttribute('y2', '0%')
    const stops = [
      { offset: '0%', color: '#00ffff' },
      { offset: '50%', color: '#ff00ff' },
      { offset: '100%', color: '#00ffff' }
    ]
    stops.forEach(s => {
      const stop = document.createElementNS('http://www.w3.org/2000/svg', 'stop')
      stop.setAttribute('offset', s.offset)
      stop.setAttribute('stop-color', s.color)
      grad.appendChild(stop)
    })
    defs.appendChild(grad)
    svg.appendChild(defs)
    document.body.appendChild(svg)
  }
})

defineExpose({ open, close })
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(16px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-container {
  position: relative;
  width: 460px;
  max-width: 90%;
  background: rgba(12, 20, 28, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 32px;
  padding: 40px 32px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  animation: modalSlideUp 0.3s ease;
}

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

/* ========== 输入框样式 ========== */
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

/* 聚焦时边框高亮 */
.input-field.focused input {
  box-shadow: 0 0 0 1px rgba(0, 255, 255, 0.3);
}

.input-field input:focus ~ .input-icon {
  color: #00ffff;
}

/* ========== SVG 光带移动动画 ========== */
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

.border-svg.animate {
  opacity: 1;
}

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
  /* 使用 pathLength="100" 后的比例：30% 长度，70% 间隙 */
  stroke-dasharray: 30 70; 
  animation: moveStroke 3s linear infinite;
}

@keyframes moveStroke {
  from {
    stroke-dashoffset: 100;
  }
  to {
    stroke-dashoffset: 0;
  }
}

/* 密码眼睛按钮 */
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

.password-eye:hover {
  color: #00ffff;
}

.input-field.focused .password-eye {
  color: #00ffff;
}

/* 表单选项 */
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

.checkbox input {
  display: none;
}

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

.forgot-link:hover {
  opacity: 0.8;
}

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

.social-login {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 24px;
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

.switch-mode a {
  color: #00ffff;
  text-decoration: none;
  font-weight: 500;
  margin-left: 4px;
}

.demo-hint {
  margin-top: 20px;
  padding: 12px;
  background: rgba(0, 255, 255, 0.05);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 11px;
  color: #666;
  border: 1px solid rgba(0, 255, 255, 0.15);
}

.demo-hint i {
  color: #00ffff;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}

@keyframes modalSlideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>