package errors

// ErrorCode 错误码结构
type ErrorCode struct {
	Code int
	Msg  string
}

// 通用错误
var (
	Success      = ErrorCode{Code: 0, Msg: "成功"}
	SystemError  = ErrorCode{Code: 100001, Msg: "系统异常"}
	ParamError   = ErrorCode{Code: 100002, Msg: "参数错误"}
	Unauthorized = ErrorCode{Code: 100003, Msg: "未授权"}
	Forbidden    = ErrorCode{Code: 100004, Msg: "无权限"}
	NotFound     = ErrorCode{Code: 100005, Msg: "资源不存在"}
)

// 认证错误
var (
	LoginFailed    = ErrorCode{Code: 101001, Msg: "登录失败"}
	AccountLocked  = ErrorCode{Code: 101002, Msg: "账户锁定"}
	AccountBanned  = ErrorCode{Code: 101003, Msg: "账户禁用"}
	SessionExpired = ErrorCode{Code: 101004, Msg: "会话过期"}
)

// 论文管理错误
var (
	ErrPaperNotFound        = ErrorCode{Code: 102001, Msg: "论文不存在"}
	ErrPaperDuplicate       = ErrorCode{Code: 102002, Msg: "论文重复"}
	ErrPaperNotAllowModify  = ErrorCode{Code: 102003, Msg: "论文不允许修改"}
	ErrPaperStatusInvalid   = ErrorCode{Code: 102004, Msg: "论文状态不允许操作"}
	ErrPaperIncomplete      = ErrorCode{Code: 102005, Msg: "论文信息不完整"}
	ErrJournalNotFound      = ErrorCode{Code: 102006, Msg: "期刊不存在"}
	ErrAuthorInvalid        = ErrorCode{Code: 102007, Msg: "作者信息无效"}
	ErrProjectNotFound      = ErrorCode{Code: 102008, Msg: "课题不存在"}
	ErrProjectLinked        = ErrorCode{Code: 102009, Msg: "课题已关联论文,不允许删除"}
	ErrDataValidationFailed = ErrorCode{Code: 102010, Msg: "数据校验失败"}
)

// 审核管理错误
var (
	ErrNoReviewPermission          = ErrorCode{Code: 103001, Msg: "无审核权限"}
	ErrAlreadyReviewed             = ErrorCode{Code: 103002, Msg: "已审核"}
	ErrPaperStatusNotAllowedReview = ErrorCode{Code: 103003, Msg: "论文状态不允许审核"}
	ErrReviewCommentRequired       = ErrorCode{Code: 103004, Msg: "驳回时必须填写审核意见"}
	ErrReviewFailed                = ErrorCode{Code: 103005, Msg: "审核失败"}
	ErrArchiveFailed               = ErrorCode{Code: 103006, Msg: "归档失败"}
	ErrAlreadyArchived             = ErrorCode{Code: 103007, Msg: "已归档"}
	ErrReviewNotFound              = ErrorCode{Code: 103008, Msg: "审核记录不存在"}
	ErrCannotArchiveUnapproved     = ErrorCode{Code: 103009, Msg: "审核未通过的论文不能归档"}
	ErrReviewerNotFound            = ErrorCode{Code: 103010, Msg: "审核人信息不存在"}
)

// 文件管理错误
var (
	ErrFileTooLarge       = ErrorCode{Code: 104001, Msg: "文件过大,最大支持100MB"}
	ErrInvalidFileType    = ErrorCode{Code: 104002, Msg: "文件格式错误,仅支持PDF、JPG、PNG"}
	ErrUploadFailed       = ErrorCode{Code: 104003, Msg: "文件上传失败"}
	ErrFileNotFound       = ErrorCode{Code: 104004, Msg: "文件不存在"}
	ErrFileAlreadyExists  = ErrorCode{Code: 104005, Msg: "文件已存在"}
	ErrDeleteFileFailed   = ErrorCode{Code: 104006, Msg: "删除文件失败"}
	ErrInvalidAttachment  = ErrorCode{Code: 104007, Msg: "附件信息无效"}
	ErrAttachmentNotFound = ErrorCode{Code: 104008, Msg: "附件不存在"}
	ErrFileTypeMismatch   = ErrorCode{Code: 104009, Msg: "文件类型不匹配"}
	ErrFileAccessDenied   = ErrorCode{Code: 104010, Msg: "文件访问被拒绝"}
)
