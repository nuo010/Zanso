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
        <h1>{{ heroTitle }}</h1>
        <p>{{ descriptionText }}</p>
      </section>

      <section v-if="isCollectionShare && detail.categories?.length" class="section section--filters">
        <div class="section-header section-header--filters">
          <h2>分类筛选</h2>
          <span>点击切换要查看的分类资源</span>
        </div>
        <el-segmented v-model="selectedFilter" :options="filterOptions" class="filter-segmented" />
      </section>

      <section class="section">
        <h2>资源展示</h2>
        <div class="media-grid">
          <article
            v-for="(item, index) in filteredResources"
            :key="item.id"
            class="media-card"
            :style="{ '--reveal-delay': `${Math.min(index % pageSize, 10) * 45}ms` }"
          >
            <div v-if="item.resourceType !== 'video' && !loadedImageMap[item.id]" class="media-skeleton"></div>
            <img
              v-if="item.resourceType !== 'video'"
              :class="{ 'is-loaded': loadedImageMap[item.id] }"
              :src="item.fileUrl"
              :alt="item.fileName"
              loading="lazy"
              decoding="async"
              @click="openPreview(item.id)"
              @load="markImageLoaded(item.id)"
              @error="markImageLoaded(item.id)"
            />
            <video
              v-else
              :src="item.fileUrl"
              :poster="item.posterUrl"
              controls
              playsinline
              preload="metadata"
            ></video>
            <div class="media-info">
              <span>{{ item.fileName }}</span>
              <span>{{ item.resourceType }}</span>
            </div>
          </article>
        </div>
        <el-empty v-if="filteredResources.length === 0 && !resourceLoading" description="当前筛选条件下没有资源" />
        <div v-if="resourceLoading" class="load-more-state">
          <span class="load-more-spinner"></span>
          <span>正在加载更多资源</span>
        </div>
        <div v-else-if="filteredResources.length && !hasMore" class="load-more-state load-more-state--done">
          已经到底了
        </div>
      </section>

      <section class="meta-inline">
        <span>发布用户：{{ detail.user.name }}</span>
        <span>所属展册：{{ detail.collection.name }}</span>
        <span>分享对象：{{ selectedCategory?.name || detail.category?.name || detail.collection.name }}</span>
        <span>浏览次数：{{ detail.shareLink.viewCount }}</span>
      </section>
    </div>

    <div v-if="previewVisible" class="image-preview" @click.self="closePreview">
      <button class="image-preview__close" type="button" aria-label="关闭预览" @click="closePreview">×</button>
      <button
        v-if="canBrowsePreview"
        class="image-preview__arrow image-preview__arrow--left"
        type="button"
        aria-label="上一张"
        @click.stop="switchPreview(-1)"
      >
        ‹
      </button>
      <div v-if="!previewImageLoaded" class="image-preview__loader">
        <span class="image-preview__loader-bar"></span>
      </div>
      <img
        :key="previewImage"
        :src="previewImage"
        :alt="previewTitle"
        class="image-preview__image"
        :class="{ 'is-loaded': previewImageLoaded }"
        decoding="async"
        @load="markPreviewImageLoaded"
        @error="markPreviewImageLoaded"
      />
      <button
        v-if="canBrowsePreview"
        class="image-preview__arrow image-preview__arrow--right"
        type="button"
        aria-label="下一张"
        @click.stop="switchPreview(1)"
      >
        ›
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { getShareLinkDetail } from '@/api/user';
import { resolveResourceURL } from '@/util/util';

const route = useRoute();
const loading = ref(false);
const detail = ref<any>(null);
const errorMessage = ref('');
const selectedFilter = ref('all');
const previewVisible = ref(false);
const previewImage = ref('');
const previewTitle = ref('');
const previewImageLoaded = ref(false);
const previewIndex = ref(-1);
const originalViewportContent = ref('');
const page = ref(1);
const pageSize = 12;
const total = ref(0);
const hasMore = ref(false);
const resourceLoading = ref(false);
const loadedImageMap = reactive<Record<string, boolean>>({});
let lastTouchEndAt = 0;
let scrollTicking = false;
let previewTouchStartX = 0;

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
  return detail.value?.resourceList || [];
});
const imageResources = computed(() => {
  return filteredResources.value.filter((item: any) => item.resourceType !== 'video' && item.fileUrl);
});
const canBrowsePreview = computed(() => imageResources.value.length > 1);

