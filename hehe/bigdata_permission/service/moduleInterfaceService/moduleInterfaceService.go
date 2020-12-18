package moduleInterfaceService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/moduleInterfaceSerializer"
	"bigdata_permission/service/common"
)

type ModuleInterfaceInfo struct {
	ModuleId int `json:"module_id"`
	InterfaceId int `json:"interface_id"`

	ModuleName string `json:"module_name"`
	InterfaceName string `json:"interface_name"`
}

type ModuleInterfaceList struct {
	PageInfo *serializer.PageInfo
	List *[]ModuleInterfaceInfo
}

func TransferModuleInterfaceInfo(moduleInterfaceDao dao.ModuleInterface) ModuleInterfaceInfo {
	return ModuleInterfaceInfo{
		ModuleId: moduleInterfaceDao.ModuleId,
		InterfaceId: moduleInterfaceDao.InterfaceId,

		ModuleName: dictionary.GetStringValue(dictionary.ModuleList, moduleInterfaceDao.ModuleId),
		InterfaceName: dictionary.GetStringValue(dictionary.InterfaceList, moduleInterfaceDao.InterfaceId),
	}
}

func CreateModuleInterface(param *moduleInterfaceSerializer.CreateModuleInterface) (ecode.Code, string, ModuleInterfaceInfo) {
	if param.ModuleId == 0 || param.InterfaceId == 0 {
		return ecode.INVALID_PARAM, "", ModuleInterfaceInfo{}
	}
	_, module := dao.GetModuleDicById(param.ModuleId)
	if module.ModuleId != param.ModuleId {
		return ecode.INVALID_PARAM, "不存在的模块", ModuleInterfaceInfo{}
	}
	_, interfaceDao := dao.GetInterfaceById(param.InterfaceId)
	if interfaceDao.InterfaceId != param.InterfaceId {
		return ecode.INVALID_PARAM, "不存在的接口", ModuleInterfaceInfo{}
	}
	checkSuccess, childNum := dao.CountChildModuleByModuleId(param.ModuleId)
	if !checkSuccess {
		return ecode.CreateErr, "", ModuleInterfaceInfo{}
	}
	if childNum > 0 {
		return ecode.NotLeafModuleCannotRelateInterface, "", ModuleInterfaceInfo{}
	}
	_, exist := dao.GetModuleInterfaceByModuleIdAndInterfaceId(param.ModuleId, param.InterfaceId)
	if exist.ModuleId == param.ModuleId && exist.InterfaceId == param.InterfaceId {
		return ecode.CreateErrExist, "", ModuleInterfaceInfo{}
	}
	insert := dao.ModuleInterface{
		ModuleId: param.ModuleId,
		InterfaceId: param.InterfaceId,
	}
	insertRes := insert.Create()
	if !insertRes {
		return ecode.CreateErr, "", ModuleInterfaceInfo{}
	}
	return ecode.OK, "", TransferModuleInterfaceInfo(insert)
}

func UpdateModuleInterface(param *moduleInterfaceSerializer.UpdateModuleInterface) (ecode.Code, string, ModuleInterfaceInfo) {
	return ecode.FunctionUnAvailable, "", ModuleInterfaceInfo{}
}

func DeleteModuleInterface(param *moduleInterfaceSerializer.DeleteModuleInterface) (ecode.Code, string) {
	if param.ModuleId == 0 || param.InterfaceId == 0 {
		return ecode.INVALID_PARAM, ""
	}
	_, exist := dao.GetModuleInterfaceByModuleIdAndInterfaceId(param.ModuleId, param.InterfaceId)
	if exist.ModuleId != param.ModuleId || exist.InterfaceId != param.InterfaceId {
		return ecode.DeleteErrNotFound, ""
	}
	deleteRes := exist.Delete()
	if !deleteRes {
		return ecode.DeleteErr, ""
	}
	return ecode.OK, ""
}

func SearchModuleInterface(param *moduleInterfaceSerializer.SearchModuleInterface) (ecode.Code, string, *ModuleInterfaceList) {
	where := &dao.ModuleInterface{}
	if param.ModuleId > 0 {
		where.ModuleId = param.ModuleId
	}
	if param.InterfaceId > 0 {
		where.InterfaceId = param.InterfaceId
	}
	page := common.GetPage(param.Page)
	offset, limit := common.GetOffsetLimit(page)
	_, total := dao.CountModuleInterfaceByWhere(where)
	result := &ModuleInterfaceList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]ModuleInterfaceInfo{},
	}
	if total == 0 || offset >= total {
		return ecode.OK, "", result
	}
	_, list := dao.GetModuleInterfaceByWhere(where, offset, limit)
	newList := []ModuleInterfaceInfo{}
	for _, item := range list {
		newList = append(newList, TransferModuleInterfaceInfo(item))
	}
	result.List = &newList
	return ecode.OK, "", result
}

func GetModuleInterfaceInfo(param *moduleInterfaceSerializer.GetModuleInterface) (ecode.Code, *ModuleInterfaceInfo) {
	if param.ModuleId == 0 || param.InterfaceId == 0 {
		return ecode.INVALID_PARAM, &ModuleInterfaceInfo{}
	}
	_, moduleInterface := dao.GetModuleInterfaceByModuleIdAndInterfaceId(param.ModuleId, param.InterfaceId)
	if moduleInterface.ModuleId != param.ModuleId || moduleInterface.InterfaceId != param.InterfaceId {
		return ecode.INVALID_PARAM, &ModuleInterfaceInfo{}
	}
	info := TransferModuleInterfaceInfo(*moduleInterface)
	return ecode.OK, &info
}

