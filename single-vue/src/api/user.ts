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

export function createCategory(data: {
  name: string;
  description?: string;
  coverUrl?: string;
  status?: string;
}) {
  return AxiosUtil({
    url: '/api/platform/categories',
    method: 'post',
    data,
  });
}

export function getUserCategories(userId: string) {
  return AxiosUtil({
    url: `/api/platform/users/${userId}/categories`,
    method: 'get',
  });
}

export function getCategoryDetail(id: string) {
  return AxiosUtil({
    url: `/api/platform/categories/${id}`,
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

export function createShareLink(data: {
  categoryId: string;
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
