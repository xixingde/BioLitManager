// 统计相关类型定义

/**
 * 论文统计数据
 */
export interface Statistics {
  /** 论文总数 */
  total: number;
  /** 待业务审核数量 */
  pendingBusiness: number;
  /** 待政工审核数量 */
  pendingPolitical: number;
  /** 已通过数量 */
  approved: number;
  /** 驳回数量 */
  rejected: number;
  /** 草稿数量 */
  draft: number;
  /** 我的草稿数量 */
  myDraft: number;
  /** 我的论文总数 */
  myTotal: number;
}

/**
 * 最近论文
 */
export interface RecentPaper {
  id: number;
  title: string;
  status: 'draft' | '待业务审核' | '待政工审核' | '审核通过' | '驳回';
  submitter_name?: string;
  submit_time: string;
  created_at: string;
}

/**
 * 待审核任务
 */
export interface PendingReviewTask {
  id: number;
  title: string;
  submitter_name: string;
  submit_time: string;
  status: string;
  days_since_submit: number;
  review_type: 'business' | 'political';
}

/**
 * 首页聚合数据
 */
export interface HomeData {
  statistics: Statistics;
  recentPapers: RecentPaper[];
  pendingReviews: PendingReviewTask[];
}

/**
 * 基础统计
 */
export interface BasicStats {
  /** 论文总数 */
  totalPapers: number;
  /** 各年份论文数量 */
  yearlyCounts: Record<number, number>;
  /** 各收录类型论文数量：SCI/EI/CI/DI/CORE */
  typeCounts: Record<string, number>;
  /** 各期刊论文数量 */
  journalCounts: Array<{ journal: string; count: number }>;
  /** 平均影响因子 */
  avgImpactFactor: number;
  /** 总引用次数 */
  totalCitations: number;
  /** 总他引次数 */
  totalSelfCitations: number;
}

/**
 * 按作者统计
 */
export interface AuthorStats {
  /** 作者ID */
  authorId: number;
  /** 作者姓名 */
  authorName: string;
  /** 论文数量 */
  paperCount: number;
  /** 第一作者数量 */
  firstAuthorCount: number;
  /** 通讯作者数量 */
  correspondingAuthorCount: number;
  /** 平均影响因子 */
  avgImpactFactor: number;
  /** 总引用次数 */
  totalCitations: number;
}

/**
 * 按课题统计
 */
export interface ProjectStats {
  /** 课题ID */
  projectId: number;
  /** 课题名称 */
  projectName: string;
  /** 课题编号 */
  projectCode: string;
  /** 论文数量 */
  paperCount: number;
  /** 高影响因子论文数量 IF>=5 */
  highImpactPaperCount: number;
  /** SCI论文数量 */
  sciPaperCount: number;
}

/**
 * 按单位统计
 */
export interface DepartmentStats {
  /** 单位名称 */
  department: string;
  /** 论文数量 */
  paperCount: number;
  /** 总影响因子 */
  totalImpactFactor: number;
  /** 总引用次数 */
  totalCitations: number;
}

/**
 * 年度统计
 */
export interface YearlyStats {
  /** 年度数据 */
  years: Array<{ year: number; count: number }>;
}

/**
 * 期刊统计
 */
export interface JournalStats {
  /** 期刊数据 */
  journals: Array<{ journal: string; count: number; avgImpactFactor: number }>;
}

/**
 * 导出请求参数
 */
export interface ExportRequest {
  /** 导出类型 */
  type: 'papers' | 'paper' | 'stats';
  /** 导出格式 */
  format: 'excel' | 'pdf' | 'word';
  /** 论文ID列表，用于查询结果导出 */
  paperIds?: number[];
  /** 单篇论文ID */
  paperId?: number;
  /** 统计类型 */
  statsType?: string;
  /** 导出字段列表 */
  fields?: string[];
}

/**
 * 导出字段定义
 */
export interface ExportField {
  /** 字段键名 */
  key: string;
  /** 字段标签 */
  label: string;
  /** 是否选中 */
  selected: boolean;
}

/**
 * 用户角色
 */
export interface UserRole {
  /** 角色类型 */
  role: 'admin' | 'user' | 'business_reviewer' | 'political_reviewer';
}
