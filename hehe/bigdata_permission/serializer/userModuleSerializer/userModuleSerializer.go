package userModuleSerializer

//新增：append模式
type CreateUserModule struct {
	UserIds []int `json:"user_ids" form:"user_ids" binding:"required"`
	ModuleIds []int `json:"module_ids" form:"module_ids" binding:"required"`//not empty, cannot insert empty record
}

//保存：与入参保持一致性
type SaveUserModule struct {
	UserIds []int `json:"user_ids" form:"user_ids" binding:"required"`
	ModuleIds []int `json:"module_ids" form:"module_ids"`//permit empty set user module empty
}