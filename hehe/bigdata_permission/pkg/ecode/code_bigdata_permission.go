package ecode

//错误码格式，五位数，128开头。eg：12800

var (
	//接口处理错误-基础错误
	CreateErr         = New(12801001)
	CreateErrExist    = New(12801002)
	UpdateErr         = New(12801003)
	UpdateErrNotFound = New(12801004)
	DeleteErr         = New(12801005)
	DeleteErrNotFound = New(12801006)
	GetListErr        = New(12801007)
	GetInfoErr        = New(12801008)
	//功能异常
	FunctionUnAvailable = New(12802001)
	//模块-接口
	NotLeafModuleCannotRelateInterface = New(12803001)
	//用户权限
	NoUserFound = New(12804001)
	NoUserRole  = New(12804002)
	//审批
	NoApprovalFound = New(12805001)
	ApprovalStatusDone = New(12805002)
	ApprovalFailed = New(12805003)
)

var BigDataPermissionCodeMsg = map[int]string{
	CreateErr.Code():                          "新增数据失败",
	CreateErrExist.Code():                     "新增数据失败:数据已存在",
	UpdateErr.Code():                          "更新数据失败",
	UpdateErrNotFound.Code():                  "更新数据失败:查询不到所要更新的数据",
	DeleteErr.Code():                          "删除数据失败",
	DeleteErrNotFound.Code():                  "删除数据失败:查询不到所要删除的数据",
	GetListErr.Code():                         "获取列表失败",
	GetInfoErr.Code():                         "获取详情失败",
	FunctionUnAvailable.Code():                "该功能不可用",
	NotLeafModuleCannotRelateInterface.Code(): "该模块拥有子模块，无法与接口关联",
	NoUserFound.Code():                        "不存在的用户",
	NoUserRole.Code():                         "未知的用户角色",
	NoApprovalFound.Code(): 				   "不存在的审批申请",
	ApprovalStatusDone.Code(): 				   "审批已完结",
	ApprovalFailed.Code(): 					   "审批失败",
}
