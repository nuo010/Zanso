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
          <article v-for="item in filteredResources" :key="item.id" class="media-card">
            <img
              v-if="item.resourceType !== 'video'"
              :src="item.fileUrl"
              :alt="item.fileName"
              loading="lazy"
              @click="openPreview(item.fileUrl, item.fileName)"
            />
            <video
              v-else
              :src="item.fileUrl"
              :poster="item.posterUrl"
              controls
              playsinline
              webkit-playsinline
              preload="metadata"
            ></video>
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

    <div v-if="previewVisible" class="image-preview" @click.self="closePreview">
      <button class="image-preview__close" type="button" aria-label="关闭预览" @click="closePreview">×</button>
      <img :src="previewImage" :alt="previewTitle" class="image-preview__image" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
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
const originalViewportContent = ref('');
let lastTouchEndAt = 0;

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
  lockShareViewport();
  document.addEventListener('gesturestart', preventGestureZoom, { passive: false });
  document.addEventListener('dblclick', preventGestureZoom, { passive: false });
  document.addEventListener('touchend', preventDoubleTapZoom, { passive: false });
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
    generateVideoPosters(detail.value.resourceList);
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

onUnmounted(() => {
  document.removeEventListener('gesturestart', preventGestureZoom);
  document.removeEventListener('dblclick', preventGestureZoom);
  document.removeEventListener('touchend', preventDoubleTapZoom);
  restoreViewport();
});

function normalizeResourceList(resourceList: any[]) {
  return resourceList.map((item) => ({
    ...item,
    fileUrl: resolveResourceURL(item.url || item.storagePath || ''),
    posterUrl: '',
  }));
}

function generateVideoPosters(resourceList: any[]) {
  resourceList
    .filter((item) => item.resourceType === 'video' && item.fileUrl)
    .forEach((item) => {
      captureVideoPoster(item.fileUrl)
        .then((posterUrl) => {
          item.posterUrl = posterUrl;
        })
        .catch(() => {
          item.posterUrl = '';
        });
    });
}

function captureVideoPoster(url: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const video = document.createElement('video');
    let settled = false;
    let captureTimes: number[] = [];
    let captureIndex = 0;
    let fallbackPoster = '';
    const timeout = window.setTimeout(() => {
      fail();
    }, 12000);
    const cleanup = () => {
      window.clearTimeout(timeout);
      video.pause();
      video.removeAttribute('src');
      video.load();
      video.remove();
    };
    const fail = () => {
      if (settled) return;
      settled = true;
      cleanup();
      reject(new Error('capture video poster failed'));
    };
    const done = (posterUrl: string) => {
      if (settled) return;
      settled = true;
      cleanup();
      resolve(posterUrl);
    };

    video.muted = true;
    video.playsInline = true;
    video.setAttribute('playsinline', 'true');
    video.setAttribute('webkit-playsinline', 'true');
    video.preload = 'metadata';
    video.crossOrigin = 'anonymous';
    video.style.position = 'fixed';
    video.style.left = '-9999px';
    video.style.top = '0';
    video.style.width = '1px';
    video.style.height = '1px';
    video.style.opacity = '0';
    video.style.pointerEvents = 'none';
    video.src = url;
    document.body.appendChild(video);

    const prepareCapture = () => {
      if (!video.videoWidth || !video.videoHeight) {
        return;
      }
      if (captureTimes.length > 0) {
        return;
      }
      captureTimes = buildPosterCaptureTimes(video.duration);
      seekNextPosterFrame(video, captureTimes, captureIndex, fail);
    };

    video.addEventListener(
      'seeked',
      () => {
        const posterFrame = drawVideoPoster(video);
        if (!posterFrame) {
          fail();
          return;
        }
        if (!fallbackPoster) {
          fallbackPoster = posterFrame.posterUrl;
        }
        if (!posterFrame.isDarkFrame) {
          done(posterFrame.posterUrl);
          return;
        }
        captureIndex += 1;
        if (captureIndex < captureTimes.length) {
          seekNextPosterFrame(video, captureTimes, captureIndex, fail);
          return;
        }
        if (fallbackPoster) {
          done(fallbackPoster);
          return;
        }
        fail();
      }
    );
    video.addEventListener('loadedmetadata', prepareCapture);
    video.addEventListener('loadeddata', prepareCapture);
    video.addEventListener('canplay', prepareCapture);
    video.addEventListener('error', fail, { once: true });
    video.load();
  });
}

function buildPosterCaptureTimes(duration: number) {
  const safeDuration = Number.isFinite(duration) && duration > 0 ? duration : 3;
  const rawTimes = [
    0.5,
    1,
    2,
    safeDuration * 0.15,
    safeDuration * 0.25,
    safeDuration * 0.5,
    Math.max(safeDuration - 0.5, 0),
  ];
  const maxTime = Math.max(safeDuration - 0.05, 0);
  return [...new Set(rawTimes.map((time) => Math.min(Math.max(time, 0), maxTime).toFixed(2)))]
    .map((time) => Number(time))
    .filter((time) => Number.isFinite(time) && time >= 0);
}

function seekNextPosterFrame(video: HTMLVideoElement, captureTimes: number[], index: number, reject: () => void) {
  const targetTime = captureTimes[index];
  if (!Number.isFinite(targetTime)) {
    reject();
    return;
  }
  try {
    video.currentTime = targetTime;
  } catch {
    reject();
  }
}

function drawVideoPoster(video: HTMLVideoElement): { posterUrl: string; isDarkFrame: boolean } | null {
  try {
    const canvas = document.createElement('canvas');
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    const context = canvas.getContext('2d');
    if (!context) {
      return null;
    }
    context.drawImage(video, 0, 0, canvas.width, canvas.height);
    const isDarkFrame = isCanvasMostlyDark(context, canvas.width, canvas.height);
    const posterUrl = canvas.toDataURL('image/jpeg', 0.82);
    return { posterUrl, isDarkFrame };
  } catch {
    return null;
  }
}

function isCanvasMostlyDark(context: CanvasRenderingContext2D, width: number, height: number) {
  const imageData = context.getImageData(0, 0, width, height).data;
  const stepX = Math.max(Math.floor(width / 48), 1);
  const stepY = Math.max(Math.floor(height / 48), 1);
  let brightPixels = 0;
  let totalBrightness = 0;
  let pixelCount = 0;

  for (let y = 0; y < height; y += stepY) {
    for (let x = 0; x < width; x += stepX) {
      const index = (y * width + x) * 4;
      const brightness = imageData[index] * 0.299 + imageData[index + 1] * 0.587 + imageData[index + 2] * 0.114;
      totalBrightness += brightness;
      pixelCount += 1;
      if (brightness > 35) {
        brightPixels += 1;
      }
    }
  }

  const averageBrightness = totalBrightness / pixelCount;
  const brightRatio = brightPixels / pixelCount;
  return averageBrightness < 28 || brightRatio < 0.08;
}

function openPreview(url: string, title: string) {
  previewImage.value = url;
  previewTitle.value = title;
  previewVisible.value = true;
}

function closePreview() {
  previewVisible.value = false;
  previewImage.value = '';
  previewTitle.value = '';
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
  const now = Date.now();
  if (now - lastTouchEndAt <= 300) {
    event.preventDefault();
  }
  lastTouchEndAt = now;
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

.media-card img {
  cursor: zoom-in;
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
}

.image-preview__image {
  display: block;
  max-width: min(100%, 960px);
  max-height: min(100vh - 80px, 90vh);
  border-radius: 18px;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.32);
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
</style>
