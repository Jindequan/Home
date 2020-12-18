package middleware

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/requests"
	"bigdata_permission/service/userService"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"bigdata_permission/conf"
	"github.com/tidwall/gjson"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code ecode.Code
		var data interface{}
		token := ""
		for k, v := range c.Request.Header {
			if k == "X-Token" {
				token = v[0]
			}
		}
		if token == "" {
			code = ecode.NO_LOGIN
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}
		//校验token
		SsoDomain := conf.BaseConf.Domain.Sso
		reqBody, err := json.Marshal(struct {
			PlatformId string `json:"platform_id"`
		}{
			PlatformId: conf.BaseConf.SsoConfig.PlatformID,
		})

		if err != nil {
			code = ecode.REQUEST_EXCEPTION
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}

		resp := requests.Post(
			SsoDomain+"/auth/verify-access-token",
			requests.Header{"X-Token": token},
			requests.JsonBodyBytes(reqBody),
		)

		if gjson.Get(resp.GetBodyString(), "code").Int() != int64(ecode.OK.Code()) {
			code = ecode.NO_LOGIN
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}

		var ssoUserInfo userService.SSOUserInfo
		err = json.Unmarshal([]byte(gjson.Get(resp.GetBodyString(), "data").String()), &ssoUserInfo)
		if err != nil {
			code = ecode.NO_LOGIN
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}

		if ssoUserInfo.Uid > 0 {
			res, userInfo := dao.GetUserById(ssoUserInfo.Uid)
			if res {
				userDetail := &userService.UserDetail{
					ssoUserInfo,
					userService.TransferUserInfo(*userInfo),
				}
				c.Set("token", token)
				c.Set("userInfo", userDetail)

			} else {
				code = ecode.PERMISSION_DENIED
				c.JSON(http.StatusOK, gin.H{
					"code": code.Code(),
					"msg":  code.Message(),
					"data": data,
				})
				c.Abort()
				return
			}
		}

		return
	}
}