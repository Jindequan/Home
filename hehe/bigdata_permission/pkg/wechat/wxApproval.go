package wechat

import (
	"bigdata_permission/conf"
	"bigdata_permission/pkg/requests"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
)

type ApprovalRequestParam struct {
	CreatorUserid       string          `json:"creator_userid"`
	TemplateId          string          `json:"template_id"`
	UseTemplateApprover int             `json:"use_template_approver"`
	Approver            []*ApproverItem `json:"approver"`
	Notifyer            []string        `json:"notifyer"`
	NotifyType          int             `json:"notify_type"`
	ApplyData           ApplyDataDetail `json:"apply_data"`
	SummaryList         []*SummaryInfo  `json:"summary_list"`
}

type ApproverItem struct {
	Attr   int      `json:"attr"`
	Userid []string `json:"userid"`
}

type ApplyDataDetail struct {
	Contents []*ContentItem `json:"contents"`
}

type ContentItem struct {
	Control string       `json:"control"`
	Id      string       `json:"id"`
	Value   ContentValue `json:"value"`
}

type ContentValue struct {
	Text      string                   `json:"text"`
	NewNumber int                      `json:"new_number"`
	Date      ContentSubtypeDate       `json:"date"`
	Selector  ContentSubtypeSelector   `json:"selector"`
	Members   []*MemberItem            `json:"members"`
	Children  []ContentSubtypeChildren `json:"children"`
}

type ContentSubtypeDate struct {
	Type       string `json:"type"`
	STimestamp string `json:"s_timestamp"`
}

type ContentSubtypeSelector struct {
	Type    string       `json:"type"`
	Options []OptionItem `json:"options"`
}

type OptionItem struct {
	Key string `json:"key"`
}

type MemberItem struct {
	Userid string `json:"userid"`
	Name   string `json:"name"`
}

type ContentSubtypeChildren struct {
	List []*ContentItem `json:"list"`
}

type SummaryInfo struct {
	SummaryInfo []SummaryItem `json:"summary_info"`
}

type SummaryItem struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

type ApprovalCallbackRequest struct {
	MsgSignature string `json:"msg_signature" form:"msg_signature"`
	Timestamp    string `json:"timestamp" form:"timestamp"`
	Nonce        string `json:"nonce" form:"nonce"`
	Echostr      string `json:"echostr" form:"echostr"`
}

//微信回调Echo会变成post
type ApprovalCallbackPostRequest struct {
	ToUserName string
	AgentID    string
	Encrypt    string
}

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
}

//微信回调正文内容
type ApprovalCallbackMsgContent struct {
	ToUsername   string       `xml:"ToUserName"`
	FromUsername string       `xml:"FromUserName"`
	CreateTime   int64        `xml:"CreateTime"`
	MsgType      string       `xml:"MsgType"`
	Agentid      int          `xml:"AgentId"`
	Event        string       `xml:"Event"`
	ApprovalInfo ApprovalInfo `xml:"ApprovalInfo"`
}

//微信查询审批详情
type ApprovalDetail struct {
	SpNo       string          `json:"sp_no"`
	SpName     string          `json:"sp_name"`
	SpStatus   int             `json:"sp_status"`
	TemplateId string          `json:"template_id"`
	ApplyTime  int64           `json:"apply_time"`
	Applyer    ApplyerItem     `json:"applyer"`
	SpRecord   []SpRecordItem  `json:"sp_record"`
	Notifyer   []ApplyerItem   `json:"notifyer"`
	ApplyData  ApplyDataDetail `json:"apply_data"`
	Comments   []CommentItem   `json:"comments"`
}

type ApprovalInfo struct {
	SpNo             int64        `xml:"SpNo" binding:"required"`
	SpName           string       `xml:"SpName"`
	SpStatus         int8         `xml:"SpStatus" binding:"required"`
	TemplateId       string       `xml:"TemplateId"`
	ApplyTime        int64        `xml:"ApplyTime"`
	Applyer          ApplyerItem  `xml:"Applyer"`
	SpRecord         SpRecordItem `xml:"SpRecord"`
	Notifyer         ApplyerItem  `xml:"Notifyer"`
	StatuChangeEvent int8         `xml:"StatuChangeEvent"`
}

type ApplyerItem struct {
	UserId string `xml:"UserId"`
}

