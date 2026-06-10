<!-- components/layout/AppHeader.vue -->
<template>
  <header class="header">
    <div class="header-content">
      <!-- Logo 区域 -->
      <div class="logo">
        <div class="logo-icon">
          <svg class="icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 7L12 12L22 7L12 2Z" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div class="logo-text">
          <h1>DesignCode AI</h1>
          <p>NEURAL CODE ENGINE</p>
        </div>
      </div>

      <!-- 导航菜单 -->
      <nav class="nav">
        <a href="javascript:void(0)" class="nav-link" :class="{'active': route.path === item.link}" v-for="item in navList" :key="item.name" 
        @click="handleNavClick(item.link)">{{item.name}}</a>
      </nav>

      <!-- 右侧区域 -->
      <div class="auth-area">
        <!-- 未登录：登录按钮 -->
        <button v-if="!isLogin" class="btn-login" @click="loginModalRef.open()">
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M20 21V19C20 16.8 18.2 15 16 15H8C5.8 15 4 16.8 4 19V21" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            <path d="M12 11C14.2091 11 16 9.20914 16 7C16 4.79086 14.2091 3 12 3C9.79086 3 8 4.79086 8 7C8 9.20914 9.79086 11 12 11Z" stroke="currentColor" stroke-width="1.5"/>
          </svg>
          登录
        </button>

        <!-- 已登录：头像 + 下拉菜单 -->
        <div v-else class="user-menu" @mouseenter="showDropdown = true" @mouseleave="showDropdown = false">
          <div class="user-avatar">
            <span class="avatar-text">{{ userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}</span>
          </div>
          <span class="user-name">{{ userInfo?.username }}</span>

          <!-- 下拉菜单 -->
          <transition name="dropdown-fade">
            <div v-if="showDropdown" class="dropdown-panel">
              <div class="dropdown-header">
                <div class="dropdown-avatar">
                  <span>{{ userInfo?.username?.charAt(0)?.toUpperCase() || 'U' }}</span>
                </div>
                <div class="dropdown-user-info">
                  <span class="dropdown-username">{{ userInfo?.username }}</span>
                  <span class="dropdown-email">{{ userInfo?.email }}</span>
                </div>
              </div>
              <div class="dropdown-divider"></div>
              <button class="dropdown-item" @click.stop="openPasswordModal">
                <i class="fas fa-key"></i>
                <span>修改密码</span>
              </button>
              <button class="dropdown-item logout" @click.stop="handleLogout">
                <i class="fas fa-sign-out-alt"></i>
                <span>退出登录</span>
              </button>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <LoginModal 
      ref="loginModalRef" 
      @login-success="handleLoginSuccess"
      @register-success="handleRegisterSuccess"
      @close="handleModalClose"
    />

    <PasswordModal ref="passwordModalRef" />

  </header>
</template>

<script setup>
import { computed } from "vue"
import LoginModal from '~/components/auth/LoginModal.vue'
import PasswordModal from '~/components/auth/PasswordModal.vue'
import { useUserStore } from '~/stores/user'
import { storeToRefs } from 'pinia'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loginModalRef = ref()
const passwordModalRef = ref()
const showDropdown = ref(false)
const { isLogin, userInfo } = storeToRefs(userStore)
const navList = ref([
  {
    name:"首页",
    link:"/",
  },
  {
    name:"代码生成",
    link:"/code",
  },
  // {
  //   name:"模型对比",
  //   link:"/compare",
  // },
  {
    name:"历史记录",
    link:"/history",
  }
])

const handleLoginSuccess = async(user) => {
  await userStore.login({ username: user.username, password: '' })
}

const handleRegisterSuccess = () => {}

const handleModalClose = () => {}

const handleLogout = () => {
  showDropdown.value = false
  userStore.logout()
}

const handleNavClick = (link) => {
  if (route.path === link) return;
  router.push(link)
}

const openPasswordModal = () => {
  showDropdown.value = false
  passwordModalRef.value?.open()
}

</script>


<style scoped>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.header {
  position: sticky;
  top: 0;
  z-index: 50;
  backdrop-filter: blur(12px);
  background: rgba(8, 8, 12, 0.85);
  border-bottom: 1px solid rgba(0, 255, 255, 0.25);
  padding: 12px 0;
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.3);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* Logo 区域 */
.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: opacity 0.2s;
}

.logo:hover {
  opacity: 0.9;
}

.logo-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #0a0a0f, #1a1a2e);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  border: 1px solid rgba(0, 255, 255, 0.3);
  box-shadow: 0 0 15px rgba(0, 255, 255, 0.1);
}

