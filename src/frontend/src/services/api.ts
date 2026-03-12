import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { message } from 'antd';
import { TOKEN_KEY } from '@/utils/constants';
import type { ApiResponse } from '@/types/api';

const baseURL = '/api';
const timeout = 30000;

const axiosInstance: AxiosInstance = axios.create({
  baseURL,
  timeout,
});

const requestInterceptor = (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
};

const responseInterceptor = (response: AxiosResponse): AxiosResponse<ApiResponse> => {
  const data = response.data as ApiResponse;
  const code = String(data.code);
  if (code !== '000000' && code !== '0' && code !== '0') {
    message.error(data.msg || '请求失败');
    throw response;
  }
  return response as AxiosResponse<ApiResponse>;
};

const responseErrorInterceptor = (error: AxiosError): Promise<never> => {
  const { response } = error;

  if (response) {
    const { status, data } = response;

    switch (status) {
      case 401:
        message.error('登录已过期，请重新登录');
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem('biolit_user_info');
        window.location.href = '/login';
        break;
      case 403:
        message.error('没有权限访问该资源');
        break;
      case 404:
        message.error('请求的资源不存在');
        break;
      case 500:
        message.error((data as ApiResponse)?.msg || '服务器错误');
        break;
      default:
        message.error((data as ApiResponse)?.msg || '网络错误');
    }
  } else {
    message.error('网络连接失败，请检查网络设置');
  }

  return Promise.reject(error);
};

axiosInstance.interceptors.request.use(requestInterceptor);
axiosInstance.interceptors.response.use(responseInterceptor, responseErrorInterceptor);

export const api = {
  get: <T = any>(url: string, params?: any): Promise<AxiosResponse<ApiResponse<T>>> => {
    return axiosInstance.get(url, { params }) as any;
  },
  post: <T = any>(url: string, data?: any): Promise<AxiosResponse<ApiResponse<T>>> => {
    return axiosInstance.post(url, data) as any;
  },
  put: <T = any>(url: string, data?: any): Promise<AxiosResponse<ApiResponse<T>>> => {
    return axiosInstance.put(url, data) as any;
  },
  delete: <T = any>(url: string, params?: any): Promise<AxiosResponse<ApiResponse<T>>> => {
    return axiosInstance.delete(url, { params }) as any;
  },
};

export default axiosInstance;
