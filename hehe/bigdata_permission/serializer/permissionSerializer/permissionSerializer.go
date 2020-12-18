package permissionSerializer

type CheckPermission struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	InterfaceType int8 `json:"interface_type" form:"interface_type"`
	Path string `json:"path" form:"path"`

	Uid int `json:"uid" form:"uid" binding:"required"`//与role_id二选一
	RoleId int `json:"role_id" form:"role_id"`
}
