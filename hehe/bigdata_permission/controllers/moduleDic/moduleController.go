package moduleDic

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/moduleDicSerializer"
	"bigdata_permission/service/moduleDicService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateModule (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleDicSerializer.CreateModuleDic{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, module := moduleDicService.CreateModuleDic(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, module)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleDicSerializer.UpdateModuleDic{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, module := moduleDicService.UpdateModuleDic(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, module)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleDicSerializer.DeleteModuleDic{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := moduleDicService.DeleteModuleDic(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchModule(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleDicSerializer.SearchModuleDic{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, moduleList := moduleDicService.SearchModuleDic(param)
	appG.JSON(http.StatusOK, code, moduleList)
	return
}

func GetModuleInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &moduleDicSerializer.GetModuleDic{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, moduleInfo := moduleDicService.GetModuleDicInfo(param)
	appG.JSON(http.StatusOK, code, moduleInfo)
	return
}
