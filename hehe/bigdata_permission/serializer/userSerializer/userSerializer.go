package userSerializer

import "bigdata_permission/serializer"

type CreateUser struct {
	Uid int `json:"uid" form:"uid" binding:"required"`
	RoleId int `json:"role_id" form:"role_id"`
	Remark string `json:"remark" form:"remark"`
}

type UpdateUser struct {
	Uid int `json:"uid" form:"uid" binding:"required"`
	RoleId int `json:"role_id" form:"role_id"`
	Remark string `json:"remark" form:"remark"`
}

type DeleteUser struct {
	Uid int `json:"uid" form:"uid" binding:"required"`
}

type SearchUser struct {
	Uid int `json:"uid" form:"uid"`
	RoleId int `json:"role_id" form:"role_id"`

	serializer.Page
}

type GetUser struct {
	Uid int `json:"uid" form:"uid"`
}