<template>
  <div class="account-form-wrap">
    <el-form
      size="large"
      ref="formRef"
      :model="state.ruleForm"
      :rules="currentRules"
      @keydown.enter="handleSubmit"
      autocomplete="on"
      class="account-form"
    >
      <div class="form-header">
        <h1 class="form-title">{{ state.isLogin ? '账号登录' : '账号注册' }}</h1>
        <div class="form-switch">
          <el-button link @click="toggleMode" class="switch-button">
            {{ state.isLogin ? '没有账号？立即注册' : '已有账号？立即登录' }}
          </el-button>
        </div>
      </div>
      <el-form-item prop="username">
      <el-input
        text
        placeholder="用户名"
        v-model="state.ruleForm.username"
        clearable
        name="username"
        autocomplete="username"
      >
        <template #prefix>
          <el-icon><User /></el-icon>
        </template>
      </el-input>
      </el-form-item>
      <el-form-item prop="password">
      <el-input
        :type="state.isShowPassword ? 'text' : 'password'"
        placeholder="密码"
        v-model="state.ruleForm.password"
        name="password"
        autocomplete="current-password"
      >
        <template #prefix>
          <el-icon><Unlock /></el-icon>
        </template>
        <template #suffix>
          <el-icon style="cursor: pointer" @click="state.isShowPassword = !state.isShowPassword">
            <View v-if="state.isShowPassword" />
            <Hide v-else />
          </el-icon>
        </template>
      </el-input>
      </el-form-item>

      <!-- 注册时显示的确认密码字段 -->
      <el-form-item v-if="!state.isLogin" prop="confirmPassword">
      <el-input
        :type="state.isShowConfirmPassword ? 'text' : 'password'"
        placeholder="确认密码"
        v-model="state.ruleForm.confirmPassword"
        name="confirmPassword"
        autocomplete="new-password"
      >
        <template #prefix>
          <el-icon><Unlock /></el-icon>
        </template>
        <template #suffix>
          <el-icon style="cursor: pointer" @click="state.isShowConfirmPassword = !state.isShowConfirmPassword">
            <View v-if="state.isShowConfirmPassword" />
            <Hide v-else />
          </el-icon>
        </template>
      </el-input>
      </el-form-item>

      <!-- 注册时显示的邮箱字段 -->
      <el-form-item v-if="!state.isLogin" prop="email">
      <div class="email-input-group">
        <el-input
          text
          placeholder="邮箱"
          v-model="state.ruleForm.email"
          name="email"
          autocomplete="email"
          class="email-input"
        >
          <template #prefix>
            <el-icon><Message /></el-icon>
          </template>
        </el-input>
        <el-button
          type="primary"
          :disabled="!state.ruleForm.email || state.countdown > 0"
          @click="sendCode"
          :loading="state.loading.sendCode"
          class="send-code-btn"
        >
          {{ state.countdown > 0 ? `${state.countdown}s后重发` : '发送验证码' }}
        </el-button>
      </div>
      </el-form-item>

      <!-- 注册时显示的验证码字段 -->
      <el-form-item v-if="!state.isLogin" prop="code">
      <el-input
        text
        placeholder="请输入验证码"
        v-model="state.ruleForm.code"
        clearable
        name="verificationCode"
        autocomplete="one-time-code"
        maxlength="6"
      >
        <template #prefix>
          <el-icon><Key /></el-icon>
        </template>
      </el-input>
      </el-form-item>

      <!--      <el-checkbox v-model="state.isRember" label="记住用户名和密码" size="large" />-->

      <el-form-item class="submit-item">
        <el-button
          type="primary"
          class="login-content-submit"
          @click="handleSubmit"
          :loading="state.isLogin ? state.loading.signIn : state.loading.signUp"
        >
          {{ state.isLogin ? '登录' : '注册' }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {doLogin, getMenuTree, getUser, registerEmail, sendEmailCode} from '@/api/user';
import { encryptAes, toast } from '@/util/util';
import { userMainStore } from '@/store';
import { setToken } from '@/util/auth';
import {tokenKey} from "@/util/constants";
import {loginRule} from "@/util/validatorRule";

// 定义变量内容

let store = userMainStore(); //接收

const formRef = ref();
const route = useRoute();
const router = useRouter();
const state = reactive({
  isLogin: true, // 默认为登录模式
  isShowPassword: false,
  isShowConfirmPassword: false,
  ruleForm: {
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    code: '',
  },
  isRember: true,
  countdown: 0, // 倒计时
  loading: {
    signIn: false,
    signUp: false,
    sendCode: false,
  }
});

// 计算属性：根据当前模式返回对应的验证规则
const currentRules = computed(() => {
  // 不要动态切换校验规则，否则组件会再切换规则的时候立马进行一次校验
  return loginRule;
});
onMounted(() => {
  if (localStorage.getItem('loginInfo') != null) {
    state.isRember = true;
    const loginInfo = JSON.parse(localStorage.getItem('loginInfo') as string);
    state.ruleForm.username = loginInfo.username;
    state.ruleForm.password = loginInfo.password;
  } else {
    state.isRember = false;
  }
});
// 切换登录/注册模式
const toggleMode = () => {
  state.isLogin = !state.isLogin;
  // 清空表单
  state.ruleForm = {
    username: '',
    password: '',
    confirmPassword: '',
    email: '',
    code: '',

  };
  // 重置倒计时
  state.countdown = 0;
  // 重置表单验证
  if (formRef.value) {
    formRef.value.clearValidate();
  }
};

// 处理提交按钮点击
const handleSubmit = () => {
  if (state.isLogin) {
    onSignIn();
  } else {
    onSignUp();
  }
};

// 登录
const onSignIn = async () => {
  // 如果已经在加载中，直接返回，防止重复提交
  if (state.loading.signIn) return;
  formRef.value.validate((valid: any) => {
    if (!valid) {
      return;
    }
    state.loading.signIn = true;
    doLogin(state.ruleForm.username, encryptAes(state.ruleForm.password))
      .then(async (res: any) => {
        // 登录成功在保存密码
        if (state.isRember) {
          localStorage.setItem('loginInfo', JSON.stringify({
            username: state.ruleForm.username,
            password: state.ruleForm.password
          }));
        } else {
          localStorage.removeItem('loginInfo');
        }
        setToken(res.data[tokenKey]);
        const { data } = await getUser();
        store.setUserStore(data);
        // const menuListRes = await getMenuTree({isTree:1,type: import.meta.env.VITE_API_SYSTEM_ID})
        // store.setMenuTreeStore(menuListRes.data)
        await store.loadMenu()
        await router.push('/home');
      })
      .catch(() => {
        state.loading.signIn = false;
      });
  });
};

// 注册
const onSignUp = async () => {
  // 如果已经在加载中，直接返回，防止重复提交
  if (state.loading.signUp) return;
  formRef.value.validate((valid: any) => {
    if (!valid) {
      return;
    }
    state.loading.signUp = true;
    registerEmail(
      state.ruleForm.username, 
      state.ruleForm.password,
      state.ruleForm.email,
      state.ruleForm.code
    )
      .then(async (res: any) => {
        toast('注册成功！请登录', 'success');
        // 注册成功后切换到登录模式
        state.isLogin = true;
        state.ruleForm.email = '';
        state.ruleForm.code = '';
        state.countdown = 0;
        state.loading.signUp = false;
      })
      .catch((error: any) => {
        state.loading.signUp = false;
        toast(error.message || '注册失败，请重试', 'error');
      });
  });
};

// 发送验证码
const sendCode = async () => {
  if (!state.ruleForm.email) {
    toast('请先输入邮箱', 'warning');
    return;
  }
  
  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(state.ruleForm.email)) {
    toast('请输入有效的邮箱地址', 'warning');
    return;
  }
  
  if (state.loading.sendCode) return;
  
  state.loading.sendCode = true;
  sendEmailCode(state.ruleForm.email)
    .then((res: any) => {
      toast('验证码已发送，请查收', 'success');
      // 开始倒计时
      startCountdown();
    })
    .catch((error: any) => {
      toast(error.message || '发送失败，请重试', 'error');
    })
    .finally(() => {
      state.loading.sendCode = false;
    });
};

