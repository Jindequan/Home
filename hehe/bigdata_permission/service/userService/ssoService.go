package userService

import (
	"bigdata_permission/conf"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/logging"
	"bigdata_permission/pkg/requests"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

type SsoUserIds struct {
	Uids []int `json:"uids"`
}

func GetSsoUserInfoByUids(uids []int) (ecode.Code, []SSOUserInfo) {
	SsoDomain := conf.BaseConf.Domain.Sso
	code := ecode.OK
	var ssoUserInfoList []SSOUserInfo
	reqBody, _ := json.Marshal(struct {
		Uids []int `json:"uids"`
	}{
		Uids: uids,
	})

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	resp := requests.Post(
		SsoDomain+"/open-api/user/get-list-by-uids?platform_id="+
			conf.BaseConf.SsoConfig.PlatformID+
			"&timestamp="+timestampStr+
			"&signature="+getSignature(timestampStr),
		requests.JsonBodyBytes(reqBody),
	)

	if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
		code = ecode.EXTERNAL_API_FAIL
	} else {
		err := json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), &ssoUserInfoList)
		if err != nil {
			code = ecode.EXTERNAL_API_JSON_ERR
		}
	}

	return code, ssoUserInfoList
}

func GetSsoUserInfoByUid(uid int) (SSOUserInfo, ecode.Code) {
	SsoDomain := conf.BaseConf.Domain.Sso
	code := ecode.OK
	var ssoUserInfo SSOUserInfo
	reqBody, _ := json.Marshal(struct {
		Uid int `json:"uid"`
	}{
		Uid: uid,
	})

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	resp := requests.Post(
		SsoDomain+"/open-api/user/get-info-by-uid?platform_id="+
			conf.BaseConf.SsoConfig.PlatformID+
			"&timestamp="+timestampStr+
			"&signature="+getSignature(timestampStr),
		requests.JsonBodyBytes(reqBody),
	)

	if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
		code = ecode.EXTERNAL_API_FAIL
	} else {
		err := json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), &ssoUserInfo)
		if err != nil {
			code = ecode.EXTERNAL_API_JSON_ERR
		}
	}

	return ssoUserInfo, code
}

//func GetSsoUserNameByUidCache(uid int) string {
//	name := ""
//	//查看缓存
//	if value, exist := dictionary.SsoNameCacheMap[uid]; exist {
//		name = value
//	} else {
//		//没有则从sso获取
//		ssoUserInfo, code := GetSsoUserInfoByUid(uid)
//		if code == ecode.OK && ssoUserInfo.Name != "" {
//			name = ssoUserInfo.Name
//			//设置缓存
//			dictionary.SsoNameCacheMap[uid] = name
//		}
//	}
//
//	return name
//}

func GetSsoUserInfoByEmail(email string) (SSOUserInfo, ecode.Code) {
	SsoDomain := conf.BaseConf.Domain.Sso
	code := ecode.OK
	var ssoUserInfo SSOUserInfo
	reqBody, _ := json.Marshal(struct {
		Email string `json:"email"`
	}{
		Email: email,
	})

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	resp := requests.Post(
		SsoDomain+"/open-api/user/get-info-by-email?platform_id="+
			conf.BaseConf.SsoConfig.PlatformID+
			"&timestamp="+timestampStr+
			"&signature="+getSignature(timestampStr),
		requests.JsonBodyBytes(reqBody),
	)

	if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
		code = ecode.EXTERNAL_API_FAIL
	} else {
		err := json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), &ssoUserInfo)
		if err != nil {
			code = ecode.EXTERNAL_API_JSON_ERR
		}
	}

	return ssoUserInfo, code
}

func GetSsoUserIdsByName(name string) (*SsoUserIds, ecode.Code) {
	SsoDomain := conf.BaseConf.Domain.Sso
	code := ecode.OK

	ssoUserIds := &SsoUserIds{}

	reqBody, _ := json.Marshal(struct {
		Name string `json:"name"`
	}{
		Name: name,
	})

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	resp := requests.Post(
		SsoDomain+"/open-api/user/get-uidlist-by-name?platform_id="+
			conf.BaseConf.SsoConfig.PlatformID+
			"&timestamp="+timestampStr+
			"&signature="+getSignature(timestampStr),
		requests.JsonBodyBytes(reqBody),
	)

	if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
		code = ecode.EXTERNAL_API_FAIL
	} else {
		err := json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), ssoUserIds)
		if err != nil {
			code = ecode.EXTERNAL_API_JSON_ERR
		}
	}

	return ssoUserIds, code
}

func getSignature(timestampStr string) string {
	h := md5.New()
	h.Write([]byte(timestampStr + "_" + conf.BaseConf.SsoConfig.SecretKey))
	return hex.EncodeToString(h.Sum(nil))
}


const (
	Male   = 0 //男性为0
	Female = 1 //女性为1
)

type SSOUserInfo struct {
	Uid                      int  `json:"uid"`
	Email                    string `json:"email"`
	Name                     string `json:"name"`
	EnName                   string `json:"en_name"`
	Mobile                   string `json:"mobile"`
	Avatar                   string `json:"avatar"`
	Status                   int  `json:"status"`
	OrganizationOid          int  `json:"organization_oid"`
	DepartmentOid            int  `json:"department_oid"`
	DepartmentName           string `json:"department_name"`
	DepartmentTreePath       string `json:"department_tree_path"`
	StraightLineManagerUid   int  `json:"straight_line_manager_uid"`
	StraightLineManagerEmail string `json:"straight_line_manager_email"`
	StraightLineManagerName  string `json:"straight_line_manager_name"`
	DottedLineManagerUid     int  `json:"dotted_line_manager_uid"`
	DottedLineManagerEmail   string `json:"dotted_line_manager_email"`
	DottedLineManagerName    string `json:"dotted_line_manager_name"`
	MentorUid                int  `json:"mentor_uid"`
	MentorEmail              string `json:"mentor_email"`
	MentorName               string `json:"mentor_name"`
	IsCharge                 int  `gorm:"column:is_charge" json:"is_charge"`
	LeaderUid                int  `json:"leader_uid"`
	LeaderEmail              string `json:"leader_email"`
	LeaderName               string `json:"leader_name"`
}

