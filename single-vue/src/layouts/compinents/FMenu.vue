<template>
  <div class="menu-wrap">
    <el-menu
      :default-active="route.path"
      :collapse="store.isCollapse"
      :collapse-transition="false"
      class="menu-panel"
      @select="handleSelect"
    >
      <el-menu-item index="/dashboard">
        <el-icon><House /></el-icon>
        <span>控制台</span>
      </el-menu-item>
      <el-menu-item index="/categories">
        <el-icon><FolderOpened /></el-icon>
        <span>分类管理</span>
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
  border-radius: 14px;
  margin-bottom: 8px;
  color: #385441;
}

.menu-panel :deep(.el-menu-item.is-active) {
  background: linear-gradient(135deg, rgba(49, 92, 69, 0.14), rgba(163, 196, 172, 0.28));
  color: #203529;
  font-weight: 700;
}
</style>
