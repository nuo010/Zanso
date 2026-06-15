<template>
  <div class="category-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>展册管理</h2>
            <p>集中管理展册、分类与资源内容，支持可见性配置、资源维护和分享发布，满足企业级内容展示与分发需求。</p>
          </div>
          <el-button type="primary" @click="openCreateCollectionDialog">新建展册</el-button>
        </div>
      </template>

      <el-table
        :data="store.categories"
        v-loading="loading"
        row-key="id"
        class="outer-table"
        @expand-change="handleCollectionExpand"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-panel" v-loading="expandedLoadingMap[row.id]">
              <template v-if="expandedDetailMap[row.id]">
                <el-table
                  v-if="expandedDetailMap[row.id].categories?.length"
                  :data="expandedDetailMap[row.id].categories"
                  size="small"
                  row-key="id"
                  class="inner-table"
                  :show-header="false"
                  @expand-change="handleInnerCategoryExpand"
                >
                  <el-table-column type="expand">
                    <template #default="{ row: categoryRow }">
                      <div class="resource-panel" v-loading="itemLoadingMap[categoryRow.id]">
                        <template v-if="expandedItemDetailMap[categoryRow.id]">
                          <template v-if="expandedItemDetailMap[categoryRow.id].resourceList?.length">
                            <TransitionGroup tag="div" class="resource-grid" name="resource-list">
                            <div
                              v-for="resource in expandedItemDetailMap[categoryRow.id].resourceList"
                              :key="resource.id"
                              class="resource-card"
                              :class="{
                                'resource-card--dragging': draggingResourceId === resource.id,
                                'resource-card--drag-over': dragOverResourceId === resource.id,
                              }"
                              draggable="true"
                              @dragstart="handleResourceDragStart(categoryRow.id, resource.id)"
                              @dragenter.prevent="handleResourceDragEnter(resource.id)"
                              @dragover.prevent
                              @drop.prevent="handleResourceDrop(categoryRow.id, resource.id)"
                              @dragend="handleResourceDragEnd"
                            >
                              <div class="resource-preview">
                                <span class="resource-sort-badge" v-if="resource.sort">{{ resource.sort }}</span>
                                <div class="resource-actions">
                                  <el-tooltip content="删除资源" placement="top" :show-after="400">
                                    <el-button
                                      :icon="Delete"
                                      size="small"
                                      text
                                      type="danger"
                                      @click.stop="confirmDeleteResource(resource.resourceId, categoryRow.categoryId, categoryRow.id)"
                                    />
                                  </el-tooltip>
                                </div>
                                <el-image
                                  v-if="resource.resourceType !== 'video'"
                                  :src="resource.fileUrl"
                                  :alt="resource.fileName"
                                  :preview-src-list="[resource.fileUrl]"
                                  :initial-index="0"
                                  fit="cover"
                                  preview-teleported
                                  hide-on-click-modal
                                  loading="lazy"
                                />
                                <button
                                  v-else-if="playingResourceId !== resource.id"
                                  class="resource-video-cover"
                                  type="button"
                                  @click.stop="playResourceVideo(resource.id)"
                                >
                                  <img
                                    v-if="resource.posterUrl"
                                    :src="resource.posterUrl"
                                    :alt="resource.fileName"
                                    loading="lazy"
                                  />
                                  <span v-else class="resource-video-placeholder">{{ resource.fileName }}</span>
                                  <span class="resource-video-play" aria-hidden="true"></span>
                                </button>
                                <video
                                  v-else
                                  :src="resource.fileUrl"
                                  :poster="resource.posterUrl"
                                  controls
                                  autoplay
                                  playsinline
                                  preload="metadata"
                                  @click.stop
                                  @ended="stopResourceVideo"
                                ></video>
                              </div>
                              <div class="resource-info">
                                <strong :title="resource.fileName">{{ resource.fileName }}</strong>
                                <span :title="`${resource.resourceType} · ${resource.mimeType || '未知类型'}`">
                                  {{ resource.resourceType }} · {{ resource.mimeType || '未知类型' }}
                                </span>
                                <div class="resource-meta-row">
                                  <span class="resource-meta">{{ formatFileSize(resource.fileSize) }}</span>
                                </div>
                              </div>
                            </div>
                            </TransitionGroup>
                          </template>
                          <el-empty v-else description="当前分类下暂无关联资源" />
                        </template>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="name" min-width="220" />
                  <el-table-column label="描述" min-width="240" show-overflow-tooltip>
                    <template #default="{ row: categoryRow }">
                      <el-tooltip
                        v-if="shouldShowDescriptionTooltip(categoryRow.description)"
                        :content="formatDescription(categoryRow.description)"
                        placement="top"
                        effect="light"
                        popper-class="description-tooltip"
                      >
                        <span class="description-cell">{{ formatDescription(categoryRow.description) }}</span>
                      </el-tooltip>
                      <span v-else class="description-cell">{{ formatDescription(categoryRow.description) }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column label="展示" width="120" align="center">
                    <template #default="{ row: categoryRow }">
                      <el-switch
                        :model-value="categoryRow.visible"
                        inline-prompt
                        active-text="可看"
                        inactive-text="隐藏"
                        @change="handleToggleInnerCategoryVisible(categoryRow, $event)"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column width="360" align="right">
                    <template #default="{ row: categoryRow }">
                      <div class="row-actions">
                        <el-button link type="primary" @click="openUploadDialog(categoryRow.categoryId, categoryRow.id)">上传资源</el-button>
                        <el-button link type="primary" @click="openEditInnerCategoryDialog(categoryRow)">修改</el-button>
                        <el-button
                          link
                          type="warning"
                          @click="openShareDialog('category', categoryRow.categoryId, categoryRow.id, categoryRow.name, categoryRow.description)"
                        >
                          分享分类
                        </el-button>
                        <el-button link type="danger" @click="confirmDeleteInnerCategory(categoryRow.id, categoryRow.categoryId)">
                          删除
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
                <div v-if="expandedDetailMap[row.id].total > pageSize" class="inner-pagination">
                  <el-pagination
                    small
                    background
                    layout="total, prev, pager, next"
                    :page-size="pageSize"
                    :total="expandedDetailMap[row.id].total"
                    :current-page="expandedDetailMap[row.id].page || 1"
                    @current-change="handleInnerCategoryPageChange(row.id, $event)"
                  />
                </div>
                <el-empty
                  v-if="!expandedDetailMap[row.id].categories?.length"
                  description="当前展册下暂无分类"
                />
              </template>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="展册名称" min-width="220" />
        <el-table-column label="描述" min-width="260" show-overflow-tooltip>
          <template #default="{ row }">
            <el-tooltip
              v-if="shouldShowDescriptionTooltip(row.description)"
              :content="formatDescription(row.description)"
              placement="top"
              effect="light"
              popper-class="description-tooltip"
            >
              <span class="description-cell">{{ formatDescription(row.description) }}</span>
            </el-tooltip>
            <span v-else class="description-cell">{{ formatDescription(row.description) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="展示" width="140" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.visible"
              inline-prompt
              active-text="可看"
              inactive-text="隐藏"
              @change="handleToggleCollectionVisible(row, $event)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="340" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button link type="primary" @click="openCreateInnerCategoryDialog(row)">新建分类</el-button>
              <el-button link type="primary" @click="openEditCollectionDialog(row)">修改描述</el-button>
              <el-button link type="warning" @click="openShareDialog('collection', row.id, '', row.name, row.description)">分享展册</el-button>
              <el-button link type="danger" @click="confirmDeleteCollection(row.id)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="store.categoryTotal > pageSize" class="table-pagination">
        <el-pagination
          background
          layout="total, prev, pager, next"
          :page-size="pageSize"
          :total="store.categoryTotal"
          :current-page="store.categoryPage"
          @current-change="handleCollectionPageChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="createCategoryDialogVisible" title="新建展册" width="460px">
      <el-form :model="categoryForm" label-position="top">
        <el-form-item label="展册名称">
          <el-input v-model="categoryForm.name" placeholder="直接输入展册名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="categoryForm.description" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="是否可看">
          <el-switch v-model="categoryForm.visible" active-text="可看" inactive-text="不可看" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingCategory" @click="handleCreateCollection">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editCategoryDialogVisible" title="修改展册描述" width="460px">
      <el-form :model="editCategoryForm" label-position="top">
        <el-form-item label="展册名称">
          <el-input :model-value="editCategoryForm.name" disabled />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editCategoryForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingCategory" @click="handleUpdateCollectionDescription">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="categoryItemDialogVisible"
      :title="categoryItemDialogMode === 'edit' ? '修改分类' : '新建分类'"
      width="460px"
    >
      <el-form :model="categoryItemForm" label-position="top">
        <el-form-item label="所属展册">
          <el-input :model-value="categoryItemParentName" disabled />
        </el-form-item>
        <el-form-item label="分类名称">
          <el-input v-model="categoryItemForm.name" placeholder="直接输入分类名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="categoryItemForm.description" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="是否可看">
          <el-switch v-model="categoryItemForm.visible" active-text="可看" inactive-text="不可看" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="categoryItemDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingCategoryItem" @click="handleSubmitInnerCategory">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="shareDialogVisible"
      :title="shareForm.targetType === 'category' ? '分享分类' : '分享展册'"
      width="460px"
    >
      <el-form :model="shareForm" label-position="top">
        <el-form-item label="分享标题">
          <el-input v-model="shareForm.title" />
        </el-form-item>
        <el-form-item label="分享描述">
          <el-input v-model="shareForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="到期时间">
          <el-date-picker
            v-model="shareForm.expiresAt"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="不选则长期有效"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <div v-if="generatedShare.shareUrl" class="share-result">
        <div class="share-result__qr">
          <qrcode-vue :value="generatedShare.shareUrl" :size="168" level="M" render-as="svg" />
        </div>
        <div class="share-result__content">
          <strong>{{ generatedShare.title || shareForm.title || '分享链接' }}</strong>
          <a :href="generatedShare.shareUrl" target="_blank" rel="noopener noreferrer">{{ generatedShare.shareUrl }}</a>
          <el-button type="primary" plain @click="copyShareUrl">一键复制分享链接</el-button>
        </div>
      </div>
      <template #footer>
        <el-button @click="shareDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingShare" @click="handleSubmitShare">生成分享链接</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="uploadDialogVisible" title="上传分类资源" width="560px" @closed="resetUploadDialog">
      <div class="upload-wall">
        <el-upload
          v-model:file-list="uploadFileList"
          list-type="picture-card"
          :auto-upload="false"
          multiple
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
        >
          <el-icon><Plus /></el-icon>
          <template #file="{ file }">
            <div class="upload-preview-card">
              <img
                v-if="getUploadPreview(file)?.type === 'image'"
                class="upload-preview-card__media"
                :src="getUploadPreview(file)?.url"
                :alt="file.name"
              />
              <template v-else-if="getUploadPreview(file)?.type === 'video'">
                <img
                  v-if="getUploadPreview(file)?.posterUrl"
                  class="upload-preview-card__media"
                  :src="getUploadPreview(file)?.posterUrl"
                  :alt="file.name"
                />
                <div v-else class="upload-preview-card__placeholder">
                  <span class="upload-preview-card__play"></span>
                </div>
              </template>
              <div v-else class="upload-preview-card__placeholder">
                <span>{{ getFileExt(file.name) }}</span>
              </div>
              <span class="upload-preview-card__name" :title="file.name">{{ file.name }}</span>
              <button
                class="upload-preview-card__remove"
                type="button"
                aria-label="移除资源"
                @click.stop="removeSelectedUploadFile(file)"
              >
                ×
              </button>
            </div>
          </template>
          <template #tip>
            <div class="upload-wall__tip">点击卡片选择资源，上传后会展示在当前分类的资源墙里。</div>
          </template>
        </el-upload>
      </div>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="submitUpload">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { UploadFile, UploadFiles, UploadUserFile } from 'element-plus';
import { ElMessageBox } from 'element-plus';
import { Delete, Plus } from '@element-plus/icons-vue';
import QrcodeVue from 'qrcode.vue';
import {
  createCategory,
  createCategoryItem,
  createShareLink,
  deleteCategory,
  deleteCategoryItem,
  deleteResource,
  getCategoryDetail,
  getCategoryItemDetail,
  updateCategoryResourceSort,
  updateCategory,
  updateCategoryItem,
  uploadCategoryResource,
} from '@/api/user';
import { userMainStore } from '@/store';
import { resolveResourceURL, toast } from '@/util/util';

const store = userMainStore();
const loading = ref(false);
const submittingCategory = ref(false);
const submittingCategoryItem = ref(false);
const submittingShare = ref(false);
const uploading = ref(false);
const createCategoryDialogVisible = ref(false);
const editCategoryDialogVisible = ref(false);
const categoryItemDialogVisible = ref(false);
const shareDialogVisible = ref(false);
const uploadDialogVisible = ref(false);
const categoryItemDialogMode = ref<'create' | 'edit'>('create');
const pageSize = 20;

const expandedDetailMap = reactive<Record<string, any>>({});
const expandedLoadingMap = reactive<Record<string, boolean>>({});
const expandedItemDetailMap = reactive<Record<string, any>>({});
const itemLoadingMap = reactive<Record<string, boolean>>({});
const playingResourceId = ref('');

const categoryForm = reactive({
  name: '',
  description: '',
  visible: true,
});

const editCategoryForm = reactive({
  id: '',
  name: '',
  description: '',
  visible: true,
  status: 'active',
});

const categoryItemForm = reactive({
  id: '',
  collectionId: '',
  collectionName: '',
  name: '',
  description: '',
  visible: true,
  status: 'active',
});

const shareForm = reactive({
  targetType: 'collection' as 'collection' | 'category',
  collectionId: '',
  categoryId: '',
  title: '',
  description: '',
  expiresAt: '',
});

const generatedShare = reactive({
  title: '',
  shareUrl: '',
});

const uploadTargetCollectionId = ref('');
const uploadTargetCategoryId = ref('');
const uploadFileList = ref<UploadUserFile[]>([]);
const uploadPreviewMap = reactive<Record<string, { type: 'image' | 'video' | 'file'; url: string; posterUrl: string }>>({});
const categoryItemParentName = ref('');
const draggingCategoryId = ref('');
const draggingResourceId = ref('');
const dragOverResourceId = ref('');

async function refreshList() {
  loading.value = true;
  try {
    await store.loadCategories(store.categoryPage || 1, pageSize);
  } finally {
    loading.value = false;
  }
}

refreshList();

function openCreateCollectionDialog() {
  categoryForm.name = '';
  categoryForm.description = '';
  categoryForm.visible = true;
  createCategoryDialogVisible.value = true;
}

function openEditCollectionDialog(collection: any) {
  editCategoryForm.id = collection.id;
  editCategoryForm.name = collection.name;
  editCategoryForm.description = collection.description || '';
  editCategoryForm.visible = collection.visible !== false;
  editCategoryForm.status = collection.status || 'active';
  editCategoryDialogVisible.value = true;
}

function openCreateInnerCategoryDialog(collection: any) {
  categoryItemDialogMode.value = 'create';
  categoryItemForm.id = '';
  categoryItemForm.collectionId = collection.id;
  categoryItemForm.collectionName = collection.name;
  categoryItemForm.name = '';
  categoryItemForm.description = '';
  categoryItemForm.visible = true;
  categoryItemForm.status = 'active';
  categoryItemParentName.value = collection.name;
  categoryItemDialogVisible.value = true;
}

function openEditInnerCategoryDialog(category: any) {
  categoryItemDialogMode.value = 'edit';
  categoryItemForm.id = category.id;
  categoryItemForm.collectionId = category.categoryId;
  categoryItemForm.collectionName = findCollectionName(category.categoryId);
  categoryItemForm.name = category.name;
  categoryItemForm.description = category.description || '';
  categoryItemForm.visible = category.visible !== false;
  categoryItemForm.status = category.status || 'active';
  categoryItemParentName.value = categoryItemForm.collectionName;
  categoryItemDialogVisible.value = true;
}

function openShareDialog(
  targetType: 'collection' | 'category',
  collectionId: string,
  categoryId: string,
  title: string,
  description?: string
) {
  shareForm.targetType = targetType;
  shareForm.collectionId = collectionId;
  shareForm.categoryId = categoryId;
  shareForm.title = title || '';
  shareForm.description = description || '';
  shareForm.expiresAt = '';
  generatedShare.title = '';
  generatedShare.shareUrl = '';
  shareDialogVisible.value = true;
}

function findCollectionName(collectionId: string) {
  return store.categories.find((item: any) => item.id === collectionId)?.name || '';
}

async function handleCreateCollection() {
  const name = categoryForm.name.trim();
  if (!name) {
    toast('展册名称不能为空', 'warning');
    return;
  }
  if (store.categories.some((item: any) => item.name === name)) {
    toast('展册名称不能重复', 'warning');
    return;
  }

  submittingCategory.value = true;
  try {
    await createCategory({
      name,
      description: categoryForm.description.trim(),
      visible: categoryForm.visible,
      status: 'active',
    });
    toast('展册创建成功', 'success');
    createCategoryDialogVisible.value = false;
    store.categoryPage = 1;
    await refreshList();
  } finally {
    submittingCategory.value = false;
  }
}

async function handleSubmitInnerCategory() {
  const name = categoryItemForm.name.trim();
  if (!categoryItemForm.collectionId || !name) {
    toast('分类名称不能为空', 'warning');
    return;
  }

  const currentItems = expandedDetailMap[categoryItemForm.collectionId]?.categories || [];
  const duplicate = currentItems.some((item: any) => item.name === name && item.id !== categoryItemForm.id);
  if (duplicate) {
    toast('分类名称不能重复', 'warning');
    return;
  }

  submittingCategoryItem.value = true;
  try {
    if (categoryItemDialogMode.value === 'edit' && categoryItemForm.id) {
      await updateCategoryItem(categoryItemForm.id, {
        name,
        description: categoryItemForm.description.trim(),
        visible: categoryItemForm.visible,
        status: categoryItemForm.status,
      });
      toast('分类修改成功', 'success');
    } else {
      await createCategoryItem({
        collectionId: categoryItemForm.collectionId,
        name,
        description: categoryItemForm.description.trim(),
        visible: categoryItemForm.visible,
        status: 'active',
      });
      toast('分类创建成功', 'success');
    }
    categoryItemDialogVisible.value = false;
    await loadCollectionDetail(categoryItemForm.collectionId, expandedDetailMap[categoryItemForm.collectionId]?.page || 1);
  } finally {
    submittingCategoryItem.value = false;
  }
}

async function handleUpdateCollectionDescription() {
  if (!editCategoryForm.id) return;
  submittingCategory.value = true;
  try {
    await updateCategory(editCategoryForm.id, {
      name: editCategoryForm.name,
      description: editCategoryForm.description.trim(),
      visible: editCategoryForm.visible,
      status: editCategoryForm.status,
    });
    toast('展册描述修改成功', 'success');
    editCategoryDialogVisible.value = false;
    const target = store.categories.find((item: any) => item.id === editCategoryForm.id);
    if (target) {
      target.description = editCategoryForm.description.trim();
    }
    await loadCollectionDetail(editCategoryForm.id, expandedDetailMap[editCategoryForm.id]?.page || 1);
  } finally {
    submittingCategory.value = false;
  }
}

async function handleSubmitShare() {
  submittingShare.value = true;
  try {
    const res = await createShareLink({
      collectionId: shareForm.collectionId,
      categoryId: shareForm.categoryId || undefined,
      targetType: shareForm.targetType,
      title: shareForm.title.trim(),
      description: shareForm.description.trim(),
      expiresAt: shareForm.expiresAt || undefined,
    });
    generatedShare.title = res.data?.shareLink?.title || shareForm.title.trim();
    generatedShare.shareUrl = res.data?.shareUrl || '';
    toast('分享链接已生成', 'success');
  } finally {
    submittingShare.value = false;
  }
}

async function copyShareUrl() {
  const shareUrl = generatedShare.shareUrl;
  if (!shareUrl) return;
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(shareUrl);
    } else {
      copyTextWithTextarea(shareUrl);
    }
    toast('分享链接已复制', 'success');
  } catch (error) {
    copyTextWithTextarea(shareUrl);
    toast('分享链接已复制', 'success');
  }
}

function copyTextWithTextarea(text: string) {
  const textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.setAttribute('readonly', 'readonly');
  textarea.style.position = 'fixed';
  textarea.style.left = '-9999px';
  document.body.appendChild(textarea);
  textarea.select();
  document.execCommand('copy');
  document.body.removeChild(textarea);
}

async function loadCollectionDetail(id: string, page = 1) {
  expandedLoadingMap[id] = true;
  try {
    const res = await getCategoryDetail(id, { page, pageSize });
    expandedDetailMap[id] = res.data;
  } finally {
    expandedLoadingMap[id] = false;
  }
}

async function loadInnerCategoryDetail(id: string) {
  itemLoadingMap[id] = true;
  try {
    const res = await getCategoryItemDetail(id);
    expandedItemDetailMap[id] = {
      ...res.data,
      resourceList: normalizeResourceList(res.data?.resourceList || []),
    };
  } finally {
    itemLoadingMap[id] = false;
  }
}

function normalizeResourceList(resourceList: any[]) {
  return resourceList.map((item) => ({
    ...item,
    fileUrl: resolveResourceURL(item.url || item.storagePath || ''),
    posterUrl: resolveVideoPosterURL(item),
  }));
}

function resolveVideoPosterURL(item: any) {
  if (item.resourceType !== 'video') return '';
  const posterPath = item.posterUrl || item.coverUrl || item.thumbnailUrl || item.previewUrl || '';
  return posterPath ? resolveResourceURL(posterPath) : '';
}

function formatFileSize(size?: number) {
  const fileSize = Number(size || 0);
  if (!fileSize) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let value = fileSize;
  let unitIndex = 0;
  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024;
    unitIndex += 1;
  }
  const fixed = value >= 100 || unitIndex === 0 ? 0 : value >= 10 ? 1 : 2;
  return `${value.toFixed(fixed)} ${units[unitIndex]}`;
}

