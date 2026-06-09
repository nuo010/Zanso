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

function trimTrailingSlash(url: string): string {
  return url.trim().replace(/\/+$/, '');
}

function getConfiguredAPIBase(): string {
  return trimTrailingSlash(import.meta.env.VITE_API_URL || '');
}

function isLocalHostName(hostname: string): boolean {
  const normalizedHost = hostname.toLowerCase();
  return normalizedHost === 'localhost' || normalizedHost === '127.0.0.1' || normalizedHost === '::1';
}

function isLocalURL(url: string): boolean {
  if (!url) return true;
  try {
    const parsedURL = new URL(url, window.location.origin);
    return isLocalHostName(parsedURL.hostname);
  } catch {
    return false;
  }
}

export function getBaseURL(): string {
  const configuredBase = getConfiguredAPIBase();
  if (configuredBase && !isLocalURL(configuredBase)) {
    return configuredBase;
  }
  return window.location.origin;
}

export function getResourceDownloadBase(): string {
  const configuredBase = trimTrailingSlash(import.meta.env.VITE_DOWNLOAD_BASE || '');
  if (configuredBase && !isLocalURL(configuredBase)) {
    return configuredBase;
  }
  return `${getBaseURL()}${resourceDownloadLocation}`;
}

export function resolveResourceURL(pathOrUrl: string): string {
  if (!pathOrUrl) return '';
  if (/^(https?:)?\/\//i.test(pathOrUrl) || pathOrUrl.startsWith('data:') || pathOrUrl.startsWith('blob:')) {
    return pathOrUrl;
  }
  const cleanPath = pathOrUrl.replace(/^\/+/, '');
  return `${getResourceDownloadBase()}/${cleanPath}`;
}