watch(selectedFilter, () => {
  if (!detail.value || !isCollectionShare.value) return;
  loadShareDetail(1, true);
});

onMounted(async () => {
  lockShareViewport();
  document.addEventListener('gesturestart', preventGestureZoom, { passive: false });
  document.addEventListener('dblclick', preventGestureZoom, { passive: false });
  document.addEventListener('touchend', preventDoubleTapZoom, { passive: false });
  document.addEventListener('keydown', handlePreviewKeydown);
  document.addEventListener('touchstart', handlePreviewTouchStart, { passive: true });
  document.addEventListener('touchend', handlePreviewTouchEnd, { passive: false });
  window.addEventListener('scroll', handleWindowScroll, { passive: true });
  const code = String(route.params.code || '');
  if (!code) {
    errorMessage.value = '分享链接不存在或分享码无效。';
    return;
  }
  loading.value = true;
  await loadShareDetail(1, true);
  loading.value = false;
});

onUnmounted(() => {
  document.removeEventListener('gesturestart', preventGestureZoom);
  document.removeEventListener('dblclick', preventGestureZoom);
  document.removeEventListener('touchend', preventDoubleTapZoom);
  document.removeEventListener('keydown', handlePreviewKeydown);
  document.removeEventListener('touchstart', handlePreviewTouchStart);
  document.removeEventListener('touchend', handlePreviewTouchEnd);
  window.removeEventListener('scroll', handleWindowScroll);
  restoreViewport();
});

async function loadShareDetail(nextPage = 1, reset = false) {
  if (resourceLoading.value) return;
  const code = String(route.params.code || '');
  if (!code) return;
  resourceLoading.value = true;
  try {
    const res = await getShareLinkDetail(code, {
      page: nextPage,
      pageSize,
      categoryId: isCollectionShare.value && selectedFilter.value !== 'all' ? selectedFilter.value : undefined,
    });
    const nextResourceList = normalizeResourceList(res.data?.resourceList || []);
    detail.value = {
      ...res.data,
      resourceList: reset ? nextResourceList : [...(detail.value?.resourceList || []), ...nextResourceList],
    };
    page.value = res.data?.page || nextPage;
    total.value = res.data?.total || detail.value.resourceList.length;
    hasMore.value = Boolean(res.data?.hasMore);
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
    resourceLoading.value = false;
  }
}

async function loadMoreResources() {
  if (!hasMore.value || resourceLoading.value || loading.value) return;
  await loadShareDetail(page.value + 1, false);
}

function normalizeResourceList(resourceList: any[]) {
  return resourceList.map((item) => ({
    ...item,
    fileUrl: resolveResourceURL(item.url || item.storagePath || ''),
    posterUrl: resolveVideoPosterURL(item),
  }));
}

function openPreview(id: string) {
  const index = imageResources.value.findIndex((item: any) => item.id === id);
  if (index < 0) return;
  setPreviewByIndex(index);
  previewVisible.value = true;
}

function setPreviewByIndex(index: number) {
  const total = imageResources.value.length;
  if (!total) return;
  const nextIndex = (index + total) % total;
  const nextImage = imageResources.value[nextIndex];
  previewIndex.value = nextIndex;
  previewImage.value = nextImage.fileUrl;
  previewTitle.value = nextImage.fileName;
  previewImageLoaded.value = false;
  preloadPreviewNeighbor(nextIndex);
}

function switchPreview(direction: number) {
  if (!canBrowsePreview.value) return;
  setPreviewByIndex(previewIndex.value + direction);
}

function markPreviewImageLoaded() {
  previewImageLoaded.value = true;
}

function markImageLoaded(id: string) {
  loadedImageMap[id] = true;
}

function closePreview() {
  previewVisible.value = false;
  previewImage.value = '';
  previewTitle.value = '';
  previewImageLoaded.value = false;
  previewIndex.value = -1;
}

