package middleware

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/service/permissionService"
	"bigdata_permission/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data interface{}
		userInfo, exists := c.Get("userInfo")
		if !exists {
			code := ecode.PERMISSION_DENIED
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}
		userDetail := userInfo.(*userService.UserDetail)
		if !permissionService.CheckPermission(dao.CURRENT_MODULE_ID, userDetail.UserInfo.RoleId, userDetail.UserInfo.Uid) {
			code := ecode.PERMISSION_DENIED
			c.JSON(http.StatusOK, gin.H{
				"code": code.Code(),
				"msg":  code.Message(),
				"data": data,
			})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}