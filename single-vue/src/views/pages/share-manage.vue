<template>
  <div class="share-manage-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>分享链接管理</h2>
            <p>统一查看所有分享链接，支持按展册和分类筛选，过期时间和访问量一眼能看明白。</p>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="filters.collectionId" clearable placeholder="展册" class="filter-select" @change="handleCollectionChange">
          <el-option v-for="item in store.categories" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-select v-model="filters.categoryId" clearable placeholder="分类" class="filter-select">
          <el-option v-for="item in categoryOptions" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-button type="primary" @click="loadShareLinks">查询</el-button>
      </div>

      <el-table :data="shareLinks" v-loading="loading" class="share-table">
        <el-table-column label="二级标题" min-width="240">
          <template #default="{ row }">
            <div class="share-title-block">
              <strong>
                {{
                  row.targetType === 'category'
                    ? `${row.collectionName} -> ${row.categoryName || '未命名分类'}`
                    : row.collectionName
                }}
              </strong>
              <span>{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="分享类型" width="140">
          <template #default="{ row }">
            <span class="share-type-tag">
              {{ row.targetType === 'category' ? '分类分享' : '展册分享' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="shareCode" label="分享码" width="140" />
        <el-table-column prop="viewCount" label="访问量" width="100" />
        <el-table-column label="到期时间" width="150">
          <template #default="{ row }">
            {{ row.expiresAt || '长期有效' }}
          </template>
        </el-table-column>
        <el-table-column label="链接" min-width="260" show-overflow-tooltip>
          <template #default="{ row }">
            <a :href="row.shareUrl" target="_blank" class="share-url">{{ row.shareUrl }}</a>
          </template>
        </el-table-column>
        <el-table-column label="二维码" width="110" align="center" header-align="center">
          <template #default="{ row }">
            <button class="qr-button" type="button" title="点击查看二维码" @click="openQrDialog(row)">
              <qrcode-vue :value="row.shareUrl" :size="44" level="M" render-as="svg" />
            </button>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="96" align="center" header-align="center" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" class="delete-link" @click="confirmDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="qrDialogVisible" title="分享二维码" width="360px" align-center>
      <div v-if="currentQrShare" class="qr-dialog-body">
        <div class="qr-large">
          <qrcode-vue :value="currentQrShare.shareUrl" :size="240" level="M" render-as="svg" />
        </div>
        <strong>{{ currentQrShare.title || currentQrShare.categoryName || '分享链接' }}</strong>
        <a :href="currentQrShare.shareUrl" target="_blank">{{ currentQrShare.shareUrl }}</a>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import QrcodeVue from 'qrcode.vue';
import { deleteShareLink, getCategoryDetail, getShareLinkList } from '@/api/user';
import { userMainStore } from '@/store';
import { toast } from '@/util/util';

const store = userMainStore();
const loading = ref(false);
const shareLinks = ref<any[]>([]);
const qrDialogVisible = ref(false);
const currentQrShare = ref<any>(null);
const categoryDetailMap = reactive<Record<string, any>>({});
const filters = reactive({
  collectionId: '',
  categoryId: '',
});

const categoryOptions = computed(() => {
  if (!filters.collectionId) return [];
  return categoryDetailMap[filters.collectionId]?.categories || [];
});

onMounted(async () => {
  if (!store.categories.length) {
    await store.loadCategories();
  }
  await loadShareLinks();
});

async function handleCollectionChange() {
  filters.categoryId = '';
  if (filters.collectionId && !categoryDetailMap[filters.collectionId]) {
    const res = await getCategoryDetail(filters.collectionId);
    categoryDetailMap[filters.collectionId] = res.data;
  }
}

async function loadShareLinks() {
  loading.value = true;
  try {
    const res = await getShareLinkList({
      collectionId: filters.collectionId || undefined,
      categoryId: filters.categoryId || undefined,
    });
    shareLinks.value = res.data || [];
  } finally {
    loading.value = false;
  }
}

async function confirmDelete(id: string) {
  await ElMessageBox.confirm('删除后该分享链接将无法继续访问，确认继续？', '二次确认', {
    type: 'warning',
    confirmButtonText: '确认',
    cancelButtonText: '取消',
  });
  await deleteShareLink(id);
  toast('分享链接已删除', 'success');
  await loadShareLinks();
}

function openQrDialog(row: any) {
  currentQrShare.value = row;
  qrDialogVisible.value = true;
}
</script>

<style scoped lang="scss">
.share-manage-page {
  display: grid;
}

.panel-card {
  border-radius: 28px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
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
  width: 220px;
}

.share-title-block strong {
  display: block;
  color: #17315f;
}

.share-title-block span {
  color: #6d82a7;
  font-size: 12px;
  margin-top: 4px;
}

.share-type-tag {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 999px;
  background: rgba(47, 107, 255, 0.08);
  color: #2f6bff;
  font-size: 12px;
  font-weight: 600;
}

.share-table {
  width: 100%;
}

.share-table :deep(.el-table__cell) {
  padding: 12px 0;
}

.share-table :deep(.el-table-fixed-column--right) {
  background: #fff;
}

.share-url {
  display: block;
  color: #2f6bff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.qr-button {
  display: inline-flex;
  padding: 6px;
  border: 1px solid rgba(47, 107, 255, 0.14);
  border-radius: 12px;
  background: #fff;
  cursor: zoom-in;
  box-shadow: 0 8px 18px rgba(60, 102, 190, 0.08);
  transition:
    transform 0.18s ease,
    box-shadow 0.18s ease;
}

.qr-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 24px rgba(60, 102, 190, 0.14);
}

.qr-dialog-body {
  display: grid;
  justify-items: center;
  gap: 12px;
  text-align: center;
}

.qr-large {
  padding: 14px;
  border-radius: 20px;
  background: #fff;
  border: 1px solid rgba(47, 107, 255, 0.12);
  box-shadow: 0 18px 38px rgba(36, 84, 170, 0.12);
}

.qr-dialog-body strong {
  color: #17315f;
}

.qr-dialog-body a {
  max-width: 100%;
  color: #2f6bff;
  word-break: break-all;
  font-size: 12px;
}

.delete-link {
  padding: 0;
  font-weight: 600;
}
</style>
