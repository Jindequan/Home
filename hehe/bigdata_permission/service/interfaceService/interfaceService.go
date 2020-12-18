package interfaceService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/interfaceSerializer"
	"bigdata_permission/service/common"
	"strings"
)

type InterfaceInfo struct {
	InterfaceId int `json:"interface_id"`
	Type int8 `json:"type"`
	ShowName string `json:"show_name"`
	Path string `json:"path"`
	Remark string `json:"remark"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type InterfaceList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]InterfaceInfo `json:"list"`
}

func TransferInterfaceInfo(interfaceDao dao.Interface) InterfaceInfo {
	times := common.BuildDaoTimeStr(interfaceDao.CreatedAt, interfaceDao.UpdatedAt, 0)
	return InterfaceInfo{
		InterfaceId: interfaceDao.InterfaceId,
		Type: interfaceDao.Type,
		ShowName: interfaceDao.ShowName,
		Remark: interfaceDao.Remark,
		CreatedAt: times.CreatedAtStr,
		UpdatedAt: times.UpdatedAtStr,
		DeletedAt: times.DeletedAtStr,
	}
}

func CreateInterface(param *interfaceSerializer.CreateInterface) (ecode.Code, string, InterfaceInfo) {
	param.Path = "/" + strings.Trim(param.Path, "/")
	if param.Type == 0 || param.ShowName == "" {
		return ecode.INVALID_PARAM, "参数错误", InterfaceInfo{}
	}
	if _, ok := dao.InterfaceTypeMap[param.Type]; !ok {
		return ecode.INVALID_PARAM, "不存在的接口类型", InterfaceInfo{}
	}
	_, exist := dao.GetInterfaceByTypeAndPath(param.Type, param.Path)
	if exist.InterfaceId > 0 {
		return ecode.CreateErrExist, "该接口已存在", InterfaceInfo{}
	}
	createInfo := dao.Interface{
		Type: param.Type,
		ShowName: param.ShowName,
		Path: param.Path,
		Remark: param.Remark,
	}
	if !createInfo.Create() {
		return ecode.MYSQL_ERR, "新增接口失败", InterfaceInfo{}
	}
	go dictionary.RefreshInterface()
	return ecode.OK, "", TransferInterfaceInfo(createInfo)
}

func UpdateInterface(param *interfaceSerializer.UpdateInterface) (ecode.Code, string, InterfaceInfo) {
	param.Path = "/" + strings.Trim(param.Path, "/")
	if param.InterfaceId == 0 || param.Type == 0 || param.ShowName == "" {
		return ecode.INVALID_PARAM, "参数错误", InterfaceInfo{}
	}
	if _, ok := dao.InterfaceTypeMap[param.Type]; !ok {
		return ecode.INVALID_PARAM, "不存在的接口类型", InterfaceInfo{}
	}
	_, interfaceInfo := dao.GetInterfaceById(param.InterfaceId)
	if interfaceInfo.InterfaceId != param.InterfaceId {
		return ecode.INVALID_PARAM, "该接口不存在", InterfaceInfo{}
	}
	hasChange := false
	updateInfo := map[string]interface{}{}
	if interfaceInfo.Type != param.Type {
		hasChange = true
		updateInfo["type"] = param.Type
	}
	if interfaceInfo.ShowName != param.ShowName {
		hasChange = true
		updateInfo["show_name"] = param.ShowName
	}
	if interfaceInfo.Path != param.Path {
		hasChange = true
		updateInfo["path"] = param.Path
	}
	if interfaceInfo.Remark != param.Remark {
		hasChange = true
		updateInfo["remark"] = param.Remark
	}
	if !hasChange {
		return ecode.OK, "没有需要的改动", TransferInterfaceInfo(*interfaceInfo)
	}
	newInterfaceExist, newInterface := dao.GetInterfaceByTypeAndPath(param.Type, param.Path)
	if newInterfaceExist && newInterface.InterfaceId != interfaceInfo.InterfaceId {
		return ecode.UpdateErr, "更新的数据已存在", TransferInterfaceInfo(*interfaceInfo)
	}

	updateRes, interfaceInfoNew := interfaceInfo.Update(updateInfo)
	if !updateRes {
		return ecode.UpdateErr, "更新接口信息失败", InterfaceInfo{}
	}
	go dictionary.RefreshInterface()
	return ecode.OK, "", TransferInterfaceInfo(*interfaceInfoNew)
}

func DeleteInterface(param *interfaceSerializer.DeleteInterface) (ecode.Code, string) {
	if param.InterfaceId < 1 {
		return ecode.INVALID_PARAM, ""
	}
	exist, interfaceInfo := dao.GetInterfaceById(param.InterfaceId)
	if !exist || interfaceInfo.InterfaceId != param.InterfaceId {
		return ecode.OK, "不存在该接口信息"
	}
	isSuccess := interfaceInfo.Delete()
	if !isSuccess {
		return ecode.DeleteErr, "删除接口信息失败"
	}
	dao.DeleteModuleInterfaceByInterfaceId(param.InterfaceId)
	go dictionary.RefreshInterface()
	return ecode.OK, ""
}

func SearchInterface(param *interfaceSerializer.SearchInterface) (ecode.Code, string, *InterfaceList) {
	where := &dao.Interface{}
	if param.InterfaceId > 0 {
		where.InterfaceId = param.InterfaceId
	}
	if param.Type > 0 {
		where.Type = param.Type
	}
	if param.Path != "" {
		where.Path = param.Path
	}
	if param.ShowName != "" {
		where.ShowName = param.ShowName
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountInterfaceByWhere(where)
	result := &InterfaceList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]InterfaceInfo{},
	}

	offset, limit := common.GetOffsetLimit(page)
	if total == 0 || total <= offset{
		return ecode.OK, "", result
	}

	_, list := dao.GetInterfaceByWhere(where, offset, limit)
	newList := []InterfaceInfo{}
	for _, interfaceInfo := range list {
		newList = append(newList, TransferInterfaceInfo(interfaceInfo))
	}
	result.List = &newList

	return ecode.OK, "", result
}

func GetInterfaceInfo(param *interfaceSerializer.GetInterface) (ecode.Code, *InterfaceInfo) {
	if param.InterfaceId == 0 {
		return ecode.INVALID_PARAM, &InterfaceInfo{}
	}
	isSuccess, interfaceInfo := dao.GetInterfaceById(param.InterfaceId)
	if !isSuccess {
		return ecode.GetInfoErr, &InterfaceInfo{}
	}
	if interfaceInfo.InterfaceId != param.InterfaceId {
		return ecode.OK, &InterfaceInfo{}
	}
	info := TransferInterfaceInfo(*interfaceInfo)
	return ecode.OK, &info
}