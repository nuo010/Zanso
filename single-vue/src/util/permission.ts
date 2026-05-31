import router from '@/router/route';
import { getToken } from './auth';
import { toast, showFullLoading, hideFullLoading } from './util';
import loadingService from "@/components/GlobalLoading";



const whiteList = ['/login', '/replace', '/555'];

let loadingShown = false

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  // 显示loading
  // showFullLoading();
  // if (!loadingShown){
  //   loadingService.show("页面加载中...")
  //   loadingShown = true
  // }
  // showAreaLoading('正在处理数据...')

  const token = getToken();
  
  // console.log('路由守卫检查:', { to: to.path, from: from.path, hasToken: !!token });
  
  // 没有token
  if (!token){
    // 没有token ，未登录，允许访问白名单
    if (whiteList.includes(to.path)){
      // console.log('白名单路径，允许访问:', to.path);
      return next();
    }else {
      // 没有token  非白名单 需要跳转到555 界面提示用户重新进入
      // console.log('无token，跳转到555页面');
      // return next({ path: '/555' });
      return next({ path: '/login' });
    }
  }

  // 有token的情况
  // console.log('有token，正常访问:', to.path);
  // 设置界面title
  document.title = String(to.meta.title ? to.meta.title : '');
  return next();
});

router.afterEach((to, from) => {
  hideFullLoading()
  // hideAreaLoading()
  router.isReady().then(() => {
    // setTimeout(() =>{
    //   if(loadingShown) {
    //     loadingService.hide()
    //     loadingShown =false
    //   }
    // },100)
  })
});

// 处理路由错误
router.onError((error) =>{
  // if(loadingShown) {
  //   loadingService.hide()
  //   loadingShown =false
  // }
  console.error('路由错误:', error)
})
