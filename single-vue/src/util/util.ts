import { ElNotification } from 'element-plus';
import { appTitle } from '@/util/constants';

export function toast(message: string, type: 'success' | 'warning' | 'info' | 'error') {
  ElNotification({
    message,
    type,
  });
}

export function showFullLoading() {}

export function hideFullLoading() {}

export function getTitle(): string {
  return appTitle;
}

export function getBaseURL(): string {
  return import.meta.env.VITE_API_URL || window.location.origin;
}
