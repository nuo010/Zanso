<template>
  <div class="category-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>展册管理</h2>
            <p>展册下面挂分类，分类下面挂资源，分享和可见性都在一套流程里管，别再整那套“分类套分类项”的绕口令了。</p>
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
                          <div v-if="expandedItemDetailMap[categoryRow.id].resourceList?.length" class="resource-grid">
                            <div
                              v-for="resource in expandedItemDetailMap[categoryRow.id].resourceList"
                              :key="resource.id"
                              class="resource-card"
                            >
                              <div class="resource-preview">
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
                                <video v-else :src="resource.fileUrl" controls playsinline preload="metadata"></video>
                              </div>
                              <div class="resource-info">
                                <strong>{{ resource.fileName }}</strong>
                                <span>{{ resource.resourceType }} · {{ resource.mimeType || '未知类型' }}</span>
                                <div class="resource-actions">
                                  <el-button
                                    link
                                    type="danger"
                                    @click="confirmDeleteResource(resource.resourceId, categoryRow.categoryId, categoryRow.id)"
                                  >
                                    删除资源
                                  </el-button>
                                </div>
                              </div>
                            </div>
                          </div>
                          <el-empty v-else description="这个分类下还没有关联资源" />
                        </template>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="name" min-width="220" />
                  <el-table-column prop="description" min-width="240" />
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
                  description="当前展册下还没有分类"
                />
              </template>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="展册名称" min-width="220" />
        <el-table-column prop="description" label="描述" min-width="260" />
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
      <template #footer>
        <el-button @click="shareDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingShare" @click="handleSubmitShare">生成分享链接</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="uploadDialogVisible" title="上传分类资源" width="460px">
      <el-upload drag :auto-upload="false" :limit="1" :on-change="handleFileChange">
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div>拖拽文件到这里，或者点击选择</div>
      </el-upload>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="uploading" @click="submitUpload">上传</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { ElMessageBox } from 'element-plus';
import { UploadFilled } from '@element-plus/icons-vue';
import {
  createCategory,
  createCategoryItem,
  createShareLink,
  deleteCategory,
  deleteCategoryItem,
  deleteResource,
  getCategoryDetail,
  getCategoryItemDetail,
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

const uploadTargetCollectionId = ref('');
const uploadTargetCategoryId = ref('');
const selectedFile = ref<File | null>(null);
const categoryItemParentName = ref('');

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
    toast(`分享链接已生成：${res.data.shareUrl}`, 'success');
    shareDialogVisible.value = false;
  } finally {
    submittingShare.value = false;
  }
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
  }));
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
  selectedFile.value = null;
  uploadDialogVisible.value = true;
}

function handleFileChange(file: any) {
  selectedFile.value = file.raw;
}

async function submitUpload() {
  if (!uploadTargetCollectionId.value || !uploadTargetCategoryId.value) {
    toast('只能给分类上传资源', 'warning');
    return;
  }
  if (!selectedFile.value) {
    toast('请选择要上传的文件', 'warning');
    return;
  }

  const formData = new FormData();
  formData.append('file', selectedFile.value);
  formData.append('categoryItemId', uploadTargetCategoryId.value);

  uploading.value = true;
  try {
    await uploadCategoryResource(uploadTargetCollectionId.value, formData);
    toast('资源上传成功', 'success');
    uploadDialogVisible.value = false;
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
  padding: 16px 20px 20px;
  background: linear-gradient(180deg, #f7fbff 0%, #eef5ff 100%);
}

.resource-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(130px, 150px));
  gap: 12px;
}

.resource-card {
  overflow: hidden;
  max-width: 150px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(116, 153, 230, 0.16);
  box-shadow: 0 10px 24px rgba(60, 102, 190, 0.07);
}

.resource-preview {
  background: linear-gradient(180deg, #edf4ff 0%, #dfeafe 100%);
}

.resource-preview :deep(.el-image),
.resource-preview video {
  display: block;
  width: 100%;
  height: 88px;
  cursor: zoom-in;
}

.resource-preview :deep(.el-image__inner),
.resource-preview video {
  object-fit: cover;
}

.resource-info {
  padding: 8px 10px 10px;
}

.resource-info strong {
  display: block;
  color: #17315f;
  font-size: 12px;
  line-height: 1.5;
  word-break: break-all;
}

.resource-info span {
  display: block;
  margin: 3px 0 6px;
  color: #6d82a7;
  font-size: 11px;
}

.resource-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  align-items: center;
}

.row-actions {
  display: flex;
  justify-content: flex-end;
  gap: 6px;
  flex-wrap: wrap;
}

.table-pagination,
.inner-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
