package dao

import "gorm.io/gorm"

const INTERFACE_TYPE_BACKEND = 1
const INTERFACE_TYPE_FRONTEND = 2

var InterfaceTypeMap = map[int8]string{
	INTERFACE_TYPE_BACKEND: "后端接口",
	INTERFACE_TYPE_FRONTEND: "前端接口",
}

type Interface struct {
	InterfaceId int `gorm:"primaryKey"`
	Type int8
	ShowName string
	Path string
	Remark string
	Base
}

func (interfaceInfo *Interface) Create () bool {
	res := DB.Create(interfaceInfo)
	return res.Error == nil && res.RowsAffected == 1
}

func (interfaceInfo *Interface) Update (data map[string]interface{}) (bool, *Interface) {
	res := DB.Model(interfaceInfo).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, interfaceInfo
}

func (interfaceInfo *Interface) Delete () bool {
	res := DB.Delete(interfaceInfo)
	return res.Error == nil && res.RowsAffected == 1
}

func GetInterfaceById (interfaceId int) (bool, *Interface) {
	interfaceInfo := &Interface{}
	res := DB.Model(interfaceInfo).Where("interface_id = ?", interfaceId).First(interfaceInfo)
	return res.Error == nil, interfaceInfo
}

func GetInterfaceByIds (interfaceIds []int) (bool, []Interface) {
	interfaceList := []Interface{}
	res := DB.Model(Interface{}).Where("interface_id in (?)", interfaceIds).Find(&interfaceList)
	return res.Error == nil, interfaceList
}

func GetInterfaceByTypeAndPath (interfaceType int8, path string) (bool, *Interface) {
	interfaceInfo := &Interface{}
	res := DB.Model(interfaceInfo).Where("type = ?", interfaceType).
		Where("path = ?", path).First(interfaceInfo)
	return res.Error == nil, interfaceInfo
}

func GetInterfaceByWhere (where *Interface, offset, limit int) (bool, []Interface) {
	query := getInterfaceQueryByWhere(where)
	list := []Interface{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountInterfaceByWhere (where *Interface) (bool, int) {
	query := getInterfaceQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil, int(total)
}

func getInterfaceQueryByWhere(where *Interface) *gorm.DB {
	query := DB.Model(Interface{})
	if where.InterfaceId > 0 {
		query = query.Where("interface_id = ?", where.InterfaceId)
	}
	if where.Type > 0 {
		query = query.Where("type = ?", where.Type)
	}
	if where.Path != "" {
		query = query.Where("path = ?", where.Path)
	}
	if where.ShowName != "" {
		query = query.Where("show_name like ?", "%" + where.ShowName + "%")
	}
	return query
}