function formatDescription(description?: string) {
  return String(description || '').trim() || '暂无描述';
}

function shouldShowDescriptionTooltip(description?: string) {
  return formatDescription(description).length > 24;
}

function playResourceVideo(resourceId: string) {
  playingResourceId.value = resourceId;
}

function stopResourceVideo() {
  playingResourceId.value = '';
}

function handleResourceDragStart(categoryId: string, resourceId: string) {
  draggingCategoryId.value = categoryId;
  draggingResourceId.value = resourceId;
  dragOverResourceId.value = resourceId;
}

function handleResourceDragEnter(resourceId: string) {
  if (!draggingResourceId.value || draggingResourceId.value === resourceId) return;
  dragOverResourceId.value = resourceId;
}

async function handleResourceDrop(categoryId: string, targetResourceId: string) {
  const sourceResourceId = draggingResourceId.value;
  if (!sourceResourceId || !categoryId || sourceResourceId === targetResourceId) {
    handleResourceDragEnd();
    return;
  }
  if (draggingCategoryId.value && draggingCategoryId.value !== categoryId) {
    handleResourceDragEnd();
    return;
  }

  const targetDetail = expandedItemDetailMap[categoryId];
  const resourceList = targetDetail?.resourceList || [];
  const sourceIndex = resourceList.findIndex((item: any) => item.id === sourceResourceId);
  const targetIndex = resourceList.findIndex((item: any) => item.id === targetResourceId);
  if (sourceIndex < 0 || targetIndex < 0 || sourceIndex === targetIndex) {
    handleResourceDragEnd();
    return;
  }

  const nextList = [...resourceList];
  const [movedItem] = nextList.splice(sourceIndex, 1);
  nextList.splice(targetIndex, 0, movedItem);
  const previousList = resourceList.slice();

  // Sync sort values with new positions after reorder
  const updatedList = nextList.map((item, index) => ({
    ...item,
    sort: index + 1,
  }));

  expandedItemDetailMap[categoryId] = {
    ...targetDetail,
    resourceList: updatedList,
  };

  try {
    await updateCategoryResourceSort(categoryId, {
      resourceRelationIds: updatedList.map((item: any) => item.id),
    });
    toast('资源顺序已更新', 'success');
  } catch (error) {
    expandedItemDetailMap[categoryId] = {
      ...targetDetail,
      resourceList: previousList,
    };
    toast('资源顺序保存失败', 'error');
  } finally {
    handleResourceDragEnd();
  }
}

