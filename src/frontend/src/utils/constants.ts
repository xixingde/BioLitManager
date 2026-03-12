import type { Role, Permission } from '@/types/user';

export const API_BASE_URL = '/api';

export const TOKEN_KEY = 'biolit_token';

export const USER_INFO_KEY = 'biolit_user_info';

export const ROLE: Record<string, Role> = {
  SUPER_ADMIN: 'super_admin',
  ADMIN: 'admin',
  DEPT_HEAD: 'dept_head',
  PROJECT_LEADER: 'project_leader',
  BUSINESS_REVIEWER: 'business_reviewer',
  POLITICAL_REVIEWER: 'political_reviewer',
  USER: 'user',
} as const;

export const PERMISSION: Record<string, Permission> = {
  PAPER_CREATE: 'paper:create',
  PAPER_EDIT: 'paper:edit',
  PAPER_VIEW: 'paper:view',
  PAPER_DELETE: 'paper:delete',
  PAPER_EXPORT: 'paper:export',
  REVIEW_BUSINESS: 'review:business',
  REVIEW_POLITICAL: 'review:political',
  SYSTEM_USER_MANAGE: 'system:user:manage',
  SYSTEM_PROJECT_MANAGE: 'system:project:manage',
  SYSTEM_JOURNAL_MANAGE: 'system:journal:manage',
  SYSTEM_CONFIG_MANAGE: 'system:config:manage',
  STATS_VIEW: 'stats:view',
  STATS_EXPORT: 'stats:export',
} as const;
