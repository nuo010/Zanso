<template>
  <div style="min-height: 100vh; display: flex; align-items: center; justify-content: center; background-color: #f5f5f5;">
    <div style="text-align: center; padding: 20px;">
      <h1 style="color: #ff4d4f; font-size: 48px; margin-bottom: 20px;">⚠️</h1>
      <h2 style="color: #333; margin-bottom: 20px;">权限错误</h2>
      <p style="color: #666; font-size: 16px; margin-bottom: 30px;">
        请在<span class="highlight-portal">门户</span>重新登录并进入该系统
      </p>
      <el-button type="primary" @click="goToLogin">前往登录</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goToLogin = () => {
  // router.push('/login');
};
</script>

<style scoped>

.highlight-portal {
  color: #ff4d4f;
  font-weight: bold;
  font-size: 22px;
  text-decoration: underline;
  padding: 0 4px;
  background-color: #fff2f0;
  border-radius: 2px;
}
</style>