function handleResourceDragEnd() {
  draggingCategoryId.value = '';
  draggingResourceId.value = '';
  dragOverResourceId.value = '';
}

function handleCollectionExpand(row: any, expandedRows: any[]) {
  if (expandedRows.some((item) => item.id === row.id)) {
    loadCollectionDetail(row.id, expandedDetailMap[row.id]?.page || 1);
  }
}

function handleInnerCategoryExpand(row: any, expandedRows: any[]) {
  if (expandedRows.some((item) => item.id === row.id)) {
    loadInnerCategoryDetail(row.id);
  }
}

async function handleToggleCollectionVisible(row: any, value: boolean) {
  await updateCategory(row.id, {
    name: row.name,
    description: row.description,
    visible: value,
    status: row.status,
  });
  row.visible = value;
  const target = store.categories.find((item: any) => item.id === row.id);
  if (target) {
    target.visible = value;
  }
  toast(`展册已切换为${value ? '可看' : '不可看'}`, 'success');
}

async function handleToggleInnerCategoryVisible(row: any, value: boolean) {
  await updateCategoryItem(row.id, {
    name: row.name,
    description: row.description,
    visible: value,
    status: row.status,
  });
  row.visible = value;
  toast(`分类已切换为${value ? '可看' : '不可看'}`, 'success');
}

