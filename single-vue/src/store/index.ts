import { defineStore } from 'pinia';
import { globalMenuAsideWidthBig, globalMenuAsideWidthLittle } from '@/util/constants';
import { getCurrentUser, getUserCategories } from '@/api/user';

export interface PlatformUser {
  id: string;
  name: string;
  loginName: string;
  contactName?: string;
  contactPhone?: string;
  status?: string;
}

export interface CategoryItem {
  id: string;
  name: string;
  description?: string;
  coverUrl?: string;
  visible?: boolean;
  status?: string;
  createdAt?: string;
}

export const userMainStore = defineStore('main', {
  persist: {
    enabled: true,
    strategies: [
      {
        key: 'zansoLocalStorage',
        storage: localStorage,
      },
    ],
  },
  state: () => ({
    user: {} as PlatformUser,
    categories: [] as CategoryItem[],
    asideWidth: globalMenuAsideWidthBig,
  }),
  getters: {
    getUserId: (state) => state.user.id || '',
    isCollapse: (state) => state.asideWidth !== globalMenuAsideWidthBig,
  },
  actions: {
    async loadProfile() {
      const res = await getCurrentUser();
      this.user = res.data;
      return res.data;
    },
    async loadCategories() {
      if (!this.user.id) return [];
      const res = await getUserCategories(this.user.id);
      this.categories = res.data || [];
      return this.categories;
    },
    setUserStore(user: PlatformUser) {
      this.user = user;
    },
    handleAsideWidth() {
      this.asideWidth = this.asideWidth === globalMenuAsideWidthBig ? globalMenuAsideWidthLittle : globalMenuAsideWidthBig;
    },
    resetAll() {
      this.$reset();
    },
  },
});
