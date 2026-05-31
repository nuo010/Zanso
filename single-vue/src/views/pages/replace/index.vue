<template>
  <div class="loading" style="background: transparent">
    <i class="el-icon-loading" />
    <!-- <div v-if="isShowError">账号登录失效，授权失败，请重新登录</div> -->
  </div>
</template>
<script setup lang="ts" name="Replace">
import { useRoute, useRouter } from 'vue-router';
import { nextTick } from 'vue';
import { setToken} from "@/util/auth";
import {userMainStore} from "@/store";
import {getUser} from "@/api/user";
// import { useStore } from '@/stores/index';
let store = userMainStore(); //接收
// const store = useStore();
const route = useRoute();
const router = useRouter();
if (route.query.token) {
  setToken(route.query.token as string);
}

nextTick(async () => {
  const { data } = await getUser();
  store.setUserStore(data);
  // 1、请注意执行顺序(存储用户信息到vuex)
  if (route.query?.redirect) {
    router.push({
      path: route.query?.redirect as string,
    });
  } else {
    router.push('/');
  }
});
</script>
<style lang="scss" scoped>
.loading {
  text-align: center;
  padding-top: 200px;

  i {
    font-size: 30px;
  }
}
</style>
