<template>
  <div class="menu-wrap">
    <el-menu
      :default-active="route.path"
      class="menu-panel"
      @select="handleSelect"
    >
      <el-menu-item index="/dashboard">
        <el-icon><House /></el-icon>
        <span>控制台</span>
      </el-menu-item>
      <el-menu-item index="/categories">
        <el-icon><FolderOpened /></el-icon>
        <span>展册管理</span>
      </el-menu-item>
      <el-menu-item index="/shares">
        <el-icon><Link /></el-icon>
        <span>分享链接管理</span>
      </el-menu-item>
      <el-menu-item v-if="store.isAdmin" index="/users">
        <el-icon><UserFilled /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
      <el-menu-item index="/user">
        <el-icon><User /></el-icon>
        <span>个人中心</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { userMainStore } from '@/store';

const router = useRouter();
const route = useRoute();
const store = userMainStore();

function handleSelect(path: string) {
  router.push(path);
}
</script>

<style scoped lang="scss">
.menu-wrap {
  padding: 18px 12px;
}

.menu-panel {
  border-right: none;
  background: transparent;
}

.menu-panel :deep(.el-menu-item),
.menu-panel :deep(.el-sub-menu__title) {
  border-radius: 16px;
  margin-bottom: 10px;
  color: #47658f;
  height: 46px;
  line-height: 46px;
}

.menu-panel :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(47, 107, 255, 0.14), rgba(120, 188, 255, 0.26));
  color: #17315f;
  font-weight: 700;
  box-shadow: 0 10px 22px rgba(47, 107, 255, 0.12);
}

.menu-panel :deep(.el-menu-item:hover),
.menu-panel :deep(.el-sub-menu__title:hover) {
  background: rgba(47, 107, 255, 0.06);
}
</style>