function openUploadDialog(collectionId: string, categoryId: string) {
  uploadTargetCollectionId.value = collectionId;
  uploadTargetCategoryId.value = categoryId;
  resetUploadDialog();
  uploadDialogVisible.value = true;
}

function handleFileChange(file: UploadFile, fileList: UploadFiles) {
  uploadFileList.value = fileList;
  prepareUploadPreview(file);
}

function handleFileRemove(file: UploadFile, fileList: UploadFiles) {
  uploadFileList.value = fileList;
  removeUploadPreview(file.uid);
}

function resetUploadDialog() {
  Object.keys(uploadPreviewMap).forEach((uid) => removeUploadPreview(Number(uid)));
  uploadFileList.value = [];
}

function getUploadPreview(file: UploadUserFile) {
  return uploadPreviewMap[String(file.uid)];
}

function getFileExt(fileName: string) {
  const ext = fileName.split('.').pop() || 'FILE';
  return ext.slice(0, 5).toUpperCase();
}

function removeSelectedUploadFile(file: UploadUserFile) {
  uploadFileList.value = uploadFileList.value.filter((item) => item.uid !== file.uid);
  removeUploadPreview(file.uid);
}

function prepareUploadPreview(file: UploadFile) {
  if (!file.raw) return;
  removeUploadPreview(file.uid);
  const url = URL.createObjectURL(file.raw);
  const fileType = file.raw.type || '';
  if (fileType.startsWith('image/')) {
    uploadPreviewMap[String(file.uid)] = { type: 'image', url, posterUrl: '' };
    return;
  }
  if (fileType.startsWith('video/')) {
    uploadPreviewMap[String(file.uid)] = { type: 'video', url, posterUrl: '' };
    captureUploadVideoPoster(file.uid, url);
    return;
  }
  uploadPreviewMap[String(file.uid)] = { type: 'file', url, posterUrl: '' };
}

