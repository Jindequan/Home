package wechat

import (
	"bigdata_permission/conf"
	"fmt"
)

//appKey
type AppKey string

var WorkWxAppApprovalModules AppKey = "approval_modules"

var workWxConf = make(map[AppKey]*WorkWxApp)

//企业微信客户端维度
type WorkWx struct {
	CorpId string
}

//企业微信应用维度  指向企业微信客户端
type WorkWxApp struct {
	*WorkWx

	//CorpSecret 应用密钥
	CorpSecret string
	//AgentId 应用ID
	AgentId int64
}

type weChatResCommon struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func init() {
	weChatConf := conf.BaseConf.WeChat
	workWxConf[WorkWxAppApprovalModules] = New(weChatConf.ApprovalApp.AppID, weChatConf.ApprovalApp.CorpSecret, weChatConf.ApprovalApp.AgentID)
}

func GetWorkWxApp(appKey AppKey) *WorkWxApp {
	fmt.Println(workWxConf)
	return workWxConf[appKey]
}

func New(corpId string, corpSecret string, agentId int64) *WorkWxApp {
	return &WorkWxApp{&WorkWx{CorpId: corpId}, corpSecret, agentId}
}
