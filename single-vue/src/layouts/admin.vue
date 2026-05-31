<template>
  <el-container>
    <el-header>
      <f-header></f-header>
    </el-header>
    <el-container>
      <el-aside :width="store.asideWidth">
        <f-menu></f-menu>
      </el-aside>
      <el-main class="main">
        <div class="main-content" :class="{ 'loading-active': isAreaLoading }">
          <router-view></router-view>
          <area-loading :visible="isAreaLoading" :text="areaLoadingText" />
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts" setup>
import FHeader from './compinents/FHeader.vue';
import FMenu from './compinents/FMenu.vue';
import { AreaLoading } from '@/components/AreaLoading';
import { isAreaLoading, areaLoadingText } from '@/util/routeLoading';
import { userMainStore } from '@/store';
import { onMounted, onUnmounted } from 'vue';
import { onlineWebSocket } from '@/util/onlineWebSocket';

const store = userMainStore();

let unsubCount: (() => void) | null = null;
let unsubStatus: (() => void) | null = null;

onMounted(() => {
  if (store.getUserId) {
    onlineWebSocket.connect(store.getUserId);
  }
  unsubCount = onlineWebSocket.onCount((count) => store.setOnlineCount(count));
  unsubStatus = onlineWebSocket.onStatus((status) => store.setWsOnlineStatus(status));
});

onUnmounted(() => {
  unsubCount?.();
  unsubStatus?.();
  onlineWebSocket.disconnect();
});
</script>

<style scoped lang="scss">
.el-container {
  height: 100vh;
  overflow: hidden;
}

.el-container .el-container {
  height: calc(100vh - 64px);
  overflow: hidden;
}

.el-header {
  height: 64px !important;
  padding: 0 !important;
}

.el-aside {
  transition: all 0.2s;
  background: #f7fafb;
  height: 100%;
}

.main {
  height: 100%;
  overflow: hidden;
  position: relative;
  //padding: 0 !important;
}

.main-content {
  height: 100%;
  position: relative;
  transition: opacity 0.3s ease;
}

.main-content.loading-active {
  opacity: 0.7;
}
</style>
