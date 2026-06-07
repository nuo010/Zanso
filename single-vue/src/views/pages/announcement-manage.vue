<template>
  <div class="announcement-manage-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>公告管理</h2>
            <p>发布平台公告，首页将展示已发布公告。</p>
          </div>
          <el-button type="primary" @click="openCreateDialog">新增公告</el-button>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="filters.status" clearable placeholder="公告状态" class="filter-select" @change="loadAnnouncements">
          <el-option label="已发布" value="active" />
          <el-option label="草稿" value="draft" />
        </el-select>
        <el-button type="primary" :loading="loading" @click="loadAnnouncements">查询</el-button>
      </div>

      <el-table :data="announcements" v-loading="loading" class="announcement-table">
        <el-table-column label="公告标题" min-width="220">
          <template #default="{ row }">
            <div class="announcement-title-block">
              <strong>{{ row.title }}</strong>
              <span>{{ row.content || '暂无内容' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" effect="light">
              {{ row.status === 'active' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="发布时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEditDialog(row)">编辑</el-button>
            <el-button link type="danger" @click="confirmDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next, total"
          @current-change="loadAnnouncements"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingAnnouncement ? '编辑公告' : '新增公告'" width="560px" align-center>
      <el-form :model="form" label-width="76px">
        <el-form-item label="标题">
          <el-input v-model="form.title" maxlength="160" show-word-limit />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="6" maxlength="2000" show-word-limit />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio-button label="active">发布</el-radio-button>
            <el-radio-button label="draft">草稿</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAnnouncement">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import {
  createAnnouncement,
  deleteAnnouncement,
  getAnnouncementManageList,
  updateAnnouncement,
} from '@/api/user';
import type { Announcement } from '@/store';
import { toast } from '@/util/util';

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const announcements = ref<Announcement[]>([]);
const editingAnnouncement = ref<Announcement | null>(null);
const page = ref(1);
const pageSize = 20;
const total = ref(0);
const filters = reactive({
  status: '',
});
const form = reactive({
  title: '',
  content: '',
  status: 'active' as 'draft' | 'active',
});

onMounted(loadAnnouncements);

async function loadAnnouncements() {
  loading.value = true;
  try {
    const res = await getAnnouncementManageList({
      page: page.value,
      pageSize,
      status: (filters.status || undefined) as 'draft' | 'active' | undefined,
    });
    announcements.value = res.data?.list || [];
    total.value = res.data?.total || 0;
    page.value = res.data?.page || page.value;
  } finally {
    loading.value = false;
  }
}

function openCreateDialog() {
  editingAnnouncement.value = null;
  form.title = '';
  form.content = '';
  form.status = 'active';
  dialogVisible.value = true;
}

function openEditDialog(row: Announcement) {
  editingAnnouncement.value = row;
  form.title = row.title;
  form.content = row.content || '';
  form.status = row.status;
  dialogVisible.value = true;
}

async function saveAnnouncement() {
  if (!form.title.trim()) {
    toast('公告标题不能为空', 'warning');
    return;
  }
  saving.value = true;
  try {
    const payload = {
      title: form.title.trim(),
      content: form.content.trim(),
      status: form.status,
    };
    if (editingAnnouncement.value) {
      await updateAnnouncement(editingAnnouncement.value.id, payload);
      toast('公告已更新', 'success');
    } else {
      await createAnnouncement(payload);
      toast('公告已创建', 'success');
    }
    dialogVisible.value = false;
    await loadAnnouncements();
  } finally {
    saving.value = false;
  }
}

async function confirmDelete(row: Announcement) {
  await ElMessageBox.confirm(`确认删除「${row.title}」？`, '删除公告', {
    type: 'warning',
    confirmButtonText: '确认',
    cancelButtonText: '取消',
  });
  await deleteAnnouncement(row.id);
  toast('公告已删除', 'success');
  await loadAnnouncements();
}

function formatDateTime(date?: string) {
  if (!date) return '--';
  return date.replace('T', ' ').slice(0, 19);
}
</script>

<style scoped lang="scss">
.announcement-manage-page {
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

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 18px;
  flex-wrap: wrap;
  align-items: center;
}

.filter-select {
  width: 180px;
}

.announcement-table {
  width: 100%;
}

.announcement-title-block strong,
.announcement-title-block span {
  display: block;
}

.announcement-title-block strong {
  color: #17315f;
}

.announcement-title-block span {
  margin-top: 4px;
  color: #6d82a7;
  font-size: 12px;
  line-height: 1.6;
}

.pagination-bar {
  padding-top: 18px;
}
</style>
