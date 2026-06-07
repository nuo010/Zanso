<template>
  <el-container class="layout-shell">
    <el-header class="layout-header">
      <f-header />
    </el-header>
    <el-container class="layout-body">
      <el-aside width="220px" class="layout-aside">
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
</script>

<style scoped lang="scss">
.layout-shell {
  height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(139, 180, 255, 0.22), transparent 28%),
    radial-gradient(circle at right center, rgba(194, 225, 255, 0.2), transparent 24%),
    linear-gradient(180deg, #f8fbff 0%, #eef4ff 100%);
}

.layout-header {
  height: 72px !important;
  padding: 0 !important;
}

.layout-body {
  height: calc(100vh - 72px);
}

.layout-aside {
  border-right: 1px solid rgba(129, 165, 235, 0.16);
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(18px);
  box-shadow: inset -1px 0 0 rgba(255, 255, 255, 0.5);
}

.layout-main {
  padding: 24px;
  overflow: auto;
}
</style>
