<template>
  <el-dialog
      v-model="dialogVisible"
      title="修改密码"
      width="400px"
      :close-on-click-modal="false"
      @closed="onClose"
  >
    <el-form
        :model="form"
        :rules="changePasswordRule"
        ref="formRef"
        size="large"
        label-width="90px"
    >
      <el-form-item label="原密码" prop="oldPassword">
        <el-input
            v-model="form.oldPassword"
            type="password"
            placeholder="请输入原密码"
            show-password
        />
      </el-form-item>

      <el-form-item label="新密码" prop="newPassword">
        <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
        />
      </el-form-item>

      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请确认新密码"
            show-password
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit" :loading="loading">
          确定
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive } from "vue";
import { ElMessage, FormInstance } from "element-plus";
import {changePassword} from "@/api/user";
import {encryptAes} from "@/util/util";

// 弹窗显示控制
const dialogVisible = ref(false);

// 表单数据
const form = reactive({
  oldPassword: "",
  newPassword: "",
  confirmPassword: "",
});

// 表单引用
const formRef = ref<FormInstance>();

// 加载状态
const loading = ref(false);

// 表单验证规则
const changePasswordRule = {
  oldPassword: [{ required: true, message: "请输入原密码", trigger: "blur" }],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    // { validator: validatePassword, trigger: "blur" }, // 可选：6位以上
  ],
  confirmPassword: [
    { required: true, message: "请确认新密码", trigger: "blur" },
    {
      validator: (_rule: any, value: string) => {
        if (value !== form.newPassword) {
          return Promise.reject("两次输入的密码不一致");
        }
        return Promise.resolve();
      },
      trigger: "blur",
    },
  ],
};
// 打开弹窗的方法
const open = () => {
  dialogVisible.value = true;
};

// 关闭弹窗时重置表单
const onClose = () => {
  formRef.value?.resetFields();
  const keys = Object.keys(form) as Array<keyof typeof form>;
  keys.forEach((key) => {
    if (typeof form[key] === 'string') {
      form[key] = "" as string;
    }
  });
};


// 提交修改
const submit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    loading.value = true;
    const {data} = await changePassword({oldPwd: encryptAes(form.oldPassword) , newPwd: form.newPassword,})

    ElMessage.success("密码修改成功，请重新登录");
    loading.value = false;
    dialogVisible.value = false;

    // 修改成功后退出登录
    // 可调用你已有的退出逻辑
    // window.location.href = "/login"; // 或调用 logout()
  } catch (err) {
    loading.value = false;
    // 验证失败或请求失败，ElMessage 会自动提示
  }
};

// 暴露 open 方法给父组件调用
defineExpose({
  open,
});
</script>