function removeUploadPreview(uid?: number) {
  if (uid === undefined) return;
  const key = String(uid);
  const preview = uploadPreviewMap[key];
  if (!preview) return;
  if (preview.url) URL.revokeObjectURL(preview.url);
  if (preview.posterUrl) URL.revokeObjectURL(preview.posterUrl);
  delete uploadPreviewMap[key];
}

function captureUploadVideoPoster(uid: number, url: string) {
  const video = document.createElement('video');
  video.preload = 'metadata';
  video.muted = true;
  video.playsInline = true;
  video.src = url;

  const cleanup = () => {
    video.removeAttribute('src');
    video.load();
  };

  video.onloadedmetadata = () => {
    const seekTime = Math.min(0.1, Math.max(video.duration || 0, 0) / 2);
    video.currentTime = Number.isFinite(seekTime) ? seekTime : 0;
  };

  video.onseeked = () => {
    const preview = uploadPreviewMap[String(uid)];
    if (!preview || preview.url !== url) {
      cleanup();
      return;
    }
    const canvas = document.createElement('canvas');
    canvas.width = video.videoWidth || 480;
    canvas.height = video.videoHeight || 270;
    const context = canvas.getContext('2d');
    if (!context) {
      cleanup();
      return;
    }
    context.drawImage(video, 0, 0, canvas.width, canvas.height);
    canvas.toBlob((blob) => {
      const currentPreview = uploadPreviewMap[String(uid)];
      if (!blob || !currentPreview || currentPreview.url !== url) {
        cleanup();
        return;
      }
      if (currentPreview.posterUrl) URL.revokeObjectURL(currentPreview.posterUrl);
      currentPreview.posterUrl = URL.createObjectURL(blob);
      cleanup();
    }, 'image/jpeg', 0.85);
  };

  video.onerror = cleanup;
}

