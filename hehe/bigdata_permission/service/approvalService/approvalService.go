package approvalService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/approvalSerializer"
	"bigdata_permission/serializer/userModuleSerializer"
	"bigdata_permission/service/common"
	"bigdata_permission/service/externalService"
	"bigdata_permission/service/userModuleService"
	"bigdata_permission/service/userService"
	"encoding/json"
	"sort"
	"strconv"
)

type ApprovalInfo struct {
	ApprovalId int `json:"approval_id"`
	RoleId int `json:"role_id"`
	ModuleIds []int `json:"module_ids"`
	WxApprovalId string `json:"wx_approval_id"`
	Status int8 `json:"status"`
	ApplyUid int `json:"apply_uid"`
	ApprovalUid []int `json:"approval_uid"`
	Reason string `json:"reason"`
	ApprovalReason string `json:"approval_reason"`
	ApprovalTime int64 `json:"approval_time"`

	RoleName string `json:"role_name"`
	ModuleList []string `json:"module_list"`
	StatusStr string `json:"status_str"`
	ApplyUserName string `json:"apply_user_name"`
	ApprovalUserName []string `json:"approval_user_name"`
	ApprovalTimeStr string `json:"approval_time_str"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type ApprovalList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]ApprovalInfo `json:"list"`
}

func SubmitApply(param *approvalSerializer.CreateApproval) (ecode.Code, string) {
	//check param
	if param.RoleId == 0 && len(param.ModuleIds) == 0 {
		return ecode.INVALID_PARAM, ""
	}
	if param.ApplyUid == 0 {
		return ecode.INVALID_PARAM, ""
	}
	//insert apply record
	moduleIdsMap := map[int]int{}
	if param.RoleId > 0 {
		_, modulesByRole := dao.GetRoleModuleByRoleId(param.RoleId)
		for _, item := range *modulesByRole {
			moduleIdsMap[item.ModuleId] = item.ModuleId
		}
	}
	if len(param.ModuleIds) > 0 {
		for _, moduleId := range param.ModuleIds {
			moduleIdsMap[moduleId] = moduleId
		}
	}
	if len(moduleIdsMap) == 0 {
		return ecode.INVALID_PARAM, "没有要申请的模块"
	}
	moduleIds := []int{}
	for _, id := range moduleIdsMap {
		moduleIds = append(moduleIds, id)
	}
	sort.Ints(moduleIds)
	insertInfo := &dao.Approval{
		RoleId: param.RoleId,
		ModuleIds: pkg.JoinIntArrToString(moduleIds, "|"),
		Status: dao.ApprovalStatusInit,
		Reason: param.Reason,
		ApplyUid: param.ApplyUid,
	}
	createRes := insertInfo.Create()
	if !createRes {
		return ecode.CreateErr, ""
	}

	//submit weChat approval
	_, modules := dao.GetModuleDicByIds(moduleIds)
	code, wxId := ApplyModuleApprovalByWeChat(insertInfo, modules)
	//update apply column
	if code != ecode.OK || wxId == "" {
		return ecode.INTERNAL_ERR, ""
	}
	updateInfo := map[string]interface{}{}
	updateInfo["wx_approval_id"] = wxId

	insertInfo.Update(updateInfo)
	return ecode.OK, ""
}

func SearchApproval(param *approvalSerializer.SearchApproval) (ecode.Code, string, *ApprovalList) {
	where := &dao.Approval{}
	if param.ApprovalId > 0 {
		where.ApprovalId = param.ApprovalId
	}
	if param.Status > 0 {
		where.Status = param.Status
	}
	if param.ApplyUid > 0 {
		where.ApplyUid = param.ApplyUid
	}
	if param.WxApprovalId != "" {
		where.WxApprovalId = param.WxApprovalId
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountApprovalByWhere(where)
	result := &ApprovalList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]ApprovalInfo{},
	}

	offset, limit := common.GetOffsetLimit(page)
	if total == 0 || total <= offset{
		return ecode.OK, "", result
	}

	_, list := dao.GetApprovalByWhere(where, offset, limit)
	newList := []ApprovalInfo{}
	for _, interfaceInfo := range list {
		newList = append(newList, TransferApprovalInfo(interfaceInfo))
	}
	result.List = &newList

	return ecode.OK, "", result
}

func GetApprovalInfo(param *approvalSerializer.ApprovalInfo) (ecode.Code, *ApprovalInfo) {
	ok, approvalDao := dao.GetApprovalByApprovalId(param.ApprovalId)
	if !ok {
		return ecode.MYSQL_ERR, &ApprovalInfo{}
	}
	info := TransferApprovalInfo(*approvalDao)
	return ecode.OK, &info
}

func GetAccessedStatus() int8 {
	return dao.ApprovalStatusAccessed
}

func GetRefusedStatus() int8 {
	return dao.ApprovalStatusReject
}

func BatchApproval(param *approvalSerializer.BatchApproval) (ecode.Code, []int) {
	if param.ApprovalUid == "" || len(param.ApprovalId) == 0 {
		return ecode.INVALID_PARAM, []int{}
	}
	if _, ok := dao.ApprovalStatusMap[param.Status]; !ok {
		return ecode.INVALID_PARAM, []int{}
	}
	errorList := []int{}
	for _, approvalId := range param.ApprovalId {
		//check base info
		ok, approvalInfo := dao.GetApprovalByApprovalId(approvalId)
		if !ok || approvalInfo.Status == dao.ApprovalStatusAccessed || approvalInfo.Status == dao.ApprovalStatusReject {
			errorList = append(errorList, approvalId)
			continue
		}
		//update approval
		updateInfo := map[string]interface{}{}
		updateInfo["status"] = param.Status
		updateInfo["approval_reason"] = param.Reason
		updateInfo["approval_uid"] = param.ApprovalUid
		updateInfo["approval_time"] = param.ApprovalTime
		isSuccess, _ := approvalInfo.Update(updateInfo)
		if !isSuccess {
			errorList = append(errorList, approvalId)
			continue
		}
		//append user module
		moduleIds := pkg.ToMultiIntArr(approvalInfo.ModuleIds, "|")
		appendUM := &userModuleSerializer.CreateUserModule{
			UserIds: []int{approvalInfo.ApplyUid},
			ModuleIds: moduleIds,
		}
		userModuleService.CreateUserModuleBatch(appendUM)
		//notify apply user
		go ApprovalNotify(approvalInfo)
	}
	code := ecode.OK
	if len(errorList) > 0 {
		paramStr, _ := json.Marshal(param)
		errStr, _ := json.Marshal(errorList)
		pkg.SendToRobot(pkg.ROBOT_ERROR, "审核处理失败：入参" + string(paramStr) + ";错误的approval_id:" + string(errStr))
		code = ecode.ApprovalFailed
	}
	return code, errorList
}

//通知申请人审批的结果:通过发送企业微信消息
func ApprovalNotify(approvalInfo *dao.Approval) {
	msg := "你申请的大数据权限："
	if approvalInfo.Reason != "" {
		msg += approvalInfo.Reason + ";"
	}
	switch approvalInfo.Status {
	case dao.ApprovalStatusAccessed:
		msg += "审核通过"
		break
	case dao.ApprovalStatusReject:
		msg += "审核拒绝"
		break
	default:
		break
	}
	if approvalInfo.ApprovalReason != "" {
		msg += ";备注：" + approvalInfo.ApprovalReason
	}

	msgParam := &externalService.SendWeChatMessageRequest{
		ToUIDs: []string{strconv.Itoa(approvalInfo.ApplyUid)},
		MsgType: externalService.MSG_TYPE_TEXT,
		Text: msg,
	}
	externalService.SendWeChatMessage(msgParam)
}

func TransferApprovalInfo(data dao.Approval) ApprovalInfo {
	times := common.BuildDaoTimeStr(data.CreatedAt, data.UpdatedAt, 0)

	moduleList := []string{}
	moduleIds := pkg.ToMultiIntArr(data.ModuleIds, "|")
	for _, moduleId := range moduleIds {
		info := ""
		moduleInfo := dictionary.ModuleIdDaoMap[moduleId]
		if moduleInfo.FirstLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.FirstLvlId, dictionary.ModuleList) + "/"
		}
		if moduleInfo.SecondLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.SecondLvlId, dictionary.ModuleList) + "/"
		}
		if moduleInfo.ThirdLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.ThirdLvlId, dictionary.ModuleList) + "/"
		}
		info += dictionary.GetDictValueOfIntKey(moduleInfo.ModuleId, dictionary.ModuleList)
		moduleList = append(moduleList, info)
	}

	approvalUIds := pkg.ToMultiIntArr(data.ApprovalUid, "|")
	approvalUserNames := []string{}
	for _, uid := range approvalUIds {
		userDetail := userService.GetUserDetailById(uid)
		approvalUserNames = append(approvalUserNames, userDetail.SSOUserInfo.Name)
	}

	applyUserDetail := userService.GetUserDetailById(data.ApplyUid)

	roleName := ""
	if data.RoleId > 0 {
		roleName = dictionary.GetStringValue(dictionary.RoleList, data.RoleId)
	}

	return ApprovalInfo {
		ApprovalId:     data.ApprovalId,
		RoleId:         data.RoleId,
		ModuleIds:      pkg.ToMultiIntArr(data.ModuleIds, "|"),
		WxApprovalId:   data.WxApprovalId,
		Status:         data.Status,
		ApplyUid:       data.ApplyUid,
		ApprovalUid:    approvalUIds,
		Reason:         data.Reason,
		ApprovalReason: data.ApprovalReason,
		ApprovalTime:   data.ApprovalTime,

		RoleName: roleName,
		ModuleList: moduleList,
		StatusStr: dao.ApprovalStatusMap[data.Status],
		ApplyUserName: applyUserDetail.SSOUserInfo.Name,
		ApprovalUserName: approvalUserNames,
		ApprovalTimeStr: pkg.TimeToString(data.ApprovalTime),
		CreatedAt: times.CreatedAtStr,
		UpdatedAt: times.UpdatedAtStr,
		DeletedAt: times.DeletedAtStr,
	}
}



