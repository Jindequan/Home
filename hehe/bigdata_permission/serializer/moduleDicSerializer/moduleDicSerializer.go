package moduleDicSerializer

import (
	"bigdata_permission/serializer"
)

type CreateModuleDic struct {
	Value string `json:"value" form:"value" binding:"required"`
	ParentId int `json:"parent_id" form:"parent_id"`
	Enable int8 `json:"enable" form:"enable"`
	Disable bool `json:"disable" form:"disable"`
	ResponsePerson int `json:"response_person" form:"response_person"`
}

type UpdateModuleDic struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
	Value string `json:"value" form:"value" binding:"required"`
	ParentId int `json:"parent_id" form:"parent_id"`
	Enable int8 `json:"enable" form:"enable"`
	Disable bool `json:"disable" form:"disable"`
	ResponsePerson int `json:"response_person" form:"response_person"`
}

type DeleteModuleDic struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}

type SearchModuleDic struct {
	ModuleId int `json:"module_id" form:"module_id"`
	Enable bool `json:"enable" form:"enable"`
	Disable bool `json:"disable" form:"disable"`
	ResponsePerson int `json:"response_person" form:"response_person"`
	ResponsePersonNull bool `json:"response_person_null" form:"response_person_null"`

	serializer.Page
}

type GetModuleDic struct {
	ModuleId int `json:"module_id" form:"module_id" binding:"required"`
}