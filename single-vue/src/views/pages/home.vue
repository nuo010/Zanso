<template>
  <div class="dashboard">
    <section class="hero-card">
      <div>
        <span class="eyebrow">Dashboard</span>
        <h1>欢迎回来，{{ store.user.name || '用户' }}</h1>
        <p>这里是资源分享平台的后台首页，主要看展册规模、资源上传和分享状态。</p>
      </div>
      <el-button type="primary" @click="router.push('/categories')">去管理展册</el-button>
    </section>

    <section class="stats-grid">
      <article class="stat-card">
        <span>当前用户</span>
        <strong>{{ store.user.loginName || '--' }}</strong>
      </article>
      <article class="stat-card">
        <span>展册数量</span>
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
          <span>最近展册</span>
          <el-button link @click="router.push('/categories')">查看全部</el-button>
        </div>
      </template>
      <el-empty v-if="store.categories.length === 0" description="还没有展册，先去创建一个" />
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
  background: linear-gradient(135deg, rgba(47, 107, 255, 0.96), rgba(120, 188, 255, 0.92));
  color: #fdfefe;
  box-shadow: 0 24px 60px rgba(47, 107, 255, 0.18);
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
  color: rgba(253, 254, 255, 0.88);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
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
}
</style>
