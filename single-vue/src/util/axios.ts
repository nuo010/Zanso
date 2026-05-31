import axios, {
  AxiosInstance,
  AxiosError,
  AxiosResponse,
  InternalAxiosRequestConfig,
} from 'axios';
import { toast, showFullLoading, hideFullLoading,getBaseURL } from './util';
import { getToken, removeToken } from './auth';
import router from '@/router/route';
import {serverLocation, tokenKey} from "@/util/constants";

// router 不可以直接在router 中引用,会无效,要在创建的实例中使用
const AxiosUtil: AxiosInstance = axios.create({
  baseURL: getBaseURL() + serverLocation,
  timeout: 60000,
});

// 请求拦截器
AxiosUtil.interceptors.request.use(function (config: InternalAxiosRequestConfig) {
    showFullLoading();
    // 如果设置了 skipAuth，则跳过 token 注入
    if (!config.skipAuth) {
      const token = getToken();
      if (token) {
        // 使用类型断言或 set 方法确保类型安全
        if (config.headers) {
          config.headers[tokenKey] = token;
        }
      }
    }
    return config;
  },
  function (error: AxiosError) {
    hideFullLoading();
    console.error('请求拦截器错误:', error);
    return Promise.reject(error);
  }
);

// 定义响应数据类型
interface ApiResponse<T = any> {
  code: number;
  msg: string;
  time: string;
  status: boolean;
  data?: T;
}

// 响应拦截器
AxiosUtil.interceptors.response.use(function (response: AxiosResponse<ApiResponse>) {
    // http状态码是2xx的
    // 默认不静默通知，根据状态码判断，如果不符合预期，就直接弹窗提示错误
    const silent = response.config.silent ?? false;
    hideFullLoading();
    
    const data = response.data;
    
    // 检查响应数据结构
    if (!data) {
      const error = new Error('响应数据格式错误');
      if (!silent) {
        toast('响应数据格式错误', 'error');
      }
      return Promise.reject(error);
    }
    
    if (data.code === 200) {
      // 业务正常返回
      return data as any;
    } else if (data.code === 500) {
      // 业务异常返回
      if (!silent) {
        toast(data.msg || '服务器错误', 'error');
      }
      // 传递错误信息给调用方
      const error = new Error(data.msg || '服务器错误');
      (error as any).code = data.code;
      (error as any).data = data.data;
      return Promise.reject(error);
    } else if (data.code === 555) {
      // 用户token不可用，需要重新登录
      if (!silent) {
        toast('请重新登录!', 'error');
      }
      removeToken();
      // router.push('/555').then(() => {});
      router.push('/login').then(() => {});
      // 传递错误信息给调用方
      const error = new Error(data.msg || '请重新登录');
      (error as any).code = data.code;
      return Promise.reject(error);
    } else {
      // 兜底逻辑，其他错误码
      if (!silent) {
        toast(data.msg || '未知错误', 'error');
      }
      const error = new Error(data.msg || '未知错误');
      (error as any).code = data.code;
      (error as any).data = data.data;
      return Promise.reject(error);
    }
  },
  function (error: AxiosError) {
    // 状态码是2xx以外的
    // http 连接不上,或者其他问题才会走这个错误
    hideFullLoading();
    console.error('axios error', error);
    
    // 检查是否有自定义错误处理器
    const config = error.config as InternalAxiosRequestConfig | undefined;
    if (config?.onError) {
      config.onError(error);
    } else {
      toast(error.message || '网络错误', 'error');
    }
    
    return Promise.reject(error);
  }
);

export default AxiosUtil;
