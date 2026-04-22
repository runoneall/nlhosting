<script setup>
import { useUserInfoStore } from '@/stores/userinfo'
import router from '@/router'
import { onMounted, ref, reactive } from 'vue'
import request from '@/request'
import { ElMessage } from 'element-plus'

const userInfo = useUserInfoStore()

const domainFetching = ref(false)
const domainList = ref([])

const domainFetch = async () => {
    if (domainList.value.length > 0) return
    domainFetching.value = true

    try {
        const resp = await request.get('/user/domains')
        if (resp.status !== 200) {
            ElMessage.error('后端错误')
            return
        }

        domainList.value = resp.data
    } finally {
        domainFetching.value = false
    }
}

const serverFetching = ref(false)
const serverList = ref([])

const serverFetch = async () => {
    if (serverList.value.length > 0) return
    serverFetching.value = true

    try {
        const resp = await request.get('/host/available')
        if (resp.status !== 200) {
            ElMessage.error('后端错误')
            return
        }

        serverList.value = resp.data
    } finally {
        serverFetching.value = false
    }
}

const form = reactive({
    domain: ref(''),
    server: ref('')
})

const requesting = ref(false)

const requestNew = async () => {
    if (requesting.value === true) return
    requesting.value = true

    try {
        const resp = await request.post('/host/new', form)
        if (resp.status === 500) {
            ElMessage.error('后端错误')
            return
        }

        if (resp.status === 400) {
            ElMessage.error('服务器拒绝了请求，因为配置不正确')
            return
        }

        ElMessage.success('请求成功，任务已加入队列，最快需 5 分钟')
    } finally {
        requesting.value = false
    }
}

onMounted(async () => {
    if (!userInfo.hasinfo) {
        router.push('/login')
        return
    }
})
</script>

<template>
    <h2 style="display: flex; justify-content: center">
        你好 &nbsp;
        <span v-if="userInfo.info.name">{{ userInfo.info.name }}</span>
        <span v-else>{{ userInfo.info.username }}</span>
    </h2>

    <el-form :model="form" label-width="auto" style="max-width: 600px; margin-left: 20px">
        <el-form-item label="域名">
            <el-select v-model="form.domain" filterable remote :loading="domainFetching" @focus="domainFetch" placeholder="请选择">
                <el-option v-for="item in domainList" :key="item" :label="item" :value="item" />
            </el-select>
        </el-form-item>
        <el-form-item label="服务器">
            <el-select v-model="form.server" filterable remote :loading="serverFetching" @focus="serverFetch" placeholder="请选择">
                <el-option v-for="item in serverList" :key="item" :label="item" :value="item" />
            </el-select>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" :loading="requesting" @click="requestNew">请求创建</el-button>
        </el-form-item>
    </el-form>
</template>

<style scoped></style>