.logo-icon svg {
  width: 22px;
  height: 22px;
  color: #00ffff;
}

.logo-text h1 {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 2px;
}

.logo-text p {
  font-size: 10px;
  color: rgba(0, 255, 255, 0.7);
  letter-spacing: 1px;
  font-weight: 500;
}

/* 导航菜单 */
.nav {
  display: flex;
  gap: 40px;
  align-items: center;
}

.nav-link {
  color: #8b8b9b;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  transition: all 0.3s;
  position: relative;
  padding: 4px 0;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}

.nav-link:hover {
  color: #00ffff;
  text-shadow: 0 0 8px rgba(0, 255, 255, 0.5);
}

.nav-link.active {
  color: #00ffff;
  text-shadow: 0 0 8px rgba(0, 255, 255, 0.5);
}

.nav-link.active::after {
  content: '';
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, #00ffff, #ff00ff);
  border-radius: 2px;
  box-shadow: 0 0 6px rgba(0, 255, 255, 0.6);
}

/* 右侧区域 */
.auth-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.btn-login {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 255, 255, 0.05);
  border: 1px solid rgba(0, 255, 255, 0.25);
  padding: 8px 20px;
  border-radius: 40px;
  color: #00ffff;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 14px;
  font-weight: 500;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  backdrop-filter: blur(4px);
}

.btn-login:hover {
  background: rgba(0, 255, 255, 0.12);
  border-color: rgba(0, 255, 255, 0.5);
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(0, 255, 255, 0.2);
}

/* ── 用户菜单 ── */
.user-menu {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 4px 12px 4px 4px;
  border-radius: 40px;
  transition: all 0.3s;
  border: 1px solid transparent;
}
.user-menu:hover {
  background: rgba(0, 255, 255, 0.06);
  border-color: rgba(0, 255, 255, 0.2);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 0 12px rgba(0, 255, 255, 0.3);
}
.avatar-text {
  font-size: 14px;
  font-weight: 700;
  color: #0a0a0f;
}
.user-name {
  font-size: 14px;
  color: #c5c5d2;
  font-weight: 500;
}

/* ── 下拉菜单 ── */
.dropdown-panel {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 240px;
  background: rgba(16, 24, 32, 0.98);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(0, 255, 255, 0.25);
  border-radius: 16px;
  padding: 8px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5), 0 0 30px rgba(0, 255, 255, 0.06);
  z-index: 100;
}

.dropdown-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
}
.dropdown-avatar {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 700;
  color: #0a0a0f;
  flex-shrink: 0;
}
.dropdown-user-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
}
.dropdown-username {
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
}
.dropdown-email {
  font-size: 12px;
  color: #6b7280;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dropdown-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.08);
  margin: 4px 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 11px 14px;
  background: none;
  border: none;
  border-radius: 10px;
  color: #c5c5d2;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
}
.dropdown-item i {
  width: 18px;
  text-align: center;
  font-size: 14px;
}
.dropdown-item:hover {
  background: rgba(0, 255, 255, 0.08);
  color: #00ffff;
}
.dropdown-item.logout:hover {
  background: rgba(255, 59, 110, 0.12);
  color: #ff3b6e;
}

/* ── 下拉过渡 ── */
.dropdown-fade-enter-active { transition: all 0.2s ease; }
.dropdown-fade-leave-active { transition: all 0.15s ease; }
.dropdown-fade-enter-from,
.dropdown-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.btn-theme {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.15);
  padding: 8px;
  border-radius: 40px;
  color: #fff;
  cursor: pointer;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
}

.btn-theme:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateY(-1px);
  border-color: rgba(0, 255, 255, 0.3);
}

.theme-icon {
  width: 18px;
  height: 18px;
  color: #aaa;
}

.btn-theme:hover .theme-icon {
  color: #00ffff;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    padding: 0 20px;
  }
  
  .nav {
    gap: 24px;
  }
  
  .nav-link {
    font-size: 13px;
  }
  
  .logo-text h1 {
    font-size: 16px;
  }
  
  .logo-text p {
    font-size: 8px;
  }
  
  .btn-login {
    padding: 6px 14px;
    font-size: 12px;
  }
  
  .btn-login span {
    display: none;
  }
  
  .btn-login .btn-icon {
    margin-right: 0;
  }
}

@media (max-width: 640px) {
  .logo-text {
    display: none;
  }
  
  .nav {
    gap: 16px;
  }
  
  .nav-link {
    font-size: 12px;
  }
}
</style>