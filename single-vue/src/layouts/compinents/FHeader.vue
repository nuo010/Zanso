<template>
  <div class="header-wrap">
    <div class="brand">
      <div class="brand-mark">Z</div>
      <div>
        <div class="brand-title">{{ getTitle() }}</div>
        <div class="brand-subtitle">数字内容管理平台</div>
      </div>
    </div>

    <div class="header-actions">
      <el-dropdown @command="handleCommand">
        <div class="user-box">
          <el-avatar :size="34">{{ store.user.name?.slice(0, 1) || 'U' }}</el-avatar>
          <div class="user-meta">
            <strong>{{ store.user.name || '未登录用户' }}</strong>
            <span>{{ store.user.loginName || '--' }}</span>
          </div>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="user">个人中心</el-dropdown-item>
            <el-dropdown-item command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';
import { userMainStore } from '@/store';
import { getTitle, toast } from '@/util/util';
import { removeToken } from '@/util/auth';
import { logoutUser } from '@/api/user';

const router = useRouter();
const store = userMainStore();

async function handleCommand(command: string) {
  if (command === 'user') {
    await router.push('/user');
    return;
  }
  if (command === 'logout') {
    try {
      await logoutUser();
    } catch (error) {
      console.warn('退出登录请求失败，直接清理本地状态', error);
    }
    removeToken();
    store.resetAll();
    toast('已退出登录', 'success');
    await router.push('/login');
  }
}
</script>

<style scoped lang="scss">
.header-wrap {
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: rgba(255, 255, 255, 0.78);
  border-bottom: 1px solid rgba(129, 165, 235, 0.16);
  backdrop-filter: blur(18px);
}

.brand {
  display: flex;
  align-items: center;
  gap: 14px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  background: linear-gradient(135deg, #2f6bff, #69b7ff);
  color: #fdfefe;
  display: grid;
  place-items: center;
  font-size: 20px;
  font-weight: 700;
  box-shadow: 0 12px 24px rgba(47, 107, 255, 0.24);
}

.brand-title {
  font-size: 18px;
  font-weight: 700;
  color: #17315f;
}

.brand-subtitle {
  font-size: 12px;
  color: #6d82a7;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.user-box {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 10px;
  border-radius: 16px;
  transition: background-color 0.2s ease;
}

.user-box:hover {
  background: rgba(47, 107, 255, 0.08);
}

.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.user-meta strong {
  font-size: 14px;
  color: #17315f;
}

.user-meta span {
  font-size: 12px;
  color: #6d82a7;
}
</style>
