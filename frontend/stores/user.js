import { defineStore } from 'pinia'
import { userApi } from '~/api/user'
import { commonApi } from '~/api/common'
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
      console.log('Initializing user store...')
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

    // ── 任务计数 ──
    async fetchUserTasks() {
      try {
        const res = await commonApi.getUserTasks()
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
      if (userInfo.token) {
        cookie.set(userInfo.token)
      }

      // 登录后拉取任务计数 + 建立用户级 SSE 连接
      await this.fetchUserTasks()
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
