package request

import (
	"order_food/model/common/request"
	"order_food/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
