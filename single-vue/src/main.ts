import { createApp, DirectiveBinding } from 'vue';
import App from './App.vue';
import ElementPlus from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import 'element-plus/dist/index.css';
import router from '@/router/route';

import 'virtual:windi.css';

import '@/util/permission';

const app = createApp(App);
import { createPinia } from 'pinia';
// pinia 持久化插件,解决刷新界面,数据不丢失
import piniaPluginPersist from 'pinia-plugin-persist';

app.use(createPinia().use(piniaPluginPersist));
// 挂载ElementPlus
app.use(ElementPlus, {
  locale: zhCn,
});
// 注册elementplus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}


// 初始化性能监控，指定后端上报地址
// const monitor = new PerfMonitor(router, {
//   reportUrl: getBaseURL() + serverLocation+ '/singleSystem/addPerf',  // 后端接收接口
//   reportInterval: 5000,             // 5秒批量上报一次
//   maxQueueSize: 500,
//   runtimeInterval: 30000,
//   debug: true,
//   sysName: "系统子应用",
// })
//
// // 包装 fetch
// PerfMonitor.wrapFetch(monitor)
// // 包装 axios
// PerfMonitor.wrapAxios(axios, monitor)


// 挂载路由
app.use(router);

// 加载样式条
import 'nprogress/nprogress.css';
import {getBaseURL} from "@/util/util";
import PerfMonitor from "@/util/perf-monitor";
import axios from "axios";
import {serverLocation} from "@/util/constants";
import {LoadingPlugin} from "@/components/GlobalLoading";



app.use(LoadingPlugin)



app.mount('#app');
// 防止el-button重复点击,直接在elbutton添加 v-preventReClick 指令
//

app.directive('preventReClick', (el, binding) => {
  function preventReClickFun(
    elValue: { disabled: boolean },
    bindingValue: DirectiveBinding<any>
  ) {
    if (!elValue.disabled) {
      elValue.disabled = true;
      setTimeout(() => {
        elValue.disabled = false;
      }, bindingValue.value || 3000);
    }
  }
  el.addEventListener('click', () => preventReClickFun(el, binding));
  binding.dir.unmounted = function () {
    el.removeEventListener('click', () => preventReClickFun(el, binding));
  };
});
