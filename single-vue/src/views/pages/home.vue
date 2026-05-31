<template>
  <div class="dashboard">
    <section class="hero-card">
      <div>
        <span class="eyebrow">Dashboard</span>
        <h1>欢迎回来，{{ store.user.name || '用户' }}</h1>
        <p>这里是资源分享平台的后台首页，主要看分类规模、资源上传和分享状态。</p>
      </div>
      <el-button type="primary" @click="router.push('/categories')">去管理分类</el-button>
    </section>

    <section class="stats-grid">
      <article class="stat-card">
        <span>当前用户</span>
        <strong>{{ store.user.loginName || '--' }}</strong>
      </article>
      <article class="stat-card">
        <span>分类数量</span>
        <strong>{{ store.categories.length }}</strong>
      </article>
      <article class="stat-card">
        <span>账号状态</span>
        <strong>{{ store.user.status || 'active' }}</strong>
      </article>
    </section>

    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="list-header">
          <span>最近分类</span>
          <el-button link @click="router.push('/categories')">查看全部</el-button>
        </div>
      </template>
      <el-empty v-if="store.categories.length === 0" description="还没有分类，先去创建一个" />
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
  await store.loadCategories();
});
</script>

<style scoped lang="scss">
.dashboard {
  display: grid;
  gap: 20px;
}

.hero-card {
  display: flex;
  justify-content: space-between;
  gap: 24px;
  padding: 28px;
  border-radius: 28px;
  background: linear-gradient(135deg, rgba(49, 92, 69, 0.94), rgba(122, 161, 131, 0.9));
  color: #fffef8;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  opacity: 0.72;
}

.hero-card h1 {
  margin: 10px 0 8px;
  font-size: 34px;
}

.hero-card p {
  max-width: 520px;
  color: rgba(255, 254, 248, 0.84);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.stat-card {
  padding: 22px;
  border-radius: 22px;
  background: rgba(255, 252, 246, 0.9);
  border: 1px solid rgba(32, 53, 41, 0.08);
}

.stat-card span {
  display: block;
  color: #6c7f73;
  margin-bottom: 12px;
}

.stat-card strong {
  font-size: 28px;
  color: #203529;
}

.list-card {
  border-radius: 24px;
  border: 1px solid rgba(32, 53, 41, 0.08);
  background: rgba(255, 252, 246, 0.92);
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 700;
  color: #203529;
}

.category-list {
  display: grid;
  gap: 12px;
}

.category-item {
  padding: 14px 16px;
  border-radius: 18px;
  background: #f5f3eb;
}

.category-item strong {
  display: block;
  color: #203529;
}

.category-item span {
  color: #6c7f73;
  font-size: 13px;
}

@media (max-width: 900px) {
  .hero-card,
  .stats-grid {
    grid-template-columns: 1fr;
    display: grid;
  }
}
</style>
