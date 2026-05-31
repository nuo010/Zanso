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
const descriptionText = computed(() => {
  if (!detail.value) return '';
  return detail.value.shareLink.description || detail.value.categoryItem?.description || detail.value.category.description || '暂无描述';
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
  background: radial-gradient(circle at top, #dfeadf 0, #f4efe6 42%, #f7f4ed 100%);
  padding: 18px 12px 40px;
}

.share-card {
  max-width: 760px;
  margin: 0 auto;
  background: rgba(255, 252, 246, 0.94);
  border: 1px solid rgba(31, 45, 36, 0.08);
  border-radius: 28px;
  overflow: hidden;
  box-shadow: 0 18px 50px rgba(31, 45, 36, 0.10);
}

.hero {
  padding: 22px 20px 18px;
  background: linear-gradient(135deg, #315c45, #7fa082);
  color: #fffef8;
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
  color: rgba(255, 254, 248, 0.84);
  font-size: 14px;
  line-height: 1.7;
}

.meta-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  padding: 18px 20px;
  background: #fcfaf4;
}

.meta-item {
  padding: 12px 14px;
  border-radius: 16px;
  background: #f3efe4;
}

.meta-item span {
  display: block;
  font-size: 12px;
  color: #7a867f;
  margin-bottom: 6px;
}

.meta-item strong {
  font-size: 14px;
  color: #203529;
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
  color: #203529;
}

.media-grid {
  display: grid;
  gap: 14px;
}

.media-card {
  overflow: hidden;
  border-radius: 20px;
  background: #f7f3ea;
  border: 1px solid rgba(31, 45, 36, 0.08);
}

.media-card img,
.media-card video {
  display: block;
  width: 100%;
  background: #e8e3d7;
}

.media-card video {
  max-height: 72vh;
}

.media-info {
  padding: 12px 14px;
  font-size: 13px;
  color: #607267;
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
    padding: 20px 28px;
  }
}
</style>
