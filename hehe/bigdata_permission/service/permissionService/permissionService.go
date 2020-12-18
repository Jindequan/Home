package permissionService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/service/userModuleService"
)

func CheckPermission(moduleId int, roleId int, uid int) bool {
	//用户无角色：无权限
	if roleId == 0 {
		return false
	}
	//超管：有权限
	if roleId == dao.SUPER_ADMIN {
		return true
	}
	//查询模块
	_, moduleInfo := dao.GetModuleDicById(moduleId)
	//无模块：无权限
	if moduleId != moduleInfo.ModuleId {
		return false
	}
	//用户模块
	code, exist := userModuleService.ExistUserModuleById(uid, moduleId)
	if code != ecode.OK {
		return false
	}

	return exist
}


//func CheckPermission(moduleId int, path string, interfaceType int8, roleId int) bool {
//	//用户无角色：无权限
//	if roleId == 0 {
//		return false
//	}
//	//超管：有权限
//	if roleId == dao.SUPER_ADMIN {
//		return true
//	}
//	//获取角色所有模块权限
//	_, roleModules := dao.GetRoleModuleByRoleId(roleId)
//	//无任何权限：无权限
//	if len(*roleModules) == 0 {
//		return false
//	}
//	//不在可用列表中：无权限
//	hasModule := false
//	for _, item := range *roleModules {
//		if item.ModuleId == moduleId {
//			hasModule = true
//			break
//		}
//	}
//	if !hasModule {
//		return false
//	}
//	//查询模块
//	_, moduleInfo := dao.GetModuleDicById(moduleId)
//	//无模块：无权限
//	if moduleId != moduleInfo.ModuleId {
//		return false
//	}
//	//获取接口信息
//	_, interfaceInfo := dao.GetInterfaceByTypeAndPath(interfaceType, path)
//	//无录入接口：无权限
//	if interfaceInfo.InterfaceId == 0 {
//		return false
//	}
//	//接口模块mapping，得到模块id
//	_, moduleInterface := dao.GetModuleInterfaceByModuleIdAndInterfaceId(moduleId, interfaceInfo.InterfaceId)
//	//无mapping：无权限
//	if moduleInterface.ModuleId != moduleId {
//		return false
//	}
//
//	return true
//}
