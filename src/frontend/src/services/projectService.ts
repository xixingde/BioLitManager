// 课题服务
import request from './api';

export interface Project {
  id: number;
  name: string;
  code: string;
  project_type: string;
  source: string;
  level: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  name: string;
  code: string;
  project_type: string;
  source: string;
  level: string;
}

export const projectService = {
  // 创建课题
  createProject: (data: CreateProjectRequest) => {
    return request.post<{ id: number }>('/projects', data);
  },

  // 获取课题详情
  getProject: (id: number) => {
    return request.get<Project>(`/projects/${id}`);
  },

  // 分页查询课题列表
  getProjects: (params?: { page?: number; size?: number; name?: string; code?: string; project_type?: string; level?: string }) => {
    return request.get<{ list: Project[]; total: number; page: number; size: number }>('/projects', { params });
  },

  // 更新课题
  updateProject: (id: number, data: Partial<CreateProjectRequest> & { status?: string }) => {
    return request.put(`/projects/${id}`, data);
  },

  // 删除课题
  deleteProject: (id: number) => {
    return request.delete(`/projects/${id}`);
  },

  // 搜索课题
  searchProjects: (keyword: string) => {
    return request.get<Project[]>('/projects', { params: { name: keyword } });
  }
};
