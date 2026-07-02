import { defineStore } from 'pinia'

export const useCommonStore = defineStore('common', {
  state: () => ({
    loginModalVisible: false,
    membershipModalVisible: false,
  }),

  actions: {
    setLoginModalVisible(visible) {
      this.loginModalVisible = visible
    },
    setMembershipModalVisible(visible) {
      this.membershipModalVisible = visible
    },
  },
})
