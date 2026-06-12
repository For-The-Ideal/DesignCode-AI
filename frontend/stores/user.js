import { defineStore } from 'pinia'
import { userApi } from '~/api/user'
import cookie from '~/utils/cookie'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
    isLogin: false,
  }),

  actions: {
    async initialize() {
      console.log('Initializing user store...')
      // 页面加载时检测 token cookie，自动恢复登录态
      const token = cookie.get()
      console.log('Token from cookie:', token)
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
      if (userInfo.token) {
        cookie.set(userInfo.token)
        console.log('Token set:', userInfo.token)
      }
      console.log('User info set:', userInfo)
    },

    async logout() {
      cookie.remove()
      this.userInfo = {}
      this.isLogin = false
    },
  },

  getters: {
    getUserInfo: (state) => state.userInfo,
  },
})