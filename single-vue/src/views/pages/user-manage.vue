<template>
  <div class="user-manage-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>用户管理</h2>
            <p>管理员可以查看平台用户，并调整用户角色。别手滑，权限这玩意儿改错了挺闹心。</p>
          </div>
          <el-button type="primary" @click="loadUsers">刷新</el-button>
        </div>
      </template>

      <el-table :data="users" v-loading="loading" class="user-table">
        <el-table-column prop="name" label="用户名称" min-width="160" />
        <el-table-column prop="loginName" label="登录账号" min-width="160" />
        <el-table-column label="联系方式" min-width="180">
          <template #default="{ row }">
            <div class="contact-block">
              <strong>{{ row.contactName || '--' }}</strong>
              <span>{{ row.contactPhone || '--' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="当前角色" width="140">
          <template #default="{ row }">
            <el-tag :type="getPrimaryRole(row) === 'admin' ? 'danger' : 'primary'" effect="light">
              {{ getPrimaryRole(row) === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="修改角色" width="180" align="center">
          <template #default="{ row }">
            <el-select
              :model-value="getPrimaryRole(row)"
              :disabled="row.id === store.user.id"
              size="small"
              class="role-select"
              @change="(roleCode: 'admin' | 'user') => handleRoleChange(row, roleCode)"
            >
              <el-option label="普通用户" value="user" />
              <el-option label="管理员" value="admin" />
            </el-select>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import { getUserList, updateUserRole } from '@/api/user';
import { userMainStore, type PlatformUser } from '@/store';
import { toast } from '@/util/util';

const store = userMainStore();
const loading = ref(false);
const users = ref<PlatformUser[]>([]);

onMounted(loadUsers);

async function loadUsers() {
  loading.value = true;
  try {
    const res = await getUserList();
    users.value = res.data || [];
  } finally {
    loading.value = false;
  }
}

function getPrimaryRole(user: PlatformUser): 'admin' | 'user' {
  return user.roleCodes?.includes('admin') ? 'admin' : 'user';
}

async function handleRoleChange(user: PlatformUser, roleCode: 'admin' | 'user') {
  const oldRoleCode = getPrimaryRole(user);
  if (oldRoleCode === roleCode) return;

  await ElMessageBox.confirm(
    `确认把「${user.name}」改为${roleCode === 'admin' ? '管理员' : '普通用户'}？`,
    '修改角色',
    {
      type: 'warning',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
    }
  );

  const res = await updateUserRole(user.id, { roleCode });
  const index = users.value.findIndex((item) => item.id === user.id);
  if (index >= 0) {
    users.value[index] = res.data;
  }
  toast('用户角色已更新', 'success');
}
</script>

<style scoped lang="scss">
.user-manage-page {
  display: grid;
}

.panel-card {
  border-radius: 28px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.panel-header h2 {
  margin: 0;
  color: #17315f;
}

.panel-header p {
  margin: 6px 0 0;
  color: #6d82a7;
}

.user-table {
  width: 100%;
}

.contact-block strong,
.contact-block span {
  display: block;
}

.contact-block strong {
  color: #17315f;
}

.contact-block span {
  margin-top: 4px;
  color: #6d82a7;
  font-size: 12px;
}

.role-select {
  width: 120px;
}
</style>
