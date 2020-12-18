package moduleInterfaceSerializer

import (
	"bigdata_permission/serializer"
)

type CreateModuleInterface struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
}

//deprecated: no need to update
type UpdateModuleInterface struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
}

type DeleteModuleInterface struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
}

type SearchModuleInterface struct {
	ModuleId int `json:"module_id" form:"module_id"`
	InterfaceId int `json:"interface_id" form:"interface_id"`

	serializer.Page
}

type GetModuleInterface struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
}

type GetModuleChildrenInterface struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}