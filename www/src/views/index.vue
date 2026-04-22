<script setup>
import { onMounted } from 'vue'
import router from '@/router'
import request from '@/request'
import { useUserInfoStore } from '@/stores/userinfo'

onMounted(async () => {
    const userInfo = useUserInfoStore()
    if (userInfo.hasinfo) {
        router.push('/dashboard')
        return
    }

    const resp = await request.get('/user/info')
    if (resp.status !== 200) {
        router.push('/login')
        return
    }

    userInfo.set(resp.data)
    router.push('/dashboard')
})
</script>

<template></template>

<style scoped></style>
