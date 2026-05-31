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
    ],
  },
  {
    path: '/login',
    meta: { title: '登录' },
    component: () => import('@/views/login/login.vue'),
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
