export type Role =
  | 'super_admin'
  | 'admin'
  | 'dept_head'
  | 'project_leader'
  | 'business_reviewer'
  | 'political_reviewer'
  | 'user';

export type Permission =
  | 'paper:create'
  | 'paper:edit'
  | 'paper:view'
  | 'paper:delete'
  | 'paper:export'
  | 'review:business'
  | 'review:political'
  | 'system:user:manage'
  | 'system:project:manage'
  | 'system:journal:manage'
  | 'system:config:manage'
  | 'stats:view'
  | 'stats:export';

export interface UserInfo {
  id: string;
  username: string;
  name: string;
  role: Role;
  department: string;
  email?: string;
  id_card?: string;
  phone?: string;
  is_disabled: boolean;
  is_locked: boolean;
}
