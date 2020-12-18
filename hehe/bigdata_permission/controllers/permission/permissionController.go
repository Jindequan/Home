package permission

import (
	"bigdata_permission/pkg/ecode"
	http2 "bigdata_permission/pkg/http"
	"bigdata_permission/serializer/permissionSerializer"
	"bigdata_permission/serializer/userSerializer"
	"bigdata_permission/service/permissionService"
	"bigdata_permission/service/userService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckPermission (c *gin.Context) {
	appG := &http2.Gin{c}
	param := &permissionSerializer.CheckPermission{}
	if err := c.ShouldBind(param); err != nil {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, "module_id, interface path, interface type required")
		return
	}
	if param.RoleId == 0 && param.Uid == 0 {
		appG.JSON(http.StatusOK, ecode.INVALID_PARAM, "check user's id or user's role id")
		return
	}
	if param.RoleId == 0 {
		_, user := userService.GetUserInfo(&userSerializer.GetUser{Uid: param.Uid})
		if user.Uid != param.Uid {
			appG.JSON(http.StatusOK, ecode.NoUserFound, false)
			return
		}
		param.RoleId = user.RoleId
	}
	if param.RoleId == 0 {
		appG.JSON(http.StatusOK, ecode.NoUserRole, false)
		return
	}
	res := permissionService.CheckPermission(param.ModuleId, param.RoleId, param.Uid)

	appG.JSON(http.StatusOK, ecode.OK, res)
}
