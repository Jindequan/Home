package common

import (
	"bigdata_permission/pkg"
	"bigdata_permission/serializer"
)

func BuildDaoTimeStr (created, updated, deleted interface{}) serializer.DaoTimeString {
	res := serializer.DaoTimeString{
		CreatedAtStr: "",
		UpdatedAtStr: "",
		DeletedAtStr: "",
	}

	if t, ok := created.(int64); ok {
		if t != 0 {
			res.CreatedAtStr = pkg.TimeToString(t)
		}
	}
	if t, ok := updated.(int64); ok {
		if t != 0 {
			res.UpdatedAtStr = pkg.TimeToString(t)
		}
	}
	if t, ok := deleted.(int64); ok {
		if t != 0 {
			res.DeletedAtStr = pkg.TimeToString(t)
		}
	}
	return res
}