type SpRecordItem struct {
	SpStatus     int8         `xml:"SpStatus"`
	ApproverAttr int          `xml:"ApproverAttr"`
	Details      []DetailItem `xml:"Details"`
}

type DetailItem struct {
	Approver ApplyerItem `xml:"Approver"`
	Speech   string      `xml:"Speech"`
	SpStatus int8        `xml:"SpStatus"`
	SpTime   int64       `xml:"SpTime"`
}

type CommentItem struct {
	CommentUserInfo ApplyerItem `json:"commentUserInfo"`
	CommentTime     int64       `json:"commenttime"`
	CommentContent  string      `json:"commentcontent"`
	CommentId       string      `json:"commentid"`
	MediaId         []string    `json:"media_id"`
}

const (
	//提单
	WX_APPROVAL_EVENT_CREATE = 1
	//同意
	WX_APPROVAL_EVENT_PASS = 2
	//驳回
	WX_APPROVAL_EVENT_REJECT = 3
	//转审
	WX_APPROVAL_EVENT_DELIVERY = 4
	//催办
	WX_APPROVAL_EVENT_HURRY = 5
	//撤销
	WX_APPROVAL_EVENT_CANCEL = 6
	//通过后撤销
	WX_APPROVAL_EVENT_CANCEL_AFTER_PASS = 8
	//添加备注
	WX_APPROVAL_EVENT_COMMENT = 10

	//审批中
	WX_APPROVAL_NODE_STATUS_PENDING = 1
	//已同意
	WX_APPROVAL_NODE_STATUS_PASS = 2
	//已驳回
	WX_APPROVAL_NODE_STATUS_REJECT = 3
	//已转审
	WX_APPROVAL_NODE_STATUS_TRANSFORM = 4

	//审批中
	WX_APPROVAL_STATUS_PENDING = 1
	//已同意
	WX_APPROVAL_STATUS_PASS = 2
	//已驳回
	WX_APPROVAL_STATUS_REJECT = 3
	//已撤销
	WX_APPROVAL_STATUS_CANCEL = 4
	//通过后撤回
	WX_APPROVAL_STATUS_CANCEL_AFTER_PASS = 6
	//已删除
	WX_APPROVAL_STATUS_DELETE = 7
	//已支付
	WX_APPROVAL_STATUS_PAY = 10

	WX_APPROVAL_STATUS_ACCESS = "agree"
	WX_APPROVAL_STATUS_REFUSE = "disagree"
)

func (w *WorkWxApp) CreateApproval(param *ApprovalRequestParam) (string, error) {
	accessToken, err := w.getWeChatAccessToken()
	if err != nil {
		return "", err
	}

	url := conf.BaseConf.WeChat.ApiHost + "/cgi-bin/oa/applyevent"
	body, _ := json.Marshal(param)

	resp := requests.Post(url, requests.QueryString{
		"access_token": accessToken,
	}, requests.JsonBodyBytes(body))

	res := resp.GetBodyBytes()

	if gjson.GetBytes(res, "errcode").Int() != 0 {
		return "", errors.New(gjson.GetBytes(res, "errmsg").String())
	}

	return gjson.GetBytes(res, "sp_no").String(), nil
}

func (w *WorkWxApp) GetApprovalDetail(approvalId int64) (ApprovalDetail, error) {
	var approvalDetail ApprovalDetail
	accessToken, err := w.getWeChatAccessToken()
	if err != nil {
		return approvalDetail, err
	}

	url := conf.BaseConf.WeChat.ApiHost + "/cgi-bin/oa/getapprovaldetail"
	body, _ := json.Marshal(struct {
		SpNo int64 `json:"sp_no"`
	}{
		SpNo: approvalId,
	})

	resp := requests.Post(url, requests.QueryString{
		"access_token": accessToken,
	}, requests.JsonBodyBytes(body))

	res := resp.GetBodyBytes()

	if gjson.GetBytes(res, "errcode").Int() != 0 {
		return approvalDetail, errors.New(gjson.GetBytes(res, "errmsg").String())
	}

	infoStr := gjson.GetBytes(res, "info").String()

	if infoStr != "" {
		err = json.Unmarshal([]byte(infoStr), &approvalDetail)
	}

	return approvalDetail, err
}