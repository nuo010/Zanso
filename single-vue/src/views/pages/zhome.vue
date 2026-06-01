<template>
  <div class="category-page">
    <el-card class="panel-card" shadow="never">
      <template #header>
        <div class="panel-header">
          <div>
            <h2>分类管理</h2>
            <p>分类和分类项都带可见性控制，分享时还能带过期时间，后台和分享治理放一块儿处理更顺手。</p>
          </div>
          <el-button type="primary" @click="openCreateCategoryDialog">新建分类</el-button>
        </div>
      </template>

      <el-table
        :data="store.categories"
        v-loading="loading"
        row-key="id"
        class="outer-table"
        @expand-change="handleCategoryExpand"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="expand-panel" v-loading="expandedLoadingMap[row.id]">
              <template v-if="expandedDetailMap[row.id]">
                <el-table
                  v-if="expandedDetailMap[row.id].categoryItems?.length"
                  :data="expandedDetailMap[row.id].categoryItems"
                  size="small"
                  row-key="id"
                  class="inner-table"
                  :show-header="false"
                  @expand-change="handleCategoryItemExpand"
                >
                  <el-table-column type="expand">
                    <template #default="{ row: itemRow }">
                      <div class="resource-panel" v-loading="itemLoadingMap[itemRow.id]">
                        <template v-if="expandedItemDetailMap[itemRow.id]">
                          <div v-if="expandedItemDetailMap[itemRow.id].resourceList?.length" class="resource-grid">
                            <div
                              v-for="resource in expandedItemDetailMap[itemRow.id].resourceList"
                              :key="resource.id"
                              class="resource-card"
                            >
                              <div class="resource-preview">
                                <el-image
                                  v-if="resource.resourceType !== 'video'"
                                  :src="resource.url"
                                  :alt="resource.fileName"
                                  :preview-src-list="[resource.url]"
                                  :initial-index="0"
                                  fit="cover"
                                  preview-teleported
                                  hide-on-click-modal
                                  loading="lazy"
                                />
                                <video v-else :src="resource.url" controls playsinline preload="metadata"></video>
                              </div>
                              <div class="resource-info">
                                <strong>{{ resource.fileName }}</strong>
                                <span>{{ resource.resourceType }} · {{ resource.mimeType || '未知类型' }}</span>
                                <div class="resource-actions">
                                  <el-button
                                    link
                                    type="danger"
                                    @click="confirmDeleteResource(resource.resourceId, itemRow.categoryId, itemRow.id)"
                                  >
                                    删除资源
                                  </el-button>
                                </div>
                              </div>
                            </div>
                          </div>
                          <el-empty v-else description="这个分类项还没有关联资源" />
                        </template>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="name" min-width="220" />
                  <el-table-column prop="description" min-width="240" />
                  <el-table-column label="展示" width="120" align="center">
                    <template #default="{ row: itemRow }">
                      <el-switch
                        :model-value="itemRow.visible"
                        inline-prompt
                        active-text="可看"
                        inactive-text="隐藏"
                        @change="handleToggleCategoryItemVisible(itemRow, $event)"
                      />
                    </template>
                  </el-table-column>
                  <el-table-column width="360" align="right">
                    <template #default="{ row: itemRow }">
                      <div class="row-actions">
                        <el-button link type="primary" @click="openUploadDialog(itemRow.categoryId, itemRow.id)">上传资源</el-button>
                        <el-button link type="primary" @click="openEditItemDialog(itemRow)">修改</el-button>
                        <el-button link type="warning" @click="openShareDialog('item', itemRow.categoryId, itemRow.id, itemRow.name, itemRow.description)">
                          分享分类项
                        </el-button>
                        <el-button link type="danger" @click="confirmDeleteCategoryItem(itemRow.id, itemRow.categoryId)">
                          删除
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
                <div v-if="expandedDetailMap[row.id].itemTotal > pageSize" class="inner-pagination">
                  <el-pagination
                    small
                    background
                    layout="total, prev, pager, next"
                    :page-size="pageSize"
                    :total="expandedDetailMap[row.id].itemTotal"
                    :current-page="expandedDetailMap[row.id].itemPage || 1"
                    @current-change="handleCategoryItemPageChange(row.id, $event)"
                  />
                </div>
                <el-empty v-else description="当前分类下还没有分类项" />
              </template>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="name" label="分类名称" min-width="220" />
        <el-table-column prop="description" label="描述" min-width="260" />
        <el-table-column label="展示" width="140" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.visible"
              inline-prompt
              active-text="可看"
              inactive-text="隐藏"
              @change="handleToggleCategoryVisible(row, $event)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="340" align="right">
          <template #default="{ row }">
            <div class="row-actions">
              <el-button link type="primary" @click="openCreateItemDialog(row)">新建分类项</el-button>
              <el-button link type="primary" @click="openEditCategoryDialog(row)">修改描述</el-button>
              <el-button link type="warning" @click="openShareDialog('category', row.id, '', row.name, row.description)">分享分类</el-button>
              <el-button link type="danger" @click="confirmDeleteCategory(row.id)">删除</el-button>
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
          @current-change="handleCategoryPageChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="createCategoryDialogVisible" title="新建分类" width="460px">
      <el-form :model="categoryForm" label-position="top">
        <el-form-item label="分类名称">
          <el-input v-model="categoryForm.name" placeholder="直接输入分类名称" />
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
        <el-button type="primary" :loading="submittingCategory" @click="handleCreateCategory">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editCategoryDialogVisible" title="修改分类描述" width="460px">
      <el-form :model="editCategoryForm" label-position="top">
        <el-form-item label="分类名称">
          <el-input :model-value="editCategoryForm.name" disabled />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editCategoryForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editCategoryDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submittingCategory" @click="handleUpdateCategoryDescription">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="categoryItemDialogVisible"
      :title="categoryItemDialogMode === 'edit' ? '修改分类项' : '新建分类项'"
      width="460px"
    >
      <el-form :model="categoryItemForm" label-position="top">
        <el-form-item label="所属分类">
          <el-input :model-value="categoryItemParentName" disabled />
        </el-form-item>
        <el-form-item label="分类项名称">
          <el-input v-model="categoryItemForm.name" placeholder="直接输入分类项名称" />
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
        <el-button type="primary" :loading="submittingCategoryItem" @click="handleSubmitCategoryItem">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="shareDialogVisible" :title="shareForm.targetType === 'item' ? '分享分类项' : '分享分类'" width="460px">
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

    <el-dialog v-model="uploadDialogVisible" :title="uploadTargetItemId ? '上传分类项资源' : '上传分类资源'" width="460px">
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
import { toast } from '@/util/util';

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
  categoryId: '',
  categoryName: '',
  name: '',
  description: '',
  visible: true,
  status: 'active',
});

