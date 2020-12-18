package dictionary

import (
	"bigdata_permission/dao"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

type GenericItem struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

type DBItem struct {
	Id          interface{}
	Value       string
	ThirdLvlId  interface{}
	SecondLvlId interface{}
	FirstLvlId  interface{}
}

type GenericTreeNode struct {
	GenericItem
	ParentId interface{}
	SubItem  []*GenericTreeNode
}

var (
	//查询dao层数据分页大小
	PageSize = 10000
	//list and tree
	RoleList      = []*GenericItem{}
	ModuleList    = []*GenericItem{}
	ModuleTree    = []*GenericTreeNode{}
	InterfaceList = []*GenericItem{}
	//ModuleInterfaceTree =

	//refresh lock
	RoleLock      sync.Mutex
	ModuleLock    sync.Mutex
	InterfaceLock sync.Mutex
)
//map
var (
	//module list
	ModuleIdDaoMap = map[int]dao.ModuleDic{}
)

type jsonDicConfig struct {
}

func Init() {
	//parseJsonConfig()
	querySqlDic()
}

func parseJsonConfig() {
	dicConfigContent, err := ioutil.ReadFile("./settings/config.json")
	if err != nil {
		panic(err)
		return
	}

	var jsonDicConfig jsonDicConfig
	err = json.Unmarshal(dicConfigContent, &jsonDicConfig)
	if err != nil {
		panic(err)
		return
	}
}

//查询dao层默认使用单次PageSize条限制
func querySqlDic() {
	queryRole()
	queryInterface()
	queryModule()
}

func queryRole() {
	offset, limit := 0, PageSize
	list := []*GenericItem{}
	for {
		_, daos := dao.GetRoleByWhere(&dao.Role{}, offset, limit)
		offset += limit
		if len(daos) == 0 {
			break
		}
		for _, item := range daos {
			list = append(list, &GenericItem{
				Key:   item.RoleId,
				Value: item.Name,
			})
		}
	}
	ModuleList = list
}

func queryInterface() {
	offset, limit := 0, PageSize
	list := []*GenericItem{}
	for {
		_, daos := dao.GetInterfaceByWhere(&dao.Interface{}, offset, limit)
		offset += limit
		if len(daos) == 0 {
			break
		}
		for _, item := range daos {
			list = append(list, &GenericItem{
				Key:   item.InterfaceId,
				Value: item.ShowName,
			})
		}
	}
	InterfaceList = list
}

func queryModule() {
	offset, limit := 0, PageSize
	dbList := []*DBItem{}
	moduleIdDaoMap := map[int]dao.ModuleDic{}
	for {
		_, daos := dao.GetModuleDicByWhere(&dao.ModuleDicForQuery{}, offset, limit)
		offset += limit
		if len(daos) == 0 {
			break
		}
		for _, item := range daos {
			dbItem := DBItem{
				Id:          item.ModuleId,
				Value:       item.Value,
				ThirdLvlId:  item.ThirdLvlId,
				SecondLvlId: item.SecondLvlId,
				FirstLvlId:  item.FirstLvlId,
			}
			byteArr, _ := json.Marshal(dbItem)

			res := &DBItem{}
			err := json.Unmarshal(byteArr, res)
			if err != nil {
				panic(err)
			}
			dbList = append(dbList, res)
			moduleIdDaoMap[item.ModuleId] = item
		}
	}
	list, tree := transform(dbList)
	ModuleList = list
	ModuleTree = getTree(tree)
	ModuleIdDaoMap = moduleIdDaoMap
}

func RefreshRole() {
	RoleLock.Lock()
	queryRole()
	RoleLock.Unlock()
}

func RefreshInterface() {
	InterfaceLock.Lock()
	queryInterface()
	InterfaceLock.Unlock()
}

func RefreshModule() {
	ModuleLock.Lock()
	queryModule()
	ModuleLock.Unlock()
}

/********
****新增list维护，与树结构修正****
********/
func AppendModule(item dao.ModuleDic) {
	if item.ModuleId == 0 {
		return
	}
	ModuleList = append(ModuleList, &GenericItem{
		Key:   item.ModuleId,
		Value: item.Value,
	})
	ModuleTree = appendTree(&DBItem{
		Id:          item.ModuleId,
		Value:       item.Value,
		FirstLvlId:  item.FirstLvlId,
		SecondLvlId: item.SecondLvlId,
		ThirdLvlId:  item.ThirdLvlId,
	}, ModuleTree)
}

func DeleteModule(item dao.ModuleDic) {
	if item.ModuleId == 0 {
		return
	}
	ModuleList = deleteNodeByKey(item.ModuleId, ModuleList)
	ModuleTree = deleteTreeNode(&DBItem{
		Id:          item.ModuleId,
		Value:       item.Value,
		FirstLvlId:  item.FirstLvlId,
		SecondLvlId: item.SecondLvlId,
		ThirdLvlId:  item.ThirdLvlId,
	}, ModuleTree)
}

func AppendInterface(item dao.Interface) {
	if item.InterfaceId == 0 {
		return
	}
	InterfaceList = append(InterfaceList, &GenericItem{
		Key:   item.InterfaceId,
		Value: item.ShowName,
	})
}

func DeleteInterface(id int) {
	if id == 0 {
		return
	}
	InterfaceList = deleteNodeByKey(id, InterfaceList)
}

func AppendRole(role dao.Role) {
	if role.RoleId == 0 {
		return
	}
	RoleList = append(RoleList, &GenericItem{
		Key:   role.RoleId,
		Value: role.Name,
	})
}

func DeleteRole(id int) {
	if id == 0 {
		return
	}
	RoleList = deleteNodeByKey(id, RoleList)
}

func appendTree(item *DBItem, tree []*GenericTreeNode) []*GenericTreeNode {
	//root
	parentIds := []int{}
	if item.FirstLvlId.(int) == 0 {
		parentIds = append(parentIds, 0)
	}
	if item.FirstLvlId.(int) > 0 {
		parentIds = append(parentIds, item.FirstLvlId.(int))
	}
	if item.SecondLvlId.(int) > 0 {
		parentIds = append(parentIds, item.SecondLvlId.(int))
	}
	if item.ThirdLvlId.(int) > 0 {
		parentIds = append(parentIds, item.ThirdLvlId.(int))
	}

	kv := &GenericItem{
		Key:   item.Id,
		Value: item.Value,
	}

	return locateAndAppend(parentIds, tree, kv)
}

func deleteTreeNode(item *DBItem, tree []*GenericTreeNode) []*GenericTreeNode {
	parentIds := []int{}
	if item.FirstLvlId.(int) > 0 {
		parentIds = append(parentIds, item.FirstLvlId.(int))
	}
	if item.SecondLvlId.(int) > 0 {
		parentIds = append(parentIds, item.SecondLvlId.(int))
	}
	if item.ThirdLvlId.(int) > 0 {
		parentIds = append(parentIds, item.ThirdLvlId.(int))
	}
	parentIds = append(parentIds, item.Id.(int))

	return locateAndDelete(parentIds, tree)
}

//搜寻所归属的节点并追加该元素
func locateAndAppend(parentIds []int, tree []*GenericTreeNode, item *GenericItem) []*GenericTreeNode {
	//初始化
	if tree == nil {
		tree = []*GenericTreeNode{}
	}
	//无父集id，无处可append
	if len(parentIds) == 0 {
		return tree
	}
	//根节点在此直接append返回
	if len(parentIds) == 1 && parentIds[0] == 0 {
		return append(tree, &GenericTreeNode{
			GenericItem: *item,
			ParentId:    parentIds[0],
		})
	}
	for k, node := range tree {
		//寻找满足条件节点
		if node.Key.(int) != parentIds[0] {
			continue
		}
		//归属于此节点，append
		if len(parentIds) == 1 {
			node.SubItem = append(node.SubItem, &GenericTreeNode{
				GenericItem: *item,
				ParentId:    parentIds[0],
			})
			tree[k] = node
			return tree
		}
		//继续追踪下一个节点
		parentIds = parentIds[1:]
		node.SubItem = locateAndAppend(parentIds, node.SubItem, item)
	}
	return tree
}

//搜寻所属节点并删除该元素
func locateAndDelete(parentIds []int, tree []*GenericTreeNode) []*GenericTreeNode {
	//无需删除
	if tree == nil || len(parentIds) == 0 {
		return tree
	}
	//非根节点
	if len(parentIds) == 1 {
		//叶子节点必须指定id
		if parentIds[0] == 0 {
			return tree
		}
		//重组
		newTree := []*GenericTreeNode{}
		for _, node := range tree {
			if node.Key.(int) != parentIds[0] {
				newTree = append(newTree, node)
			}
		}
		return newTree
	}
	//len > 1 需要深入下一层
	for k, node := range tree {
		//寻找满足条件节点
		if node.Key.(int) != parentIds[0] {
			continue
		}
		//继续追踪下一个节点
		parentIds = parentIds[1:]
		node.SubItem = locateAndDelete(parentIds, node.SubItem)
		tree[k] = node
		break
	}
	return tree
}

func deleteNodeByKey(key int, nodeList []*GenericItem) []*GenericItem {
	if key <= 0 {
		return nodeList
	}
	newNodeList := []*GenericItem{}
	for _, node := range nodeList {
		if node.Key.(int) == key {
			continue
		}
		newNodeList = append(newNodeList, node)
	}
	return newNodeList
}

// 会将float64转成int类型
func transform(dbItemList []*DBItem) ([]*GenericItem, []*GenericTreeNode) {
	var itemList []*GenericItem
	var nodeList []*GenericTreeNode
	for _, dbItem := range dbItemList {
		// jsonUnmarshal 把数字赋给interface{}时，默认使用float64类型
		// dbItem的属性id（和lvlId）属于int或者string两者之一
		if _, ok := dbItem.Id.(float64); ok {
			dbItem.Id = int(dbItem.Id.(float64))
		}
		if _, ok := dbItem.FirstLvlId.(float64); ok {
			dbItem.FirstLvlId = int(dbItem.FirstLvlId.(float64))
		}
		if _, ok := dbItem.SecondLvlId.(float64); ok {
			dbItem.SecondLvlId = int(dbItem.SecondLvlId.(float64))
		}
		if _, ok := dbItem.ThirdLvlId.(float64); ok {
			dbItem.ThirdLvlId = int(dbItem.ThirdLvlId.(float64))
		}
		//存储字典kvmap
		item := &GenericItem{
			Key:   dbItem.Id,
			Value: dbItem.Value,
		}
		itemList = append(itemList, item)

		//计算parentid
		// 没有层级的dao可以不使用transform方法，参考roleDic
		zeroValue := zeroValueOf(dbItem.FirstLvlId)
		parentId := zeroValue
		if !isZeroValue(dbItem.FirstLvlId) {
			parentId = dbItem.FirstLvlId
		}
		if !isZeroValue(dbItem.SecondLvlId) {
			parentId = dbItem.SecondLvlId
		}
		if !isZeroValue(dbItem.ThirdLvlId) {
			parentId = dbItem.ThirdLvlId
		}

		//存储标准格式的节点列表
		nodeList = append(nodeList, &GenericTreeNode{
			GenericItem: *item,
			ParentId:    parentId,
		})
	}
	return itemList, nodeList
}

func isZeroValue(id interface{}) bool {
	if id == "" || id == float64(0) || id == int(0) || id == nil {
		return true
	}
	return false
}

// 默认只有int(0)或者""两种零值
func zeroValueOf(id interface{}) interface{} {
	if _, ok := id.(int); ok {
		return int(0)
	} else if _, ok := id.(string); ok {
		return ""
	} else {
		return nil
	}
}

func getTree(nodeList []*GenericTreeNode) []*GenericTreeNode {
	if len(nodeList) == 0 {
		return []*GenericTreeNode{}
	}
	firstNode := nodeList[0]
	zeroValueOfKey := zeroValueOf(firstNode.Key)
	rootNode := &GenericTreeNode{
		GenericItem: GenericItem{
			Key:   zeroValueOfKey,
			Value: "root",
		},
		ParentId: zeroValueOfKey,
		SubItem:  nil,
	}

	makeTree(nodeList, rootNode)

	return rootNode.SubItem
}

func makeTree(allNode []*GenericTreeNode, node *GenericTreeNode) {
	if childs, has := hasChild(allNode, node); has {
		node.SubItem = append(node.SubItem, childs[0:]...)
		for _, v := range childs {
			if _, has = hasChild(allNode, v); has {
				makeTree(allNode, v)
			}
		}
	}
}

func hasChild(allNode []*GenericTreeNode, node *GenericTreeNode) (childs []*GenericTreeNode, yes bool) {
	for _, v := range allNode {
		if v.ParentId == node.Key {
			childs = append(childs, v)
		}
	}

	if childs != nil {
		yes = true
	} else {
		yes = false
	}

	return
}

func GetUintCodeOf(originStr string, itemDict []*GenericItem) (int, bool) {
	resCode := int(0)
	found := false

	for _, item := range itemDict {
		if strings.ToLower(originStr) == strings.ToLower(item.Value.(string)) {
			if _, ok := item.Key.(int); ok {
				resCode = int(item.Key.(int))
			} else if _, ok := item.Key.(float64); ok {
				resCode = int(item.Key.(float64))
			}
			found = true
			break
		}
	}
	return resCode, found
}

func GetDictValueOfIntKey(code int, dict []*GenericItem) string {
	for _, item := range dict {
		key, err := GetUintKey(item.Key)
		if err != nil {
			return ""
		}
		if code == key {
			return item.Value.(string)
		}
	}
	return ""
}

func GetUintKey(originKey interface{}) (int, error) {
	switch originKey.(type) {
	case int:
		return originKey.(int), nil
	case float64:
		return int(originKey.(float64)), nil
	case uint:
		return int(originKey.(uint)), nil
	case int8:
		return int(originKey.(int8)), nil
	default:
		return int(0), fmt.Errorf("unknown key type")
	}
}

func GetStringValue(dic []*GenericItem, key interface{}) string {
	var result string

	if _, ok := key.(string); ok {
		for _, v := range dic {
			if v.Key.(string) == key {
				result = v.Value.(string)
				break
			}
		}
	} else {
		keyUint, _ := GetUintKey(key)
		for _, v := range dic {
			if v.Key.(int) == keyUint {
				result = v.Value.(string)
				break
			}
		}
	}

	return result
}

func transferKeyFloatToInt(dic []*GenericItem) []*GenericItem {
	for _, item := range dic {
		item.Key = int(item.Key.(float64))
	}

	return dic
}

// 将key(int), value(string)的dic转成map
func intStrDic2Map(dic []*GenericItem) map[int]string {
	resMap := map[int]string{}
	for _, item := range dic {
		if _, ok := item.Key.(int); !ok {
			return map[int]string{}
		}
		if _, ok := item.Value.(string); !ok {
			return map[int]string{}
		}
		resMap[item.Key.(int)] = item.Value.(string)
	}

	return resMap
}

// 本方法只支持字符串的过滤查找
func FilterTreeByStringValue(tree []*GenericTreeNode, keyword string) ([]*GenericTreeNode, error) {
	if len(tree) == 0 {
		return tree, nil
	}
	// 本方法只支持字符串的过滤查找
	if _, ok := tree[0].Value.(string); !ok {
		return tree, fmt.Errorf("invalid tree, value is not string type")
	}

	resTree := []*GenericTreeNode{}

	for _, subTree := range tree {
		resSubTree := filterTreeByStringValue(subTree, keyword)
		if resSubTree != nil {
			resTree = append(resTree, resSubTree)
		}
	}
	return resTree, nil
}

// tree的value需要是string字符串，才可使用本方法
func filterTreeByStringValue(tree *GenericTreeNode, keyword string) *GenericTreeNode {
	if keyword == "" {
		return tree
	}

	if tree == nil {
		return nil
	}

	resTree := &GenericTreeNode{
		GenericItem: GenericItem{
			Key:   tree.Key,
			Value: tree.Value,
		},
		ParentId: tree.ParentId,
		SubItem:  []*GenericTreeNode{},
	}

	matchFlag := false

	for _, node := range tree.SubItem {
		subTree := filterTreeByStringValue(node, keyword)
		if subTree != nil {
			matchFlag = true
			resTree.SubItem = append(resTree.SubItem, subTree)
		}
	}

	if valueStr, ok := tree.Value.(string); !ok {
		return nil
	} else {
		if !matchFlag && strings.Contains(valueStr, keyword) {
			matchFlag = true
		}
	}

	if matchFlag {
		return resTree
	} else {
		return nil
	}
}
