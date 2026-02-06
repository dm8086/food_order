package api

import (
	"context"
	"fmt"
	"order_food/global"
	"order_food/model/common/response"
	"order_food/model/sea/request"
	"order_food/service"
	"time"

	"github.com/gin-gonic/gin"
)

var OrderService = service.OrderService{}

type OrderApi struct{}

// Add 添加订单
// @Tags 订单管理
// @Summary 添加订单
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.OrderAddReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/add [post]
func (OrderApi) Add(c *gin.Context) {
	req := request.OrderAddReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	if req.RequestId == "" {
		response.FailWithMessage("请求id不能为空", c)
		return
	}

	if len(req.Dishes) == 0 {
		response.FailWithMessage("添加菜品不能为空", c)
		return
	}

	redisKey := fmt.Sprintf("order:request:%s", req.RequestId)
	res, _ := global.GVA_REDIS.SetNX(context.Background(), redisKey, req.RequestId, 15*time.Second).Result()
	if !res {
		response.FailWithMessage("请求已失效", c)
		return
	}

	orderInfo, _, err := OrderService.OrderAdd(req)
	if err != nil {
		response.FailWithMessage("菜品添加失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(map[string]any{
		"orderId":     orderInfo.OrderId,
		"tableId":     req.TableId,
		"peopleCount": orderInfo.PeopleCount,
	}, "订单添加成功", c)
}

// Sub 订单减少商品
// @Tags 订单管理
// @Summary 订单减少商品
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.OrderSubReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/sub [post]
func (OrderApi) Sub(c *gin.Context) {
	req := request.OrderSubReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	if req.RequestId == "" {
		response.FailWithMessage("请求id不能为空", c)
		return
	}

	if len(req.GoodsList) == 0 {
		response.FailWithMessage("减少菜品不能为空", c)
		return
	}

	redisKey := fmt.Sprintf("order:request:%s", req.RequestId)
	res, _ := global.GVA_REDIS.SetNX(context.Background(), redisKey, req.RequestId, 15*time.Second).Result()
	if !res {
		response.FailWithMessage("请求已失效", c)
		return
	}

	_, _, err := OrderService.OrderSub(req)
	if err != nil {
		response.FailWithMessage("菜品移除失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]any{
		"orderId": "",
		"tableId": 0,
	}, "订单菜品移除成功", c)

}

// Sub 订单更新
// @Tags 订单管理
// @Summary 订单更新
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.OrderUpdateReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/update [post]
func (OrderApi) Update(c *gin.Context) {
	req := request.OrderUpdateReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	if req.RequestId == "" {
		response.FailWithMessage("请求id不能为空", c)
		return
	}
	err := OrderService.Update(req)
	if err != nil {
		response.FailWithMessage("修改订单失败"+err.Error(), c)
		return
	}

	response.OkWithData("修改订单成功", c)
}

// Sub 订单菜品扣减
// @Tags 订单管理
// @Summary 订单菜品扣减
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.CombWriteoffReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/sku/update [post]
func (OrderApi) OrderSkuUpdate(c *gin.Context) {
	req := request.CombWriteoffReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	if req.RequestId == "" {
		response.FailWithMessage("请求id不能为空", c)
		return
	}
	err := OrderService.DishSkuUpdate(req)
	if err != nil {
		response.FailWithMessage("扣减失败"+err.Error(), c)
		return
	}

	response.OkWithData("扣减成功成功", c)
}

// Info 订单详情
// @Tags 订单管理
// @Summary 订单详情
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   orderId 	query  string   true "订单id"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/info [get]
func (OrderApi) Info(c *gin.Context) {
	orderId := c.Query("orderId")
	if orderId == "" {
		response.FailWithMessage("订单id不能为空", c)
		return
	}
	info, err := OrderService.Info(orderId)
	if err != nil {
		response.FailWithMessage("获取详情失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(info, "获取详情成功", c)
}

// List 订单列表
// @Tags 订单管理
// @Summary 订单列表
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   areaId 	query  string   false "区域id"
// @Param   status 	query  string   false "状态"
// @Param   sort 	query  string   false "开始时间"
// @Param   page 	query  string   false "分页页码"
// @Param   pageSize 	query  string   false "分页大小"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/list [get]
func (OrderApi) List(c *gin.Context) {
	req := request.OrderListReq{}
	_ = c.ShouldBindQuery(&req)

	if req.Page <= 1 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	list, count, err := OrderService.List(req)
	if err != nil {
		response.FailWithMessage("获取详情失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// SkuRemove 订单商品移除
// @Tags 订单管理
// @Summary 订单商品移除
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.OrderSkuRemoveReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/sku/remove [post]
func (OrderApi) SkuRemove(c *gin.Context) {
	req := request.OrderSkuRemoveReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	err := OrderService.SkuRemove(req)
	if err != nil {
		response.FailWithMessage("送达商品失败"+err.Error(), c)
		return
	}

	response.OkWithData("送达商品成功", c)
}

// SoupRemove 锅底移除
// @Tags 订单管理
// @Summary 锅底移除
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.SoupRemoveReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /order/soup/remove [post]
func (OrderApi) SoupRemove(c *gin.Context) {
	req := request.SoupRemoveReq{}
	_ = c.ShouldBindBodyWithJSON(&req)

	err := OrderService.SoupRemove(req)
	if err != nil {
		response.FailWithMessage("操作失败"+err.Error(), c)
		return
	}

	response.OkWithData("操作成功", c)
}
