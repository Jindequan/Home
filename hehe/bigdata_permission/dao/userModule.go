package dao

type UserModule struct {
	UserId int `gorm:"primaryKey"`
	ModuleId int `gorm:"primaryKey"`
	Base
}

func (userModule *UserModule) Create () bool {
	res := DB.Create(userModule)
	return res.Error == nil && res.RowsAffected == 1
}

func (userModule *UserModule) Update (data map[string]interface{}) (bool, *UserModule) {
	res := DB.Model(userModule).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, userModule
}

func (userModule *UserModule) Delete () bool {
	res := DB.Delete(userModule)
	return res.Error == nil && res.RowsAffected == 1
}

func BatchInsertUserModule(userModules []UserModule) bool {
	res := DB.Create(&userModules)
	return res.Error == nil
}

func DeleteUserModuleByUserId(userId int) bool {
	if userId == 0 {
		return true
	}
	res := DB.Where("user_id = ?", userId).Delete(UserModule{})
	return res.Error == nil
}

func DeleteUserModuleByModuleId(moduleId int) bool {
	if moduleId == 0 {
		return true
	}
	res := DB.Where("module_id = ?", moduleId).Delete(UserModule{})
	return res.Error == nil
}

func DeleteUserModuleByUserIds(userIds []int) bool {
	res := DB.Where("user_id in (?)", userIds).Delete(UserModule{})
	return res.Error == nil
}

func GetUserModuleByUserIds(userIds []int) (bool, []UserModule) {
	list := []UserModule{}
	res := DB.Where("user_id in (?)", userIds).Find(&list)
	return res.Error == nil, list
}

func GetUserModuleByUserIdAndModuleId(userId, moduleId int) (bool, UserModule) {
	ret := UserModule{}
	res := DB.Where("user_id = ? and module_id = ?", userId, moduleId).First(&ret)
	return res.Error == nil, ret
}