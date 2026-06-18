<template>
  <div class="login-page">
    <div class="hero-panel">
      <span class="hero-tag">Zanso</span>
      <h1>数字内容管理平台</h1>
      <p>面向企业内容展示与资源分发场景，提供用户、展册、资源和分享链接的一体化管理能力。</p>
    </div>

    <el-card class="login-card" shadow="never">
      <template #header>
        <div class="card-header">
          <h2>{{ authMode === 'login' ? '账号登录' : '注册账号' }}</h2>
          <span>{{ authMode === 'login' ? '使用已授权账号访问管理平台' : '创建平台用户账号并完善基础信息' }}</span>
        </div>
      </template>

      <el-form v-if="authMode === 'login'" :model="form" label-position="top" @submit.prevent="handleLogin">
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

      <el-form v-else :model="registerForm" label-position="top" @submit.prevent="handleRegister">
        <el-form-item label="用户名称">
          <el-input v-model="registerForm.name" placeholder="请输入用户名称" />
        </el-form-item>
        <el-form-item label="登录账号">
          <el-input v-model="registerForm.loginName" placeholder="仅支持字母和数字，至少 6 位" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="registerForm.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="registerForm.email" placeholder="可选" />
        </el-form-item>
        <el-button type="primary" class="submit-btn" :loading="submitting" @click="handleRegister">
          注册
        </el-button>
      </el-form>

      <p class="helper-text">
        {{ authMode === 'login' ? '暂无账号？' : '已有账号？' }}
        <el-button link type="primary" class="mode-switch" @click="toggleAuthMode">
          {{ authMode === 'login' ? '立即注册' : '返回登录' }}
        </el-button>
      </p>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { createUser, loginUser, getCurrentUser } from '@/api/user';
import { setToken } from '@/util/auth';
import { toast } from '@/util/util';
import { userMainStore } from '@/store';

const router = useRouter();
const store = userMainStore();
const submitting = ref(false);
const authMode = ref<'login' | 'register'>('login');
const form = reactive({
  loginName: '',
  password: '',
});
const registerForm = reactive({
  name: '',
  loginName: '',
  password: '',
  email: '',
});

const loginNamePattern = /^[A-Za-z0-9]{6,}$/;
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function toggleAuthMode() {
  authMode.value = authMode.value === 'login' ? 'register' : 'login';
}

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

async function handleRegister() {
  const name = registerForm.name.trim();
  const loginName = registerForm.loginName.trim();
  const password = registerForm.password.trim();
  const email = registerForm.email.trim();
  if (!name || !loginName || !password) {
    toast('用户名称、登录账号和密码不能为空', 'warning');
    return;
  }
  if (!loginNamePattern.test(loginName)) {
    toast('登录账号只能使用字母和数字，且至少 6 位', 'warning');
    return;
  }
  if (email && !emailPattern.test(email)) {
    toast('请输入正确的邮箱格式', 'warning');
    return;
  }
  submitting.value = true;
  try {
    await createUser({
      name,
      loginName,
      password,
      email,
    });
    form.loginName = loginName;
    form.password = '';
    authMode.value = 'login';
    toast('注册成功，请登录', 'success');
  } catch (error) {
    console.error('注册失败:', error);
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
    radial-gradient(circle at top left, rgba(139, 180, 255, 0.18), transparent 32%),
    radial-gradient(circle at right center, rgba(194, 225, 255, 0.18), transparent 24%),
    linear-gradient(180deg, #f8fbff 0%, #eef4ff 100%);
}

.hero-panel {
  padding: 72px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  color: #17315f;
}

.hero-tag {
  display: inline-flex;
  width: fit-content;
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(47, 107, 255, 0.10);
  color: #2f6bff;
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
  color: #6d82a7;
}

.login-card {
  width: min(460px, calc(100% - 40px));
  align-self: center;
  justify-self: center;
  border-radius: 28px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(14px);
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
}

.card-header h2 {
  margin: 0;
  color: #17315f;
}

.card-header span {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  color: #6d82a7;
}

.submit-btn {
  width: 100%;
  margin-top: 12px;
  height: 44px;
}

.helper-text {
  margin-top: 18px;
  color: #7f90af;
  font-size: 12px;
  line-height: 1.6;
}

.mode-switch {
  padding: 0 2px;
  vertical-align: baseline;
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
