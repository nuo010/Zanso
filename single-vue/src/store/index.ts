import { defineStore } from 'pinia';
import {globalMenuAsideWidthBig, globalMenuAsideWidthLittle} from "@/util/constants";
import {getMenuList, getMenuTree} from "@/api/user";


interface User {
  userId: string;
  nickName: string;
  userName?: string;
  email?: string;
  ipaddr?: string;
}

export const userMainStore = defineStore('main', {
  persist: {
    enabled: true, //开启数据持久化
    strategies: [
      {
        key: 'singleLocalStorage', //给一个要保存的名称
        storage: localStorage, //sessionStorage / localStorage 存储方式
      },
    ],
  },
  state: () => {
    return {
      user: {} as User,
      menuTree: [],
      // 侧边宽度
      asideWidth: globalMenuAsideWidthBig,
      // 实时在线人数（WebSocket）
      onlineCount: 0,
      wsOnlineStatus: 'closed' as 'connecting' | 'connected' | 'closed' | 'error',
    };
  },
  getters: {
    getUserId: (state) => {
      return state.user.userId || '';
    },
  },
  actions: {
    async loadMenu() {
      const menuListRes = await getMenuTree({isTree:1,type: import.meta.env.VITE_API_SYSTEM_ID})
      this.menuTree =  menuListRes.data
    },
    setUserStore(user: any) {
      this.user = user;
    },
    setMenuTreeStore(menuTree: any) {
      console.log('store:menuTree', menuTree);
      this.menuTree = menuTree;
    },
    handleAsideWidth() {
      this.asideWidth = this.asideWidth == globalMenuAsideWidthBig ? globalMenuAsideWidthLittle : globalMenuAsideWidthBig;
    },
    resetAll() {
      this.$reset();
    },
    setOnlineCount(count: number) {
      this.onlineCount = count;
    },
    setWsOnlineStatus(status: 'connecting' | 'connected' | 'closed' | 'error') {
      this.wsOnlineStatus = status;
    },
  },
});
