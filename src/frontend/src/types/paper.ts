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

export interface PaperListResponse extends PageResponse<Paper> {}