type SSOUserInfoWithWxId struct {
	*SSOUserInfo
	WechatId string
}

type CheckLoginReq struct {
	PlatformId string `json:"platform_id"`
}

type CheckLoginResp struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data *SSOUserInfo `json:"data"`
}
type UsersResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data []*SSOUserInfo `json:"data"`
}

func CheckLogin(token,platformId string) (*SSOUserInfo, error) {
	url := conf.BaseConf.Domain.Sso + "/auth/verify-access-token"

	var req = map[string]string{
		"platform_id": platformId,
	}
	ret, err := json.Marshal(req)

	if err != nil {
		logging.Error("转化json失败", err)
		return nil, err
	}
	resp := requests.Post(url, requests.Header{"X-Token": token}, requests.JsonBodyBytes(ret))
	if resp.Err != nil {
		logging.Error("请求失败", resp.Err)
		return nil, resp.Err
	}

	var res = resp.GetBodyBytes()

	var respStruct CheckLoginResp
	err = json.Unmarshal(res, &respStruct)
	if err != nil {
		logging.Error("无法解析sso平台回复的内容", err)
		return nil, err
	}

	return respStruct.Data, nil
}

type OtherIdItem struct {
	Uid       int    `json:"uid"`
	OaUid     int    `json:"oa_uid"`
	WeChatUid string `json:"wechat_uid"`
}

func GetUserOtherIdByUIds(uidArr []int, typeStr string) (ecode.Code, []OtherIdItem) {
	resultArr := []OtherIdItem{}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := pkg.Md5(timestamp + "_" + conf.BaseConf.SsoConfig.SecretKey)
	url := conf.BaseConf.Domain.Sso + "/open-api/user/uid-map?timestamp=" + timestamp +
		"&signature=" + signature + "&platform_id=" + conf.BaseConf.SsoConfig.PlatformID

	reqBody := struct {
		Uids []int
		Type string
	}{
		Uids: uidArr,
		Type: typeStr,
	}
	jsonBody, _ := json.Marshal(reqBody)
	resp := requests.Post(url,requests.JsonBodyBytes(jsonBody))

	var res = string(resp.GetBodyBytes())

	if 0 != gjson.Get(res, "code").Int() {
		return ecode.EXTERNAL_API_FAIL, resultArr
	}

	err := json.Unmarshal([]byte(gjson.Get(res, "data").Raw), &resultArr)
	if err != nil {
		return ecode.EXTERNAL_API_JSON_ERR, resultArr
	}

	return ecode.OK, resultArr
}

func GetSsoUserInfoWithWeChatIdByUIds(uidArr []int) (ecode.Code, map[int]*SSOUserInfoWithWxId) {
	//uid去重
	queryUidArr, uidMap := []int{}, map[int]int{}
	for _, uid := range uidArr {
		if uid > 0 {
			if _, exist := uidMap[uid]; !exist {
				queryUidArr = append(queryUidArr, uid)
				uidMap[uid] = 1
			}
		}
	}

	resultMap := map[int]*SSOUserInfoWithWxId{}
	code, userInfoList := GetSsoUserInfoByUids(queryUidArr)
	if code != ecode.OK {
		return code, resultMap
	}

	code, wxIdArr := GetUserOtherIdByUIds(queryUidArr, "wechat")
	if code != ecode.OK {
		return code, resultMap
	}
	wxIdMap := map[int]string{}
	for _, resultItem := range wxIdArr {
		wxIdMap[resultItem.Uid] = resultItem.WeChatUid
	}

	for _, ssoUserInfo := range userInfoList {
		wxId := ""
		if wxIdInfo, exist := wxIdMap[ssoUserInfo.Uid]; exist {
			wxId = wxIdInfo
		}
		resultMap[ssoUserInfo.Uid] = &SSOUserInfoWithWxId{
			WechatId:    wxId,
			SSOUserInfo: &ssoUserInfo,
		}
	}
	return ecode.OK, resultMap
}

func GetUserUidInfoByWxId(wxUidArr []string) (ecode.Code, []OtherIdItem) {
	resultArr := []OtherIdItem{}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := pkg.Md5(timestamp + "_" + conf.BaseConf.SsoConfig.SecretKey)
	url := conf.BaseConf.Domain.Sso + "/open-api/user/other-uid-to-sso-uid?timestamp=" + timestamp +
		"&signature=" + signature + "&platform_id=" + conf.BaseConf.SsoConfig.PlatformID

	reqBody := struct {
		WechatUids []string
	}{
		WechatUids: wxUidArr,
	}
	jsonBody, _ := json.Marshal(reqBody)
	resp := requests.Post(url,requests.JsonBodyBytes(jsonBody))

	if resp.Err != nil {
		return ecode.EXTERNAL_API_FAIL, resultArr
	}

	var res = string(resp.GetBodyBytes())

	if 0 != gjson.Get(res, "code").Int() {
		return ecode.EXTERNAL_API_FAIL, resultArr
	}

	err := json.Unmarshal([]byte(gjson.Get(res, "data").Raw), &resultArr)
	code := ecode.OK
	if err != nil {
		code = ecode.EXTERNAL_API_FAIL
	}
	return code, resultArr
}
