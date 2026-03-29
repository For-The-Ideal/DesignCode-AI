<!-- components/auth/LoginModal.vue -->
<template>
  <Teleport to="body">
    <div v-if="visible" class="modal-overlay" @click.self="close">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isLogin ? '登录账户' : '注册账户' }}</h2>
          <button class="close-btn" @click="close">&times;</button>
        </div>
        
        <!-- 登录表单 -->
        <form v-if="isLogin" @submit.prevent="handleLogin">
          <input 
            v-model="loginForm.username" 
            type="text" 
            placeholder="用户名"
            class="input-field"
          />
          <input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="密码"
            class="input-field"
          />
          <button type="submit" class="submit-btn">登录</button>
          <p class="switch-text">
            还没有账户？ 
            <a @click="switchToRegister">立即注册</a>
          </p>
        </form>
        
        <!-- 注册表单 -->
        <form v-else @submit.prevent="handleRegister">
          <input 
            v-model="registerForm.username" 
            type="text" 
            placeholder="用户名"
            class="input-field"
          />
          <input 
            v-model="registerForm.email" 
            type="email" 
            placeholder="邮箱"
            class="input-field"
          />
          <input 
            v-model="registerForm.password" 
            type="password" 
            placeholder="密码"
            class="input-field"
          />
          <input 
            v-model="registerForm.confirm" 
            type="password" 
            placeholder="确认密码"
            class="input-field"
          />
          <button type="submit" class="submit-btn">注册</button>
          <p class="switch-text">
            已有账户？ 
            <a @click="switchToLogin">立即登录</a>
          </p>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'

const visible = ref(false)
const isLogin = ref(true)

const loginForm = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirm: ''
})

const open = () => {
  visible.value = true
  isLogin.value = true
  resetForms()
}

const close = () => {
  visible.value = false
}

const resetForms = () => {
  loginForm.username = ''
  loginForm.password = ''
  registerForm.username = ''
  registerForm.email = ''
  registerForm.password = ''
  registerForm.confirm = ''
}

const switchToRegister = () => {
  isLogin.value = false
  resetForms()
}

const switchToLogin = () => {
  isLogin.value = true
  resetForms()
}

const handleLogin = () => {
 
}

const handleRegister = () => {
 
}

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
  backdrop-filter: blur(12px);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: rgba(15, 25, 35, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 255, 255, 0.4);
  border-radius: 24px;
  padding: 32px;
  width: 400px;
  max-width: 90%;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h2 {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 28px;
  cursor: pointer;
  transition: color 0.3s;
}

.close-btn:hover {
  color: #fff;
}

.input-field {
  width: 100%;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  color: white;
  margin-bottom: 16px;
  transition: all 0.3s;
}

.input-field:focus {
  outline: none;
  border-color: #00ffff;
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.2);
}

.submit-btn {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  border-radius: 12px;
  color: white;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  margin-top: 8px;
}

.submit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 20px rgba(0, 255, 255, 0.3);
}

.switch-text {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #888;
}

.switch-text a {
  color: #00ffff;
  text-decoration: none;
  cursor: pointer;
}

.switch-text a:hover {
  text-decoration: underline;
}
</style>