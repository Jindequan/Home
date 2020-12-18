package common

import (
	"bigdata_permission/serializer"
	"math"
)



func GetPage(param serializer.Page) serializer.Page {
	page := serializer.Page{
		PageIndex: 1,
		PageSize: 15,
	}
	if param.PageIndex > 0 {
		page.PageIndex = param.PageIndex
	}
	if param.PageSize > 0 && param.PageSize < 10000 {
		page.PageSize = param.PageSize
	}
	return page
}

func GetOffsetLimit(page serializer.Page) (int, int) {
	return page.PageSize * (page.PageIndex - 1), page.PageSize
}

func BuildPageInfo(total int, pageSize int, pageIndex int) *serializer.PageInfo{
	totalPage := math.Ceil(float64(total) / float64(pageSize))
	return &serializer.PageInfo{
		PageIndex: pageIndex,
		PageSize: pageSize,
		Total: total,
		TotalPage: int(totalPage),
	}
}
