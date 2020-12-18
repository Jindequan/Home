package user

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/userSerializer"
	"bigdata_permission/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userSerializer.CreateUser{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, user := userService.CreateUser(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, user)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateUser(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userSerializer.UpdateUser{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, user := userService.UpdateUser(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, user)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteUser(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userSerializer.DeleteUser{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := userService.DeleteUser(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchUser(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userSerializer.SearchUser{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, userList := userService.SearchUser(param)
	appG.JSON(http.StatusOK, code, userList)
	return
}

func GetUserInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userSerializer.GetUser{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, userInfo := userService.GetUserInfo(param)
	appG.JSON(http.StatusOK, code, userInfo)
	return
}
