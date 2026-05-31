<template>
  <div class="share-page" v-loading="loading">
    <div v-if="detail" class="share-card">
      <section class="hero">
        <span class="eyebrow">Zanso Share</span>
        <h1>{{ detail.shareLink.title }}</h1>
        <p>{{ descriptionText }}</p>
      </section>

      <section class="meta-grid">
        <div class="meta-item">
          <span>发布用户</span>
          <strong>{{ detail.user.name }}</strong>
        </div>
        <div class="meta-item">
          <span>所属分类</span>
          <strong>{{ detail.category.name }}</strong>
        </div>
        <div class="meta-item">
          <span>分享对象</span>
          <strong>{{ detail.categoryItem?.name || detail.category.name }}</strong>
        </div>
        <div class="meta-item">
          <span>浏览次数</span>
          <strong>{{ detail.shareLink.viewCount }}</strong>
        </div>
      </section>

      <section v-if="isCategoryShare && detail.categoryItems?.length" class="section">
        <div class="section-header">
          <h2>分类项筛选</h2>
          <el-segmented v-model="selectedFilter" :options="filterOptions" />
        </div>
        <p v-if="selectedItemDescription" class="filter-description">{{ selectedItemDescription }}</p>
      </section>

      <section class="section">
        <h2>资源展示</h2>
        <div class="media-grid">
          <article v-for="item in filteredResources" :key="item.id" class="media-card">
            <img v-if="item.resourceType !== 'video'" :src="item.url" :alt="item.fileName" loading="lazy" />
            <video v-else :src="item.url" controls playsinline preload="metadata"></video>
            <div class="media-info">
              <span>{{ item.fileName }}</span>
              <span>{{ item.resourceType }}</span>
            </div>
          </article>
        </div>
        <el-empty v-if="filteredResources.length === 0" description="当前筛选条件下没有资源" />
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { getShareLinkDetail } from '@/api/user';
import { toast } from '@/util/util';

const route = useRoute();
const loading = ref(false);
const detail = ref<any>(null);
const selectedFilter = ref('all');

const isCategoryShare = computed(() => detail.value?.shareLink?.targetType === 'category');
const selectedCategoryItem = computed(() => {
  if (!detail.value?.categoryItems?.length || selectedFilter.value === 'all') return null;
  return detail.value.categoryItems.find((item: any) => item.id === selectedFilter.value) || null;
});
const descriptionText = computed(() => {
  if (!detail.value) return '';
  if (isCategoryShare.value) {
    return detail.value.category.description || detail.value.shareLink.description || '暂无描述';
  }
  return detail.value.categoryItem?.description || detail.value.category.description || detail.value.shareLink.description || '暂无描述';
});
const selectedItemDescription = computed(() => {
  if (!selectedCategoryItem.value) return '';
  return selectedCategoryItem.value.description || '';
});

const filterOptions = computed(() => {
  if (!detail.value?.categoryItems?.length) {
    return [{ label: '全部', value: 'all' }];
  }
  return [
    { label: '全部', value: 'all' },
    ...detail.value.categoryItems.map((item: any) => ({
      label: item.name,
      value: item.id,
    })),
  ];
});

const filteredResources = computed(() => {
  const list = detail.value?.resourceList || [];
  if (!isCategoryShare.value) return list;
  if (selectedFilter.value === 'all') return list;
  return list.filter((item: any) => item.categoryItemId === selectedFilter.value);
});

onMounted(async () => {
  const code = String(route.params.code || '');
  if (!code) {
    toast('分享码不存在', 'error');
    return;
  }
  loading.value = true;
  try {
    const res = await getShareLinkDetail(code);
    detail.value = res.data;
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped lang="scss">
.share-page {
  min-height: 100vh;
  background:
    radial-gradient(circle at top left, rgba(139, 180, 255, 0.18), transparent 30%),
    radial-gradient(circle at right center, rgba(194, 225, 255, 0.18), transparent 24%),
    linear-gradient(180deg, #f8fbff 0%, #eef4ff 100%);
  padding: 18px 12px 40px;
}

.share-card {
  max-width: 760px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(123, 162, 255, 0.16);
  border-radius: 28px;
  overflow: hidden;
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
}

.hero {
  padding: 22px 20px 18px;
  background: linear-gradient(135deg, #2f6bff, #78bcff);
  color: #fdfefe;
}

.eyebrow {
  display: inline-block;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.hero h1 {
  margin: 14px 0 8px;
  font-size: 28px;
  line-height: 1.2;
}

.hero p {
  margin: 0;
  color: rgba(253, 254, 255, 0.88);
  font-size: 14px;
  line-height: 1.7;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  padding: 14px 20px;
  background: #f8fbff;
}

.meta-item {
  padding: 10px 12px;
  border-radius: 14px;
  background: #eef5ff;
}

.meta-item span {
  display: block;
  font-size: 11px;
  color: #6d82a7;
  margin-bottom: 4px;
}

.meta-item strong {
  font-size: 13px;
  color: #17315f;
  word-break: break-all;
}

.section {
  padding: 20px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.section h2 {
  margin: 0 0 14px;
  font-size: 18px;
  color: #17315f;
}

.filter-description {
  margin: 14px 0 0;
  padding: 12px 14px;
  border-radius: 14px;
  background: #f4f8ff;
  color: #5f76a0;
  line-height: 1.7;
  font-size: 13px;
}

.media-grid {
  display: grid;
  gap: 14px;
}

.media-card {
  overflow: hidden;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(116, 153, 230, 0.16);
  box-shadow: 0 14px 32px rgba(60, 102, 190, 0.08);
}

.media-card img,
.media-card video {
  display: block;
  width: 100%;
  background: linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
}

.media-card video {
  max-height: 72vh;
}

.media-info {
  padding: 12px 14px;
  font-size: 13px;
  color: #6d82a7;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

@media (min-width: 768px) {
  .share-page {
    padding: 30px 24px 56px;
  }

  .hero {
    padding: 30px 28px 22px;
  }

  .hero h1 {
    font-size: 34px;
  }

  .section {
    padding: 24px 28px;
  }

  .meta-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
    padding: 16px 28px;
  }
}
</style>
