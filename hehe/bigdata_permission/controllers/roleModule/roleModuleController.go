package roleModule

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/roleModuleSerializer"
	"bigdata_permission/service/roleModuleService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRoleModule (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.CreateRoleModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, roleModule := roleModuleService.CreateRoleModule(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, roleModule)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func CreateRoleModuleBatch (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.CreateRoleModuleBatch{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := roleModuleService.CreateRoleModuleBatch(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, true)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateRoleModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.UpdateRoleModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, roleModule := roleModuleService.UpdateRoleModule(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, roleModule)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteRoleModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.DeleteRoleModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := roleModuleService.DeleteRoleModule(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchRoleModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.SearchRoleModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, roleModuleList := roleModuleService.SearchRoleModule(param)
	appG.JSON(http.StatusOK, code, roleModuleList)
	return
}

func GetRoleModuleInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleModuleSerializer.GetRoleModule{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, roleModuleInfo := roleModuleService.GetRoleModuleInfo(param)
	appG.JSON(http.StatusOK, code, roleModuleInfo)
	return
}
