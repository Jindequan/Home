package moduleDicService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/moduleDicSerializer"
	"bigdata_permission/service/common"
)

type ModuleDicInfo struct {
	ModuleId int `json:"module_id"`
	Value string `json:"value"`
	FirstLvlId int `json:"first_lvl_id"`
	SecondLvlId int `json:"second_lvl_id"`
	ThirdLvlId int `json:"third_lvl_id"`
	Level int `json:"level"`
	Enable int8 `json:"enable"`
	EnableStr string `json:"enable_str"`
	ResponsePerson int `json:"response_person"`
	ResponsePersonName string `json:"response_person_name"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

type ModuleDicList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]ModuleDicInfo `json:"list"`
}

func TransferModuleDic(moduleDao dao.ModuleDic) ModuleDicInfo {
	times := common.BuildDaoTimeStr(moduleDao.CreatedAt, moduleDao.UpdatedAt, 0)
	level := 0
	if moduleDao.FirstLvlId > 0 {
		level = 1
	}
	if moduleDao.SecondLvlId > 0 {
		level = 2
	}
	if moduleDao.ThirdLvlId > 0 {
		level = 3
	}
	return ModuleDicInfo{
		ModuleId: moduleDao.ModuleId,
		Value: moduleDao.Value,
		FirstLvlId: moduleDao.FirstLvlId,
		SecondLvlId: moduleDao.SecondLvlId,
		ThirdLvlId: moduleDao.ThirdLvlId,
		Level: level,
		CreatedAt: times.CreatedAtStr,
		UpdatedAt: times.UpdatedAtStr,
		DeletedAt: times.DeletedAtStr,
	}
}

func CreateModuleDic(param *moduleDicSerializer.CreateModuleDic) (ecode.Code, string, ModuleDicInfo) {
	moduleDicInfo := ModuleDicInfo{}
	if param.Value == "" {
		return ecode.INVALID_PARAM, "", moduleDicInfo
	}
	firstId, secondId, thirdId, isSuccess := GetLevelIdsByParentId(param.ParentId)
	if !isSuccess {
		return ecode.INVALID_PARAM, "", moduleDicInfo
	}

	if firstId == 0 {
		return insertRootLevel(param.Value)
	}
	if secondId == 0 {
		return insertFirstLevel(firstId, param.Value)
	}
	if thirdId == 0 {
		return insertSecondLevel(firstId, secondId, param.Value)
	}
	return insertThirdLevel(firstId, secondId, thirdId, param.Value)
}

func insertRootLevel(value string) (ecode.Code, string, ModuleDicInfo) {
	_, exist := dao.GetModuleDicByLevelAndValue(0, 0, 0, value)
	if exist.ModuleId > 0 {
		return ecode.CreateErrExist, "", ModuleDicInfo{}
	}
	insertInfo := dao.ModuleDic{
		Value: value,
		Enable: 1,
	}
	isSuccess := insertInfo.Create()
	if !isSuccess {
		return ecode.CreateErr, "", ModuleDicInfo{}
	}
	dictionary.RefreshModule()
	return ecode.OK, "", TransferModuleDic(insertInfo)
}

func insertFirstLevel(firstId int, value string) (ecode.Code, string, ModuleDicInfo) {
	_, exist := dao.GetModuleDicByLevelAndValue(firstId, 0, 0, value)
	if exist.ModuleId > 0 {
		return ecode.CreateErrExist, "", ModuleDicInfo{}
	}
	insertInfo := dao.ModuleDic{
		FirstLvlId: firstId,
		Value: value,
		Enable: 1,
	}
	isSuccess := insertInfo.Create()
	if !isSuccess {
		return ecode.CreateErr, "", ModuleDicInfo{}
	}
	dictionary.RefreshModule()
	return ecode.OK, "", TransferModuleDic(insertInfo)
}

func insertSecondLevel(firstId, secondId int, value string) (ecode.Code, string, ModuleDicInfo) {
	_, exist := dao.GetModuleDicByLevelAndValue(firstId, secondId, 0, value)
	if exist.ModuleId > 0 {
		return ecode.CreateErrExist, "", ModuleDicInfo{}
	}
	insertInfo := dao.ModuleDic{
		FirstLvlId: firstId,
		SecondLvlId: secondId,
		Value: value,
		Enable: 1,
	}
	isSuccess := insertInfo.Create()
	if !isSuccess {
		return ecode.CreateErr, "", ModuleDicInfo{}
	}
	dictionary.RefreshModule()
	return ecode.OK, "", TransferModuleDic(insertInfo)
}

func insertThirdLevel(firstId, secondId, thirdId int, value string) (ecode.Code, string, ModuleDicInfo) {
	_, exist := dao.GetModuleDicByLevelAndValue(firstId, secondId, thirdId, value)
	if exist.ModuleId > 0 {
		return ecode.CreateErrExist, "", ModuleDicInfo{}
	}
	insertInfo := dao.ModuleDic{
		FirstLvlId: firstId,
		SecondLvlId: secondId,
		ThirdLvlId: thirdId,
		Value: value,
		Enable: 1,
	}
	isSuccess := insertInfo.Create()
	if !isSuccess {
		return ecode.CreateErr, "", ModuleDicInfo{}
	}
	dictionary.RefreshModule()
	return ecode.OK, "", TransferModuleDic(insertInfo)
}

func GetLevelIdsByParentId(parentId int) (int, int, int, bool) {
	firstId, secondId, thirdId := 0, 0, 0
	if parentId == 0 {
		return firstId, secondId, thirdId, true
	}
	res, parent := dao.GetModuleDicById(parentId)
	if !res || parent.ModuleId != parentId {
		return firstId, secondId, thirdId, false
	}
	if parent.ThirdLvlId != 0 {//已达到最多4级
		return firstId, secondId, thirdId, false
	}
	if parent.SecondLvlId != 0 {
		return parent.FirstLvlId, parent.SecondLvlId, parent.ModuleId, true
	}
	if parent.FirstLvlId != 0 {
		return parent.FirstLvlId, parent.ModuleId, 0, true
	}
	return parent.ModuleId, 0, 0, true
}

func GetChildModule(moduleId int) (ecode.Code, string, []ModuleDicInfo){
	_, list := dao.GetChildModuleByModuleId(moduleId)
	newList := []ModuleDicInfo{}
	for _, item := range list {
		newList = append(newList, TransferModuleDic(item))
	}
	return ecode.OK, "", newList
}

func UpdateModuleDic(param *moduleDicSerializer.UpdateModuleDic) (ecode.Code, string, ModuleDicInfo) {
	if param.ModuleId == 0 {
		return ecode.INVALID_PARAM, "", ModuleDicInfo{}
	}
	_, moduleInfo := dao.GetModuleDicById(param.ModuleId)
	if moduleInfo.ModuleId != param.ModuleId {
		return ecode.UpdateErrNotFound, "", ModuleDicInfo{}
	}
	firstId, secondId, thirdId, isSuccess := GetLevelIdsByParentId(param.ParentId)
	if !isSuccess {
		return ecode.INVALID_PARAM, "", ModuleDicInfo{}
	}
	hasChange := false
	levelChange := false
	updateInfo := map[string]interface{}{}
	if moduleInfo.Value != param.Value {
		hasChange = true
		updateInfo["value"] = param.Value
	}
	if moduleInfo.Enable == 1 && param.Enable == 0 {
		hasChange = true
		updateInfo["enable"] = 0
	}
	if moduleInfo.Enable == 0 && param.Enable == 1 {
		hasChange = true
		updateInfo["enable"] = 1
	}
	if moduleInfo.ResponsePerson != param.ResponsePerson {
		hasChange = true
		updateInfo["response_person"] = param.ResponsePerson
	}
	if moduleInfo.FirstLvlId != firstId {
		hasChange = true
		levelChange =true
		updateInfo["first_lvl_id"] = firstId
	}
	if moduleInfo.SecondLvlId != secondId {
		hasChange = true
		levelChange = true
		updateInfo["second_lvl_id"] = secondId
	}
	if moduleInfo.ThirdLvlId != thirdId {
		hasChange = true
		levelChange = true
		updateInfo["third_lvl_id"] = thirdId
	}
	if !hasChange {
		return ecode.OK, "没有需要的改动", TransferModuleDic(*moduleInfo)
	}
	if levelChange {
		check, childrenNum := dao.CountChildModuleByModuleId(moduleInfo.ModuleId)
		if !check {
			return ecode.UpdateErr, "数据库出错，暂时无法更新", TransferModuleDic(*moduleInfo)
		}
		if childrenNum > 0 {
			return ecode.UpdateErr, "已经存在子模块，无法修改层级结构", TransferModuleDic(*moduleInfo)
		}
	}

	updateRes, moduleDicNew := moduleInfo.Update(updateInfo)
	if !updateRes {
		return ecode.UpdateErr, "", TransferModuleDic(*moduleInfo)
	}
	go dictionary.RefreshModule()
	return ecode.OK, "", TransferModuleDic(*moduleDicNew)
}

func DeleteModuleDic(param *moduleDicSerializer.DeleteModuleDic) (ecode.Code, string) {
	if param.ModuleId < 1 {
		return ecode.INVALID_PARAM, ""
	}
	exist, moduleInfo := dao.GetModuleDicById(param.ModuleId)
	if !exist || moduleInfo.ModuleId != param.ModuleId {
		return ecode.OK, "不存在该模块信息"
	}
	check, childrenNum := dao.CountChildModuleByModuleId(moduleInfo.ModuleId)
	if !check {
		return ecode.DeleteErr, "数据库出错，暂时无法删除"
	}
	if childrenNum > 0 {
		return ecode.DeleteErr, "已经存在子模块，无法删除"
	}
	isSuccess := moduleInfo.Delete()
	if !isSuccess {
		return ecode.DeleteErr, "删除模块信息失败"
	}
	dao.DeleteUserModuleByModuleId(param.ModuleId)
	dao.DeleteRoleModuleByModuleId(param.ModuleId)
	dao.DeleteModuleInterfaceByModuleId(param.ModuleId)
	go dictionary.RefreshModule()
	return ecode.OK, ""
}

func SearchModuleDic(param *moduleDicSerializer.SearchModuleDic) (ecode.Code, string, *ModuleDicList) {
	where := &dao.ModuleDicForQuery{}
	if param.ModuleId > 0 {
		where.ModuleId = param.ModuleId
	}
	if param.Enable &&  param.Disable{
		return ecode.OK, "", &ModuleDicList{}
	}
	if param.Enable {
		where.Enable = 1
	}
	if param.Disable {
		where.EnableZero = true
	}
	if param.ResponsePerson > 0 {
		where.ResponsePerson = param.ResponsePerson
	}
	if param.ResponsePersonNull {
		where.ResponsePersonZero = true
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountModuleDicByWhere(where)
	pageInfo := common.BuildPageInfo(total, page.PageSize, page.PageIndex)
	result := &ModuleDicList{
		PageInfo: pageInfo,
		List: &[]ModuleDicInfo{},
	}
	offset, limit := common.GetOffsetLimit(page)
	if total == 0 || total <= offset{
		return ecode.OK, "", result
	}
	_, list := dao.GetModuleDicByWhere(where, offset, limit)
	newList := []ModuleDicInfo{}
	for _, item := range list {
		newList = append(newList, TransferModuleDic(item))
	}
	result.List = &newList
	return ecode.OK, "", result
}

func GetModuleDicInfo(param *moduleDicSerializer.GetModuleDic) (ecode.Code, *ModuleDicInfo) {
	if param.ModuleId == 0 {
		return ecode.INVALID_PARAM, &ModuleDicInfo{}
	}
	isSuccess, moduleDicInfo := dao.GetModuleDicById(param.ModuleId)
	if !isSuccess {
		return ecode.GetInfoErr, &ModuleDicInfo{}
	}
	if moduleDicInfo.ModuleId != param.ModuleId {
		return ecode.OK, &ModuleDicInfo{}
	}
	info := TransferModuleDic(*moduleDicInfo)
	return ecode.OK, &info
}