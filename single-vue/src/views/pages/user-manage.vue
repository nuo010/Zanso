<template>
  <div class="user-manage-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>用户管理</h2>
            <p>集中维护平台用户信息与角色权限，支持管理员查看账号状态并进行权限配置。</p>
          </div>
          <el-button type="primary" @click="loadUsers">刷新</el-button>
        </div>
      </template>

      <el-table :data="users" v-loading="loading" class="user-table">
        <el-table-column prop="name" label="用户名称" min-width="160" />
        <el-table-column prop="loginName" label="登录账号" min-width="160" />
        <el-table-column label="邮箱" min-width="200">
          <template #default="{ row }">
            <span class="email-text">{{ row.email || '--' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="当前角色" width="140">
          <template #default="{ row }">
            <el-tag :type="getPrimaryRole(row) === 'admin' ? 'danger' : 'primary'" effect="light">
              {{ getPrimaryRole(row) === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="账号状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" effect="light">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="资源数量" width="120" align="center">
          <template #default="{ row }">
            <span class="resource-count">{{ row.resourceCount || 0 }}</span>
          </template>
        </el-table-column>
        <el-table-column label="占用空间" width="140">
          <template #default="{ row }">
            <span class="resource-size">{{ formatFileSize(row.fileSizeTotal) }}</span>
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
        <el-table-column label="账号操作" width="140" align="center">
          <template #default="{ row }">
            <el-button
              link
              :type="row.status === 'active' ? 'danger' : 'success'"
              :disabled="row.id === store.user.id"
              @click="handleStatusChange(row)"
            >
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadUsers"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import { getUserList, updateUserRole, updateUserStatus } from '@/api/user';
import { userMainStore, type PlatformUser } from '@/store';
import { toast } from '@/util/util';

const store = userMainStore();
const loading = ref(false);
const users = ref<PlatformUser[]>([]);
const page = ref(1);
const pageSize = 20;
const total = ref(0);

onMounted(loadUsers);

async function loadUsers() {
  loading.value = true;
  try {
    const res = await getUserList({
      page: page.value,
      pageSize,
    });
    users.value = res.data?.list || [];
    total.value = res.data?.total || 0;
    page.value = res.data?.page || page.value;
  } finally {
    loading.value = false;
  }
}

function getPrimaryRole(user: PlatformUser): 'admin' | 'user' {
  return user.roleCodes?.includes('admin') ? 'admin' : 'user';
}

function formatFileSize(size?: number) {
  const fileSize = Number(size || 0);
  if (!fileSize) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = fileSize;
  let unitIndex = 0;
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }
  const fixed = value >= 100 || unitIndex === 0 ? 0 : value >= 10 ? 1 : 2;
  return `${value.toFixed(fixed)} ${units[unitIndex]}`;
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

async function handleStatusChange(user: PlatformUser) {
  const nextStatus = user.status === 'active' ? 'inactive' : 'active';
  await ElMessageBox.confirm(
    `确认${nextStatus === 'active' ? '启用' : '禁用'}「${user.name}」？`,
    '修改账号状态',
    {
      type: nextStatus === 'active' ? 'success' : 'warning',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
    }
  );

  const res = await updateUserStatus(user.id, { status: nextStatus });
  const index = users.value.findIndex((item) => item.id === user.id);
  if (index >= 0) {
    users.value[index] = res.data;
  }
  toast(`账号已${nextStatus === 'active' ? '启用' : '禁用'}`, 'success');
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

.pagination-bar {
  display: flex;
  justify-content: flex-end;
  padding-top: 18px;
}

.email-text {
  color: #6d82a7;
  word-break: break-all;
}

.role-select {
  width: 120px;
}

.resource-count {
  font-weight: 700;
  color: #17315f;
}

.resource-size {
  color: #47658f;
  font-weight: 600;
}
</style>
