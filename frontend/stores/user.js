import { defineStore } from 'pinia'
import { userApi } from '~/api/user'
import { taskApi } from '~/api/task'
import cookie from '~/utils/cookie'
import { useSSE } from '~/composables/useSSE'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {},
    isLogin: false,
    // ── 任务统计 ──
    taskCounts: { pending: 0, running: 0, success: 0, failed: 0 },
  }),

  actions: {
    async initialize() {
      if (!import.meta.client) return;
      console.log('Initializing user store...')
      const token = cookie.get()
      if(!token) return;
      if (import.meta.client && token && !this.isLogin) {
        try {
          const res = await userApi.userInfo()
          if (res.code === 200) {
            await this.setUserInfo(res.data)
            await this.fetchUserTasks()
          }
        } catch (e) {
          // token 过期或无效，静默忽略
        }
      }
    },

    // ── 任务计数 ──
    async fetchUserTasks() {
      try {
        const res = await taskApi.getTaskStatus()
        if (res.code !== 200 || !res.data) return

        this.taskCounts = {
          pending: res.data.pending || 0,
          running: res.data.running || 0,
          success: res.data.success || 0,
          failed: res.data.failed || 0,
        }
      } catch {
        // 静默
      }
    },

    // ── 用户级 SSE ──
    startUserSSE() {
      const { connectUser } = useSSE()
      connectUser(({ task_id, status }) => {
        console.log('[UserSSE] 收到状态变更:', task_id, status)
        this.fetchUserTasks()
      })
    },

    stopUserSSE() {
      const { disconnectUser } = useSSE()
      disconnectUser()
    },

    async setUserInfo(userInfo) {
      this.userInfo = userInfo
      this.isLogin = true
      const isNewLogin = !!userInfo.token
      if (isNewLogin) {
        cookie.set(userInfo.token)
        // 等待 cookie 生效后再建 SSE，避免 Nuxt 代理层转发时丢失认证
        await new Promise(r => setTimeout(r, 300))
      }
      this.startUserSSE()
    },

    async logout() {
      // 先断开 SSE 再清除状态
      this.stopUserSSE()

      cookie.remove()
      this.userInfo = {}
      this.isLogin = false
      this.taskCounts = { pending: 0, running: 0, success: 0, failed: 0 }
    },
  },

  getters: {
    getUserInfo: (state) => state.userInfo,
  },
})
