<template>
  <div class="dashboard">
    <section class="stats-grid">
      <article class="stat-card">
        <span class="stat-label">当前用户</span>
        <strong class="stat-value">{{ store.user.loginName || '--' }}</strong>
      </article>
      <article class="stat-card">
        <span class="stat-label">展册数量</span>
        <strong class="stat-value">{{ store.dashboardStats.collectionCount }}</strong>
      </article>
      <article class="stat-card stat-card--media">
        <span class="stat-label">图片资源</span>
        <div class="stat-pair">
          <div class="stat-metric">
            <strong>{{ store.dashboardStats.imageCount }}</strong>
            <small>数量</small>
          </div>
          <div class="stat-metric">
            <strong>{{ formatFileSize(store.dashboardStats.imageSizeTotal) }}</strong>
            <small>占用大小</small>
          </div>
        </div>
      </article>
      <article class="stat-card stat-card--media">
        <span class="stat-label">视频资源</span>
        <div class="stat-pair">
          <div class="stat-metric">
            <strong>{{ store.dashboardStats.videoCount }}</strong>
            <small>数量</small>
          </div>
          <div class="stat-metric">
            <strong>{{ formatFileSize(store.dashboardStats.videoSizeTotal) }}</strong>
            <small>占用大小</small>
          </div>
        </div>
      </article>
    </section>

    <el-card class="list-card" shadow="never">
      <template #header>
        <div class="list-header">
          <span>系统公告</span>
          <el-button v-if="store.isAdmin" link @click="router.push('/announcements')">公告管理</el-button>
        </div>
      </template>
      <el-empty v-if="announcements.length === 0" description="暂无公告" />
      <div v-else class="announcement-list">
        <div v-for="item in announcements" :key="item.id" class="announcement-item">
          <div class="announcement-title">
            <strong>{{ item.title }}</strong>
            <span>{{ formatDate(item.createdAt) }}</span>
          </div>
          <p>{{ item.content || '暂无内容' }}</p>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { getAnnouncementList } from '@/api/user';
import { userMainStore, type Announcement } from '@/store';

const store = userMainStore();
const router = useRouter();
const announcements = ref<Announcement[]>([]);

onMounted(async () => {
  if (!store.user.id) {
    await store.loadProfile();
  }
  await Promise.all([store.loadDashboardStats(), loadAnnouncements()]);
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

async function loadAnnouncements() {
  const res = await getAnnouncementList({ page: 1, pageSize: 6 });
  announcements.value = res.data?.list || [];
}

function formatDate(date?: string) {
  if (!date) return '--';
  return date.slice(0, 10);
}
</script>

<style scoped lang="scss">
.dashboard {
  display: grid;
  gap: 18px;
  max-width: 1180px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
}

.stat-card {
  min-height: 118px;
  padding: 18px 20px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid rgba(117, 151, 210, 0.18);
  box-shadow: 0 10px 24px rgba(48, 86, 150, 0.06);
}

.stat-label {
  display: block;
  color: #6f82a5;
  margin-bottom: 14px;
  font-size: 14px;
  font-weight: 600;
}

.stat-value {
  display: block;
  min-width: 0;
  overflow-wrap: anywhere;
  font-size: 26px;
  line-height: 1.15;
  color: #17315f;
}

.stat-card--media {
  padding-bottom: 16px;
}

.stat-pair {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.stat-metric {
  min-width: 0;
  padding: 10px 12px;
  border-radius: 10px;
  background: #f5f8fd;
  border: 1px solid rgba(128, 159, 214, 0.12);
}

.stat-metric strong {
  display: block;
  color: #17315f;
  font-size: 22px;
  line-height: 1.15;
  overflow-wrap: anywhere;
}

.stat-metric small {
  display: block;
  margin-top: 8px;
  color: #7f91b1;
  font-size: 12px;
}

.list-card {
  width: min(760px, 100%);
  border-radius: 14px;
  border: 1px solid rgba(117, 151, 210, 0.16);
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 10px 24px rgba(48, 86, 150, 0.06);
}

.list-card :deep(.el-card__header) {
  padding: 16px 20px;
}

.list-card :deep(.el-card__body) {
  padding: 16px 20px 20px;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 700;
  color: #17315f;
}

.announcement-list {
  display: grid;
  gap: 10px;
}

.announcement-item {
  padding: 12px 14px;
  border-radius: 10px;
  background: #f5f8fd;
  border: 1px solid rgba(123, 162, 255, 0.12);
}

.announcement-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.announcement-title strong {
  color: #17315f;
}

.announcement-title span {
  color: #8a9abb;
  font-size: 12px;
  white-space: nowrap;
}

.announcement-item p {
  margin: 6px 0 0;
  color: #6d82a7;
  font-size: 13px;
  line-height: 1.7;
}

@media (max-width: 900px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 520px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
