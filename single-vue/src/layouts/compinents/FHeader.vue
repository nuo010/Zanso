<template>
    <div class="f-header">
        <span class="logo">
            <el-icon class="mr-1"><MostlyCloudy /></el-icon>
            {{ getTitle() }}
        </span>
<!--   菜单缩放功能   -->
<!--        <el-icon class="icon-btn" @click="handleAsideWidth()">-->
<!--          <fold v-if = "store.asideWidth == globalMenuAsideWidthBig"/>-->
<!--          <Expand v-else/>-->
<!--        </el-icon>-->

<!--        <el-tooltip effect="dark" content="刷新" persistent="bottom">-->
<!--            <el-icon class="icon-btn" @click="handleRefresh"><refresh/></el-icon>-->
<!--        </el-tooltip>-->
<!--        时间: {{ nowTime }}-->

        <div class="ml-auto flex items-center" style="margin-left: auto; margin-right: 20px;">
            <!-- 实时在线人数 -->
            <el-tooltip :content="wsStatusText" placement="bottom">
                <span class="online-count">
                    <el-icon class="online-icon"><User /></el-icon>
                    <span class="online-num">{{ store.onlineCount }}</span>
                    <span class="online-label">人在线</span>
                </span>
            </el-tooltip>





            <el-dropdown class="dropdown" @command="handleCommand">
                <span class="flex justify-center items-center text-light-50">
                    <el-avatar class="mr-2" :size="25" src="https://qiniu.ligl.top/1682046831918.jpg" />
                    {{ store.user.nickName }}
                    <el-icon class="el-icon--right">
                        <arrow-down />
                    </el-icon>
                </span>
                <template #dropdown>
                <el-dropdown-menu>
                    <el-dropdown-item command="user">个人中心</el-dropdown-item>
                    <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                    <el-dropdown-item command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
                </template>
            </el-dropdown>
        </div>
    </div>
  <ChangePasswordDialog ref="passwordDialog" />
</template>


<script lang="ts" setup>
import { useRouter } from "vue-router";
import {getToken,removeToken} from "@/util/auth"
import { logout } from "@/api/user";
import { toast } from "@/util/util";
import {userMainStore} from "@/store/"
import {onMounted, onUnmounted, Ref, ref} from "vue";
import dayjs from "dayjs";
import {getTitle} from "@/util/util";
const token = getToken()
const router = useRouter()
import ChangePasswordDialog from "@/views/user/ChangePasswordDialog.vue";
import {globalMenuAsideWidthBig} from "@/util/constants";
import { computed } from "vue";
import { User } from "@element-plus/icons-vue";

const store = userMainStore()

const wsStatusText = computed(() => {
  const s = store.wsOnlineStatus;
  if (s === 'connected') return '实时在线人数';
  if (s === 'connecting') return '连接中...';
  if (s === 'error') return '连接异常，将自动重连';
  return '已断开，将自动重连';
})
const addressData = ref('')
let nowTime:Ref<string> = ref("");
const handleAsideWidth = () =>{
  console.log("修改侧边菜单宽度")
  store.handleAsideWidth()

}
const passwordDialog = ref<InstanceType<typeof ChangePasswordDialog>>();

const handleCommand = (c:any)=>{
    switch (c){
        case "logout":
            exitLogout()
            break
        case "user":
            router.push('/user')
            break
        case "changePassword":
          passwordDialog.value?.open() // 调用弹窗的 open 方法
          break
    }
}

const handleRefresh = ()=>location.reload()


const exitLogout = ()=>{
  removeToken()
  toast("退出登录成功","success")
  router.push('/login')
  // logout().then(()=>{
  //   removeToken()
  //   toast("退出登录成功","success")
  //   router.push('/login')
  // })
}


let timer: number | null = null;

onMounted(() => {
  timer = window.setInterval(() => {
    nowTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss');
  }, 1000);
});

onUnmounted(() => {
  if (timer) {
    window.clearInterval(timer);
  }
});





</script>

<style scoped>
.f-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #ffffff; /* 简约白底 */
  color: #1f2933;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  padding: 0 24px;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  box-sizing: border-box;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #1f2933;
}

.logo .el-icon {
  color: #2196f3; /* 品牌主色，只强调图标 */
  font-size: 20px;
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 36px;
  width: 36px;
  border-radius: 999px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.icon-btn:hover {
  background-color: rgba(148, 163, 184, 0.15);
}

.f-header .dropdown {
  height: 40px;
  cursor: pointer;
  display: flex;
  align-items: center;
}

/* 头像右侧用户名：覆盖 text-light-50，让文字变深可见 */
.f-header .dropdown > span {
  color: #1f2933 !important;
  font-weight: 500;
}

.online-count {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-right: 16px;
  padding: 6px 12px;
  background: rgba(33, 150, 243, 0.08);
  border-radius: 8px;
  color: #1f2933;
  font-size: 14px;
}

.online-icon {
  font-size: 16px;
  color: #2196f3;
}

.online-num {
  font-weight: 600;
  color: #2196f3;
}

.online-label {
  color: #64748b;
  font-size: 13px;
}
</style>
