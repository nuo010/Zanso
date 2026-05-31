<template>
  <div class="user-page">
    <el-card class="profile-card" shadow="never">
      <template #header>
        <div class="profile-header">
          <div>
            <h2>个人中心</h2>
            <p>当前登录用户的基础信息。</p>
          </div>
        </div>
      </template>

      <div class="profile-grid">
        <div class="profile-item">
          <span>用户名称</span>
          <strong>{{ store.user.name || '--' }}</strong>
        </div>
        <div class="profile-item">
          <span>登录账号</span>
          <strong>{{ store.user.loginName || '--' }}</strong>
        </div>
        <div class="profile-item">
          <span>联系人</span>
          <strong>{{ store.user.contactName || '--' }}</strong>
        </div>
        <div class="profile-item">
          <span>联系电话</span>
          <strong>{{ store.user.contactPhone || '--' }}</strong>
        </div>
        <div class="profile-item">
          <span>用户 ID</span>
          <strong>{{ store.user.id || '--' }}</strong>
        </div>
        <div class="profile-item">
          <span>Token</span>
          <strong>{{ getToken() || '--' }}</strong>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { userMainStore } from '@/store';
import { getToken } from '@/util/auth';

const store = userMainStore();
</script>

<style scoped lang="scss">
.profile-card {
  border-radius: 24px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 20px 46px rgba(36, 84, 170, 0.08);
}

.profile-header h2 {
  margin: 0;
  color: #17315f;
}

.profile-header p {
  margin: 6px 0 0;
  color: #6d82a7;
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.profile-item {
  padding: 16px;
  border-radius: 18px;
  background: #f4f8ff;
  border: 1px solid rgba(123, 162, 255, 0.12);
}

.profile-item span {
  display: block;
  color: #6d82a7;
  margin-bottom: 10px;
}

.profile-item strong {
  color: #17315f;
  word-break: break-all;
}
</style>
