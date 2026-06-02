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
  roleCodes?: string[];
  roleNames?: string[];
}

export interface Collection {
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
    categories: [] as Collection[],
    categoryTotal: 0,
    categoryPage: 1,
    categoryPageSize: 20,
    asideWidth: globalMenuAsideWidthBig,
  }),
  getters: {
    getUserId: (state) => state.user.id || '',
    isCollapse: (state) => state.asideWidth !== globalMenuAsideWidthBig,
    isAdmin: (state) => state.user.roleCodes?.includes('admin') || false,
  },
  actions: {
    async loadProfile() {
      const res = await getCurrentUser();
      this.user = res.data;
      return res.data;
    },
    async loadCategories(page = 1, pageSize = 20) {
      if (!this.user.id) return [];
      const res = await getUserCategories(this.user.id, { page, pageSize });
      this.categories = res.data?.list || [];
      this.categoryTotal = res.data?.total || 0;
      this.categoryPage = res.data?.page || page;
      this.categoryPageSize = res.data?.pageSize || pageSize;
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
