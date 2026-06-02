import AxiosUtil from '@/util/axios';

export function createUser(data: {
  name: string;
  loginName: string;
  password: string;
  contactName?: string;
  contactPhone?: string;
}) {
  return AxiosUtil({
    url: '/api/platform/users',
    method: 'post',
    data,
    skipAuth: true,
  });
}

export function loginUser(data: { loginName: string; password: string }) {
  return AxiosUtil({
    url: '/api/platform/auth/login',
    method: 'post',
    data,
    skipAuth: true,
  });
}

export function logoutUser() {
  return AxiosUtil({
    url: '/api/platform/auth/logout',
    method: 'post',
  });
}

export function getCurrentUser() {
  return AxiosUtil({
    url: '/api/platform/auth/profile',
    method: 'get',
  });
}

export function getUserList() {
  return AxiosUtil({
    url: '/api/platform/users',
    method: 'get',
  });
}

export function updateUserRole(id: string, data: { roleCode: 'admin' | 'user' }) {
  return AxiosUtil({
    url: `/api/platform/users/${id}/role`,
    method: 'post',
    data,
  });
}

export function createCategory(data: {
  name: string;
  description?: string;
  coverUrl?: string;
  visible?: boolean;
  status?: string;
}) {
  return AxiosUtil({
    url: '/api/platform/collections',
    method: 'post',
    data,
  });
}

export function deleteCategory(id: string) {
  return AxiosUtil({
    url: `/api/platform/collections/${id}/delete`,
    method: 'post',
  });
}

export function getCurrentUserCategories(params?: { page?: number; pageSize?: number }) {
  return AxiosUtil({
    url: '/api/platform/collections',
    method: 'get',
    params,
  });
}

export function getDashboardStats() {
  return AxiosUtil({
    url: '/api/platform/dashboard/stats',
    method: 'get',
  });
}

export function getCategoryDetail(id: string, params?: { page?: number; pageSize?: number }) {
  return AxiosUtil({
    url: `/api/platform/collections/${id}`,
    method: 'get',
    params,
  });
}

export function createCategoryItem(data: {
  collectionId: string;
  name: string;
  description?: string;
  coverUrl?: string;
  visible?: boolean;
  status?: string;
}) {
  return AxiosUtil({
    url: '/api/platform/categories',
    method: 'post',
    data,
  });
}

export function updateCategoryItem(
  id: string,
  data: {
    name: string;
    description?: string;
    visible?: boolean;
    status?: string;
  }
) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}/update`,
    method: 'post',
    data,
  });
}

export function deleteCategoryItem(id: string) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}/delete`,
    method: 'post',
  });
}

export function getCategoryItemDetail(id: string) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}`,
    method: 'get',
  });
}

export function uploadCategoryResource(id: string, formData: FormData) {
  return AxiosUtil({
    url: `/api/platform/collections/${id}/resources`,
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

export function deleteResource(id: string) {
  return AxiosUtil({
    url: `/api/platform/resources/${id}/delete`,
    method: 'post',
  });
}

export function updateCategoryResourceSort(id: string, data: { resourceRelationIds: string[] }) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}/resources/sort`,
    method: 'post',
    data,
  });
}

export function createShareLink(data: {
  collectionId: string;
  categoryId?: string;
  targetType?: 'collection' | 'category';
  title?: string;
  description?: string;
  expiresAt?: string;
}) {
  return AxiosUtil({
    url: '/api/platform/share-links',
    method: 'post',
    data,
  });
}

export function getShareLinkDetail(code: string) {
  return AxiosUtil({
    url: `/api/platform/share-links/${code}`,
    method: 'get',
    skipAuth: true,
  });
}

export function updateCategory(
  id: string,
  data: {
    name: string;
    description?: string;
    visible?: boolean;
    status?: string;
  }
) {
  return AxiosUtil({
    url: `/api/platform/collections/${id}/update`,
    method: 'post',
    data,
  });
}

export function getShareLinkList(params?: { collectionId?: string; categoryId?: string }) {
  return AxiosUtil({
    url: '/api/platform/share-links',
    method: 'get',
    params,
  });
}

export function deleteShareLink(id: string) {
  return AxiosUtil({
    url: `/api/platform/share-links/${id}/delete`,
    method: 'post',
  });
}
