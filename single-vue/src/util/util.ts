// @ts-ignore
import { ElNotification } from 'element-plus';
// @ts-ignore
import nprogress from 'nprogress';
// @ts-ignore
import QrCode from 'qrcode-decoder';
import CryptoJS from 'crypto-js';
import {useLoading} from "@/components/GlobalLoading";

// 成功提示
// 参数默认值
export function toast(message: string, type: any) {
  ElNotification({
    message,
    type,
  });
}

// 显示全局loading
export function showFullLoading() {
  // console.log('showFullLoading');
  // nprogress.start();
}
// 关闭全局loading
export function hideFullLoading() {
  // console.log('hideFullLoading');
  // nprogress.done();
}

// 传入file对象，返回promise
export function getQrUrl(file: Blob | MediaSource) {
  // 获取临时路径 chrome有效，其他浏览器的方法请自行查找
  const url = window.webkitURL.createObjectURL(file);
  // 初始化
  const qr = new QrCode();
  // 解析二维码，返回promise
  return qr.decodeFromImage(url);
}

export function encryptAes(str: string) {
  const key = CryptoJS.enc.Utf8.parse('single7812345678'); //abcdefghigkliopk密码，16位
  const encryptResult = CryptoJS.AES.encrypt(str, key, {
    mode: CryptoJS.mode.ECB, //aes加密模式cbc
  });
  //把object转化为string
  return String(encryptResult);
}








// 直接读取环境变量

// 从浏览器地址栏动态获取IP和端口
export function getTitle(): string {

  return import.meta.env.VITE_APP_TITLE;
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