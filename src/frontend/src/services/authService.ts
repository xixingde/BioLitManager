import { api } from './api';
import type { UserInfo } from '@/types/user';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },

  getProfile: async (): Promise<UserInfo> => {
    const response = await api.get<UserInfo>('/auth/profile');
    return response.data.data;
  },

  searchUsers: async (params: { keyword?: string; page?: number; size?: number }) => {
    const response = await api.get<{ list: UserInfo[]; total: number }>('/users', { params });
    return response.data;
  },
};
