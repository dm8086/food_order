package api

import (
	"order_food/model/common/response"
	"order_food/model/sea/request"
	"order_food/service"

	"github.com/gin-gonic/gin"
)

type ButtonApi struct{}

// Click 点击
// @Tags 按钮服务
// @Summary   点击
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     req body      request.ButtonServiceReq true "请求体"
// @Success   200   {object}  response.Response "成功"
// @Router    /button/click [post]
func (*ButtonApi) Click(c *gin.Context) {
	req := request.ButtonServiceReq{}
	_ = c.ShouldBindJSON(&req)

	if req.EventType == 0 {
		response.FailWithMessage("事件不能为空", c)
		return
	}
	if req.TableId == 0 {
		response.FailWithMessage("桌台不能为空", c)
		return
	}
	err := service.BS.Click(req)
	if err != nil {
		response.FailWithMessage("按钮点击错误:"+err.Error(), c)
		return
	}
	response.OkWithDetailed("成功", "", c)

}