function preloadPreviewNeighbor(index: number) {
  const total = imageResources.value.length;
  if (total <= 1) return;
  [-1, 1].forEach((offset) => {
    const neighbor = imageResources.value[(index + offset + total) % total];
    if (!neighbor?.fileUrl) return;
    const image = new Image();
    image.decoding = 'async';
    image.src = neighbor.fileUrl;
  });
}

function resolveVideoPosterURL(item: any) {
  if (item.resourceType !== 'video') return '';
  const posterPath = item.posterUrl || item.coverUrl || item.thumbnailUrl || item.previewUrl || '';
  return posterPath ? resolveResourceURL(posterPath) : '';
}

function preventGestureZoom(event: Event) {
  event.preventDefault();
}

function lockShareViewport() {
  const viewport = document.querySelector('meta[name="viewport"]');
  if (viewport) {
    originalViewportContent.value = viewport.getAttribute('content') || '';
    viewport.setAttribute(
      'content',
      'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no, viewport-fit=cover'
    );
  }
}

function restoreViewport() {
  const viewport = document.querySelector('meta[name="viewport"]');
  if (viewport && originalViewportContent.value) {
    viewport.setAttribute('content', originalViewportContent.value);
  }
}

function preventDoubleTapZoom(event: TouchEvent) {
  if (previewVisible.value) return;
  const now = Date.now();
  if (now - lastTouchEndAt <= 300) {
    event.preventDefault();
  }
  lastTouchEndAt = now;
}

function handlePreviewKeydown(event: KeyboardEvent) {
  if (!previewVisible.value) return;
  if (event.key === 'Escape') {
    closePreview();
  } else if (event.key === 'ArrowLeft') {
    event.preventDefault();
    switchPreview(-1);
  } else if (event.key === 'ArrowRight') {
    event.preventDefault();
    switchPreview(1);
  }
}

function handlePreviewTouchStart(event: TouchEvent) {
  if (!previewVisible.value) return;
  previewTouchStartX = event.changedTouches[0]?.clientX || 0;
}

function handlePreviewTouchEnd(event: TouchEvent) {
  if (!previewVisible.value || !previewTouchStartX) return;
  const endX = event.changedTouches[0]?.clientX || 0;
  const distance = endX - previewTouchStartX;
  previewTouchStartX = 0;
  if (Math.abs(distance) < 56) return;
  event.preventDefault();
  switchPreview(distance > 0 ? -1 : 1);
}

function handleWindowScroll() {
  if (scrollTicking) return;
  scrollTicking = true;
  window.requestAnimationFrame(() => {
    scrollTicking = false;
    const distanceToBottom = document.documentElement.scrollHeight - window.scrollY - window.innerHeight;
    if (distanceToBottom < 420) {
      loadMoreResources();
    }
  });
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
  touch-action: manipulation;
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
  padding: 18px 20px 16px;
  background: linear-gradient(135deg, #2f6bff, #78bcff);
  color: #fdfefe;
}

.hero h1 {
  margin: 0 0 6px;
  font-size: 24px;
  line-height: 1.2;
}

.hero p {
  margin: 0;
  color: rgba(253, 254, 255, 0.88);
  font-size: 13px;
  line-height: 1.6;
}

.section {
  padding: 20px;
}

.section--filters {
  padding-top: 16px;
  padding-bottom: 16px;
}

.section-header {
  display: block;
}

.section-header--filters {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 10px;
}

.section-header--filters span {
  color: #7f90af;
  font-size: 12px;
}

.section h2 {
  margin: 0 0 14px;
  font-size: 18px;
  color: #17315f;
}

.filter-segmented {
  display: inline-flex;
  padding: 6px;
  border-radius: 16px;
  background: #f3f7ff;
  border: 1px solid rgba(123, 162, 255, 0.16);
}

.filter-segmented :deep(.el-segmented__item) {
  min-height: 38px;
  padding: 0 16px;
  font-size: 13px;
  font-weight: 600;
}

.media-grid {
  display: grid;
  gap: 14px;
}

.media-card {
  position: relative;
  overflow: hidden;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(116, 153, 230, 0.16);
  box-shadow: 0 14px 32px rgba(60, 102, 190, 0.08);
  animation: card-reveal 0.5s ease both;
  animation-delay: var(--reveal-delay, 0ms);
}

.media-card img,
.media-card video {
  display: block;
  width: 100%;
  min-height: 220px;
  background: linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
}

.media-card img {
  cursor: zoom-in;
  opacity: 0;
  filter: blur(10px);
  transform: scale(1.012);
  transition:
    opacity 0.34s ease,
    filter 0.34s ease,
    transform 0.34s ease;
}

.media-card img.is-loaded {
  opacity: 1;
  filter: blur(0);
  transform: scale(1);
}

.media-card video {
  max-height: 72vh;
}

.media-skeleton {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 220px;
  background:
    linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.68) 46%, transparent 72%),
    linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
  background-size: 220px 100%, 100% 100%;
  animation: image-loading 1.1s ease-in-out infinite;
  pointer-events: none;
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

