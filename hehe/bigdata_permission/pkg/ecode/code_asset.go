package ecode

//asset错误码格式，五位数，20开头。eg：20xxx

var (
	OrderExistErr         = New(20001) //工单不存在
	NoStaffAvailableErr   = New(20002) //无人员可分配
	ServerNotExistErr     = New(20003) //服务器不存在
	DomainNotExistErr     = New(20004) //域名不存在
	RoleNotExistErr       = New(20005) //角色不存在
	PermissionNotExistErr = New(20006) //权限不存在
	NotInApprovalStageErr = New(20007) //非法操作，该申请不在指定的审批阶段
)

var AssetCodeMsg = map[int]string{
	OrderExistErr.Code():         "工单不存在",
	NoStaffAvailableErr.Code():   "无人员可分配",
	ServerNotExistErr.Code():     "服务器不存在",
	DomainNotExistErr.Code():     "域名不存在",
	RoleNotExistErr.Code():       "角色不存在",
	PermissionNotExistErr.Code(): "权限不存在",
	NotInApprovalStageErr.Code(): "非法操作，该申请不在指定的审批阶段",
}