const shareForm = reactive({
  targetType: 'category' as 'category' | 'item',
  categoryId: '',
  categoryItemId: '',
  title: '',
  description: '',
  expiresAt: '',
});

const uploadTargetCategoryId = ref('');
const uploadTargetItemId = ref('');
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

function openCreateCategoryDialog() {
  categoryForm.name = '';
  categoryForm.description = '';
  categoryForm.visible = true;
  createCategoryDialogVisible.value = true;
}

function openEditCategoryDialog(category: any) {
  editCategoryForm.id = category.id;
  editCategoryForm.name = category.name;
  editCategoryForm.description = category.description || '';
  editCategoryForm.visible = category.visible !== false;
  editCategoryForm.status = category.status || 'active';
  editCategoryDialogVisible.value = true;
}

function openCreateItemDialog(category: any) {
  categoryItemDialogMode.value = 'create';
  categoryItemForm.id = '';
  categoryItemForm.categoryId = category.id;
  categoryItemForm.categoryName = category.name;
  categoryItemForm.name = '';
  categoryItemForm.description = '';
  categoryItemForm.visible = true;
  categoryItemForm.status = 'active';
  categoryItemParentName.value = category.name;
  categoryItemDialogVisible.value = true;
}

function openEditItemDialog(item: any) {
  categoryItemDialogMode.value = 'edit';
  categoryItemForm.id = item.id;
  categoryItemForm.categoryId = item.categoryId;
  categoryItemForm.categoryName = findCategoryName(item.categoryId);
  categoryItemForm.name = item.name;
  categoryItemForm.description = item.description || '';
  categoryItemForm.visible = item.visible !== false;
  categoryItemForm.status = item.status || 'active';
  categoryItemParentName.value = categoryItemForm.categoryName;
  categoryItemDialogVisible.value = true;
}

