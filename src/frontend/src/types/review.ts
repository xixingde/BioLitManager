// 审核相关类型定义

export interface ReviewLog {
  id: number;
  paper_id: number;
  review_type: 'business' | 'political';
  result: 'approved' | 'rejected';
  comment: string;
  reviewer?: User;
  review_time: string;
  created_at: string;
}

export interface ReviewForm {
  result: '通过' | '驳回';
  comment: string;
}

export interface PendingReview {
  id: number;
  title: string;
  submitter_name: string;
  submit_time: string;
  status: string;
  days_since_submit: number;
}

export interface User {
  id: number;
  username: string;
  name: string;
  role: string;
}
