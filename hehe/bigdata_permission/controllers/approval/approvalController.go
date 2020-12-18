package approval

import (
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/approvalSerializer"
	"bigdata_permission/service/approvalService"
	"bigdata_permission/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


func SubmitApply(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.CreateApproval{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	userInfo, exist := c.Get("userInfo")
	if !exist {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	userDetail, ok := userInfo.(*userService.UserDetail)
	if !ok {
		appG.JSON(http.StatusOK, ecode.NoUserFound, false)
		return
	}
	param.ApplyUid = userDetail.UserInfo.Uid
	code, reason := approvalService.SubmitApply(param)
	if code != ecode.OK {
		appG.JSON(http.StatusOK, code, reason)
		return
	}
	appG.JSON(http.StatusOK, ecode.OK, true)
}

func GetApprovalList(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.SearchApproval{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	code, _, list := approvalService.SearchApproval(param)
	appG.JSON(http.StatusOK, code, list)
}

func GetApplyInfo(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.ApprovalInfo{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	code, info := approvalService.GetApprovalInfo(param)
	appG.JSON(http.StatusOK, code, info)
}

func GetWxApprovalInfo(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.WxApprovalInfo{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	appG.JSON(http.StatusOK, ecode.OK, approvalService.GetWxApprovalDetail(param.WxApprovalId))
}

//批量审批（通过）
func BatchApprove(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.BatchApproval{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	userInfo, ok := c.Get("userInfo")
	if !ok {
		appG.JSON(http.StatusOK, ecode.NO_LOGIN, false)
		return
	}
	userDetail := userInfo.(*userService.UserDetail)
	param.ApprovalUid = strconv.Itoa(userDetail.SSOUserInfo.Uid)
	param.ApprovalTime = pkg.NowUnixMs()
	param.Status = approvalService.GetAccessedStatus()

	_, errList := approvalService.BatchApproval(param)
	if len(errList) > 0 {
		appG.JSON(http.StatusOK, ecode.UpdateErr, errList)
		return
	}
	appG.JSON(http.StatusOK, ecode.OK, true)
}

//拒绝
func BatchRefuse(c *gin.Context) {
	appG := http2.Gin{c}
	param := &approvalSerializer.BatchApproval{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, false)
		return
	}
	userInfo, ok := c.Get("userInfo")
	if !ok {
		appG.JSON(http.StatusOK, ecode.NO_LOGIN, false)
		return
	}
	userDetail := userInfo.(*userService.UserDetail)
	param.ApprovalUid = strconv.Itoa(userDetail.SSOUserInfo.Uid)
	param.ApprovalTime = pkg.NowUnixMs()
	param.Status = approvalService.GetRefusedStatus()

	_, errList := approvalService.BatchApproval(param)
	if len(errList) > 0 {
		appG.JSON(http.StatusOK, ecode.UpdateErr, errList)
		return
	}
	appG.JSON(http.StatusOK, ecode.OK, true)
}