function openShareDialog(
  targetType: 'category' | 'item',
  categoryId: string,
  categoryItemId: string,
  title: string,
  description?: string
) {
  shareForm.targetType = targetType;
  shareForm.categoryId = categoryId;
  shareForm.categoryItemId = categoryItemId;
  shareForm.title = title || '';
  shareForm.description = description || '';
  shareForm.expiresAt = '';
  shareDialogVisible.value = true;
}

function findCategoryName(categoryId: string) {
  return store.categories.find((item: any) => item.id === categoryId)?.name || '';
}

async function handleCreateCategory() {
  const name = categoryForm.name.trim();
  if (!name) {
    toast('分类名称不能为空', 'warning');
    return;
  }
  if (store.categories.some((item: any) => item.name === name)) {
    toast('分类名称不能重复', 'warning');
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
    toast('分类创建成功', 'success');
    createCategoryDialogVisible.value = false;
    store.categoryPage = 1;
    await refreshList();
  } finally {
    submittingCategory.value = false;
  }
}

async function handleSubmitCategoryItem() {
  const name = categoryItemForm.name.trim();
  if (!categoryItemForm.categoryId || !name) {
    toast('分类项名称不能为空', 'warning');
    return;
  }

  const currentItems = expandedDetailMap[categoryItemForm.categoryId]?.categoryItems || [];
  const duplicate = currentItems.some((item: any) => item.name === name && item.id !== categoryItemForm.id);
  if (duplicate) {
    toast('分类项名称不能重复', 'warning');
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
      toast('分类项修改成功', 'success');
    } else {
      await createCategoryItem({
        categoryId: categoryItemForm.categoryId,
        name,
        description: categoryItemForm.description.trim(),
        visible: categoryItemForm.visible,
        status: 'active',
      });
      toast('分类项创建成功', 'success');
    }
    categoryItemDialogVisible.value = false;
    await loadCategoryDetail(categoryItemForm.categoryId, expandedDetailMap[categoryItemForm.categoryId]?.itemPage || 1);
  } finally {
    submittingCategoryItem.value = false;
  }
}

async function handleUpdateCategoryDescription() {
  if (!editCategoryForm.id) return;
  submittingCategory.value = true;
  try {
    await updateCategory(editCategoryForm.id, {
      name: editCategoryForm.name,
      description: editCategoryForm.description.trim(),
      visible: editCategoryForm.visible,
      status: editCategoryForm.status,
    });
    toast('分类描述修改成功', 'success');
    editCategoryDialogVisible.value = false;
    const target = store.categories.find((item: any) => item.id === editCategoryForm.id);
    if (target) {
      target.description = editCategoryForm.description.trim();
    }
    await loadCategoryDetail(editCategoryForm.id, expandedDetailMap[editCategoryForm.id]?.itemPage || 1);
  } finally {
    submittingCategory.value = false;
  }
}

