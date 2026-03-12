// 审核服务
import request from './api';
import type { ReviewLog, PendingReview } from '../types/review';

export const reviewService = {
  // 业务审核
  businessReview: (paperId: number, data: { result: '通过' | '驳回'; comment: string }) => {
    return request.post(`/reviews/business/${paperId}`, data);
  },

  // 政工审核
  politicalReview: (paperId: number, data: { result: '通过' | '驳回'; comment: string }) => {
    return request.post(`/reviews/political/${paperId}`, data);
  },

  // 获取审核记录
  getReviewLogs: (paperId: number) => {
    return request.get<ReviewLog[]>(`/reviews/${paperId}/logs`);
  },

  // 获取待业务审核列表
  getPendingBusinessReviews: () => {
    return request.get<PendingReview[]>('/reviews/pending/business');
  },

  // 获取待政工审核列表
  getPendingPoliticalReviews: () => {
    return request.get<PendingReview[]>('/reviews/pending/political');
  },

  // 获取我的审核记录
  getMyReviews: () => {
    return request.get<ReviewLog[]>('/reviews/my');
  }
};
