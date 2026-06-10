import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
    token: '',
    isLogin: false,
  }),

  actions: {
    async initialize() {
      console.log("第一次访问");
      // 这里可以添加初始化逻辑
    },
    
    async setUserInfo(userInfo) {
      this.userInfo = userInfo
      this.isLogin = true
      console.log('User info set:', this.userInfo)
    },

    async logout() {
      this.userInfo = {}
      this.token = ''
      this.isLogin = false
    },
  },

  getters: {
    getUserInfo: (state) => state.userInfo,
  },
})