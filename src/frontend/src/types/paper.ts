// 论文相关类型定义

export interface Journal {
  id: number;
  full_name: string;
  short_name: string;
  issn: string;
  impact_factor: number;
  publisher: string;
}

export interface Author {
  id: number;
  name: string;
  author_type: 'first_author' | 'co_first_author' | 'corresponding_author' | 'author';
  rank: number;
  department: string;
  user_id?: number;
}

export interface Project {
  id: number;
  name: string;
  code: string;
  project_type: string;
  source: string;
  level: string;
}

export interface Attachment {
  id: number;
  file_type: string;
  file_name: string;
  file_path: string;
  file_size: number;
  mime_type: string;
  created_at: string;
}

export interface User {
  id: number;
  username: string;
  name: string;
  role: string;
}

export interface Paper {
  id: number;
  title: string;
  abstract: string;
  journal?: Journal;
  doi: string;
  impact_factor: number;
  publish_date?: string;
  status: 'draft' | '待业务审核' | '待政工审核' | '审核通过' | '驳回';
  submitter?: User;
  submit_time: string;
  authors?: Author[];
  projects?: Project[];
  attachments?: Attachment[];
  created_at: string;
  updated_at: string;
}

export interface PaperForm {
  title: string;
  abstract: string;
  journal_id: number;
  doi: string;
  impact_factor: number;
  publish_date?: string;
  authors: Author[];
  projects: number[];
}

export interface PageResponse<T> {
  list: T[];
  total: number;
  page: number;
  size: number;
}

// 逻辑类型
export type LogicType = 'AND' | 'OR' | 'NOT';

// 作者类型筛选
export type AuthorTypeFilter = 'first_author' | 'co_first_author' | 'corresponding_author' | 'all';

// 排序字段
export type SortField = 'publish_date' | 'impact_factor' | 'created_at' | 'title' | 'journal_name';

// 排序方向
export type SortOrder = 'asc' | 'desc';

// 单个查询条件
export interface QueryCondition {
  field: string;
  operator: 'eq' | 'ne' | 'gt' | 'gte' | 'lt' | 'lte' | 'like' | 'in' | 'between';
  value: string | number | string[] | number[];
}

// 分页信息
export interface PaginationInfo {
  page: number;
  size: number;
  total?: number;
}

// 排序信息
export interface SortInfo {
  field: SortField;
  order: SortOrder;
}

// 搜索请求
export interface SearchRequest {
  query?: QueryCondition[];
  logic?: LogicType;
  pagination?: PaginationInfo;
  sort?: SortInfo;
  author_type_filter?: AuthorTypeFilter;
}

// 归档信息
export interface ArchiveInfo {
  is_archived: boolean;
  archive_time?: string;
  archive_user?: User;
  archive_reason?: string;
}

// 搜索结果项
export interface PaperSearchResult extends Paper {
  archive_info?: ArchiveInfo;
  match_score?: number;
}

export interface PaperListResponse extends PageResponse<Paper> {}

// 搜索响应
export interface PaperSearchResponse extends PageResponse<PaperSearchResult> {}
