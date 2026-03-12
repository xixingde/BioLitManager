export interface ApiResponse<T = any> {
  code: string;
  msg: string;
  data: T;
}

export interface PageResult<T = any> {
  list: T[];
  total: number;
  page: number;
  size: number;
}
