<template>
  <div class="login-page">
    <div class="hero-panel">
      <span class="hero-tag">Zanso</span>
      <h1>资源分享后台</h1>
      <p>统一管理用户、分类、资源和分享链接，给外链展示留个正经后台。</p>
    </div>

    <el-card class="login-card" shadow="never">
      <template #header>
        <div class="card-header">
          <h2>账号登录</h2>
          <span>使用已创建的用户账号进入后台</span>
        </div>
      </template>

      <el-form :model="form" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="登录账号">
          <el-input v-model="form.loginName" placeholder="请输入 loginName" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-button type="primary" class="submit-btn" :loading="submitting" @click="handleLogin">
          登录
        </el-button>
      </el-form>

      <p class="helper-text">
        如果还没有账号，先调用后端 `POST /api/platform/users` 创建一个用户。
      </p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { loginUser, getCurrentUser } from '@/api/user';
import { setToken } from '@/util/auth';
import { toast } from '@/util/util';
import { userMainStore } from '@/store';

const router = useRouter();
const store = userMainStore();
const submitting = ref(false);
const form = reactive({
  loginName: '',
  password: '',
});

async function handleLogin() {
  if (!form.loginName || !form.password) {
    toast('登录账号和密码不能为空', 'warning');
    return;
  }
  submitting.value = true;
  try {
    const loginRes = await loginUser({
      loginName: form.loginName,
      password: form.password,
    });
    setToken(loginRes.data.token);
    const profileRes = await getCurrentUser();
    store.setUserStore(profileRes.data);
    await store.loadCategories();
    toast('登录成功', 'success');
    await router.push('/dashboard');
  } catch (error) {
    console.error('登录失败:', error);
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  background:
    radial-gradient(circle at top left, rgba(183, 217, 189, 0.9), transparent 32%),
    linear-gradient(135deg, #f4efe2, #e8f1eb);
}

.hero-panel {
  padding: 72px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: #1e3428;
}

.hero-tag {
  display: inline-flex;
  width: fit-content;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(49, 92, 69, 0.12);
  color: #315c45;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.hero-panel h1 {
  font-size: 56px;
  line-height: 1.05;
  margin: 22px 0 16px;
}

.hero-panel p {
  max-width: 480px;
  font-size: 18px;
  color: #587061;
}

.login-card {
  width: min(460px, calc(100% - 40px));
  align-self: center;
  justify-self: center;
  border-radius: 28px;
  border: 1px solid rgba(32, 53, 41, 0.08);
  background: rgba(255, 253, 249, 0.92);
  backdrop-filter: blur(14px);
}

.card-header h2 {
  margin: 0;
  color: #203529;
}

.card-header span {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  color: #69806f;
}

.submit-btn {
  width: 100%;
  margin-top: 12px;
  height: 44px;
}

.helper-text {
  margin-top: 18px;
  color: #7a877f;
  font-size: 12px;
  line-height: 1.6;
}

@media (max-width: 960px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .hero-panel {
    padding: 40px 28px 16px;
  }

  .hero-panel h1 {
    font-size: 40px;
  }

  .login-card {
    margin: 0 0 32px;
  }
}
</style>
