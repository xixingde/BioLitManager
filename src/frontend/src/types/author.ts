// 作者相关类型定义

/**
 * 作者类型枚举
 */
export type AuthorType = 'first_author' | 'co_first_author' | 'corresponding_author' | 'author';

/**
 * 作者接口
 */
export interface Author {
  id: number;
  name: string;
  author_type: AuthorType;
  rank: number;
  department: string;
  user_id?: number;
}

/**
 * 作者表单接口
 */
export interface AuthorForm {
  name: string;
  author_type: AuthorType;
  rank: number;
  department: string;
  user_id?: number;
}

/**
 * 作者类型选项
 */
export const AUTHOR_TYPE_OPTIONS = [
  { label: '第一作者', value: 'first_author' as AuthorType },
  { label: '共同第一作者', value: 'co_first_author' as AuthorType },
  { label: '通讯作者', value: 'corresponding_author' as AuthorType },
  { label: '普通作者', value: 'author' as AuthorType },
];

/**
 * 检查作者类型是否互斥
 */
export const AUTHOR_TYPE_CONFLICTS: Record<AuthorType, AuthorType[]> = {
  first_author: ['co_first_author', 'corresponding_author'],
  co_first_author: ['first_author'],
  corresponding_author: ['first_author'],
  author: [],
};
