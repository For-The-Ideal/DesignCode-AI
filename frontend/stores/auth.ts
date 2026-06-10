// stores/auth.ts
import { defineStore } from 'pinia'

interface UserInfo {
  username: string
  email?: string
  avatar?: string
}

interface LoginParams {
  username: string
  password: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    userInfo: null as UserInfo | null,
    token: null as string | null,
    isLoggedIn: false
  }),

  actions: {
    login(params: LoginParams) {
      // 演示：模拟登录成功
      this.userInfo = { username: params.username, email: `${params.username}@example.com` }
      this.token = 'mock-token-' + Date.now()
      this.isLoggedIn = true
    },

    register(username: string, email: string, password: string) {
      // 预留扩展
    },

    logout() {
      this.userInfo = null
      this.token = null
      this.isLoggedIn = false
      localStorage.removeItem('user')
    },

    changePassword(oldPwd: string, newPwd: string): { success: boolean; message: string } {
      // 演示：旧密码固定为 123456
      if (oldPwd !== '123456') return { success: false, message: '当前密码错误' }
      if (newPwd.length < 6) return { success: false, message: '新密码长度至少6位' }
      return { success: true, message: '密码修改成功！' }
    },

    checkAuth() {
      const saved = localStorage.getItem('user')
      if (saved) {
        try {
          this.userInfo = JSON.parse(saved)
          this.isLoggedIn = true
          this.token = 'restored-token'
        } catch {
          localStorage.removeItem('user')
        }
      }
    }
  }
})