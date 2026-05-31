import { createApp, DirectiveBinding } from 'vue';
import { createPinia } from 'pinia';
import piniaPluginPersist from 'pinia-plugin-persist';
import ElementPlus from 'element-plus';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import 'element-plus/dist/index.css';
import 'nprogress/nprogress.css';
import 'virtual:windi.css';

import App from './App.vue';
import router from '@/router/route';
import '@/util/permission';
import { LoadingPlugin } from '@/components/GlobalLoading';

const app = createApp(App);

app.use(createPinia().use(piniaPluginPersist));
app.use(ElementPlus, { locale: zhCn });
app.use(router);
app.use(LoadingPlugin);

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

app.mount('#app');

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
