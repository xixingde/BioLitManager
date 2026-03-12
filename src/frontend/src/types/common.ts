// 通用类型定义

/**
 * 分页响应接口
 */
export interface PageResponse<T> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

/**
 * 分页参数接口
 */
export interface PaginationParams {
  page?: number;
  size?: number;
}

/**
 * 通用响应接口
 */
export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data: T;
}

/**
 * 操作结果接口
 */
export interface OperationResult {
  success: boolean;
  message?: string;
  error?: string;
}

/**
 * 批量操作结果接口
 */
export interface BatchOperationResult {
  success_count: number;
  failed_count: number;
  errors: string[];
}

/**
 * 搜索参数接口
 */
export interface SearchParams extends PaginationParams {
  keyword?: string;
}

/**
 * 日期范围接口
 */
export interface DateRange {
  start_date?: string;
  end_date?: string;
}
