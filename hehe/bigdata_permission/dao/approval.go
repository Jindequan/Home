package dao

import "gorm.io/gorm"

type Approval struct {
	ApprovalId int `gorm:"primaryKey"`
	RoleId int
	ModuleIds string
	WxApprovalId string
	Status int8
	ApplyUid int
	ApprovalUid string
	Reason string
	ApprovalReason string
	ApprovalTime int64

	Base
}

const (
	ApprovalStatusInit     = 1 //审批中
	ApprovalStatusAccessed = 2 //已通过
	ApprovalStatusReject   = 3 //已驳回
)

var ApprovalStatusMap = map[int8]string{
	ApprovalStatusInit:     "审批中",
	ApprovalStatusAccessed: "已通过",
	ApprovalStatusReject:   "已驳回",
}

func (approval *Approval) Create () bool {
	res := DB.Create(approval)
	return res.Error == nil && res.RowsAffected == 1
}

func (approval *Approval) Update (data map[string]interface{}) (bool, *Approval) {
	res := DB.Model(approval).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, approval
}

func (approval *Approval) Delete () bool {
	res := DB.Delete(approval)
	return res.Error == nil && res.RowsAffected == 1
}

func GetApprovalByApprovalId (approvalId int) (bool, *Approval) {
	approval := &Approval{}
	res := DB.Model(approval).Where("approval_id = ?", approvalId).First(approval)
	return res.Error == nil, approval
}

func GetApprovalByWxId (wxApprovalId string) (bool, *Approval) {
	approval := &Approval{}
	res := DB.Model(approval).Where("wx_approval_id = ?", wxApprovalId).First(approval)
	return res.Error == nil, approval
}

func GetApprovalByWhere (where *Approval, offset, limit int) (bool, []Approval) {
	query := getApprovalQueryByWhere(where)
	list := []Approval{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountApprovalByWhere (where *Approval) (bool, int) {
	query := getApprovalQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil, int(total)
}

func getApprovalQueryByWhere(where *Approval) *gorm.DB {
	query := DB.Model(Approval{})
	if where.ApprovalId > 0 {
		query = query.Where("approval_id = ?", where.ApprovalId)
	}
	if where.ApplyUid > 0 {
		query = query.Where("apply_uid = ?", where.ApplyUid)
	}
	if where.Status > 0 {
		query = query.Where("status = ?", where.Status)
	}
	if where.WxApprovalId != "" {
		query = query.Where("wx_approval_uid = ?", where.WxApprovalId)
	}
	return query
}
