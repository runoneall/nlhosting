<script setup>
import router from '@/router'
router.push('/')
</script>

<template>
    <router-view />
</template>

<style>
body {
    margin: 0;
    border: 0;
    padding: 0;
}

#app {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
}
</style>
