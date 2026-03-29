import { defineNuxtPlugin } from '#app'

export default defineNuxtPlugin((nuxtApp) => {
    if (!nuxtApp.$pinia) {
        console.error('Pinia not installed!')
        return
    }
    
    // try {
    //     const userStore = useUserStore(nuxtApp.$pinia)
    //     console.log('Store initialized:', userStore)
        
    //     // 初始化 store
    //     userStore.initialize()
    // } catch (error) {
    //     console.error('Error initializing store:', error)
    // }
})