import { defineNuxtPlugin } from '#app'

export default defineNuxtPlugin((nuxtApp) => {
    if (!nuxtApp.$pinia) {
        console.error('Pinia not installed!')
        return
    }
})