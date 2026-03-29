// stores/auth.ts
import { defineStore } from 'pinia'

interface UserInfo {
  userName: string
  email?: string
  avatar?: string
}

interface LoginParams{
    userName: string
    password:string 
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    userInfo: null,
    token: null,
    isLoggedIn: false
  }) as {
    userInfo: UserInfo | null
    token: string | null
    isLoggedIn: boolean
  },
  
  actions: {
    login(params:LoginParams) {
     
    },
    
    register(username: string, email: string, password: string) {
    
    },
    
    logout() {
      this.userInfo = null
      this.token = null
      this.isLoggedIn = false
    },
    
    checkAuth() {
     
    }
  }
})