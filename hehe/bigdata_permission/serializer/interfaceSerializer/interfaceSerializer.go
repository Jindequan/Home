package interfaceSerializer

import "bigdata_permission/serializer"

type CreateInterface struct {
	Type int8 `json:"type" form:"type" binding:"required"`
	ShowName string `json:"show_name" form:"show_name" binding:"required"`
	Path string `json:"path" form:"path" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

type UpdateInterface struct {
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
	Type int8 `json:"type" form:"type" binding:"required"`
	ShowName string `json:"show_name" form:"show_name" binding:"required"`
	Path string `json:"path" form:"path" binding:"required"`
	Remark string `json:"remark" form:"remark"`
}

type DeleteInterface struct {
	InterfaceId int `json:"interface_id" form:"interface_id" binding:"required"`
}

type SearchInterface struct {
	InterfaceId int `json:"interface_id" form:"interface_id"`
	Type int8 `json:"type" form:"type"`
	ShowName string `json:"show_name" form:"show_name"`
	Path string `json:"path" form:"path"`

	serializer.Page
}

type GetInterface struct {
	InterfaceId int `json:"interface_id" form:"interface_id"`
}