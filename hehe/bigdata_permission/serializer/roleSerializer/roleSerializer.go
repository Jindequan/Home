package roleSerializer

import (
	"bigdata_permission/serializer"
)

type CreateRole struct {
	Name string `json:"name" form:"name" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

type UpdateRole struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
	Name string `json:"name" form:"name"`
	Remark string `json:"remark" form:"remark"`
}

type DeleteRole struct {
	RoleId int `json:"role_id" form:"role_id" binding:"required"`
}

type SearchRole struct {
	RoleId int `json:"role_id" form:"role_id"`
	Name string `json:"name" form:"name"`

	serializer.Page
}

type GetRole struct {
	RoleId int `json:"role_id" form:"role_id"`
}