package roleModuleService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/roleModuleSerializer"
	"bigdata_permission/service/common"
)

type RoleModuleInfo struct {
	RoleId int `json:"role_id"`
	ModuleId int `json:"module_id"`

	RoleName string `json:"role_name"`
	ModuleName string `json:"module_name"`
}

type RoleModuleList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]RoleModuleInfo `json:"list"`
}

func TransferRoleModuleInfo (roleModule dao.RoleModule) RoleModuleInfo {
	return RoleModuleInfo{
		RoleId: roleModule.RoleId,
		ModuleId: roleModule.ModuleId,

		RoleName: dictionary.GetStringValue(dictionary.RoleList, roleModule.RoleId),
		ModuleName: dictionary.GetStringValue(dictionary.ModuleList, roleModule.ModuleId),
	}
}

func CreateRoleModule(param *roleModuleSerializer.CreateRoleModule) (ecode.Code, string, RoleModuleInfo) {
	if param.RoleId == 0 || param.ModuleId == 0 {
		return ecode.INVALID_PARAM, "", RoleModuleInfo{}
	}
	_, exist := dao.GetRoleModuleByRoleIdAndModuleId(param.RoleId, param.ModuleId)
	if exist.RoleId == param.RoleId && exist.ModuleId == param.ModuleId {
		return ecode.CreateErrExist, "", RoleModuleInfo{}
	}
	insert := dao.RoleModule{
		RoleId: param.RoleId,
		ModuleId: param.ModuleId,
	}
	insertRes := insert.Create()
	if !insertRes {
		return ecode.CreateErr, "", RoleModuleInfo{}
	}
	return ecode.OK, "", TransferRoleModuleInfo(insert)
}

func CreateRoleModuleBatch(param *roleModuleSerializer.CreateRoleModuleBatch) (ecode.Code, string) {
	if len(param.RoleIds) == 0 || len(param.ModuleIds) == 0 {
		return ecode.INVALID_PARAM, ""
	}

	for _, roleId := range param.RoleIds {
		for _, moduleId := range param.ModuleIds {
			CreateRoleModule(&roleModuleSerializer.CreateRoleModule{
				RoleId: roleId,
				ModuleId: moduleId,
			})
		}
	}

	return ecode.OK, ""
}

func UpdateRoleModule(param *roleModuleSerializer.UpdateRoleModule) (ecode.Code, string, RoleModuleInfo) {
	return ecode.FunctionUnAvailable, "", RoleModuleInfo{}
}

func DeleteRoleModule(param *roleModuleSerializer.DeleteRoleModule) (ecode.Code, string) {
	if param.RoleId == 0 || param.ModuleId == 0 {
		return ecode.INVALID_PARAM, ""
	}
	_, exist := dao.GetRoleModuleByRoleIdAndModuleId(param.RoleId, param.ModuleId)
	if exist.RoleId != param.RoleId || exist.ModuleId != param.ModuleId {
		return ecode.DeleteErrNotFound, ""
	}
	deleteRes := exist.Delete()
	if !deleteRes {
		return ecode.DeleteErr, ""
	}
	return ecode.OK, ""
}

func SearchRoleModule(param *roleModuleSerializer.SearchRoleModule) (ecode.Code, string, *RoleModuleList) {
	where := &dao.RoleModule{}
	if param.RoleId > 0 {
		where.RoleId = param.RoleId
	}
	if param.ModuleId > 0 {
		where.ModuleId = param.ModuleId
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountRoleModuleByWhere(where)
	result := &RoleModuleList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]RoleModuleInfo{},
	}
	offset, limit := common.GetOffsetLimit(page)
	if total == 0 || total <= offset {
		return ecode.OK, "", result
	}
	_, list := dao.GetRoleModuleByWhere(where, offset, limit)
	newList := []RoleModuleInfo{}
	for _, item := range list {
		newList = append(newList, TransferRoleModuleInfo(item))
	}
	result.List = &newList
	return ecode.OK, "", result
}

func GetRoleModuleInfo(param *roleModuleSerializer.GetRoleModule) (ecode.Code, string, *RoleModuleInfo) {
	if param.RoleId == 0 || param.ModuleId == 0 {
		return ecode.INVALID_PARAM, "", &RoleModuleInfo{}
	}
	_, exist := dao.GetRoleModuleByRoleIdAndModuleId(param.RoleId, param.ModuleId)
	if exist.RoleId != param.RoleId || exist.ModuleId != param.ModuleId {
		return ecode.GetInfoErr, "", &RoleModuleInfo{}
	}
	info := TransferRoleModuleInfo(*exist)
	return ecode.OK, "", &info
}
