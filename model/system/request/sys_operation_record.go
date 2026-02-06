package request

import (
	"order_food/model/common/request"
	"order_food/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
