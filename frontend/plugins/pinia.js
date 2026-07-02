import { defineNuxtPlugin } from '#app'
import { useUserStore } from '~/stores/user'

export default defineNuxtPlugin(async () => {
  const userStore = useUserStore()
  await userStore.initialize()
  console.log('User store initialized')
})