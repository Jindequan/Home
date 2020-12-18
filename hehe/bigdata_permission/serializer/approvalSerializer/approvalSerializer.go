package approvalSerializer

import "bigdata_permission/serializer"

type CreateApproval struct {
	RoleId int `json:"role_id" form:"role_id"`
	ModuleIds []int `json:"module_ids" form:"module_ids"`
	Reason string `json:"reason" form:"reason"`

	ApplyUid int
}

type SearchApproval struct {
	ApprovalId int `json:"approval_id" form:"approval_id"`
	Status int8 `json:"status" form:"status"`
	ApplyUid int `json:"apply_uid" form:"apply_uid"`
	WxApprovalId string `json:"wx_approval_id" form:"wx_approval_id"`

	serializer.Page
}

type ApprovalInfo struct {
	ApprovalId int `json:"approval_id" form:"approval_id" binding:"required"`
}

type WxApprovalInfo struct {
	WxApprovalId string `json:"wx_approval_id" form:"wx_approval_id" binding:"required"`
}

type BatchApproval struct {
	ApprovalId []int `json:"approval_id" form:"approval_id" binding:"required"`
	Status int8 `json:"status" form:"status"`
	Reason string `json:"reason" form:"reason"`

	ApprovalUid string `json:"approval_uid" form:"approval_uid"`
	ApprovalTime int64 `json:"approval_time" form:"approval_time"`
}