<template>
  <div class="login-bg">
    <div class="login-header">
      <div class="logo"></div>
    </div>
    <div class="login-main">
      <div class="login-left">
        <div class="quote-block">singlevue</div>
      </div>
      <div class="login-right">
        <Account />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue';
import { getTitle } from '@/util/util';

const Account = defineAsyncComponent(() => import('@/views/login/component/account.vue'));
</script>

<style scoped lang="scss">
.login-bg {
  min-height: 100vh;
  width: 100vw;
  background: #f5f6f8;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.login-header {
  position: absolute;
  top: 32px;
  left: 48px;
  z-index: 2;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}


.login-main {
  display: flex;
  align-items: stretch;
  justify-content: center;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.08);
  min-width: 880px;
  min-height: 520px;
  overflow: hidden;
}

.login-left {
  flex: 1;
  min-width: 380px;
  background: linear-gradient(135deg, #5eb8b0 0%, #7ec8c2 40%, #b8a8d4 80%, #d4c4e8 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.quote-block {
  color: #fff;
  font-size: 42px;
  font-weight: 700;
  letter-spacing: 2px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.login-right {
  flex: 1;
  min-width: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 56px;
}

@media (max-width: 900px) {
  .login-header {
    left: 24px;
    top: 24px;
  }

  .login-main {
    flex-direction: column;
    min-width: 0;
    width: calc(100vw - 32px);
    max-width: 420px;
    min-height: auto;
  }

  .login-left {
    min-width: 0;
    padding: 32px 24px;
  }

  .login-right {
    min-width: 0;
    padding: 32px 24px;
  }
}
</style>
