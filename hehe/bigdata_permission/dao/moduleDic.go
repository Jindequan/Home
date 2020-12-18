package dao

import "gorm.io/gorm"

const CURRENT_MODULE_ID = 1

type ModuleDic struct {
	ModuleId int `gorm:"primaryKey"`
	Value string
	FirstLvlId int
	SecondLvlId int
	ThirdLvlId int
	Enable int8
	ResponsePerson int
	Base
}
//用来查询，包含零值检索
type ModuleDicForQuery struct {
	ModuleId int
	FirstLvlId int
	FirstLvlIdZero bool
	SecondLvlId int
	SecondLvlIdZero bool
	ThirdLvlId int
	ThirdLvlIdZero bool
	Enable int8
	EnableZero bool
	ResponsePerson int
	ResponsePersonZero bool
}

func (moduleDic *ModuleDic) Create () bool {
	res := DB.Create(moduleDic)
	return res.Error == nil && res.RowsAffected == 1
}

func (moduleDic *ModuleDic) Update (data map[string]interface{}) (bool, *ModuleDic) {
	res := DB.Model(moduleDic).Updates(data)
	return res.Error ==nil && res.RowsAffected == 1, moduleDic
}

func (moduleDic *ModuleDic) Delete () bool {
	res := DB.Delete(moduleDic)
	return res.Error == nil && res.RowsAffected == 1
}

func GetModuleDicById (moduleId int) (bool, *ModuleDic) {
	moduleDic := &ModuleDic{}
	res := DB.Model(moduleDic).Where("module_id = ?", moduleId).First(moduleDic)
	return res.Error == nil, moduleDic
}

func GetModuleDicByIds (moduleIds []int) (bool, *[]ModuleDic) {
	moduleDic := &[]ModuleDic{}
	res := DB.Model(moduleDic).Where("module_id in (?)", moduleIds).Find(moduleDic)
	return res.Error == nil, moduleDic
}

func GetModuleDicByLevelAndValue (firstId, secondId, thirdId int, value string) (bool, *ModuleDic) {
	module := &ModuleDic{}
	res := DB.Model(ModuleDic{}).
		Where("first_lvl_id = ? and second_lvl_id = ? and third_lvl_id = ? and value = ?", firstId, secondId, thirdId, value).
		First(module)
	return res.Error == nil, module
}

func CountChildModuleByModuleId (moduleId int) (bool, int) {
	count := int64(0)
	res := DB.Model(ModuleDic{}).
		Where("first_lvl_id = ? or second_lvl_id = ? or third_lvl_id = ?", moduleId, moduleId, moduleId).
		Count(&count)
	return res.Error == nil, int(count)
}

func GetChildModuleByModuleId(moduleId int) (bool, []ModuleDic) {
	list := []ModuleDic{}
	res := DB.Model(ModuleDic{}).
		Where("first_lvl_id = ? or second_lvl_id = ? or third_lvl_id = ?", moduleId, moduleId, moduleId).
		Find(&list)
	return res.Error == nil, list
}

func GetRelateModuleByModuleId(moduleId int) (bool, []ModuleDic) {
	list := []ModuleDic{}
	res := DB.Model(ModuleDic{}).
		Where("module_id = ? or first_lvl_id = ? or second_lvl_id = ? or third_lvl_id = ?", moduleId, moduleId, moduleId, moduleId).
		Find(&list)
	return res.Error == nil, list
}

func GetModuleDicByWhere (where *ModuleDicForQuery, offset, limit int) (bool, []ModuleDic) {
	query := getModuleDicQueryByWhere(where)
	list := []ModuleDic{}
	res := query.Offset(offset).Limit(limit).Find(&list)
	return res.Error == nil, list
}

func CountModuleDicByWhere (where *ModuleDicForQuery) (bool, int) {
	query := getModuleDicQueryByWhere(where)
	total := int64(0)
	res := query.Count(&total)
	return res.Error == nil, int(total)
}

func getModuleDicQueryByWhere(where *ModuleDicForQuery) *gorm.DB {
	query := DB.Model(ModuleDic{})

	if where.ModuleId > 0 {
		query = query.Where("module_id = ?", where.ModuleId)
	}
	if where.FirstLvlId > 0 {
		query = query.Where("first_lvl_id = ?", where.FirstLvlId)
	} else if where.FirstLvlIdZero {
		query = query.Where("first_lvl_id = ?", 0)
	}
	if where.SecondLvlId > 0 {
		query = query.Where("second_lvl_id = ?", where.SecondLvlId)
	} else if where.SecondLvlIdZero {
		query = query.Where("second_lvl_id = ?", 0)
	}
	if where.ThirdLvlId > 0 {
		query = query.Where("third_lvl_id = ?", where.ThirdLvlId)
	} else if where.ThirdLvlIdZero {
		query = query.Where("third_lvl_id = ?", 0)
	}
	if where.Enable > 0 {
		query = query.Where("enable = ?", where.Enable)
	} else if where.EnableZero {
		query = query.Where("enable = ?", 0)
	}
	if where.ResponsePerson > 0 {
		query = query.Where("response_person = ?", where.ResponsePerson)
	} else if where.ResponsePersonZero {
		query = query.Where("response_person = ?", 0)
	}
	return query
}