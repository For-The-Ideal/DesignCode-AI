import { defineStore } from 'pinia'

export const useMenuListStore = defineStore('menuList', {
  state: () => ({
    menuList: []
  }),

  actions: {
    
  },

  getters: {
    getMenuList: (state) => state.menuList
  }
})
