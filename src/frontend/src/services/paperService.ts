// 论文服务
import request from './api';
import type { Paper, PaperListResponse, PaperForm } from '../types/paper';

export const paperService = {
  // 创建论文
  createPaper: (data: PaperForm) => {
    return request.post<{ id: number }>('/papers', data);
  },

  // 获取论文详情
  getPaper: (id: number) => {
    return request.get<Paper>(`/papers/${id}`);
  },

  // 分页查询论文列表
  getPapers: (params?: { page?: number; size?: number; status?: string; keyword?: string }) => {
    return request.get<PaperListResponse>('/papers', { params });
  },

  // 更新论文
  updatePaper: (id: number, data: Partial<PaperForm>) => {
    return request.put(`/papers/${id}`, data);
  },

  // 删除论文
  deletePaper: (id: number) => {
    return request.delete(`/papers/${id}`);
  },

  // 提交审核
  submitForReview: (id: number) => {
    return request.post(`/papers/${id}/submit`);
  },

  // 保存草稿
  saveDraft: (id: number, data: Partial<PaperForm>) => {
    return request.post(`/papers/${id}/save-draft`, data);
  },

  // 检查重复
  checkDuplicate: (data: { title: string; doi: string }) => {
    return request.post<{ count: number; papers: Paper[] }>('/papers/check-duplicate', data);
  },

  // 批量导入
  batchImport: (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return request.post<{ success: number; failed: number; errors: string[] }>('/papers/batch-import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
  },

  // 获取我的论文
  getMyPapers: (params?: { page?: number; size?: number }) => {
    return request.get<PaperListResponse>('/papers/my', { params });
  },

  // 下载导入模板
  downloadImportTemplate: () => {
    return request.get('/papers/import-template', {
      responseType: 'blob'
    });
  }
};
