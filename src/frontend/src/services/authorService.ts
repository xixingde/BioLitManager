// 作者服务
import request from './api';
import type { Author, AuthorForm } from '../types/author';
import type { PageResponse } from '../types/common';

export const authorService = {
  /**
   * 获取论文的作者列表
   * @param paperId 论文ID
   */
  getAuthorsByPaperID: (paperId: number) => {
    return request.get<Author[]>(`/papers/${paperId}/authors`);
  },

  /**
   * 创建作者
   * @param data 作者数据
   */
  createAuthor: (data: AuthorForm) => {
    return request.post<{ id: number }>('/authors', data);
  },

  /**
   * 更新作者
   * @param id 作者ID
   * @param data 作者数据
   */
  updateAuthor: (id: number, data: Partial<AuthorForm>) => {
    return request.put(`/authors/${id}`, data);
  },

  /**
   * 删除作者
   * @param id 作者ID
   */
  deleteAuthor: (id: number) => {
    return request.delete(`/authors/${id}`);
  },

  /**
   * 批量创建作者
   * @param data 作者列表数据
   */
  batchCreateAuthors: (data: AuthorForm[]) => {
    return request.post('/authors/batch', { authors: data });
  },

  /**
   * 搜索用户（用于人员库）
   * @param params 搜索参数
   */
  searchUsers: (params: { keyword?: string; page?: number; size?: number }) => {
    return request.get<PageResponse<{ id: number; username: string; name: string }>>('/users', { params });
  }
};