async function submitUpload() {
  if (!uploadTargetCollectionId.value || !uploadTargetCategoryId.value) {
    toast('只能给分类上传资源', 'warning');
    return;
  }
  const rawFiles = uploadFileList.value
    .map((file) => file.raw)
    .filter((file): file is File => Boolean(file));

  if (!rawFiles.length) {
    toast('请选择要上传的文件', 'warning');
    return;
  }

  uploading.value = true;
  try {
    for (const file of rawFiles) {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('categoryId', uploadTargetCategoryId.value);
      await uploadCategoryResource(uploadTargetCollectionId.value, formData);
    }
    toast(`成功上传 ${rawFiles.length} 个资源`, 'success');
    uploadDialogVisible.value = false;
    resetUploadDialog();
    await loadCollectionDetail(uploadTargetCollectionId.value, expandedDetailMap[uploadTargetCollectionId.value]?.page || 1);
    await loadInnerCategoryDetail(uploadTargetCategoryId.value);
  } finally {
    uploading.value = false;
  }
}

async function confirmDeleteCollection(collectionId: string) {
  await confirmAction('删除后该展册下的分类、资源和分享链接都会一起删除，确认继续？');
  await deleteCategory(collectionId);
  toast('展册已删除', 'success');
  delete expandedDetailMap[collectionId];
  if (store.categories.length === 1 && store.categoryPage > 1) {
    store.categoryPage -= 1;
  }
  await refreshList();
}

