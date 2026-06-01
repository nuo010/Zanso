import { ElNotification } from 'element-plus';
import { appTitle, resourceDownloadLocation } from '@/util/constants';

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
  return window.location.origin;
}

export function getResourceDownloadBase(): string {
  const configuredBase = import.meta.env.VITE_DOWNLOAD_BASE;
  if (configuredBase) {
    return configuredBase.replace(/\/$/, '');
  }
  return `${window.location.origin}${resourceDownloadLocation}`;
}

export function resolveResourceURL(pathOrUrl: string): string {
  if (!pathOrUrl) return '';
  if (/^(https?:)?\/\//i.test(pathOrUrl) || pathOrUrl.startsWith('data:') || pathOrUrl.startsWith('blob:')) {
    return pathOrUrl;
  }
  const cleanPath = pathOrUrl.replace(/^\/+/, '');
  return `${getResourceDownloadBase()}/${cleanPath}`;
}
