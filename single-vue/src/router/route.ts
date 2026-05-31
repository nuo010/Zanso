import {createRouter, createWebHashHistory, createWebHistory, RouteRecordRaw} from 'vue-router';





// 所有的路由定义，只有定义了，才可以打开界面
const routes: RouteRecordRaw[] = [
    // 一级路由，只是一个壳子
  {
    path: '/',
    component: () => import('@/layouts/admin.vue'),
    redirect: '/home',
    children: [
        // 二级路由，二级路由才是项目真实需要配置的地址，所有的业务界面都是二级及二级以下的路由
        // 定义每个path 对应的组件
      {
        path: '/home',
        meta: { title: '首页' },
        component: () => import('@/views/pages/home.vue'),
      },
      {
        path: '/user',
        meta: { title: '子首页' },
        component: () => import('@/views/user/user.vue'),
      },
      {
        path: '/zhome',
        meta: { title: '子首页' },
        component: () => import('@/views/pages/zhome.vue'),
      },

    ],
  },
  {
    path: '/login',
    meta: { title: '登录页' },
    component: () => import('@/views/login/login.vue'),
  },
  {
    path: '/replace',
    meta: { title: '跳转中间页' },
    component: () => import('@/views/pages/replace/index.vue'),
  },
  {
    path: '/555',
    meta: { title: '555' },
    component: () => import('@/views/pages/555.vue'),
  },
  {
    path: '/gateway',
    meta: { title: '555' },
    component: () => import('@/views/pages/gateway.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    meta: { title: '404' },
    name: 'NotFound',
    component: () => import('@/views/pages/404.vue'),
  },
];

const router = createRouter({
  // 路由模式 createWebHashHistory 和 createWebHistory
  // hash 模式简单，history 需要服务端在nginx中对路径特殊处理，如果是多级nginx映射，有点麻烦
  history: createWebHashHistory(),
  routes,
});

export default router;
