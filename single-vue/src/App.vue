<template>
  <router-view />
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import { userMainStore } from '@/store'
import { useRoute, useRouter } from 'vue-router'

const store = userMainStore()
const route = useRoute()
const router = useRouter()

onMounted(async () => {
  // 等路由完成初始导航后再取 path，否则可能一直是 /
  await router.isReady()
  const path = route.path
  console.log('app.vue onmounted 当前路由:', path)
  if (path === '/login' || path === '/') return
  await store.loadMenu()
})




</script>



<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

#app {
  height: 100%;
}

/* 全局loading */
#nprogress .bar {
  background-color: rgb(232, 80, 10) !important;
  height: 5px !important;
}
/* 全局分页 组件居中*/
.el-pagination {
  justify-content: center;
}
</style>