async function handleSubmitShare() {
  submittingShare.value = true;
  try {
    const res = await createShareLink({
      categoryId: shareForm.categoryId,
      categoryItemId: shareForm.categoryItemId || undefined,
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

async function loadCategoryDetail(id: string, page = 1) {
  expandedLoadingMap[id] = true;
  try {
    const res = await getCategoryDetail(id, { page, pageSize });
    expandedDetailMap[id] = res.data;
  } finally {
    expandedLoadingMap[id] = false;
  }
}

async function loadCategoryItemDetail(id: string) {
  itemLoadingMap[id] = true;
  try {
    const res = await getCategoryItemDetail(id);
    expandedItemDetailMap[id] = res.data;
  } finally {
    itemLoadingMap[id] = false;
  }
}

function handleCategoryExpand(row: any, expandedRows: any[]) {
  if (expandedRows.some((item) => item.id === row.id)) {
    loadCategoryDetail(row.id, expandedDetailMap[row.id]?.itemPage || 1);
  }
}

function handleCategoryItemExpand(row: any, expandedRows: any[]) {
  if (expandedRows.some((item) => item.id === row.id)) {
    loadCategoryItemDetail(row.id);
  }
}

async function handleToggleCategoryVisible(row: any, value: boolean) {
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
  toast(`分类已切换为${value ? '可看' : '不可看'}`, 'success');
}

async function handleToggleCategoryItemVisible(row: any, value: boolean) {
  await updateCategoryItem(row.id, {
    name: row.name,
    description: row.description,
    visible: value,
    status: row.status,
  });
  row.visible = value;
  toast(`分类项已切换为${value ? '可看' : '不可看'}`, 'success');
}

function openUploadDialog(categoryId: string, categoryItemId: string) {
  uploadTargetCategoryId.value = categoryId;
  uploadTargetItemId.value = categoryItemId;
  selectedFile.value = null;
  uploadDialogVisible.value = true;
}

function handleFileChange(file: any) {
  selectedFile.value = file.raw;
}

async function submitUpload() {
  if (!uploadTargetCategoryId.value || !uploadTargetItemId.value) {
    toast('只能给分类项上传资源', 'warning');
    return;
  }
  if (!selectedFile.value) {
    toast('请选择要上传的文件', 'warning');
    return;
  }

  const formData = new FormData();
  formData.append('file', selectedFile.value);
  formData.append('categoryItemId', uploadTargetItemId.value);

  uploading.value = true;
  try {
    await uploadCategoryResource(uploadTargetCategoryId.value, formData);
    toast('资源上传成功', 'success');
    uploadDialogVisible.value = false;
    await loadCategoryDetail(uploadTargetCategoryId.value, expandedDetailMap[uploadTargetCategoryId.value]?.itemPage || 1);
    await loadCategoryItemDetail(uploadTargetItemId.value);
  } finally {
    uploading.value = false;
  }
}

async function confirmDeleteCategory(categoryId: string) {
  await confirmAction('删除后该分类下的分类项、资源和分享链接都会一起删除，确认继续？');
  await deleteCategory(categoryId);
  toast('分类已删除', 'success');
  delete expandedDetailMap[categoryId];
  if (store.categories.length === 1 && store.categoryPage > 1) {
    store.categoryPage -= 1;
  }
  await refreshList();
}

async function confirmDeleteCategoryItem(itemId: string, categoryId: string) {
  await confirmAction('删除后该分类项下的资源和分享链接都会一起删除，确认继续？');
  await deleteCategoryItem(itemId);
  toast('分类项已删除', 'success');
  delete expandedItemDetailMap[itemId];
  const currentPage = expandedDetailMap[categoryId]?.itemPage || 1;
  await loadCategoryDetail(categoryId, currentPage);
}

async function confirmDeleteResource(resourceId: string, categoryId: string, itemId: string) {
  await confirmAction('删除资源后不可恢复，确认继续？');
  await deleteResource(resourceId);
  toast('资源已删除', 'success');
  await loadCategoryDetail(categoryId, expandedDetailMap[categoryId]?.itemPage || 1);
  await loadCategoryItemDetail(itemId);
}

async function confirmAction(message: string) {
  await ElMessageBox.confirm(message, '二次确认', {
    type: 'warning',
    confirmButtonText: '确认',
    cancelButtonText: '取消',
  });
}

function handleCategoryPageChange(page: number) {
  store.categoryPage = page;
  refreshList();
}

function handleCategoryItemPageChange(categoryId: string, page: number) {
  loadCategoryDetail(categoryId, page);
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
