<template>
  <div class="share-page" v-loading="loading">
    <div v-if="errorMessage" class="share-card share-empty-card">
      <section class="empty-state">
        <span class="eyebrow">Zanso Share</span>
        <h1>分享不可用</h1>
        <p>{{ errorMessage }}</p>
      </section>
    </div>
    <div v-if="detail" class="share-card">
      <section class="hero">
        <span class="eyebrow">Zanso Share</span>
        <h1>{{ heroTitle }}</h1>
        <p>{{ descriptionText }}</p>
      </section>

      <section v-if="isCollectionShare && detail.categories?.length" class="section">
        <div class="section-header">
          <h2>分类筛选</h2>
        </div>
        <el-segmented v-model="selectedFilter" :options="filterOptions" class="filter-segmented" />
      </section>

      <section class="section">
        <h2>资源展示</h2>
        <div class="media-grid">
          <article v-for="item in filteredResources" :key="item.id" class="media-card">
            <img v-if="item.resourceType !== 'video'" :src="item.fileUrl" :alt="item.fileName" loading="lazy" />
            <video v-else :src="item.fileUrl" controls playsinline preload="metadata"></video>
            <div class="media-info">
              <span>{{ item.fileName }}</span>
              <span>{{ item.resourceType }}</span>
            </div>
          </article>
        </div>
        <el-empty v-if="filteredResources.length === 0" description="当前筛选条件下没有资源" />
      </section>

      <section class="meta-inline">
        <span>发布用户：{{ detail.user.name }}</span>
        <span>所属展册：{{ detail.collection.name }}</span>
        <span>分享对象：{{ selectedCategory?.name || detail.category?.name || detail.collection.name }}</span>
        <span>浏览次数：{{ detail.shareLink.viewCount }}</span>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';
import { getShareLinkDetail } from '@/api/user';
import { resolveResourceURL } from '@/util/util';

const route = useRoute();
const loading = ref(false);
const detail = ref<any>(null);
const errorMessage = ref('');
const selectedFilter = ref('all');

const isCollectionShare = computed(() => detail.value?.shareLink?.targetType === 'collection');
const selectedCategory = computed(() => {
  if (!detail.value?.categories?.length || selectedFilter.value === 'all') return null;
  return detail.value.categories.find((item: any) => item.id === selectedFilter.value) || null;
});
const heroTitle = computed(() => {
  if (!detail.value) return '';
  if (selectedCategory.value?.name) {
    return selectedCategory.value.name;
  }
  return detail.value.shareLink.title || detail.value.category?.name || detail.value.collection.name || '';
});
const descriptionText = computed(() => {
  if (!detail.value) return '';
  if (selectedCategory.value?.description) {
    return selectedCategory.value.description;
  }
  if (isCollectionShare.value) {
    return detail.value.collection.description || detail.value.shareLink.description || '暂无描述';
  }
  return detail.value.category?.description || detail.value.collection.description || detail.value.shareLink.description || '暂无描述';
});

const filterOptions = computed(() => {
  if (!detail.value?.categories?.length) {
    return [{ label: '全部', value: 'all' }];
  }
  return [
    { label: '全部', value: 'all' },
    ...detail.value.categories.map((item: any) => ({
      label: item.name,
      value: item.id,
    })),
  ];
});

const filteredResources = computed(() => {
  const list = detail.value?.resourceList || [];
  if (!isCollectionShare.value) return list;
  if (selectedFilter.value === 'all') return list;
  return list.filter((item: any) => item.categoryId === selectedFilter.value);
});

onMounted(async () => {
  const code = String(route.params.code || '');
  if (!code) {
    errorMessage.value = '分享链接不存在或分享码无效。';
    return;
  }
  loading.value = true;
  try {
    const res = await getShareLinkDetail(code);
    detail.value = {
      ...res.data,
      resourceList: normalizeResourceList(res.data?.resourceList || []),
    };
    errorMessage.value = '';
  } catch (error: any) {
    const message = error?.message || '分享链接不存在';
    if (message.includes('过期')) {
      errorMessage.value = '这个分享链接已过期，暂时无法查看。';
    } else if (message.includes('不可查看') || message.includes('不可看')) {
      errorMessage.value = '当前分享内容已被关闭查看权限。';
    } else {
      errorMessage.value = '这个分享链接不存在，可能已被删除。';
    }
  } finally {
    loading.value = false;
  }
});

function normalizeResourceList(resourceList: any[]) {
  return resourceList.map((item) => ({
    ...item,
    fileUrl: resolveResourceURL(item.url || item.storagePath || ''),
  }));
}
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

.share-empty-card {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 320px;
}

.empty-state {
  padding: 42px 28px;
  text-align: center;
}

.empty-state h1 {
  margin: 18px 0 10px;
  font-size: 32px;
  color: #17315f;
}

.empty-state p {
  margin: 0;
  color: #6d82a7;
  font-size: 15px;
  line-height: 1.8;
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

.section {
  padding: 20px;
}

.section-header {
  display: block;
}

.section h2 {
  margin: 0 0 14px;
  font-size: 18px;
  color: #17315f;
}

.filter-segmented {
  display: inline-flex;
  margin-top: 4px;
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

.meta-inline {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 18px;
  padding: 0 20px 20px;
  color: #7b8dad;
  font-size: 12px;
  line-height: 1.7;
}

.meta-inline span {
  white-space: nowrap;
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

  .meta-inline {
    padding: 0 28px 24px;
  }
}
</style>
