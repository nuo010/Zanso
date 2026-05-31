import router from '@/router/route';
import { getToken } from './auth';
import { hideFullLoading } from './util';

const whiteList = ['/login'];

router.beforeEach((to, _from, next) => {
  const token = getToken();
  if (!token && !whiteList.includes(to.path)) {
    return next({ path: '/login' });
  }
  document.title = String(to.meta.title || '');
  next();
});

router.afterEach(() => {
  hideFullLoading();
});

router.onError((error) => {
  console.error('路由错误:', error);
});
