import { defineStore } from 'pinia';
import { getCurrentUser, getCurrentUserCategories, getDashboardStats } from '@/api/user';

export interface PlatformUser {
  id: string;
  name: string;
  loginName: string;
  email?: string;
  status?: string;
  roleCodes?: string[];
  roleNames?: string[];
  resourceCount?: number;
  fileSizeTotal?: number;
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

export interface DashboardStats {
  collectionCount: number;
  resourceCount: number;
  imageCount: number;
  videoCount: number;
  fileSizeTotal: number;
  imageSizeTotal: number;
  videoSizeTotal: number;
}

export interface Announcement {
  id: string;
  title: string;
  content?: string;
  status: 'draft' | 'active';
  createdBy?: string;
  createdAt?: string;
  updatedAt?: string;
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
    dashboardStats: {
      collectionCount: 0,
      resourceCount: 0,
      imageCount: 0,
      videoCount: 0,
      fileSizeTotal: 0,
      imageSizeTotal: 0,
      videoSizeTotal: 0,
    } as DashboardStats,
  }),
  getters: {
    getUserId: (state) => state.user.id || '',
    isAdmin: (state) => state.user.roleCodes?.includes('admin') || false,
  },
  actions: {
    async loadProfile() {
      const res = await getCurrentUser();
      this.user = res.data;
      return res.data;
    },
    async loadCategories(page = 1, pageSize = 20) {
      const res = await getCurrentUserCategories({ page, pageSize });
      this.categories = res.data?.list || [];
      this.categoryTotal = res.data?.total || 0;
      this.categoryPage = res.data?.page || page;
      this.categoryPageSize = res.data?.pageSize || pageSize;
      return this.categories;
    },
    async loadDashboardStats() {
      const res = await getDashboardStats();
      this.dashboardStats = {
        collectionCount: res.data?.collectionCount || 0,
        resourceCount: res.data?.resourceCount || 0,
        imageCount: res.data?.imageCount || 0,
        videoCount: res.data?.videoCount || 0,
        fileSizeTotal: res.data?.fileSizeTotal || 0,
        imageSizeTotal: res.data?.imageSizeTotal || 0,
        videoSizeTotal: res.data?.videoSizeTotal || 0,
      };
      return this.dashboardStats;
    },
    setUserStore(user: PlatformUser) {
      this.user = user;
    },
    resetAll() {
      this.$reset();
    },
  },
});
