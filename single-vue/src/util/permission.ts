import router from '@/router/route';
import { getToken } from './auth';
import { hideFullLoading, toast } from './util';
import { userMainStore } from '@/store';

const whiteList = ['/login'];

router.beforeEach(async (to, _from, next) => {
  if (to.meta.public) {
    document.title = String(to.meta.title || '');
    return next();
  }
  const token = getToken();
  if (!token && !whiteList.includes(to.path)) {
    return next({ path: '/login' });
  }
  const store = userMainStore();
  if (token && !store.user.id) {
    try {
      await store.loadProfile();
    } catch (error) {
      return next({ path: '/login' });
    }
  }
  if (to.meta.adminOnly && !store.isAdmin) {
    toast('只有管理员可以访问用户管理', 'warning');
    return next({ path: '/dashboard' });
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
