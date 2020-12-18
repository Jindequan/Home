package roleModuleSerializer

import (
	"bigdata_permission/serializer"
)

type CreateRoleModule struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}

type CreateRoleModuleBatch struct {
	RoleIds []int `json:"role_ids" form:"role_ids" binding:"required"`
	ModuleIds []int `json:"module_ids" form:"module_ids" binding:"required"`
}

//deprecated: no need to update
type UpdateRoleModule struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}

type DeleteRoleModule struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}

type SearchRoleModule struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	ModuleId int `json:"module_id" form:"module_id"`

	serializer.Page
}

type GetRoleModule struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}