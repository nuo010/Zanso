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
  border: 1px solid rgba(32, 53, 41, 0.08);
  background: rgba(255, 252, 246, 0.92);
}

.profile-header h2 {
  margin: 0;
  color: #203529;
}

.profile-header p {
  margin: 6px 0 0;
  color: #6c7f73;
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 16px;
}

.profile-item {
  padding: 16px;
  border-radius: 18px;
  background: #f5f3eb;
}

.profile-item span {
  display: block;
  color: #6c7f73;
  margin-bottom: 10px;
}

.profile-item strong {
  color: #203529;
  word-break: break-all;
}
</style>
