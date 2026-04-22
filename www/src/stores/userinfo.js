import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserInfoStore = defineStore('userinfo', () => {
    const info = ref({})
    const hasinfo = ref(false)

    const set = data => {
        info.value = data
        hasinfo.value = true
    }

    return { info, hasinfo, set }
})
