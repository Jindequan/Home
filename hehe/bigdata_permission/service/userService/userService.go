package userService

import (
	"bigdata_permission/dao"
	"bigdata_permission/pkg/dictionary"
	"bigdata_permission/pkg/ecode"
	"bigdata_permission/pkg/redis"
	"bigdata_permission/serializer"
	"bigdata_permission/serializer/userSerializer"
	"bigdata_permission/service/common"
	"encoding/json"
	"strconv"
)

type UserDetail struct {
	SSOUserInfo SSOUserInfo `json:"sso_user_info"`
	UserInfo UserInfo `json:"user_info"`
}

type UserInfo struct {
	Uid      int    `json:"uid"`
	RoleId   int    `json:"role_id"`
	RoleName string `json:"role_name"`
	Remark   string `json:"remark"`
}

type UserList struct {
	PageInfo *serializer.PageInfo `json:"page_info"`
	List *[]UserInfo `json:"list"`
}

func TransferUserInfo(user dao.User) UserInfo {
	return UserInfo{
		Uid:      user.Uid,
		RoleId:   user.RoleId,
		RoleName: dictionary.GetStringValue(dictionary.RoleList, user.RoleId),
		Remark:   user.Remark,
	}
}

func CreateUser(param *userSerializer.CreateUser) (ecode.Code, string, UserInfo) {
	if param.Uid < 1 {
		return ecode.INVALID_PARAM, "", UserInfo{}
	}
	exist, userInfo := dao.GetUserById(param.Uid)
	if exist && userInfo.Uid == param.Uid {
		return ecode.MYSQL_RECORD_EXIST, "已存在该用户", UserInfo{}
	}
	user := &dao.User{
		Uid: param.Uid,
		RoleId: param.RoleId,
		Remark: param.Remark,
	}
	isSuccess := user.Create()
	if !isSuccess {
		return ecode.MYSQL_ERR, "新增用户失败", UserInfo{}
	}
	return ecode.OK, "", TransferUserInfo(*user)
}

func UpdateUser(param *userSerializer.UpdateUser) (ecode.Code, string, UserInfo) {
	if param.Uid < 1 {
		return ecode.INVALID_PARAM, "", UserInfo{}
	}
	exist, userInfo := dao.GetUserById(param.Uid)
	if !exist || userInfo.Uid != param.Uid {
		return ecode.UpdateErrNotFound, "不存在该用户", UserInfo{}
	}
	if userInfo.Remark == param.Remark && userInfo.RoleId == param.RoleId {
		return ecode.OK, "", TransferUserInfo(*userInfo)
	}
	updateData := map[string]interface{}{}

	if userInfo.Remark != param.Remark {
		updateData["remark"] = param.Remark
	}
	if userInfo.RoleId != param.RoleId {
		updateData["role_id"] = param.RoleId
	}

	isSuccess, user := userInfo.Update(updateData)
	if !isSuccess {
		return ecode.MYSQL_ERR, "更新用户失败", UserInfo{}
	}
	return ecode.OK, "", TransferUserInfo(*user)
}

func DeleteUser(param *userSerializer.DeleteUser) (ecode.Code, string) {
	if param.Uid < 1 {
		return ecode.INVALID_PARAM, ""
	}
	exist, userInfo := dao.GetUserById(param.Uid)
	if !exist || userInfo.Uid != param.Uid {
		return ecode.OK, "不存在该用户"
	}
	isSuccess := userInfo.Delete()
	if !isSuccess {
		return ecode.MYSQL_ERR, "删除用户失败"
	}
	//联动删除用户的module关联数据
	dao.DeleteUserModuleByUserId(param.Uid)
	return ecode.OK, ""
}

func SearchUser(param *userSerializer.SearchUser) (ecode.Code, string, *UserList) {
	where := &dao.User{}
	if param.Uid > 0 {
		where.Uid = param.Uid
	}
	if param.RoleId > 0 {
		where.RoleId = param.RoleId
	}
	page := common.GetPage(param.Page)
	_, total := dao.CountUserByWhere(where)
	result := &UserList{
		PageInfo: common.BuildPageInfo(total, page.PageSize, page.PageIndex),
		List: &[]UserInfo{},
	}
	if total == 0 {
		return ecode.OK, "", result
	}

	offset, limit := common.GetOffsetLimit(page)
	_, list := dao.GetUserByWhere(where, offset, limit)
	newList := []UserInfo{}
	for _, user := range list {
		newList = append(newList, TransferUserInfo(user))
	}
	result.List = &newList

	return ecode.OK, "", result
}

func GetUserInfo(param *userSerializer.GetUser) (ecode.Code, *UserInfo) {
	if param.Uid == 0 {
		return ecode.INVALID_PARAM, &UserInfo{}
	}
	isSuccess, userInfo := dao.GetUserById(param.Uid)
	if !isSuccess {
		return ecode.NoUserFound, &UserInfo{}
	}
	if userInfo.Uid != param.Uid {
		return ecode.OK, &UserInfo{}
	}
	info := TransferUserInfo(*userInfo)
	return ecode.OK, &info
}

func GetUserDetailById(uid int) *UserDetail {
	userDetail := getUserDetailFromRedis(uid)
	if userDetail.UserInfo.Uid > 0 && userDetail.UserInfo.Uid == userDetail.SSOUserInfo.Uid {
		return userDetail
	}
	result := &UserDetail{}
	userInfo, code := GetSsoUserInfoByUid(uid)
	if code != ecode.OK {
		return result
	}
	result.SSOUserInfo = userInfo
	ok, user := dao.GetUserById(uid)
	if !ok {
		return result
	}
	result.UserInfo = TransferUserInfo(*user)
	setUserDetailToRedis(uid, result)
	return result
}

func getUserDetailKey(uid int) string {
	return "user-detail-" + strconv.Itoa(uid)
}

func getUserDetailFromRedis(uid int) *UserDetail {
	userDetail := &UserDetail{}
	ok, res := redis.Get(getUserDetailKey(uid))
	if ok {
		json.Unmarshal([]byte(res), userDetail)
	}
	return userDetail
}

func setUserDetailToRedis(uid int, data *UserDetail) {
	value, _ := json.Marshal(data)
	redis.Set(getUserDetailKey(uid), string(value), 3600 * 16)
}