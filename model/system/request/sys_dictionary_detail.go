package request

import (
	"order_food/model/common/request"
	"order_food/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
