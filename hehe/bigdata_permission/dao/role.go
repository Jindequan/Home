package dao

import "gorm.io/gorm"

const SUPER_ADMIN = 1

type Role struct {
	RoleId int `gorm:"primaryKey"`
	Name string
	Remark string
	Base
}

func (role *Role) Create () bool {
	res := DB.Create(role)
	return res.Error == nil && res.RowsAffected == 1
}

func (role *Role) Update (data map[string]interface{}) (bool, *Role) {
	res := DB.Model(role).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, role
}

func (role *Role) Delete () bool {
	res := DB.Delete(role)
	return res.Error == nil && res.RowsAffected == 1
}

func GetRoleById (roleId int) (bool, *Role) {
	role := &Role{}
	res := DB.Model(role).Where("role_id = ?", roleId).First(role)
	return res.Error == nil, role
}

func GetRoleByName (name string) (bool, *Role) {
	role := &Role{}
	res := DB.Model(role).Where("name = ?", name).First(role)
	return res.Error == nil, role
}

func GetRoleByWhere (where *Role, offset, limit int) (bool, []Role) {
	query := getRoleQueryByWhere(where)
	list := []Role{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountRoleByWhere (where *Role) (bool, int) {
	query := getRoleQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil, int(total)
}

func getRoleQueryByWhere(where *Role) *gorm.DB {
	query := DB.Model(Role{})
	if where.RoleId > 0 {
		query = query.Where("role_id = ?", where.RoleId)
	}
	if where.Name != "" {
		query = query.Where("name like ?", "%" + where.Name + "%")
	}
	return query
}