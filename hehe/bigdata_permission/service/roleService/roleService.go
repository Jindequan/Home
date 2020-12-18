package roleService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/roleSerializer"
	"bigdata_permission/service/common"
)

type RoleInfo struct {
	RoleId int   `json:"role_id"`
	Name string `json:"name"`
	Remark string `json:"remark"`
}

type RoleList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]RoleInfo `json:"list"`
}

func TransferRoleInfo(role dao.Role) RoleInfo {
	return RoleInfo{
		RoleId: role.RoleId,
		Name: role.Name,
		Remark: role.Remark,
	}
}

func CreateRole(param *roleSerializer.CreateRole) (ecode.Code, string, RoleInfo) {
	if param.Name == "" {
		return ecode.INVALID_PARAM, "", RoleInfo{}
	}
	exist, roleInfo := dao.GetRoleByName(param.Name)
	if exist && roleInfo.Name == param.Name {
		return ecode.MYSQL_RECORD_EXIST, "已存在该角色", RoleInfo{}
	}
	role := &dao.Role{
		Name: param.Name,
		Remark: param.Remark,
	}
	isSuccess := role.Create()
	if !isSuccess {
		return ecode.MYSQL_ERR, "新增角色失败", RoleInfo{}
	}
	go dictionary.RefreshRole()
	return ecode.OK, "", TransferRoleInfo(*role)
}

func UpdateRole(param *roleSerializer.UpdateRole) (ecode.Code, string, RoleInfo) {
	if param.RoleId < 1 {
		return ecode.INVALID_PARAM, "", RoleInfo{}
	}
	exist, roleInfo := dao.GetRoleById(param.RoleId)
	if !exist || roleInfo.RoleId != param.RoleId {
		return ecode.INVALID_PARAM, "不存在该角色", RoleInfo{}
	}
	if roleInfo.Remark == param.Remark && roleInfo.Name == param.Name {
		return ecode.OK, "", TransferRoleInfo(*roleInfo)
	}

	existNew, roleInfoExist := dao.GetRoleByName(param.Name)
	if existNew && roleInfoExist.Name == param.Name {
		return ecode.MYSQL_RECORD_EXIST, "已存在该角色", RoleInfo{}
	}

	updateData := map[string]interface{}{}

	if roleInfo.Remark != param.Remark {
		updateData["remark"] = param.Remark
	}
	if roleInfo.Name != param.Name {
		updateData["name"] = param.Name
	}

	isSuccess, role := roleInfo.Update(updateData)
	if !isSuccess {
		return ecode.MYSQL_ERR, "更新角色失败", RoleInfo{}
	}
	go dictionary.RefreshRole()
	return ecode.OK, "", TransferRoleInfo(*role)
}

func DeleteRole(param *roleSerializer.DeleteRole) (ecode.Code, string) {
	if param.RoleId < 1 {
		return ecode.INVALID_PARAM, ""
	}
	exist, roleInfo := dao.GetRoleById(param.RoleId)
	if !exist || roleInfo.RoleId != param.RoleId {
		return ecode.OK, "不存在该角色"
	}
	isSuccess := roleInfo.Delete()
	if !isSuccess {
		return ecode.MYSQL_ERR, "删除角色失败"
	}
	go dictionary.RefreshRole()
	return ecode.OK, ""
}

func SearchRole(param *roleSerializer.SearchRole) (ecode.Code, string, *RoleList) {
	where := &dao.Role{}
	if param.RoleId > 0 {
		where.RoleId = param.RoleId
	}
	if param.Name != "" {
		where.Name = param.Name
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountRoleByWhere(where)
	result := &RoleList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]RoleInfo{},
	}
	if total == 0 {
		return ecode.OK, "", result
	}

	offset, limit := common.GetOffsetLimit(page)
	_, list := dao.GetRoleByWhere(where, offset, limit)
	newList := []RoleInfo{}
	for _, role := range list {
		newList = append(newList, TransferRoleInfo(role))
	}
	result.List = &newList

	return ecode.OK, "", result
}

func GetRoleInfo(param *roleSerializer.GetRole) (ecode.Code, *RoleInfo) {
	if param.RoleId == 0 {
		return ecode.INVALID_PARAM, &RoleInfo{}
	}
	isSuccess, role := dao.GetRoleById(param.RoleId)
	if !isSuccess {
		return ecode.MYSQL_ERR, &RoleInfo{}
	}
	if role.RoleId != param.RoleId {
		return ecode.OK, &RoleInfo{}
	}
	info := TransferRoleInfo(*role)
	return ecode.OK, &info
}