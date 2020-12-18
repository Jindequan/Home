package moduleInterface

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/moduleInterfaceSerializer"
	"bigdata_permission/service/moduleInterfaceService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateModuleInterface (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.CreateModuleInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, moduleInterface := moduleInterfaceService.CreateModuleInterface(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, moduleInterface)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateModuleInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.UpdateModuleInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, moduleInterface := moduleInterfaceService.UpdateModuleInterface(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, moduleInterface)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteModuleInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.DeleteModuleInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := moduleInterfaceService.DeleteModuleInterface(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchModuleInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.SearchModuleInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, moduleInterfaceList := moduleInterfaceService.SearchModuleInterface(param)
	appG.JSON(http.StatusOK, code, moduleInterfaceList)
	return
}

func GetModuleInterfaceInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.GetModuleInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, moduleInterfaceInfo := moduleInterfaceService.GetModuleInterfaceInfo(param)
	appG.JSON(http.StatusOK, code, moduleInterfaceInfo)
	return
}

func GetModuleAllInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleInterfaceSerializer.GetModuleChildrenInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	if param.ModuleId < 1 {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	tree := moduleInterfaceService.GetChildModuleInterface(param.ModuleId)
	appG.JSON(http.StatusOK, ecode.OK, tree)
}