type InterfaceInfo struct {
	Key int `json:"key"`
	Type int8 `json:"type"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type ModuleInterfaceTree struct {
	Key int `json:"key"`
	Name string `json:"name"`
	ParentId int `json:"parent_id"`
	InterfaceList []InterfaceInfo `json:"interface_list"`
	Children []ModuleInterfaceTree `json:"children"`
}

//获取任何层级的模块下的子模块与接口
func GetChildModuleInterface(moduleId int) ModuleInterfaceTree {
	_, moduleInfo := dao.GetModuleDicById(moduleId)
	if moduleInfo.ModuleId != moduleId {
		return ModuleInterfaceTree{}
	}
	_, list := dao.GetRelateModuleByModuleId(moduleId)
	moduleIds := []int{}
	for _, module := range list {
		moduleIds = append(moduleIds, module.ModuleId)
	}
	if len(moduleIds) == 0 {
		return ModuleInterfaceTree{}
	}
	isSuccess, mapping := dao.GetModuleInterfaceByModuleIds(moduleIds)
	if !isSuccess {
		return ModuleInterfaceTree{}
	}

	return GetTreeByModuleAndInterface(list, mapping)
}

func GetTreeByModuleAndInterface(moduleList []dao.ModuleDic, mapList *[]dao.ModuleInterface) ModuleInterfaceTree {
	interfaceIds := map[int]int{}
	interfaceIdList := []int{}

	mapListByModule := map[int]map[int]int{}
	//模块接口
	for _, v := range *mapList {
		if _, ok := interfaceIds[v.InterfaceId]; !ok {
			interfaceIdList = append(interfaceIdList, v.InterfaceId)
			interfaceIds[v.InterfaceId] = v.InterfaceId
		}
		if _, ok := mapListByModule[v.ModuleId]; !ok {
			mapListByModule[v.ModuleId] = map[int]int{}
		}
		if _, ok := mapListByModule[v.ModuleId][v.InterfaceId]; !ok {
			mapListByModule[v.ModuleId][v.InterfaceId] = v.InterfaceId
		}
	}

	//接口
	_, interfaceList := dao.GetInterfaceByIds(interfaceIdList)
	interfaceById := map[int]dao.Interface{}
	for _, v := range interfaceList {
		interfaceById[v.InterfaceId] = v
	}

	//模块与接口信息绑定
	moduleInterfaceList := map[int][]InterfaceInfo{}
	for moduleId, interfaces := range mapListByModule {
		moduleInterfaceList[moduleId] = []InterfaceInfo{}
		if len(interfaces) == 0 {
			continue
		}
		for _, interfaceId := range interfaces {
			interfaceInfo, ok := interfaceById[interfaceId]
			if !ok {
				continue
			}
			moduleInterfaceList[moduleId] = append(moduleInterfaceList[moduleId], InterfaceInfo{
				Key: interfaceId,
				Name: interfaceInfo.ShowName,
				Type: interfaceInfo.Type,
				Path: interfaceInfo.Path,
			})
		}
	}

	nodeList := []ModuleInterfaceTree{}
	for _, moduleInfo := range moduleList {
		parentId := 0
		if moduleInfo.ThirdLvlId > 0 {
			parentId = moduleInfo.ThirdLvlId
		} else if moduleInfo.SecondLvlId > 0 {
			parentId = moduleInfo.SecondLvlId
		} else if moduleInfo.FirstLvlId > 0 {
			parentId = moduleInfo.FirstLvlId
		}
		interfaceList, _ := moduleInterfaceList[moduleInfo.ModuleId]

		nodeList = append(nodeList, ModuleInterfaceTree{
			Key: moduleInfo.ModuleId,
			Name: moduleInfo.Value,
			ParentId: parentId,
			InterfaceList: interfaceList,
		})
	}

	moduleTree := BuildModuleInterfaceTree(nodeList)
	return moduleTree
}

func BuildModuleInterfaceTree(nodeList []ModuleInterfaceTree) ModuleInterfaceTree {
	if len(nodeList) == 0 {
		return ModuleInterfaceTree{}
	}
	rootModule := GetRootParentModule(nodeList, nodeList[0])
	return SetModuleInterfaceChildren(nodeList, rootModule)
}

func GetRootParentModule(nodeList []ModuleInterfaceTree, nodeInfo ModuleInterfaceTree) ModuleInterfaceTree {
	if nodeInfo.ParentId == 0 {
		return nodeInfo
	}
	for _, node := range nodeList {
		if node.Key == nodeInfo.ParentId {
			nodeInfo = node
			return GetRootParentModule(nodeList, nodeInfo)
		}
	}
	return nodeInfo
}

func SetModuleInterfaceChildren(nodeList []ModuleInterfaceTree, baseNode ModuleInterfaceTree) ModuleInterfaceTree {
	baseNode.Children = GetModuleInterfaceChildren(nodeList, baseNode.Key)
	if len(baseNode.Children) == 0 {
		return baseNode
	}
	for k, node := range baseNode.Children {
		baseNode.Children[k] = SetModuleInterfaceChildren(nodeList, node)
	}
	return baseNode
}

func GetModuleInterfaceChildren(nodeList []ModuleInterfaceTree, parentId int) []ModuleInterfaceTree {
	children := []ModuleInterfaceTree{}
	for _, node := range nodeList {
		if node.ParentId == parentId {
			children = append(children, node)
		}
	}
	return children
}