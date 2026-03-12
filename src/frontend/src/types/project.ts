// 课题相关类型定义

/**
 * 项目类型枚举
 */
export type ProjectType = 'vertical' | 'horizontal';

/**
 * 项目级别枚举
 */
export type ProjectLevel = 'national' | 'provincial' | 'municipal';

/**
 * 课题接口
 */
export interface Project {
  id: number;
  name: string;
  code: string;
  project_type: ProjectType;
  source: string;
  level: ProjectLevel;
  status: string;
}

/**
 * 课题表单接口
 */
export interface ProjectForm {
  name: string;
  code: string;
  project_type: ProjectType;
  source: string;
  level: ProjectLevel;
}

/**
 * 项目类型选项
 */
export const PROJECT_TYPE_OPTIONS = [
  { label: '纵向项目', value: 'vertical' as ProjectType },
  { label: '横向项目', value: 'horizontal' as ProjectType },
];

/**
 * 项目级别选项
 */
export const PROJECT_LEVEL_OPTIONS = [
  { label: '国家级', value: 'national' as ProjectLevel },
  { label: '省部级', value: 'provincial' as ProjectLevel },
  { label: '市级', value: 'municipal' as ProjectLevel },
];