async function confirmDeleteInnerCategory(categoryId: string, collectionId: string) {
  await confirmAction('删除后该分类下的资源和分享链接都会一起删除，确认继续？');
  await deleteCategoryItem(categoryId);
  toast('分类已删除', 'success');
  delete expandedItemDetailMap[categoryId];
  const currentPage = expandedDetailMap[collectionId]?.page || 1;
  await loadCollectionDetail(collectionId, currentPage);
}

async function confirmDeleteResource(resourceId: string, collectionId: string, categoryId: string) {
  await confirmAction('删除资源后不可恢复，确认继续？');
  await deleteResource(resourceId);
  toast('资源已删除', 'success');
  await loadCollectionDetail(collectionId, expandedDetailMap[collectionId]?.page || 1);
  await loadInnerCategoryDetail(categoryId);
}

async function confirmAction(message: string) {
  await ElMessageBox.confirm(message, '二次确认', {
    type: 'warning',
    confirmButtonText: '确认',
    cancelButtonText: '取消',
  });
}

function handleCollectionPageChange(page: number) {
  store.categoryPage = page;
  refreshList();
}

function handleInnerCategoryPageChange(collectionId: string, page: number) {
  loadCollectionDetail(collectionId, page);
}
</script>

<style scoped lang="scss">
.category-page {
  display: grid;
}

.panel-card {
  border-radius: 28px;
  border: 1px solid rgba(123, 162, 255, 0.16);
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 24px 60px rgba(36, 84, 170, 0.08);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
}

.panel-header h2 {
  margin: 0;
  color: #17315f;
}

.panel-header p {
  margin: 6px 0 0;
  color: #6d82a7;
}

.outer-table :deep(.el-table__expanded-cell),
.inner-table :deep(.el-table__expanded-cell) {
  padding: 0;
  background: linear-gradient(180deg, #f7faff 0%, #f2f7ff 100%);
}

.outer-table :deep(.el-table__inner-wrapper::before),
.inner-table :deep(.el-table__inner-wrapper::before) {
  display: none;
}

.outer-table :deep(th.el-table__cell) {
  background: #f8fbff;
  color: #6f85ad;
}

.outer-table :deep(td.el-table__cell),
.inner-table :deep(td.el-table__cell) {
  border-bottom-color: rgba(144, 174, 230, 0.18);
}

.expand-panel {
  padding: 12px 18px 18px;
}

.inner-table {
  border-radius: 22px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid rgba(123, 162, 255, 0.12);
}

.resource-panel {
  padding: 16px 18px 20px;
  background: #f5f8fd;
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(164px, 184px));
  gap: 14px;
}

.resource-card {
  overflow: hidden;
  max-width: 184px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid #e4ebf7;
  box-shadow: 0 8px 18px rgba(43, 77, 130, 0.06);
  cursor: grab;
  transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
}

.resource-sort-badge {
  position: absolute;
  bottom: 7px;
  left: 7px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 22px;
  padding: 0 7px;
  border: 1px solid rgba(255, 255, 255, 0.46);
  border-radius: 999px;
  background: rgba(17, 31, 54, 0.5);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  line-height: 1;
  pointer-events: none;
  user-select: none;
}

.resource-card:hover .resource-sort-badge {
  background: rgba(17, 31, 54, 0.68);
}

.resource-card--dragging {
  opacity: 0.62;
  transform: rotate(1.5deg) scale(0.98);
  box-shadow: 0 18px 32px rgba(60, 102, 190, 0.16);
}

.resource-card--drag-over {
  border-color: rgba(47, 107, 255, 0.52);
  box-shadow: 0 0 0 2px rgba(47, 107, 255, 0.12);
}

.resource-card:hover {
  border-color: rgba(47, 107, 255, 0.24);
  box-shadow: 0 12px 26px rgba(43, 77, 130, 0.1);
}

.resource-list-move {
  transition: all 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.resource-preview {
  position: relative;
  overflow: hidden;
  background: #eef3fb;
}

.resource-preview :deep(.el-image),
.resource-video-cover,
.resource-preview video {
  display: block;
  width: 100%;
  height: 124px;
}

.resource-preview :deep(.el-image__inner),
.resource-video-cover img,
.resource-preview video {
  object-fit: cover;
}

.resource-preview :deep(.el-image) {
  cursor: zoom-in;
}

.resource-video-cover {
  position: relative;
  overflow: hidden;
  padding: 0;
  border: 0;
  background: linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
  cursor: pointer;
}

.resource-video-cover img {
  display: block;
  width: 100%;
  height: 100%;
}

.resource-video-cover::after {
  position: absolute;
  inset: 0;
  content: '';
  background: linear-gradient(180deg, rgba(8, 16, 28, 0.04) 0%, rgba(8, 16, 28, 0.22) 100%);
  pointer-events: none;
}

.resource-video-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  padding: 18px;
  color: #486487;
  font-size: 13px;
  font-weight: 700;
  line-height: 1.35;
  text-align: center;
  word-break: break-word;
}

