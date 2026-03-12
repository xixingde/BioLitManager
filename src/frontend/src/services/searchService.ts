// 搜索服务
import request from './api';
import type { SearchRequest, PaperSearchResponse, Paper } from '../types/paper';

export const searchService = {
  // 高级搜索 - 多维度组合查询
  advancedSearch: (params: SearchRequest) => {
    return request.get<PaperSearchResponse>('/search', { params });
  },

  // 获取论文详情
  getPaperDetail: (paperId: number) => {
    return request.get<Paper>(`/search/papers/${paperId}`);
  }
};
