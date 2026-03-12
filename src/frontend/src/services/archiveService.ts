// 归档服务
import request from './api';
import type { ArchiveInfo, Paper } from '../types/paper';

// 归档列表查询参数
export interface ArchiveListParams {
  year?: number;
  paperType?: string;
  author?: string;
  projectCode?: string;
  page?: number;
  pageSize?: number;
}

// 归档列表响应
export interface ArchiveListResponse {
  list: ArchiveInfo[];
  total: number;
  page: number;
  size: number;
}

// 修改申请请求
export interface ModifyRequest {
  requestType: string;
  requestReason: string;
  requestData?: Record<string, any>;
}

// 操作结果响应
export interface OperationResult {
  success: boolean;
  message?: string;
}

export const archiveService = {
  // 获取归档列表
  getArchiveList: (params: ArchiveListParams) => {
    return request.get<ArchiveListResponse>('/archives', { params });
  },

  // 获取归档详情
  getArchiveByPaperId: (paperId: number) => {
    return request.get<ArchiveInfo>(`/archives/paper/${paperId}`);
  },

  // 隐藏归档论文
  hideArchive: (paperId: number) => {
    return request.put<OperationResult>(`/archives/${paperId}/hide`);
  },

  // 提交修改申请
  submitModifyRequest: (paperId: number, data: ModifyRequest) => {
    return request.post<OperationResult>(`/archives/${paperId}/modify`, data);
  }
};
