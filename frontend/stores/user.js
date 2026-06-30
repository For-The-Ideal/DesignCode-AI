import { defineStore } from 'pinia'
import { userApi } from '~/api/user'
import cookie from '~/utils/cookie'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
    credits: 100,
    creditsUsed: 0,
    isLogin: false,
  }),

  actions: {
    async initialize() {
      console.log('Initializing user store...')
      // 页面加载时检测 token cookie，自动恢复登录态
      const token = cookie.get()
      if (import.meta.client && token && !this.isLogin) {
        try {
          const res = await userApi.userInfo()
          if (res.code === 200) {
            this.setUserInfo(res.data)
          }
        } catch (e) {
          // token 过期或无效，静默忽略
        }
      }
    },

    async setUserInfo(userInfo) {
      this.userInfo = userInfo
      this.isLogin = true
      if (userInfo.credits !== undefined) {
        this.credits = userInfo.credits
        this.creditsUsed = userInfo.credits_used || 0
      }
      if (userInfo.token) {
        cookie.set(userInfo.token)
      }
    },

    async logout() {
      cookie.remove()
      this.userInfo = {}
      this.credits = 0
      this.creditsUsed = 0
      this.isLogin = false
    },
  },

  getters: {
    getUserInfo: (state) => state.userInfo,
  },
})