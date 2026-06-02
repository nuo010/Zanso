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
    url: '/api/platform/categories',
    method: 'post',
    data,
  });
}

export function deleteCategory(id: string) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}/delete`,
    method: 'post',
  });
}

export function getUserCategories(userId: string, params?: { page?: number; pageSize?: number }) {
  return AxiosUtil({
    url: `/api/platform/users/${userId}/categories`,
    method: 'get',
    params,
  });
}

export function getCategoryDetail(id: string, params?: { page?: number; pageSize?: number }) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}`,
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
    url: '/api/platform/category-items',
    method: 'post',
    data: {
      ...data,
      categoryId: data.collectionId,
    },
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
    url: `/api/platform/category-items/${id}/update`,
    method: 'post',
    data,
  });
}

export function deleteCategoryItem(id: string) {
  return AxiosUtil({
    url: `/api/platform/category-items/${id}/delete`,
    method: 'post',
  });
}

export function getCategoryItemDetail(id: string) {
  return AxiosUtil({
    url: `/api/platform/category-items/${id}`,
    method: 'get',
  });
}

export function uploadCategoryResource(id: string, formData: FormData) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}/resources`,
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
    data: {
      ...data,
      categoryItemId: data.categoryId,
    },
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
    url: `/api/platform/categories/${id}/update`,
    method: 'post',
    data,
  });
}

export function getShareLinkList(params?: { collectionId?: string; categoryId?: string }) {
  return AxiosUtil({
    url: '/api/platform/share-links',
    method: 'get',
    params: {
      categoryId: params?.collectionId,
      categoryItemId: params?.categoryId,
    },
  });
}

export function deleteShareLink(id: string) {
  return AxiosUtil({
    url: `/api/platform/share-links/${id}/delete`,
    method: 'post',
  });
}
