// 统计服务
import { api } from './api';
import type {
  BasicStats,
  AuthorStats,
  ProjectStats,
  DepartmentStats,
  YearlyStats,
  JournalStats
} from '../types/statistics';

export const statsService = {
  /**
   * 获取基础统计
   */
  getBasicStats() {
    return api.get<BasicStats>('/stats/basic');
  },

  /**
   * 获取作者统计
   */
  getAuthorStats(authorId: number) {
    return api.get<AuthorStats>(`/stats/author/${authorId}`);
  },

  /**
   * 获取课题统计
   */
  getProjectStats(projectId: number) {
    return api.get<ProjectStats>(`/stats/project/${projectId}`);
  },

  /**
   * 获取单位统计
   */
  getDepartmentStats(department: string) {
    return api.get<DepartmentStats>('/stats/department', { name: department });
  },

  /**
   * 获取年度统计
   */
  getYearlyStats() {
    return api.get<YearlyStats>('/stats/yearly');
  },

  /**
   * 获取期刊统计
   */
  getJournalStats() {
    return api.get<JournalStats>('/stats/journal');
  }
};
