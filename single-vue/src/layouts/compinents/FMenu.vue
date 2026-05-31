<template>
  <div class="f-menu" :style="{ width: store.asideWidth }">
    <el-menu
      :default-active="defaultActive"
      class="border-0"
      :collapse="isCollapse"
      :collapse-transition="false"
      unique-opened
      @select="handleSelect"
    >
      <template v-for="(item, index) in asideMenus" :key="index">
        <el-sub-menu
          :index="item.title"
          v-if="item.children && item.children.length > 0"
        >
<!--     有子菜单的一级菜单     -->
          <!-- 第一级别  -->
          <template #title>
            <el-icon>
              <component :is="item.meta.icon"></component>
            </el-icon>

            <span v-if="!isCollapse">{{ item.title }}</span>

          </template>
          <!-- 第二级别 -->
          <template v-for="item2 in item.children">
            <el-menu-item :index="item2.path">
              <el-icon>
                <component :is="item2.meta.icon"></component>
              </el-icon>
              <span v-if="!isCollapse">{{ item2.title }}</span>
            </el-menu-item>
          </template>
        </el-sub-menu>
<!--    没有子菜单的 普通菜单    -->
        <el-menu-item v-else :index="item.path">
          <el-icon>
            <component :is="item.meta.icon"></component>
          </el-icon>
          <span v-if="!isCollapse">{{ item.title }}</span>
        </el-menu-item>
      </template>
    </el-menu>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { userMainStore } from '@/store/';
import {globalMenuAsideWidthBig} from "@/util/constants";
const store = userMainStore();
const router = useRouter();
const route = useRoute();
const isCollapse = computed(() => !(store.asideWidth == globalMenuAsideWidthBig));
// 默认选中
const defaultActive = ref(route.path);
const asideMenus: Array<any> = store.menuTree;
const handleSelect = (e: any) => {
  router.push(e);
};
</script>

<style scoped lang="scss">
/* 悬浮效果 */
.f-menu .el-menu-item:hover,
.f-menu .el-sub-menu__title:hover {
  background: #f5f7fa;
  color: #409eff;
}

/* 选中时统一加背景色和左侧高亮条 - 增强视觉效果 */
.f-menu .el-menu-item.is-active,
.f-menu .el-sub-menu__title.is-active {
  background: #b3d8ff !important; /* 更深的蓝色背景，更明显 */
  color: #1890ff !important; /* 更深的蓝色文字 */
  font-weight: 600 !important; /* 加粗文字 */
  border-left-color: #1890ff !important; /* 更深的蓝色边框 */
  border-left-width: 4px !important; /* 加粗左侧边框 */
}
</style>
