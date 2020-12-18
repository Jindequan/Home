package interfaceDic

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/interfaceSerializer"
	"bigdata_permission/service/interfaceService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateInterface (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &interfaceSerializer.CreateInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, interfaceInfo := interfaceService.CreateInterface(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, interfaceInfo)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func UpdateInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &interfaceSerializer.UpdateInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg, interfaceInfo := interfaceService.UpdateInterface(param)
	if code == ecode.OK {
		appG.JSON(http.StatusOK, code, interfaceInfo)
		return
	}
	appG.JSON(http.StatusOK, code, errMsg)
}

func DeleteInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &interfaceSerializer.DeleteInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, errMsg := interfaceService.DeleteInterface(param)
	appG.JSON(http.StatusOK, code, errMsg)
}

func SearchInterface(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &interfaceSerializer.SearchInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, _, interfaceList := interfaceService.SearchInterface(param)
	appG.JSON(http.StatusOK, code, interfaceList)
	return
}

func GetInterfaceInfo(c *gin.Context) {
	appG := &http2.Gin{c}
	param := &interfaceSerializer.GetInterface{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, err.Error())
		return
	}
	code, interfaceInfo := interfaceService.GetInterfaceInfo(param)
	appG.JSON(http.StatusOK, code, interfaceInfo)
	return
}
