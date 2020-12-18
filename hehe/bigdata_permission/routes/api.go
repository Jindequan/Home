package routes

import (
	"bigdata_permission/controllers/approval"
	"bigdata_permission/controllers/common"
	"bigdata_permission/controllers/interfaceDic"
	"bigdata_permission/controllers/moduleDic"
	"bigdata_permission/controllers/moduleInterface"
	"bigdata_permission/controllers/permission"
	"bigdata_permission/controllers/role"
	"bigdata_permission/controllers/roleModule"
	"bigdata_permission/controllers/test"
	"bigdata_permission/controllers/user"
	"bigdata_permission/controllers/userModule"
	"bigdata_permission/middleware"
	"github.com/gin-gonic/gin"
)

func routeListForApi(r *gin.Engine) *gin.Engine {

	userGroup := r.Group("user")
	{
		userGroup.Use(middleware.Auth())
		userGroup.Use(middleware.CheckPermission())
		userGroup.POST("create", user.CreateUser)
		userGroup.POST("update", user.UpdateUser)
		userGroup.POST("delete", user.DeleteUser)
		userGroup.GET("search", user.SearchUser)
		userGroup.GET("info", user.GetUserInfo)
	}

	roleGroup := r.Group("role")
	{
		roleGroup.Use(middleware.Auth())
		roleGroup.Use(middleware.CheckPermission())
		roleGroup.POST("create", role.CreateRole)
		roleGroup.POST("update", role.UpdateRole)
		roleGroup.POST("delete", role.DeleteRole)
		roleGroup.GET("search", role.SearchRole)
		roleGroup.GET("info", role.GetRoleInfo)
	}

	interfaceGroup := r.Group("interface")
	{
		interfaceGroup.Use(middleware.Auth())
		interfaceGroup.Use(middleware.CheckPermission())
		interfaceGroup.POST("create", interfaceDic.CreateInterface)
		interfaceGroup.POST("update", interfaceDic.UpdateInterface)
		interfaceGroup.POST("delete", interfaceDic.DeleteInterface)
		interfaceGroup.GET("search", interfaceDic.SearchInterface)
		interfaceGroup.GET("info", interfaceDic.GetInterfaceInfo)
	}

	moduleGroup := r.Group("module")
	{
		moduleGroup.Use(middleware.Auth())
		moduleGroup.Use(middleware.CheckPermission())
		moduleGroup.POST("create", moduleDic.CreateModule)
		moduleGroup.POST("update", moduleDic.UpdateModule)
		moduleGroup.POST("delete", moduleDic.DeleteModule)
		moduleGroup.GET("search", moduleDic.SearchModule)
		moduleGroup.GET("info", moduleDic.GetModuleInfo)
	}

	moduleInterfaceGroup := r.Group("module-interface")
	{
		moduleInterfaceGroup.Use(middleware.Auth())
		moduleInterfaceGroup.Use(middleware.CheckPermission())
		moduleInterfaceGroup.POST("create", moduleInterface.CreateModuleInterface)
		moduleInterfaceGroup.POST("update", moduleInterface.UpdateModuleInterface)
		moduleInterfaceGroup.POST("delete", moduleInterface.DeleteModuleInterface)
		moduleInterfaceGroup.GET("search", moduleInterface.SearchModuleInterface)
		moduleInterfaceGroup.GET("info", moduleInterface.GetModuleInterfaceInfo)

		moduleInterfaceGroup.GET("module-all-interface", moduleInterface.GetModuleAllInterface)
	}

	roleModuleGroup := r.Group("role-module")
	{
		roleModuleGroup.Use(middleware.Auth())
		roleModuleGroup.Use(middleware.CheckPermission())
		roleModuleGroup.POST("create", roleModule.CreateRoleModule)
		roleModuleGroup.POST("batch-create", roleModule.CreateRoleModuleBatch)
		roleModuleGroup.POST("update", roleModule.UpdateRoleModule)
		roleModuleGroup.POST("delete", roleModule.DeleteRoleModule)
		roleModuleGroup.GET("search", roleModule.SearchRoleModule)
		roleModuleGroup.GET("info", roleModule.GetRoleModuleInfo)
	}

	userModuleGroup := r.Group("user-module")
	{
		userModuleGroup.Use(middleware.Auth())
		userModuleGroup.Use(middleware.CheckPermission())
		userModuleGroup.POST("batch-create", userModule.CreateUserModule)
		userModuleGroup.POST("batch-save", userModule.SaveUserModule)
	}

	permissionGroup := r.Group("permission")
	{
		permissionGroup.Use(middleware.Auth())
		permissionGroup.POST("check", permission.CheckPermission)
	}

	commonGroup := r.Group("common")
	{
		commonGroup.GET("dictionary", common.GetDictionary)
	}

	approvalGroup := r.Group("approval")
	{
		approvalGroup.POST("wx-callback", approval.WxApprovalCallback)
		approvalGroup.Use(middleware.Auth())
		approvalGroup.Use(middleware.CheckPermission())
		approvalGroup.POST("apply", approval.SubmitApply)
		approvalGroup.POST("approve", approval.BatchApprove)
		approvalGroup.POST("refuse", approval.BatchRefuse)
		approvalGroup.GET("search", approval.GetApprovalList)
		approvalGroup.GET("info", approval.GetApplyInfo)
		approvalGroup.GET("info-wx", approval.GetWxApprovalInfo)
	}

	testGroup := r.Group("test")
	{
		testGroup.GET("no-auth", test.TestNoAuth)
		testGroup.Use(middleware.Auth())
		testGroup.GET("auth", test.Test)
	}
	return r
}
