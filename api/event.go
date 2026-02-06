package api

import (
	"order_food/model/common/response"
	"order_food/model/sea/request"

	"github.com/gin-gonic/gin"
)

type EventApi struct{}

// Add /api/event/action
func (EventApi) Action(c *gin.Context) {
	req := request.TableEventReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	res, err := TableService.DeskHandler(req)
	if err != nil {
		response.FailWithMessage("操作失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(res, "操作成功", c)
}

// // Occupy /api/event/occupy 占桌
// func (EventApi) Occupy(c *gin.Context) {
// 	req := request.TableEventReq{}
// 	_ = c.ShouldBindBodyWithJSON(&req)

// 	if err != nil {
// 		response.FailWithMessage("操作失败"+err.Error(), c)
// 		return
// 	}

// 	// 处理完

// 	response.OkWithDetailed("", "操作成功", c)
// }

// // Unoccupy /api/event/unoccupy 解除占桌
// func (EventApi) Unoccupy(c *gin.Context) {
// 	req := request.TableEventReq{}
// 	_ = c.ShouldBindBodyWithJSON(&req)

// 	if err != nil {
// 		response.FailWithMessage("操作失败"+err.Error(), c)
// 		return
// 	}

// 	// 处理完

// 	response.OkWithDetailed("", "操作成功", c)
// }

// // Union /api/event/unoccupy 合桌
// func (EventApi) Union(c *gin.Context) {
// 	req := request.TableEventReq{}
// 	_ = c.ShouldBindBodyWithJSON(&req)

// 	if err != nil {
// 		response.FailWithMessage("操作失败"+err.Error(), c)
// 		return
// 	}

// 	// 处理完

// 	response.OkWithDetailed("", "操作成功", c)
// }

// // Ununion /api/event/unoccupy 解除合桌
// func (EventApi) Ununion(c *gin.Context) {
// 	req := request.TableEventReq{}
// 	_ = c.ShouldBindBodyWithJSON(&req)

// 	if err != nil {
// 		response.FailWithMessage("操作失败"+err.Error(), c)
// 		return
// 	}

// 	// 处理完

// 	response.OkWithDetailed("", "操作成功", c)
// }
