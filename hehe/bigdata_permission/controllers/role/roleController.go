package role

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/roleSerializer"
	"bigdata_permission/service/roleService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRole (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleSerializer.CreateRole{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, role := roleService.CreateRole(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, role)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateRole(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleSerializer.UpdateRole{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, role := roleService.UpdateRole(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, role)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteRole(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleSerializer.DeleteRole{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := roleService.DeleteRole(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchRole(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleSerializer.SearchRole{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, roleList := roleService.SearchRole(param)
	appG.JSON(http.StatusOK, code, roleList)
	return
}

func GetRoleInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &roleSerializer.GetRole{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, roleList := roleService.GetRoleInfo(param)
	appG.JSON(http.StatusOK, code, roleList)
	return
}
