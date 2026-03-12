// 期刊服务
import request from './api';

export interface Journal {
  id: number;
  full_name: string;
  short_name: string;
  issn: string;
  impact_factor: number;
  publisher: string;
  created_at: string;
  updated_at: string;
}

export interface CreateJournalRequest {
  full_name: string;
  short_name: string;
  issn: string;
  impact_factor: number;
  publisher: string;
}

export const journalService = {
  // 创建期刊
  createJournal: (data: CreateJournalRequest) => {
    return request.post<{ id: number }>('/journals', data);
  },

  // 获取期刊详情
  getJournal: (id: number) => {
    return request.get<Journal>(`/journals/${id}`);
  },

  // 分页查询期刊列表
  getJournals: (params?: { page?: number; size?: number }) => {
    return request.get<{ list: Journal[]; total: number; page: number; size: number }>('/journals', { params });
  },

  // 更新期刊
  updateJournal: (id: number, data: Partial<CreateJournalRequest>) => {
    return request.put(`/journals/${id}`, data);
  },

  // 更新影响因子
  updateImpactFactor: (id: number, impact_factor: number) => {
    return request.put(`/journals/${id}/impact-factor`, { impact_factor });
  },

  // 搜索期刊
  searchJournals: (keyword: string) => {
    return request.get<{ list: Journal[] }>('/journals/search', { params: { keyword } });
  }
};
