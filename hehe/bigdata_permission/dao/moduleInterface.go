package dao

import "gorm.io/gorm"

type ModuleInterface struct {
	ModuleId int `gorm:"primaryKey"`
	InterfaceId int `gorm:"primaryKey"`
	Base
}

func (moduleInterface *ModuleInterface) Create () bool {
	res := DB.Create(moduleInterface)
	return res.Error == nil && res.RowsAffected == 1
}

func (moduleInterface *ModuleInterface) Update (data map[string]interface{}) (bool, *ModuleInterface) {
	res := DB.Model(moduleInterface).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, moduleInterface
}

func (moduleInterface *ModuleInterface) Delete () bool {
	res := DB.Delete(moduleInterface)
	return res.Error == nil && res.RowsAffected == 1
}

func GetModuleInterfaceByModuleId (moduleId int) (bool, *[]ModuleInterface) {
	moduleInterfaces := &[]ModuleInterface{}
	res := DB.Model(ModuleInterface{}).Where("module_id = ?", moduleId).Find(moduleInterfaces)
	return res.Error == nil, moduleInterfaces
}

func GetModuleInterfaceByModuleIds (moduleIds []int) (bool, *[]ModuleInterface) {
	moduleInterfaces := &[]ModuleInterface{}
	res := DB.Model(ModuleInterface{}).Where("module_id in (?)", moduleIds).Find(moduleInterfaces)
	return res.Error == nil, moduleInterfaces
}

func GetModuleInterfaceByInterfaceId (interfaceId int) (bool, *ModuleInterface) {
	moduleInterfaces := &ModuleInterface{}
	res := DB.Model(ModuleInterface{}).Where("interface_id = ?", interfaceId).First(moduleInterfaces)
	return res.Error == nil, moduleInterfaces
}

func GetModuleInterfaceByModuleIdAndInterfaceId (moduleId, interfaceId int) (bool, *ModuleInterface) {
	moduleInterfaces := &ModuleInterface{}
	res := DB.Model(ModuleInterface{}).
		Where("module_id = ? and interface_id = ?", moduleId, interfaceId).
		First(moduleInterfaces)
	return res.Error == nil, moduleInterfaces
}

func GetModuleInterfaceByWhere (where *ModuleInterface, offset, limit int) (bool, []ModuleInterface) {
	query := getModuleInterfaceQueryByWhere(where)
	list := []ModuleInterface{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountModuleInterfaceByWhere (where *ModuleInterface) (bool, int) {
	query := getModuleInterfaceQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil || res.RowsAffected == 0, int(total)
}

func getModuleInterfaceQueryByWhere(where *ModuleInterface) *gorm.DB {
	query := DB.Model(ModuleInterface{})
	if where.ModuleId > 0 {
		query = query.Where("module_id = ?", where.ModuleId)
	}
	if where.ModuleId > 0 {
		query = query.Where("interface_id = ?", where.InterfaceId)
	}
	return query
}

func DeleteModuleInterfaceByModuleId(moduleId int) bool {
	res := DB.Delete(&ModuleInterface{ModuleId: moduleId})
	return res.Error == nil
}

func DeleteModuleInterfaceByInterfaceId(interfaceId int) bool {
	res := DB.Delete(&ModuleInterface{InterfaceId: interfaceId})
	return res.Error == nil
}