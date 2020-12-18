package wechat

import (
	"bigdata_permission/conf"
	"bigdata_permission/pkg/redis"
	"bigdata_permission/pkg/requests"
	"errors"
	"github.com/tidwall/gjson"
	"strconv"
)

//!!!企业微信可能会出于运营需要，提前使access_token失效，开发者应实现access_token失效时重新获取的逻辑。
func (w *WorkWxApp) getWeChatAccessToken() (string, error) {
	cacheKey := "asset_wechat_access_token_" + strconv.FormatInt(w.AgentId, 10)

	ok, accessToken := redis.Get(cacheKey)
	if !ok {
		//logging.Error("Redis get err:", err.Error())
	} else if accessToken != "" {
		return accessToken, nil
	}

	url := conf.BaseConf.WeChat.ApiHost + "/cgi-bin/gettoken"
	getParam := requests.QueryString{
		"corpid":     w.WorkWx.CorpId,
		"corpsecret": w.CorpSecret,
	}

	//尝试两次，直到成功
	var resp *requests.Response
	for i := 0; i < 2; i++ {
		resp = requests.Get(url, getParam)
		if resp.Err == nil {
			break
		}
	}

	res := resp.GetBodyBytes()

	if gjson.GetBytes(res, "errcode").Int() != 0 {
		return "", errors.New(gjson.GetBytes(res, "errmsg").String())
	}

	accessToken = gjson.GetBytes(res, "access_token").String()
	expiresTime := gjson.GetBytes(res, "expires_in").Int() - 120
	redis.Set(cacheKey, accessToken, int(expiresTime))

	return accessToken, nil
}
