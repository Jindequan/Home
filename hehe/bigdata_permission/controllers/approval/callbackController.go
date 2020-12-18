package approval

import (
	"bigdata_permission/pkg"
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/pkg/logging"
	"bigdata_permission/pkg/wechat"
	"bigdata_permission/service/approvalService"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func WxApprovalCallback(c *gin.Context) {
	appG := &http2.Gin{c}
	params := &wechat.ApprovalCallbackMsgContent{}
	data, _ := ioutil.ReadAll(c.Request.Body)

	if err := c.BindQuery(params); err != nil {
		pkg.SendToRobot(pkg.ROBOT_ERROR, "微信审批回调失败：缺少必要参数" + string(data))
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, nil)
		return
	}

	code, msg := approvalService.WxApprovalCallback(params)

	if code != ecode.OK {
		logging.Error(fmt.Sprintf("微信审批变更回调失败:%s;%s", msg, string(data)))
		pkg.SendToRobot(pkg.ROBOT_ERROR, fmt.Sprintf("微信审批变更回调失败:%s;%s", msg, string(data)))
	} else {
		logging.Info(fmt.Sprintf("微信审批变更回调成功:%s", string(data)))
	}

	appG.JSON(http.StatusOK, code, true)
}