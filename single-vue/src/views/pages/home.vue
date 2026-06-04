<template>
  <div class="dashboard">
    <section class="hero-card">
      <div class="hero-accent"></div>
      <div class="hero-body">
        <div class="hero-content">
          <span class="eyebrow">Dashboard</span>
          <h1>欢迎回来，{{ store.user.name || '用户' }}</h1>
          <p>汇总展示平台核心运营数据，帮助管理员快速掌握展册规模、资源容量与内容管理概况。</p>
        </div>
        <el-button type="primary" @click="router.push('/categories')">去管理展册</el-button>
      </div>
    </section>

    <section class="stats-grid">
      <article class="stat-card">
        <span>当前用户</span>
        <strong>{{ store.user.loginName || '--' }}</strong>
      </article>
      <article class="stat-card">
        <span>展册数量</span>
        <strong>{{ store.dashboardStats.collectionCount }}</strong>
      </article>
      <article class="stat-card">
        <span>资源数量</span>
        <strong>{{ store.dashboardStats.resourceCount }}</strong>
      </article>
      <article class="stat-card">
        <span>文件占用大小</span>
        <strong>{{ formatFileSize(store.dashboardStats.fileSizeTotal) }}</strong>
      </article>
    </section>

    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="list-header">
          <span>最近展册</span>
          <el-button link @click="router.push('/categories')">查看全部</el-button>
        </div>
      </template>
      <el-empty v-if="store.categories.length === 0" description="暂无展册数据" />
      <div v-else class="category-list">
        <div v-for="item in store.categories.slice(0, 6)" :key="item.id" class="category-item">
          <strong>{{ item.name }}</strong>
          <span>{{ item.description || '暂无描述' }}</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { userMainStore } from '@/store';

const store = userMainStore();
const router = useRouter();

onMounted(async () => {
  if (!store.user.id) {
    await store.loadProfile();
  }
  await Promise.all([store.loadCategories(), store.loadDashboardStats()]);
});

function formatFileSize(size: number) {
  if (!size) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = size;
  let unitIndex = 0;
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }
  const fixed = value >= 100 || unitIndex === 0 ? 0 : value >= 10 ? 1 : 2;
  return `${value.toFixed(fixed)} ${units[unitIndex]}`;
}
</script>

<style scoped lang="scss">
.dashboard {
  display: grid;
  gap: 20px;
}

.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 0;
  padding: 0;
  border-radius: 28px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(123, 162, 255, 0.16);
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
  overflow: hidden;
}

.hero-accent {
  width: 4px;
  flex-shrink: 0;
  background: linear-gradient(180deg, #2f6bff, #69b7ff);
}

.hero-body {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  padding: 28px;
  flex: 1;
}

.eyebrow {
  color: #2f6bff;
  font-size: 12px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  font-weight: 600;
}

.hero-body h1 {
  margin: 10px 0 8px;
  font-size: 34px;
  color: #17315f;
}

.hero-body p {
  max-width: 520px;
  color: #6d82a7;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.stat-card {
  padding: 22px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(123, 162, 255, 0.16);
  box-shadow: 0 14px 32px rgba(60, 102, 190, 0.08);
}

.stat-card span {
  display: block;
  color: #6d82a7;
  margin-bottom: 12px;
}

.stat-card strong {
  font-size: 28px;
  color: #17315f;
}

.list-card {
  border-radius: 24px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 20px 46px rgba(36, 84, 170, 0.08);
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 700;
  color: #17315f;
}

.category-list {
  display: grid;
  gap: 12px;
}

.category-item {
  padding: 14px 16px;
  border-radius: 18px;
  background: #f4f8ff;
  border: 1px solid rgba(123, 162, 255, 0.12);
}

.category-item strong {
  display: block;
  color: #17315f;
}

.category-item span {
  color: #6d82a7;
  font-size: 13px;
}

@media (max-width: 900px) {
  .hero-card,
  .stats-grid {
    grid-template-columns: 1fr;
    display: grid;
  }

  .hero-body {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
