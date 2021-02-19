package warningJob

import (
	"encoding/json"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
	"time"
	"web_leads_backend/conf"
	"web_leads_backend/dao"
	"web_leads_backend/pkg"
	"web_leads_backend/pkg/dictionary"
	"web_leads_backend/pkg/ecode"
	"web_leads_backend/service/externalApiService"
	"web_leads_backend/service/leads"
	"web_leads_backend/service/user"
)

/**
 * 通过dispatch拿到励销的分配
 */
func (job *LixiaoUpdateJob) Run(task dao.WarningTask) (bool, string) {
	assignUids := []uint{}
	err := json.Unmarshal([]byte(task.AssignUids), &assignUids)
	if err != nil || len(assignUids) < 1 {
		return true, "无可通知的人员"
	}
	//两周未更新跟进记录|check exist
	timeStamp := (time.Now().Unix() - 14*24*3600) * 1000
	exist := dao.FindDispatchForLiXiaoUpdateJob(conf.BaseConf.SelfWorkspaceId, []uint{dao.SYSTEM_ID_LIXIAO_MIANXIAO}, timeStamp, 0, 1)
	if len(exist) == 0 {
		return true, "无需要处理的线索"
	}
	//refresh dispatch and record
	offset, limit := uint(0), uint(1000)
	//check again
	exist = dao.FindDispatchForLiXiaoUpdateJob(conf.BaseConf.SelfWorkspaceId, []uint{dao.SYSTEM_ID_LIXIAO_MIANXIAO}, timeStamp, 0, 1)
	if len(exist) == 0 {
		return true, "无需要处理的线索"
	}

	//clean file before created
	fileNameLocal := getFilePath("lixiao_update_warning.xlsx")
	if !checkAndCleanFile(fileNameLocal) {
		return false, "excel历史文件已存在，且删除文件失败"
	}
	//初始化excel文件与工作区
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("励销逾期未更新线索")
	if err != nil {
		return false, "初始化excel文件失败"
	}
	sheet.SetColWidth(1, 6, 15)
	sheet.SetColWidth(8, 8, 15)
	//写入头部行
	pkg.WriteXlsxByRow(sheet, &[]string{
		"励销线索id",   //0
		"线索名称",     //1
		"线索状态",     //2
		"线索创建时间",   //3
		"线索创建人",    //4
		"公司名称",     //5
		"当前系统",     //6
		"当前负责人",    //7
		"上次操作时间",   //8
		"未更新时长（天）", //9
		"LC线索id",   //10
		"备注",       //11
	})

	//no updated over 2 weeks
	offset, limit = uint(0), uint(1000)
	count, insertNum := 0, 0
	for {
		list := dao.FindDispatchForLiXiaoUpdateJob(conf.BaseConf.SelfWorkspaceId, []uint{dao.SYSTEM_ID_LIXIAO_MIANXIAO}, timeStamp, offset, limit)
		if len(list) == 0 {
			break
		}
		leadsArr := getLeadsByDispatch(list)
		count += len(list)
		for _, dispatchInfo := range list {
			lastRecord, isOk := dao.GetLastRecordListByLeadsId(dispatchInfo.WorkspaceId, dispatchInfo.LeadsId)
			if !isOk {
				continue
			}
			if !(lastRecord.FeedbackSubType == 71 || lastRecord.FeedbackSubType == 83 || lastRecord.FeedbackSubType == 0) {
				continue
			}
			leadsInfo, ok := leadsArr[dispatchInfo.LeadsId]
			insert := []string{}
			lastRecordTime := "无"
			gapDay := "无"
			day := 0
			if dispatchInfo.LastRecordTime > 0 {
				gap := pkg.NowUnixMs() - dispatchInfo.LastRecordTime
				day = int(gap / 1000 / 3600 / 24)
				gapDay = strconv.Itoa(day)
				lastRecordTime = pkg.TimeToDateString(dispatchInfo.LastRecordTime)
			}
			if !ok {
				insert = []string{
					dispatchInfo.OutLeadsId, //10
					"",                      //1
					"",                      //2
					"",                      //3
					"",                      //4
					"",                      //5
					dictionary.GetStringValue(dictionary.SystemDic, dispatchInfo.SystemId), //6
					dispatchInfo.CurrentFollower,            //7
					lastRecordTime,                          //8
					gapDay,                                  //9
					strconv.Itoa(int(dispatchInfo.LeadsId)), //10
					"查询不到线索信息，请检查是否被删除", //11
				}
			} else {
				creatorName := ""
				if leadsInfo.CreatorId > 0 {
					userListMap := user.GetUserListMap()
					if userInfo, exist := userListMap[leadsInfo.CreatorId]; exist {
						creatorName = userInfo.Name
					}
				}

				insert = []string{
					dispatchInfo.OutLeadsId, //0
					leadsInfo.Name,          //1
					leads.GetStatusStr(leadsInfo.Status, leadsInfo.SystemId), //2
					pkg.TimeToDateString(leadsInfo.CreateTime),               //3
					creatorName,       //4
					leadsInfo.Company, //5
					dictionary.GetStringValue(dictionary.SystemDic, dispatchInfo.SystemId), //6
					dispatchInfo.CurrentFollower,            //7
					lastRecordTime,                          //8
					gapDay,                                  //9
					strconv.Itoa(int(dispatchInfo.LeadsId)), //10
					"",                                      //11
				}
			}

			insertNum += 1
			pkg.WriteXlsxByRow(sheet, &insert)
		}
		if len(list) < int(limit) {
			break
		}
		offset += limit
	}
	//写入文件
	saveErr := file.Save(fileNameLocal)
	if saveErr != nil {
		return false, "共查询出" + strconv.Itoa(count) + "条数据；导出" + strconv.Itoa(insertNum) + "条有效数据。但是写入文件失败：" + saveErr.Error()
	}
	//清除文件
	defer cleanFile(fileNameLocal)
	//上传文件到sso
	fileName := "励销未处理工单-" + pkg.DateNow() + ".xlsx"
	uploadSuccess, msg, fileResp := externalApiService.UploadFileSso(fileName, fileNameLocal)
	if !uploadSuccess {
		fileRespStr, _ := json.Marshal(fileResp)
		return false, "上传文件失败:" + msg + "::" + string(fileRespStr)
	}
	//发送邮件
	params := externalApiService.SsoEmailParams{
		Subject:       "【CRM】【" + strings.ToUpper(pkg.GetEnv()) + "】励销面销未更新预警",
		Content:       "励销跟进记录已超过14天未更新",
		FromName:      "线索管理系统预警中心",
		ToUids:        assignUids,
		AttachFileIds: []string{fileResp.FileId},
	}
	code, resp := externalApiService.SendSsoEmail(params)
	if code == ecode.OK {
		return true, ""
	}
	req, _ := json.Marshal(params)
	res, _ := json.Marshal(resp)
	return false, "发送预警邮件失败; request:" + string(req) + "; response:" + string(res)
}

func getLeadsByDispatch(dispatchList []*dao.Dispatch) map[uint]*dao.Leads {
	wlMapping := map[uint][]uint{} //workspace => leads_id

	for _, v := range dispatchList {
		_, ok := wlMapping[v.WorkspaceId]
		if !ok {
			wlMapping[v.WorkspaceId] = []uint{v.LeadsId}
			continue
		}

		wlMapping[v.WorkspaceId] = append(wlMapping[v.WorkspaceId], v.LeadsId)
	}

	leadsList := map[uint]*dao.Leads{}
	for workspaceId, leadsIds := range wlMapping {
		leadsArr, _ := dao.FindLeadsByIdArr(workspaceId, leadsIds)
		for _, daoLeads := range leadsArr {
			leadsList[daoLeads.LeadsId] = daoLeads
		}
	}

	return leadsList
}