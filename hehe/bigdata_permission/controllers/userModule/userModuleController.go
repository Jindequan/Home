package userModule

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/userModuleSerializer"
	"bigdata_permission/service/userModuleService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUserModule (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userModuleSerializer.CreateUserModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := userModuleService.CreateUserModuleBatch(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, true)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func SaveUserModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &userModuleSerializer.SaveUserModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := userModuleService.SaveUserModule(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, true)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

