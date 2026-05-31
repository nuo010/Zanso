import * as axios from 'axios'

declare module 'axios' {
    interface AxiosInstance {
        (config: AxiosRequestConfig): Promise<any>
    }
}



declare module 'axios' {
    interface InternalAxiosRequestConfig {
        /** 是否静默请求（不自动 toast 错误） */
        silent?: boolean;
        /** 是否跳过 token 注入（如登录接口） */
        skipAuth?: boolean;
        /** 自定义错误处理器 */
        onError?: (error: any) => void;
    }
}
