package dao

import "gorm.io/gorm"

type User struct {
	Uid int `gorm:"primaryKey"`
	RoleId int
	Remark string
	Base
}

func (user *User) Create () bool {
	res := DB.Create(user)
	return res.Error == nil && res.RowsAffected == 1
}

func (user *User) Update (data map[string]interface{}) (bool, *User) {
	res := DB.Model(user).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, user
}

func (user *User) Delete () bool {
	res := DB.Delete(user)
	return res.Error == nil && res.RowsAffected == 1
}

func GetUserById (uid int) (bool, *User) {
	user := &User{}
	res := DB.Model(user).Where("uid = ?", uid).First(user)
	return res.Error == nil, user
}

func GetUserByWhere (where *User, offset, limit int) (bool, []User) {
	query := getUserQueryByWhere(where)
	list := []User{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountUserByWhere (where *User) (bool, int) {
	query := getUserQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil, int(total)
}

func getUserQueryByWhere(where *User) *gorm.DB {
	query := DB.Model(User{})
	if where.Uid > 0 {
		query = query.Where("uid = ?", where.Uid)
	}
	if where.RoleId > 0 {
		query = query.Where("role_id = ?", where.RoleId)
	}
	return query
}