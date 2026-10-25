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
    // ── 实时进度（SSE task_progress 推送）：{ [task_id]: { progress, status, current_step } }
    taskProgress: {},
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

    // ── 任务计数（登录后拉一次做初始值；之后由 SSE 推送 counts 保持最新，见 applyTaskCounts）──
    async fetchUserTasks() {
        const res = await taskApi.getTaskStatus()
        if (res && res.code === 200 && res.data) {
          this.taskCounts = {
            pending: res.data.pending || 0,
            running: res.data.running || 0,
            success: res.data.success || 0,
            failed: res.data.failed || 0,
          }
        }
    },

    // 直接应用 SSE 推送的权威计数（事件驱动，替代每事件回拉 /task/status）
    applyTaskCounts(counts) {
      if (!counts) return
      this.taskCounts = {
        pending: counts.pending || 0,
        running: counts.running || 0,
        success: counts.success || 0,
        failed: counts.failed || 0,
      }
    },

    // 将 SSE 推来的局部字段合并到指定任务的实时进度（触发响应式更新）
    patchTaskRow(taskId, fields) {
      if (!taskId) return
      this.taskProgress = {
        ...this.taskProgress,
        [taskId]: { ...(this.taskProgress[taskId] || {}), ...fields },
      }
    },

    // ── 用户级 SSE ──
    startUserSSE() {
      const { connectUser } = useSSE()
      connectUser(
        // task_status：携带权威 counts → 直接更新计数；同时同步该行状态
        ({ task_id, status, counts }) => {
          this.applyTaskCounts(counts)
          if (task_id) this.patchTaskRow(task_id, { status })
        },
        // task_progress：中间阶段进度 → 更新对应行
        ({ task_id, status, progress, current_step }) => {
          this.patchTaskRow(task_id, { progress, status, current_step })
        }
      )
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
      this.taskProgress = {}
    },
  },

  getters: {
    getUserInfo: (state) => state.userInfo,
  },
})
