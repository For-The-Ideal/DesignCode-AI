import { defineStore } from 'pinia'

export const useCommonStore = defineStore('common', {
  state: () => ({
    loginModalVisible: false,
    isGenerating: false,
  }),

  actions: {
    openLoginModal() {
      this.loginModalVisible = true
    },
    closeLoginModal() {
      this.loginModalVisible = false
    },
    startGenerating() {
      this.isGenerating = true
    },
    finishGenerating() {
      this.isGenerating = false
    },
  },
})
