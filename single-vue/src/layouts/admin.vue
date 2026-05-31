<template>
  <el-container class="layout-shell">
    <el-header class="layout-header">
      <f-header />
    </el-header>
    <el-container class="layout-body">
      <el-aside :width="store.asideWidth" class="layout-aside">
        <f-menu />
      </el-aside>
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import FHeader from './compinents/FHeader.vue';
import FMenu from './compinents/FMenu.vue';
import { userMainStore } from '@/store';

const store = userMainStore();
</script>

<style scoped lang="scss">
.layout-shell {
  height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(210, 233, 214, 0.95), transparent 35%),
    linear-gradient(180deg, #f6f2e8 0%, #edf4ef 100%);
}

.layout-header {
  height: 72px !important;
  padding: 0 !important;
}

.layout-body {
  height: calc(100vh - 72px);
}

.layout-aside {
  border-right: 1px solid rgba(32, 53, 41, 0.08);
  background: rgba(255, 252, 246, 0.84);
  backdrop-filter: blur(10px);
}

.layout-main {
  padding: 24px;
  overflow: auto;
}
</style>
