package externalService

import (
	"bigdata_permission/conf"
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/logging"
	"bigdata_permission/pkg/requests"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

type WeChatMessageBtn struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	ReplaceName string `json:"replace_name"`
	Color       string `json:"color"`
	IsBold      bool   `json:"is_bold"`
}

type SendWeChatMessageRequest struct {
	ToUIDs      []string           `json:"to_uids"`
	MsgType     int64              `json:"msg_type"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Url         string             `json:"url"`
	IsSafe      int64              `json:"is_safe"`
	Text        string             `json:"text"`
	BtnText     string             `json:"btn_text"`
	PicUrl      string             `json:"pic_url"`
	CallbackUrl string             `json:"callback_url"`
	Extra       string             `json:"extra"`
	Btn         []WeChatMessageBtn `json:"btn"`
}

type SendMessageResponse struct {
	InvalidUIDs   []int64 `json:"invalid_uids"`
	InvalidEmails []int64 `json:"invalid_emails"`
	MsgID         int64   `json:"msg_id"`
}

const (
	MSG_TYPE_TEXT      int64 = 1
	MSG_TYPE_PIC       int64 = 2
	MSG_TYPE_AUDIO     int64 = 3
	MSG_TYPE_VIDEO     int64 = 4
	MSG_TYPE_FILE      int64 = 5
	MSG_TYPE_CARD      int64 = 6
	MSG_TYPE_PIC_TEXT  int64 = 7
	MSG_TYPE_CARD_TASK int64 = 8
)

func SendWeChatMessage(req *SendWeChatMessageRequest) (bool, string, *SendMessageResponse) {
	ret := &SendMessageResponse{}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := getSignature(timestamp)
	url := conf.BaseConf.Domain.Sso + "/open-api/message/send-wechat"
	queryString := requests.QueryString{
		"timestamp":   timestamp,
		"signature":   signature,
		"platform_id": conf.BaseConf.SsoConfig.PlatformID,
	}
	btnJson := ""
	if len(req.Btn) > 0 {
		btn, _ := json.Marshal(req.Btn)
		btnJson = string(btn)
	}

	formData := requests.FormData{
		"to_uids":      req.ToUIDs,
		"msg_type":     req.MsgType,
		"title":        req.Title,
		"description":  req.Description,
		"url":          req.Url,
		"is_safe":      req.IsSafe,
		"text":         req.Text,
		"btn_text":     req.BtnText,
		"pic_url":      req.PicUrl,
		"callback_url": req.CallbackUrl,
		"extra":        req.Extra,
		"btn":          btnJson,
	}
	response := requests.Post(url, queryString, formData)
	if response.Err != nil {
		return false, "sso发送微信消息错误" + response.Err.Error(), ret
	}

	errCode := gjson.GetBytes(response.GetBodyBytes(), "code").Int()
	errMsg := gjson.GetBytes(response.GetBodyBytes(), "msg").String()
	if errCode != 0 {
		logging.Error(errMsg)
		return false, "sso发送微信消息错误" + errMsg, ret
	}
	data := gjson.GetBytes(response.GetBodyBytes(), "data").String()
	jsonErr := json.Unmarshal([]byte(data), ret)
	if jsonErr != nil {
		logging.Error(data)
		return false, "sso发送微信消息错误" + jsonErr.Error(), ret
	}
	return true, "", ret
}

type SsoEmailParams struct {
	Subject  string   `json:"subject"`
	Content  string   `json:"content"`
	FromName string   `json:"from_name"`
	ToUids   []uint   `json:"to_uids"`   //sso_uid
	ToEmails []string `json:"to_emails"` //emails
	AttachFileIds []string `json:"attach_file_ids"` //文件id
}

type SsoEmailResponse struct {
	InvalidUids   []uint   `json:"invalid_uids"`
	InvalidEmails []string `json:"invalid_emails"`
	MsgId         uint     `json:"msg_id"`
}

func SendSsoEmail(params SsoEmailParams) (ecode.Code, *SsoEmailResponse) {
	SsoDomain := conf.BaseConf.Domain.Sso
	code := ecode.OK
	ssoEmailResponse := &SsoEmailResponse{}
	reqBody, _ := json.Marshal(params)

	timestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	resp := requests.Post(
		SsoDomain+"/open-api/message/send-email?platform_id="+
			conf.BaseConf.SsoConfig.PlatformID+
			"&timestamp="+timestampStr+
			"&signature="+getSignature(timestampStr),
		requests.JsonBodyBytes(reqBody),
	)

	if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
		code = ecode.EXTERNAL_API_FAIL
	} else {
		err := json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), ssoEmailResponse)
		if err != nil {
			code = ecode.EXTERNAL_API_JSON_ERR
		}
	}

	return code, ssoEmailResponse
}

type SsoFileResponse struct {
	FileUrl string `json:"file_url"`
	FileId string `json:"file_id"`
}

func UploadFileSso(fileName, filePath string) (bool, string, *SsoFileResponse) {
	uploadedFile := &SsoFileResponse{}
	isSuccess, token := getSsoToken()
	if !isSuccess {
		return false, "获取token失败", uploadedFile
	}
	url := conf.BaseConf.Domain.Sso + "/file/upload?platform_id=" + conf.BaseConf.SsoConfig.PlatformID +
		"&access_token=" + token

	if _, err := os.Lstat(filePath); err != nil {
		return false, "获取文件信息失败", uploadedFile
	}
	file, err := os.Open(filePath)
	if err != nil {
		return false, "读取文件内容失败", uploadedFile
	}
	defer file.Close()

	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	fileWriter, err := writer.CreateFormFile("file", fileName)
	f, err := os.Open(filePath)
	io.Copy(fileWriter, f)
	writer.WriteField("file_name", fileName)
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), bodyBuf)
	if err != nil || res.Body == nil{
		return false, "接口response数据为空", uploadedFile
	}
	resp, err := ioutil.ReadAll(res.Body)

	if err != nil || gjson.Get(string(resp), "code").String() != "0" {
		return false, "读取接口response数据失败", uploadedFile
	}
	data := gjson.Get(string(resp), "data").String()
	err = json.Unmarshal([]byte(data), uploadedFile)
	fmt.Println(data,uploadedFile)
	return err == nil, "", uploadedFile
}

type ssoTokenParam struct {
	PlatformId string `json:"platform_id"`
	Sign string `json:"sign"`
}
func getSsoToken() (bool, string) {
	url := conf.BaseConf.Domain.Sso + "/open-api/get-access-token"
	sign := pkg.Md5(conf.BaseConf.SsoConfig.SecretKey)
	reqBody, _ := json.Marshal(ssoTokenParam{
		PlatformId: conf.BaseConf.SsoConfig.PlatformID,
		Sign: sign,
	})

	resp := requests.Post(url, requests.JsonBodyBytes(reqBody))
	code := gjson.Get(resp.GetBodyString(), "code").String()
	token := gjson.Get(resp.GetBodyString(), "data").String()
	return code == "0", token
}

func getSignature(timestampStr string) string {
	h := md5.New()
	h.Write([]byte(timestampStr + "_" + conf.BaseConf.SsoConfig.SecretKey))
	return hex.EncodeToString(h.Sum(nil))
}