.load-more-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 52px;
  color: #7b8dad;
  font-size: 13px;
}

.load-more-state--done {
  color: #9aa8bf;
}

.load-more-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(47, 107, 255, 0.18);
  border-top-color: #2f6bff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.image-preview {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(7, 17, 40, 0.88);
}

.image-preview__close {
  position: absolute;
  top: 18px;
  right: 18px;
  width: 42px;
  height: 42px;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  color: #fff;
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
  z-index: 3;
}

.image-preview__arrow {
  position: absolute;
  top: 50%;
  z-index: 3;
  width: 46px;
  height: 64px;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  color: #fff;
  font-size: 48px;
  line-height: 1;
  cursor: pointer;
  transform: translateY(-50%);
  transition:
    background 0.2s ease,
    transform 0.2s ease;
}

.image-preview__arrow:hover {
  background: rgba(255, 255, 255, 0.24);
}

.image-preview__arrow--left {
  left: 16px;
}

.image-preview__arrow--right {
  right: 16px;
}

.image-preview__loader {
  position: absolute;
  inset: 20px;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 18px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.09) 0%, rgba(255, 255, 255, 0.04) 100%),
    rgba(255, 255, 255, 0.05);
}

.image-preview__loader-bar {
  width: min(280px, 62vw);
  height: 6px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
}

.image-preview__loader-bar::after {
  content: '';
  display: block;
  width: 42%;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, rgba(120, 188, 255, 0.1), rgba(120, 188, 255, 0.95));
  animation: preview-loading 1s ease-in-out infinite;
}

.image-preview__image {
  position: relative;
  z-index: 2;
  display: block;
  max-width: min(100%, 960px);
  max-height: min(100vh - 80px, 90vh);
  border-radius: 18px;
  opacity: 0;
  filter: blur(12px);
  transform: scale(0.99);
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.32);
  transition:
    opacity 0.35s ease,
    filter 0.35s ease,
    transform 0.35s ease;
}

.image-preview__image.is-loaded {
  opacity: 1;
  filter: blur(0);
  transform: scale(1);
}

@keyframes card-reveal {
  from {
    opacity: 0;
    transform: translateY(12px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes image-loading {
  0% {
    background-position: -220px 0, 0 0;
  }

  100% {
    background-position: calc(100% + 220px) 0, 0 0;
  }
}

@keyframes preview-loading {
  0% {
    transform: translateX(-110%);
  }

  100% {
    transform: translateX(260%);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (min-width: 768px) {
  .share-page {
    padding: 30px 24px 56px;
  }

  .hero {
    padding: 24px 28px 18px;
  }

  .hero h1 {
    font-size: 30px;
  }

  .section {
    padding: 24px 28px;
  }

  .meta-inline {
    padding: 0 28px 24px;
  }
}

@media (max-width: 640px) {
  .image-preview {
    padding: 14px;
  }

  .image-preview__arrow {
    width: 38px;
    height: 54px;
    font-size: 38px;
    background: rgba(255, 255, 255, 0.18);
  }

  .image-preview__arrow--left {
    left: 8px;
  }

  .image-preview__arrow--right {
    right: 8px;
  }

  .image-preview__image {
    max-height: min(100vh - 72px, 88vh);
  }
}
</style>
