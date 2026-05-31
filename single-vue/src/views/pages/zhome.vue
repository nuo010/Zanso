<template>
  <div class="category-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>分类管理</h2>
            <p>创建分类、查看分类下的资源，并生成分享链接。</p>
          </div>
          <el-button type="primary" @click="createDialog = true">新建分类</el-button>
        </div>
      </template>

      <el-table :data="store.categories" v-loading="loading">
        <el-table-column prop="name" label="分类名称" min-width="180" />
        <el-table-column prop="description" label="描述" min-width="220" />
        <el-table-column prop="status" label="状态" width="120" />
        <el-table-column label="操作" width="280">
          <template #default="{ row }">
            <el-button link type="primary" @click="loadDetail(row.id)">查看详情</el-button>
            <el-button link type="success" @click="openUpload(row.id)">上传资源</el-button>
            <el-button link type="warning" @click="handleCreateShare(row.id, row.name, row.description)">生成分享</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card v-if="detail" class="panel-card detail-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h3>{{ detail.category.name }}</h3>
            <p>{{ detail.category.description || '暂无描述' }}</p>
          </div>
        </div>
      </template>

      <div class="detail-meta">
        <span>分享链接数量：{{ detail.shareLinks?.length || 0 }}</span>
        <span>资源数量：{{ detail.resourceList?.length || 0 }}</span>
      </div>

      <div class="resource-grid">
        <div v-for="item in detail.resourceList" :key="item.id" class="resource-card">
          <strong>{{ item.fileName }}</strong>
          <span>{{ item.resourceType }} · {{ item.mimeType || '未知类型' }}</span>
          <a :href="item.url" target="_blank">查看资源</a>
        </div>
      </div>

      <div v-if="detail.shareLinks?.length" class="share-list">
        <h4>分享链接</h4>
        <div v-for="item in detail.shareLinks" :key="item.id" class="share-item">
          <div>
            <strong>{{ item.title }}</strong>
            <span>{{ item.shareCode }}</span>
          </div>
          <a :href="buildShareLink(item.shareCode)" target="_blank">{{ buildShareLink(item.shareCode) }}</a>
        </div>
      </div>
    </el-card>

    <el-dialog v-model="createDialog" title="新建分类" width="460px">
      <el-form :model="createForm" label-position="top">
        <el-form-item label="分类名称">
          <el-input v-model="createForm.name" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleCreateCategory">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="uploadDialog" title="上传分类资源" width="460px">
      <el-upload
        drag
        :auto-upload="false"
        :limit="1"
        :on-change="handleFileChange"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div>拖拽文件到这里，或者点击选择</div>
      </el-upload>
      <template #footer>
        <el-button @click="uploadDialog = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="submitUpload">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { createCategory, createShareLink, getCategoryDetail, uploadCategoryResource } from '@/api/user';
import { userMainStore } from '@/store';
import { toast, getBaseURL } from '@/util/util';

const store = userMainStore();
const loading = ref(false);
const submitting = ref(false);
const uploading = ref(false);
const createDialog = ref(false);
const uploadDialog = ref(false);
const selectedCategoryId = ref('');
const selectedFile = ref<File | null>(null);
const detail = ref<any>(null);
const createForm = reactive({
  name: '',
  description: '',
});

async function refreshList() {
  loading.value = true;
  try {
    await store.loadCategories();
  } finally {
    loading.value = false;
  }
}

refreshList();

async function handleCreateCategory() {
  if (!createForm.name) {
    toast('分类名称不能为空', 'warning');
    return;
  }
  submitting.value = true;
  try {
    await createCategory({
      name: createForm.name,
      description: createForm.description,
      status: 'active',
    });
    toast('分类创建成功', 'success');
    createDialog.value = false;
    createForm.name = '';
    createForm.description = '';
    await refreshList();
  } finally {
    submitting.value = false;
  }
}

async function loadDetail(id: string) {
  const res = await getCategoryDetail(id);
  detail.value = res.data;
}

function openUpload(id: string) {
  selectedCategoryId.value = id;
  selectedFile.value = null;
  uploadDialog.value = true;
}

function handleFileChange(file: any) {
  selectedFile.value = file.raw;
}

async function submitUpload() {
  if (!selectedCategoryId.value || !selectedFile.value) {
    toast('请选择要上传的文件', 'warning');
    return;
  }
  const formData = new FormData();
  formData.append('file', selectedFile.value);
  uploading.value = true;
  try {
    await uploadCategoryResource(selectedCategoryId.value, formData);
    toast('资源上传成功', 'success');
    uploadDialog.value = false;
    await loadDetail(selectedCategoryId.value);
    await refreshList();
  } finally {
    uploading.value = false;
  }
}

async function handleCreateShare(categoryId: string, name: string, description?: string) {
  const res = await createShareLink({
    categoryId,
    title: name,
    description,
  });
  toast(`分享链接已生成：${res.data.shareUrl}`, 'success');
  await loadDetail(categoryId);
}

function buildShareLink(code: string) {
  return `${getBaseURL()}/share/${code}`;
}
</script>

<style scoped lang="scss">
.category-page {
  display: grid;
  gap: 20px;
}

.panel-card {
  border-radius: 24px;
  border: 1px solid rgba(32, 53, 41, 0.08);
  background: rgba(255, 252, 246, 0.92);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.panel-header h2,
.panel-header h3 {
  margin: 0;
  color: #203529;
}

.panel-header p {
  margin: 6px 0 0;
  color: #6c7f73;
}

.detail-meta {
  display: flex;
  gap: 18px;
  color: #587061;
  margin-bottom: 16px;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 14px;
}

.resource-card,
.share-item {
  padding: 14px 16px;
  border-radius: 18px;
  background: #f5f3eb;
}

.resource-card strong,
.share-item strong {
  display: block;
  color: #203529;
}

.resource-card span,
.share-item span {
  display: block;
  margin: 6px 0;
  color: #6c7f73;
  font-size: 13px;
}

.resource-card a,
.share-item a {
  color: #315c45;
  font-size: 13px;
}

.share-list {
  margin-top: 20px;
  display: grid;
  gap: 12px;
}
</style>