.resource-video-play {
  position: absolute;
  top: 50%;
  left: 50%;
  z-index: 2;
  width: 34px;
  height: 34px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 8px 18px rgba(17, 31, 54, 0.2);
  transform: translate(-50%, -50%);
}

.resource-video-play::before {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-top: 8px solid transparent;
  border-bottom: 8px solid transparent;
  border-left: 12px solid #2f6fed;
  content: '';
  transform: translate(-36%, -50%);
}

.resource-info {
  padding: 7px 10px 8px;
}

.resource-info strong {
  display: block;
  overflow: hidden;
  color: #17315f;
  font-size: 12px;
  font-weight: 700;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resource-info span {
  display: block;
  overflow: hidden;
  margin: 0;
  color: #6d82a7;
  font-size: 10px;
  line-height: 1.25;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resource-meta {
  color: #8a99b8;
  font-size: 10px;
  line-height: 1.2;
}

.resource-meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 4px;
}

.resource-actions {
  position: absolute;
  top: 6px;
  right: 6px;
  z-index: 3;
  display: flex;
  align-items: center;
  margin: 0;
  opacity: 0;
  transform: translateY(-2px);
  transition: opacity 0.18s ease, transform 0.18s ease;
}


.resource-actions .el-button {
  width: 24px;
  height: 24px;
  padding: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 6px 14px rgba(21, 37, 64, 0.16);
  transition: opacity 0.18s ease;
}

.resource-card:hover .resource-actions {
  opacity: 1;
  transform: translateY(0);
}
.row-actions {
  display: flex;
  justify-content: flex-end;
  gap: 6px;
  flex-wrap: wrap;
}

.description-cell {
  display: block;
  max-width: 100%;
  overflow: hidden;
  color: #606266;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:global(.description-tooltip) {
  max-width: 360px;
  color: #4f5b6b;
  line-height: 1.7;
  white-space: normal;
  word-break: break-word;
}

.share-result {
  display: grid;
  grid-template-columns: 184px 1fr;
  gap: 16px;
  align-items: center;
  margin-top: 8px;
  padding: 16px;
  border: 1px solid rgba(47, 107, 255, 0.16);
  border-radius: 14px;
  background: #f8fbff;
}

.share-result__qr {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  border: 1px solid #e1e8f5;
  border-radius: 10px;
  background: #fff;
}

.share-result__content {
  display: grid;
  gap: 10px;
  min-width: 0;
}

.share-result__content strong {
  overflow: hidden;
  color: #17315f;
  font-size: 15px;
  line-height: 1.4;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.share-result__content a {
  overflow: hidden;
  color: #2f6fed;
  font-size: 13px;
  line-height: 1.5;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.table-pagination,
.inner-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.upload-wall {
  display: flex;
  justify-content: center;
  padding: 10px 0 4px;
}

.upload-wall :deep(.el-upload-list--picture-card) {
  display: flex;
  justify-content: center;
  gap: 14px;
}

.upload-wall :deep(.el-upload--picture-card),
.upload-wall :deep(.el-upload-list--picture-card .el-upload-list__item) {
  width: 148px;
  height: 148px;
  border-radius: 18px;
}

.upload-wall :deep(.el-upload--picture-card) {
  border: 1px dashed rgba(47, 107, 255, 0.42);
  background: linear-gradient(180deg, #f8fbff 0%, #eef5ff 100%);
}

.upload-preview-card {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  border-radius: 18px;
  background: linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
}

.upload-preview-card__media {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-preview-card__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: #526b91;
  font-size: 13px;
  font-weight: 800;
}

.upload-preview-card__placeholder::after {
  position: absolute;
  inset: 0;
  content: '';
  background: linear-gradient(180deg, rgba(8, 16, 28, 0.02) 0%, rgba(8, 16, 28, 0.16) 100%);
}

.upload-preview-card__play {
  position: relative;
  z-index: 1;
  width: 38px;
  height: 38px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 8px 18px rgba(17, 31, 54, 0.18);
}

.upload-preview-card__play::before {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-top: 9px solid transparent;
  border-bottom: 9px solid transparent;
  border-left: 13px solid #2f6fed;
  content: '';
  transform: translate(-35%, -50%);
}

.upload-preview-card__name {
  position: absolute;
  right: 8px;
  bottom: 8px;
  left: 8px;
  z-index: 2;
  overflow: hidden;
  padding: 5px 7px;
  border-radius: 8px;
  background: rgba(17, 31, 54, 0.58);
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  line-height: 1.2;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.upload-preview-card__remove {
  position: absolute;
  top: 7px;
  right: 7px;
  z-index: 3;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: rgba(17, 31, 54, 0.58);
  color: #fff;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
}

.upload-wall__tip {
  width: 100%;
  margin-top: 12px;
  color: #6d82a7;
  font-size: 13px;
  text-align: center;
}

@media (max-width: 560px) {
  .share-result {
    grid-template-columns: 1fr;
  }
}
</style>
