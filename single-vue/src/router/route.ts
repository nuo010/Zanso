import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/admin.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        meta: { title: '控制台' },
        component: () => import('@/views/pages/home.vue'),
      },
      {
        path: '/categories',
        meta: { title: '分类管理' },
        component: () => import('@/views/pages/zhome.vue'),
      },
      {
        path: '/user',
        meta: { title: '个人中心' },
        component: () => import('@/views/user/user.vue'),
      },
      {
        path: '/shares',
        meta: { title: '分享链接管理' },
        component: () => import('@/views/pages/share-manage.vue'),
      },
      {
        path: '/users',
        meta: { title: '用户管理', adminOnly: true },
        component: () => import('@/views/pages/user-manage.vue'),
      },
    ],
  },
  {
    path: '/login',
    meta: { title: '登录' },
    component: () => import('@/views/login/login.vue'),
  },
  {
    path: '/share/:code',
    meta: { title: '分享详情', public: true },
    component: () => import('@/views/pages/share.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    meta: { title: '404' },
    component: () => import('@/views/pages/404.vue'),
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

export default router;
