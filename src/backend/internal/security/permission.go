package security

// Role 角色类型
type Role string

const (
	// RoleSuperAdmin 超级管理员
	RoleSuperAdmin Role = "super_admin"
	// RoleAdmin 管理员
	RoleAdmin Role = "admin"
	// RoleDeptHead 部门主管
	RoleDeptHead Role = "dept_head"
	// RoleProjectLeader 项目负责人
	RoleProjectLeader Role = "project_leader"
	// RoleBusinessReviewer 业务审核人
	RoleBusinessReviewer Role = "business_reviewer"
	// RolePoliticalReview 政治审核人
	RolePoliticalReview Role = "political_reviewer"
	// RoleUser 普通用户
	RoleUser Role = "user"
)

// Permission 权限类型
type Permission string

const (
	// PermissionPaperCreate 创建论文
	PermissionPaperCreate Permission = "paper:create"
	// PermissionPaperEdit 编辑论文
	PermissionPaperEdit Permission = "paper:edit"
	// PermissionPaperView 查看论文
	PermissionPaperView Permission = "paper:view"
	// PermissionPaperDelete 删除论文
	PermissionPaperDelete Permission = "paper:delete"
	// PermissionPaperExport 导出论文
	PermissionPaperExport Permission = "paper:export"
	// PermissionReviewBusiness 业务审核
	PermissionReviewBusiness Permission = "review:business"
	// PermissionReviewPolitical 政治审核
	PermissionReviewPolitical Permission = "review:political"
	// PermissionSystemUserManage 用户管理
	PermissionSystemUserManage Permission = "system:user:manage"
	// PermissionSystemProjectManage 项目管理
	PermissionSystemProjectManage Permission = "system:project:manage"
	// PermissionSystemJournalManage 期刊管理
	PermissionSystemJournalManage Permission = "system:journal:manage"
	// PermissionSystemConfigManage 配置管理
	PermissionSystemConfigManage Permission = "system:config:manage"
	// PermissionStatsView 查看统计
	PermissionStatsView Permission = "stats:view"
	// PermissionStatsExport 导出统计
	PermissionStatsExport Permission = "stats:export"
)

// RolePermissions 角色到权限的映射
var RolePermissions = map[Role][]Permission{
	RoleSuperAdmin: {
		PermissionPaperCreate,
		PermissionPaperEdit,
		PermissionPaperView,
		PermissionPaperDelete,
		PermissionPaperExport,
		PermissionReviewBusiness,
		PermissionReviewPolitical,
		PermissionSystemUserManage,
		PermissionSystemProjectManage,
		PermissionSystemJournalManage,
		PermissionSystemConfigManage,
		PermissionStatsView,
		PermissionStatsExport,
	},
	RoleAdmin: {
		PermissionPaperCreate,
		PermissionPaperEdit,
		PermissionPaperView,
		PermissionPaperDelete,
		PermissionPaperExport,
		PermissionReviewBusiness,
		PermissionReviewPolitical,
		PermissionSystemUserManage,
		PermissionSystemProjectManage,
		PermissionSystemJournalManage,
		PermissionStatsView,
		PermissionStatsExport,
	},
	RoleDeptHead: {
		PermissionPaperCreate,
		PermissionPaperEdit,
		PermissionPaperView,
		PermissionPaperDelete,
		PermissionPaperExport,
		PermissionReviewBusiness,
		PermissionReviewPolitical,
		PermissionStatsView,
		PermissionStatsExport,
	},
	RoleProjectLeader: {
		PermissionPaperCreate,
		PermissionPaperEdit,
		PermissionPaperView,
		PermissionPaperDelete,
		PermissionPaperExport,
		PermissionReviewBusiness,
		PermissionReviewPolitical,
		PermissionStatsView,
		PermissionStatsExport,
	},
	RoleBusinessReviewer: {
		PermissionPaperView,
		PermissionPaperExport,
		PermissionReviewBusiness,
		PermissionStatsView,
	},
	RolePoliticalReview: {
		PermissionPaperView,
		PermissionPaperExport,
		PermissionReviewPolitical,
		PermissionStatsView,
	},
	RoleUser: {
		PermissionPaperCreate,
		PermissionPaperEdit,
		PermissionPaperView,
		PermissionPaperExport,
		PermissionStatsView,
	},
}

// GetPermissionsByRole 根据角色获取权限列表
func GetPermissionsByRole(role Role) []Permission {
	permissions, exists := RolePermissions[role]
	if !exists {
		return []Permission{}
	}
	return permissions
}

// HasPermission 检查角色是否具有指定权限
func HasPermission(role Role, permission Permission) bool {
	permissions := GetPermissionsByRole(role)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasAnyPermission 检查角色是否具有任一指定权限
func HasAnyPermission(role Role, permissions ...Permission) bool {
	for _, permission := range permissions {
		if HasPermission(role, permission) {
			return true
		}
	}
	return false
}

// HasAllPermissions 检查角色是否具有所有指定权限
func HasAllPermissions(role Role, permissions ...Permission) bool {
	for _, permission := range permissions {
		if !HasPermission(role, permission) {
			return false
		}
	}
	return true
}
