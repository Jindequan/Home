package userModuleService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/serializer/userModuleSerializer"
)

//追加用户-模块关联关系，已有的关系保持不变
func CreateUserModuleBatch(param *userModuleSerializer.CreateUserModule) (ecode.Code, string) {
	if len(param.UserIds) == 0 || len(param.ModuleIds) == 0 {
		return ecode.INVALID_PARAM, ""
	}

	//find and remove repeat part
	_, exist := dao.GetUserModuleByUserIds(param.UserIds)
	userModuleMap := map[int]map[int]int{}
	for _, userModule := range exist {
		if _, ok := userModuleMap[userModule.UserId]; !ok {
			userModuleMap[userModule.UserId] = map[int]int{}
		}
		userModuleMap[userModule.UserId][userModule.ModuleId] = userModule.ModuleId
	}

	insertArr := []dao.UserModule{}
	for _, userId := range param.UserIds {
		for _, moduleId := range param.ModuleIds {
			//already exist
			if _, ok := userModuleMap[userId][moduleId]; ok {
				continue
			}
			insertArr = append(insertArr, dao.UserModule{
				UserId: userId,
				ModuleId: moduleId,
			})
		}
	}

	if len(insertArr) == 0 {
		return ecode.OK, ""
	}
	res := dao.BatchInsertUserModule(insertArr)
	if !res {
		return ecode.CreateErr, ""
	}
	return ecode.OK, ""
}

func SaveUserModule(param *userModuleSerializer.SaveUserModule) (ecode.Code, string) {
	if len(param.UserIds) == 0 {
		return ecode.INVALID_PARAM, ""
	}
	dao.DeleteUserModuleByUserIds(param.UserIds)
	insertArr := []dao.UserModule{}
	for _, userId := range param.UserIds {
		for _, moduleId := range param.ModuleIds {
			insertArr = append(insertArr, dao.UserModule{
				UserId: userId,
				ModuleId: moduleId,
			})
		}
	}

	if len(insertArr) == 0 {
		return ecode.OK, ""
	}
	res := dao.BatchInsertUserModule(insertArr)
	if !res {
		return ecode.CreateErr, ""
	}
	return ecode.OK, ""
}

func ExistUserModuleById(uid, moduleId int) (ecode.Code, bool) {
	ok, res := dao.GetUserModuleByUserIdAndModuleId(uid, moduleId)
	if !ok {
		return ecode.MYSQL_ERR, false
	}
	return ecode.OK, res.UserId == uid && res.ModuleId == moduleId
}


