// 文件服务
import request from './api';

export interface Attachment {
  id: number;
  file_type: string;
  file_name: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

export const fileService = {
  // 上传文件
  uploadFile: (paperId: number, fileType: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('paper_id', paperId.toString());
    formData.append('file_type', fileType);
    return request.post<{ id: number; file_name: string; file_size: number; file_type: string }>('/files/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
  },

  // 获取文件信息
  getFile: (id: number) => {
    return request.get<Attachment>(`/files/${id}`);
  },

  // 下载文件
  downloadFile: (id: number) => {
    return request.get(`/files/${id}/download`, {
      responseType: 'blob'
    });
  },

  // 删除文件
  deleteFile: (id: number) => {
    return request.delete(`/files/${id}`);
  }
};
