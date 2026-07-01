import { defineStore } from 'pinia'

export const useCommonStore = defineStore('common', {
  state: () => ({
    loginModalVisible: false,
  }),

  actions: {
    setLoginModalVisible(visible) {
      this.loginModalVisible = visible
    },
  },
})
