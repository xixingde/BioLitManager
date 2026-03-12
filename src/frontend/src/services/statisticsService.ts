// 统计服务
import { paperService } from './paperService';
import { reviewService } from './reviewService';
import type { Statistics, RecentPaper, PendingReviewTask, HomeData } from '../types/statistics';
import type { Paper, PaperListResponse } from '../types/paper';
import type { PendingReview } from '../types/review';
import type { AxiosResponse } from 'axios';
import type { ApiResponse } from '../types/api';

/**
 * 统计论文列表中的各状态数量
 */
const countPaperStatus = (papers: Paper[]): Omit<Statistics, 'myDraft' | 'myTotal'> => {
  let total = 0;
  let pendingBusiness = 0;
  let pendingPolitical = 0;
  let approved = 0;
  let rejected = 0;
  let draft = 0;

  papers.forEach((paper) => {
    total++;
    switch (paper.status) {
      case 'draft':
        draft++;
        break;
      case '待业务审核':
        pendingBusiness++;
        break;
      case '待政工审核':
        pendingPolitical++;
        break;
      case '审核通过':
        approved++;
        break;
      case '驳回':
        rejected++;
        break;
    }
  });

  return {
    total,
    pendingBusiness,
    pendingPolitical,
    approved,
    rejected,
    draft,
  };
};

/**
 * 获取论文统计数据
 */
export const getStatistics = async (userId?: string): Promise<Statistics> => {
  // 获取所有论文列表进行统计
  const allPapersResponse = await paperService.getPapers({ page: 1, size: 1000 }) as unknown as AxiosResponse<ApiResponse<PaperListResponse>>;
  const allPapers = (allPapersResponse.data as ApiResponse<PaperListResponse>).data?.list || [];
  
  const baseStats = countPaperStatus(allPapers);

  // 如果有用户ID，获取该用户的论文统计
  let myDraft = 0;
  let myTotal = 0;

  if (userId) {
    const myPapersResponse = await paperService.getMyPapers({ page: 1, size: 1000 }) as unknown as AxiosResponse<ApiResponse<PaperListResponse>>;
    const myPapers = (myPapersResponse.data as ApiResponse<PaperListResponse>).data?.list || [];
    myTotal = myPapers.length;
    myDraft = myPapers.filter((p: Paper) => p.status === 'draft').length;
  }

  return {
    ...baseStats,
    myDraft,
    myTotal,
  };
};

/**
 * 获取待审核数据
 */
export const getPendingReviews = async (): Promise<PendingReviewTask[]> => {
  try {
    const [businessReviewsResponse, politicalReviewsResponse] = await Promise.all([
      reviewService.getPendingBusinessReviews(),
      reviewService.getPendingPoliticalReviews(),
    ]) as unknown as [AxiosResponse<ApiResponse<PendingReview[]>>, AxiosResponse<ApiResponse<PendingReview[]>>];

    const businessReviews = (businessReviewsResponse.data as ApiResponse<PendingReview[]>)?.data || [];
    const politicalReviews = (politicalReviewsResponse.data as ApiResponse<PendingReview[]>)?.data || [];

    const businessTasks: PendingReviewTask[] = businessReviews.map((item: PendingReview) => ({
      ...item,
      review_type: 'business' as const,
    }));

    const politicalTasks: PendingReviewTask[] = politicalReviews.map((item: PendingReview) => ({
      ...item,
      review_type: 'political' as const,
    }));

    // 合并并按提交时间倒序排列
    return [...businessTasks, ...politicalTasks].sort(
      (a, b) => new Date(b.submit_time).getTime() - new Date(a.submit_time).getTime()
    );
  } catch (error) {
    console.error('获取待审核数据失败:', error);
    return [];
  }
};

/**
 * 获取最近录入的论文
 */
export const getRecentPapers = async (limit: number = 5): Promise<RecentPaper[]> => {
  try {
    const response = await paperService.getPapers({ page: 1, size: limit * 2 }) as unknown as AxiosResponse<ApiResponse<PaperListResponse>>;
    const papers = (response.data as ApiResponse<PaperListResponse>)?.data?.list || [];

    // 转换为 RecentPaper 格式并按时间倒序
    const recentPapers: RecentPaper[] = papers
      .map((paper: Paper) => ({
        id: paper.id,
        title: paper.title,
        status: paper.status,
        submitter_name: paper.submitter?.name || paper.submitter?.username || '未知',
        submit_time: paper.submit_time,
        created_at: paper.created_at,
      }))
      .sort((a: RecentPaper, b: RecentPaper) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      .slice(0, limit);

    return recentPapers;
  } catch (error) {
    console.error('获取最近论文失败:', error);
    return [];
  }
};

/**
 * 获取首页聚合数据
 */
export const getHomeData = async (userId?: string): Promise<HomeData> => {
  const [statistics, recentPapers, pendingReviews] = await Promise.all([
    getStatistics(userId),
    getRecentPapers(5),
    getPendingReviews(),
  ]);

  return {
    statistics,
    recentPapers,
    pendingReviews,
  };
};
