<template>
  <router-view />
</template>

<script lang="ts" setup>
import { onMounted } from 'vue';
import { userMainStore } from '@/store';
import { useRoute } from 'vue-router';
import { getToken } from '@/util/auth';

const store = userMainStore();
const route = useRoute();

onMounted(async () => {
  if (route.meta.public) return;
  const token = getToken();
  if (!token) return;
  if (!store.user.id) {
    await store.loadProfile();
  }
  await store.loadCategories();
});
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
