package dao

import "gorm.io/gorm"

type RoleModule struct {
	RoleId int `gorm:"primaryKey"`
	ModuleId int `gorm:"primaryKey"`
	Base
}

func (roleModule *RoleModule) Create () bool {
	res := DB.Create(roleModule)
	return res.Error == nil && res.RowsAffected == 1
}

func (roleModule *RoleModule) Update (data map[string]interface{}) (bool, *RoleModule) {
	res := DB.Model(roleModule).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, roleModule
}

func (roleModule *RoleModule) Delete () bool {
	res := DB.Delete(roleModule)
	return res.Error == nil && res.RowsAffected == 1
}

func GetRoleModuleByRoleId (roleId int) (bool, *[]RoleModule) {
	roleModules := &[]RoleModule{}
	res := DB.Model(RoleModule{}).Where("role_id = ?", roleId).Find(roleModules)
	return res.Error == nil, roleModules
}

func GetRoleModuleByModuleId (moduleId int) (bool, *[]RoleModule) {
	roleModules := &[]RoleModule{}
	res := DB.Model(RoleModule{}).Where("module_id = ?", moduleId).Find(roleModules)
	return res.Error == nil, roleModules
}

func GetRoleModuleByRoleIdAndModuleId (roleId, moduleId int) (bool, *RoleModule) {
	roleModule := &RoleModule{}
	res := DB.Model(RoleModule{}).Where("role_id = ? and module_id = ?", roleId, moduleId).First(roleModule)
	return res.Error == nil, roleModule
}

func GetRoleModuleByWhere (where *RoleModule, offset, limit int) (bool, []RoleModule) {
	query := getRoleModuleQueryByWhere(where)
	list := []RoleModule{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountRoleModuleByWhere (where *RoleModule) (bool, int) {
	query := getRoleModuleQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil || res.RowsAffected == 0, int(total)
}

func getRoleModuleQueryByWhere(where *RoleModule) *gorm.DB {
	query := DB.Model(RoleModule{})
	if where.RoleId > 0 {
		query = query.Where("role_id = ?", where.RoleId)
	}
	if where.ModuleId > 0 {
		query = query.Where("module_id = ?", where.ModuleId)
	}
	return query
}

func DeleteRoleModuleByModuleId(moduleId int) bool {
	res := DB.Delete(&RoleModule{ModuleId: moduleId})
	return res.Error == nil
}