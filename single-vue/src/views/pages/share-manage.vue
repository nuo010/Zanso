<template>
  <div class="share-manage-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>分享链接管理</h2>
            <p>统一查看所有分享链接，支持按分类和分类项筛选，过期时间和访问量一眼能看明白。</p>
          </div>
        </div>
      </template>

      <div class="filter-bar">
        <el-select v-model="filters.categoryId" clearable placeholder="按分类筛选" @change="handleCategoryChange">
          <el-option v-for="item in store.categories" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-select v-model="filters.categoryItemId" clearable placeholder="按分类项筛选">
          <el-option v-for="item in categoryItemOptions" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
        <el-button type="primary" @click="loadShareLinks">查询</el-button>
      </div>

      <el-table :data="shareLinks" v-loading="loading" class="share-table">
        <el-table-column prop="title" label="分享标题" min-width="200" />
        <el-table-column label="分享对象" min-width="220">
          <template #default="{ row }">
            <div class="share-target">
              <strong>{{ row.categoryName }}</strong>
              <span>{{ row.targetType === 'item' ? row.categoryItemName || '分类项分享' : '分类分享' }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="shareCode" label="分享码" width="140" />
        <el-table-column prop="viewCount" label="访问量" width="100" />
        <el-table-column label="到期时间" min-width="180">
          <template #default="{ row }">
            {{ row.expiresAt || '长期有效' }}
          </template>
        </el-table-column>
        <el-table-column label="链接" min-width="280">
          <template #default="{ row }">
            <a :href="row.shareUrl" target="_blank" class="share-url">{{ row.shareUrl }}</a>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { getCategoryDetail, getShareLinkList } from '@/api/user';
import { userMainStore } from '@/store';

const store = userMainStore();
const loading = ref(false);
const shareLinks = ref<any[]>([]);
const categoryDetailMap = reactive<Record<string, any>>({});
const filters = reactive({
  categoryId: '',
  categoryItemId: '',
});

const categoryItemOptions = computed(() => {
  if (!filters.categoryId) return [];
  return categoryDetailMap[filters.categoryId]?.categoryItems || [];
});

onMounted(async () => {
  if (!store.categories.length) {
    await store.loadCategories();
  }
  await loadShareLinks();
});

async function handleCategoryChange() {
  filters.categoryItemId = '';
  if (filters.categoryId && !categoryDetailMap[filters.categoryId]) {
    const res = await getCategoryDetail(filters.categoryId);
    categoryDetailMap[filters.categoryId] = res.data;
  }
}

async function loadShareLinks() {
  loading.value = true;
  try {
    const res = await getShareLinkList({
      categoryId: filters.categoryId || undefined,
      categoryItemId: filters.categoryItemId || undefined,
    });
    shareLinks.value = res.data || [];
  } finally {
    loading.value = false;
  }
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
}

.share-target strong {
  display: block;
  color: #17315f;
}

.share-target span {
  color: #6d82a7;
  font-size: 12px;
}

.share-url {
  color: #2f6bff;
  word-break: break-all;
}
</style>