// 开始倒计时
const startCountdown = () => {
  state.countdown = 60; // 60秒倒计时
  const timer = setInterval(() => {
    state.countdown--;
    if (state.countdown <= 0) {
      clearInterval(timer);
    }
  }, 1000);
};
</script>

<style scoped lang="scss">
.account-form-wrap {
  width: 100%;
  max-width: 360px;
}

.form-header {
  margin-bottom: 28px;
}

.form-title {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a2e;
  margin: 0 0 12px;
  letter-spacing: -0.3px;
}

.form-switch {
  margin: 0;
}

.switch-button {
  color: #2563eb;
  font-size: 14px;
  padding: 0;

  &:hover {
    color: #3b82f6;
  }
}

:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-input__wrapper) {
  background-color: #f5f6f8;
  border-radius: 10px;
  box-shadow: none;
  padding: 4px 14px;
  min-height: 48px;

  &:hover {
    background-color: #eef0f2;
  }

  &.is-focus {
    background-color: #fff;
    box-shadow: 0 0 0 1px #2563eb;
  }
}

:deep(.el-input__inner) {
  color: #1a1a2e;
}

:deep(.el-input__inner::placeholder) {
  color: #9ca3af;
}

:deep(.el-form-item__error) {
  font-size: 12px;
}

.submit-item {
  margin-top: 28px;
  margin-bottom: 0;
}

.login-content-submit {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.5px;
  border-radius: 10px;
  background: #2563eb;
  border: none;

  &:hover {
    background: #3b82f6;
  }
}

.email-input-group {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.email-input {
  flex: 1;
}

.send-code-btn {
  white-space: nowrap;
  min-width: 110px;
  height: 48px;
  border-radius: 10px;
  font-weight: 500;
}
</style>
