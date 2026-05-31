import axios, {
  AxiosError,
  AxiosInstance,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';
import { hideFullLoading, showFullLoading, toast, getBaseURL } from './util';
import { getToken, removeToken } from './auth';
import router from '@/router/route';
import { serverLocation, tokenKey } from '@/util/constants';

const AxiosUtil: AxiosInstance = axios.create({
  baseURL: getBaseURL() + serverLocation,
  timeout: 60000,
});

AxiosUtil.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    showFullLoading();
    if (!config.skipAuth) {
      const token = getToken();
      if (token && config.headers) {
        config.headers[tokenKey] = token;
      }
    }
    return config;
  },
  (error: AxiosError) => {
    hideFullLoading();
    return Promise.reject(error);
  }
);

interface ApiResponse<T = any> {
  code: number;
  msg: string;
  time: string;
  status: boolean;
  data?: T;
}

AxiosUtil.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    hideFullLoading();
    const data = response.data;
    if (!data) {
      toast('响应数据格式错误', 'error');
      return Promise.reject(new Error('响应数据格式错误'));
    }
    if (data.code === 200) {
      return data as any;
    }
    toast(data.msg || '请求失败', 'error');
    return Promise.reject(new Error(data.msg || '请求失败'));
  },
  (error: AxiosError) => {
    hideFullLoading();
    if (error.response?.status === 401) {
      removeToken();
      router.push('/login');
    }
    toast(error.message || '网络错误', 'error');
    return Promise.reject(error);
  }
);

export default AxiosUtil;
