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

// 从浏览器地址栏动态获取IP和端口
export function getBaseURL(): string {
    // import.meta.env.VITE_APP_API
    const protocol = window.location.protocol;
    const hostname = window.location.hostname;
    const port = window.location.port;

    if (hostname === 'localhost' || hostname === '127.0.0.1') {
        return import.meta.env.VITE_API_URL
    }

    // 如果端口存在且不是80/443，则包含端口
    if (port && port !== '80' && port !== '443') {
        return `${protocol}//${hostname}:${port}`;
    }
    return `${protocol}//${hostname}`;
}