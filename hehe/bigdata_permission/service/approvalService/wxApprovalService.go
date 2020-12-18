package approvalService

import (
	"bigdata_permission/conf"
	"bigdata_permission/dao"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/wechat"
	"bigdata_permission/serializer/approvalSerializer"
	"bigdata_permission/service/userService"
	"encoding/json"
	"strconv"
)

var wxContentIdMap = map[string]string{
	"apply_modules": "Textarea-1607482913017",
	"apply_reason":  "Text-1607504584106",
	"apply_user":    "Contact-1607504642320",
}

func ApplyModuleApprovalByWeChat(approval *dao.Approval, modules *[]dao.ModuleDic) (ecode.Code, string) {
	//处理审批内容
	moduleContent := []string{}
	for _, moduleInfo := range *modules {
		info := ""
		if moduleInfo.FirstLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.FirstLvlId, dictionary.ModuleList) + "/"
		}
		if moduleInfo.SecondLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.SecondLvlId, dictionary.ModuleList) + "/"
		}
		if moduleInfo.ThirdLvlId > 0 {
			info += dictionary.GetDictValueOfIntKey(moduleInfo.ThirdLvlId, dictionary.ModuleList) + "/"
		}
		info += dictionary.GetDictValueOfIntKey(moduleInfo.ModuleId, dictionary.ModuleList) + "\n"
		moduleContent = append(moduleContent, info)
	}
	params := []*wechat.ContentItem{}
	params = append(params, &wechat.ContentItem{
		Control: "text",
		Id: wxContentIdMap["apply_modules"],
		Value: wechat.ContentValue{
			Text: pkg.JoinStrArrToString(moduleContent, ""),
		},
	})
	params = append(params, &wechat.ContentItem{
		Control: "text",
		Id: wxContentIdMap["apply_reason"],
		Value: wechat.ContentValue{
			Text: approval.Reason,
		},
	})
	summaryInfo := []*wechat.SummaryInfo{}
	summaryInfo = append(summaryInfo, &wechat.SummaryInfo{
		[]wechat.SummaryItem{
			{
				Text: "共" + strconv.Itoa(len(moduleContent)) + "条权限申请",
				Lang: "zh_CN",
			},
		},
	})
	summaryInfo = append(summaryInfo, &wechat.SummaryInfo{
		[]wechat.SummaryItem{
			{
				Text: "原因：" + approval.Reason,
				Lang: "zh_CN",
			},
		},
	})

	code, userInfo := userService.GetSsoUserInfoWithWeChatIdByUIds([]int{approval.ApplyUid})
	if code != ecode.OK {
		pkg.SendToRobot(pkg.ROBOT_ERROR, code.Message() + "sso获取用户微信id失败：" + strconv.Itoa(approval.ApplyUid))
		return code, ""
	}
	wxId := userInfo[approval.ApplyUid].WechatId
	wxName := userInfo[approval.ApplyUid].Name
	params = append(params, &wechat.ContentItem{
		Control: "Contact",
		Id: wxContentIdMap["apply_user"],
		Value: wechat.ContentValue{
			Members: []*wechat.MemberItem{
				{
					Userid: wxId,
					Name:   wxName,
				},
			},
		},
	})

	useTemplateApprover := 1
	approverList := []*wechat.ApproverItem{}
	if !pkg.IsProduction() {
		wxId = "JinDeQuan"
		useTemplateApprover = 0
		approverList = []*wechat.ApproverItem{
			{
				Attr:   1,
				Userid: []string{"JinDeQuan"},
			},
		}
	}
	//拼接微信请求体
	wxApprovalRequestParam := &wechat.ApprovalRequestParam{
		CreatorUserid:       wxId,
		TemplateId:          conf.BaseConf.WeChat.ApprovalTemplateId.ModuleApprovalTemplateId,
		UseTemplateApprover: useTemplateApprover,
		Approver:            approverList,
		Notifyer:            []string{wxId},
		NotifyType:          1,
		ApplyData: wechat.ApplyDataDetail{
			Contents: params,
		},
		SummaryList: summaryInfo,
	}

	weChatApp := wechat.GetWorkWxApp(wechat.WorkWxAppApprovalModules)
	wxApprovalId, err := weChatApp.CreateApproval(wxApprovalRequestParam)
	if err != nil {
		pkg.SendToRobot(pkg.ROBOT_ERROR, strconv.Itoa(approval.ApprovalId) + "：发起微信审批流程失败：" + err.Error())
		return ecode.EXTERNAL_API_FAIL, ""
	}

	return ecode.OK, wxApprovalId
}

func WxApprovalCallback(param *wechat.ApprovalCallbackMsgContent) (ecode.Code, string) {
	isExist, approvalInfo := dao.GetApprovalByWxId(strconv.Itoa(int(param.ApprovalInfo.SpNo)))
	if !isExist  {
		return ecode.NoApprovalFound, "错误的微信审批号"
	}
	if approvalInfo.Status != dao.ApprovalStatusInit {
		return ecode.ApprovalStatusDone, ""
	}
	wechatApp := wechat.GetWorkWxApp(wechat.WorkWxAppApprovalModules)
	wxId, _ := strconv.Atoi(approvalInfo.WxApprovalId)
	approvalDetail, err := wechatApp.GetApprovalDetail(int64(wxId))
	if err != nil || approvalDetail.SpNo != approvalInfo.WxApprovalId {
		return ecode.EXTERNAL_API_FAIL, "微信审批详情获取失败"
	}

	//查询审批人UID
	needQueryUserWxId := []string{}
	for _, spRecordItem := range approvalDetail.SpRecord {
		for _, detailItem := range spRecordItem.Details {
			if detailItem.Approver.UserId != "" {
				needQueryUserWxId = append(needQueryUserWxId, detailItem.Approver.UserId)
			}
		}
	}

	//查询微信id对应的uid
	wxIdMap, uidArr := map[string]int{}, []int{}
	if len(needQueryUserWxId) > 0 {
		code, uidInfoArr := userService.GetUserUidInfoByWxId(needQueryUserWxId)
		if code != ecode.OK && len(uidInfoArr) > 0 {
			for _, uidInfo := range uidInfoArr {
				wxIdMap[uidInfo.WeChatUid] = uidInfo.Uid
				uidArr = append(uidArr, uidInfo.Uid)
			}
		}
	}
	status := int8(0)
	switch approvalDetail.SpStatus {
	case wechat.WX_APPROVAL_STATUS_PASS:
		status = dao.ApprovalStatusAccessed
		break
	case wechat.WX_APPROVAL_STATUS_REJECT:
		status = dao.ApprovalStatusReject
		break
	default:
		break
	}
	//回调成功，仍需等待下一次确切的通过或拒绝
	if status == 0 {
		return ecode.OK, ""
	}

	//更新审批信息
	operateParam := &approvalSerializer.BatchApproval{
		Status: status,
		ApprovalId:  []int{approvalInfo.ApprovalId},
		Reason:      "微信审批",
		ApprovalUid: pkg.JoinIntArrToString(uidArr, "|"),
		ApprovalTime: pkg.NowUnixMs(),
	}
	BatchApproval(operateParam)
	
	return ecode.OK, ""
}

func GetWxApprovalDetail(wxId string) string {
	wechatApp := wechat.GetWorkWxApp(wechat.WorkWxAppApprovalModules)
	wxIdInt, _ := strconv.Atoi(wxId)
	approvalDetail, err := wechatApp.GetApprovalDetail(int64(wxIdInt))
	if err != nil {

	}
	//todo transfer to visible struct
	detail, _ := json.Marshal(approvalDetail)
	return string(detail)
}