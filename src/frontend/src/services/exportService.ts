// 导出服务
import request from './api';
import type { ExportField } from '../types/statistics';

/**
 * 导出论文参数
 */
export interface ExportPapersParams {
  /** 搜索参数 */
  searchParams?: Record<string, any>;
  /** 导出字段 */
  fields?: string[];
}

/**
 * 导出服务
 */
export const exportService = {
  /**
   * 导出查询结果（Excel格式）
   * @param params 导出参数
   * @returns 导出文件路径
   */
  exportPapers: (params: ExportPapersParams) => {
    return request.post<{ filePath: string }>('/export/papers', params);
  },

  /**
   * 导出单篇论文
   * @param paperId 论文ID
   * @param format 导出格式
   * @returns 导出文件路径
   */
  exportPaper: (paperId: number, format: 'pdf' | 'word') => {
    return request.get<{ filePath: string }>(`/export/paper/${paperId}`, { params: { format } });
  },

  /**
   * 导出统计结果
   * @param statsType 统计类型
   * @param format 导出格式
   * @returns 导出文件路径
   */
  exportStats: (statsType: string, format: 'excel' | 'pdf') => {
    return request.post<{ filePath: string }>('/export/stats', { statsType, format });
  },

  /**
   * 获取可导出字段列表
   * @returns 导出字段列表
   */
  getExportFields: () => {
    return request.get<ExportField[]>('/export/fields');
  },

  /**
   * 下载文件
   * @param url 文件URL
   * @param filename 文件名
   */
  downloadFile: (url: string, filename: string) => {
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  },

  /**
   * 检查导出权限
   * @param paperId 论文ID
   * @returns 是否有权限
   */
  checkExportPermission: (paperId: number) => {
    return request.get<{ hasPermission: boolean }>(`/export/paper/${paperId}`);
  }
};

export default exportService;
