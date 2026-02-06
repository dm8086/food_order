package api

import (
	"order_food/model/common/response"
	"order_food/model/sea/request"
	"order_food/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var StallService = service.StallService{}

type StallApi struct{}

// Add 明档添加
// @Tags 明档管理
// @Summary 明档添加
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.StallAddReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/add [post]
func (StallApi) Add(c *gin.Context) {
	req := []request.StallAddReq{}
	_ = c.ShouldBindJSON(&req)

	err := StallService.StoreStall.Add(req)
	if err != nil {
		response.FailWithMessage("添加明档失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed("", "添加明档成功", c)
}

// Edit 明档修改
// @Tags 明档管理
// @Summary 明档修改
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.StallEditReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/edit [post]
func (StallApi) Edit(c *gin.Context) {
	req := []request.StallEditReq{}
	_ = c.ShouldBindJSON(&req)

	err := StallService.StoreStall.Edit(req)
	if err != nil {
		response.FailWithMessage("修改明档失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed("", "修改明档成功", c)
}

// StallList 明档列表
// @Tags 明档管理
// @Summary 明档列表
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.StallListReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/list [get]
func (StallApi) StallList(c *gin.Context) {
	req := request.StallListReq{}
	_ = c.ShouldBindQuery(&req)

	if req.Page <= 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 30
	}

	list, count, err := StallService.StoreStall.List(req)
	if err != nil {
		response.FailWithMessage("获取列表失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取列表成功", c)
}

func (StallApi) StallListNoPage(c *gin.Context) {

	storeId := c.Query("storeId")
	if storeId == "" {
		response.FailWithMessage("门店id不能为空", c)
		return
	}

	list, err := StallService.StoreStall.ListNopage(storeId)
	if err != nil {
		response.FailWithMessage("获取列表失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(list, "获取列表成功", c)
}

// StallBingAdd 明档绑定商品添加
// @Tags 明档管理
// @Summary 明档绑定商品添加
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.StallBingFoodAddReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/bing/add [post]
func (StallApi) StallBingAdd(c *gin.Context) {
	req := []request.StallBingFoodAddReq{}
	_ = c.ShouldBindJSON(&req)
	err := StallService.StoreBindFood.Add(req)
	if err != nil {
		response.FailWithMessage("添加绑定成功"+err.Error(), c)
		return
	}

	response.OkWithData("添加绑定失败", c)
}

// StallBingEdit 明档绑定商品修改
// @Tags 明档管理
// @Summary 明档绑定商品修改
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.StallBingFoodEditReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/bing/edit [post]
func (StallApi) StallBingEdit(c *gin.Context) {
	req := []request.StallBingFoodEditReq{}
	_ = c.ShouldBindJSON(&req)
	err := StallService.StoreBindFood.Edit(req)
	if err != nil {
		response.FailWithMessage("修改绑定失败"+err.Error(), c)
		return
	}

	response.OkWithData("修改绑定成功", c)
}

// StallBingInfo 明档绑定商品信息
// @Tags 明档管理
// @Summary 明档绑定商品信息
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   id 	query  string   true "id"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/bing/info [get]
func (StallApi) StallBingInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	info, err := StallService.StoreBindFood.Info(id)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(info, "获取成功", c)
}

// StallBingList 明档绑定商品列表
// @Tags 明档管理
// @Summary 明档绑定商品列表
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   storeId 	query  string   false "门店id"
// @Param   stallId 	query  string   false "明档id"
// @Param   sort 	query  string   false "排序"
// @Param   page 	query  string   false "分页页码"
// @Param   pageSize 	query  string   false "分页大小"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/bing/list [get]
func (StallApi) StallBingList(c *gin.Context) {
	req := request.StallBingFoodListReq{}
	_ = c.ShouldBindQuery(&req)

	if req.Page <= 1 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 30
	}

	list, count, err := StallService.StoreBindFood.List(req)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// StallSkusInfo 明档绑定商品详情
// @Tags 明档管理
// @Summary 明档绑定商品详情
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   id 	query  string   false "id"
// @Success   200   {object} response.Response  "成功"
// @Router    /stall/skus/info [get]
func (StallApi) StallSkusInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	if idStr == "" {
		response.FailWithMessage("id不能为空", c)
		return
	}
	res, err := StallService.StallSkus.Info(id)
	if err != nil {
		response.FailWithMessage("获取详情错误", c)
		return
	}
	response.OkWithDetailed(res, "获取成功", c)
}

func (StallApi) StallSkusList(c *gin.Context) {
	req := request.StallSkuListReq{}
	_ = c.ShouldBindQuery(&req)
	list, err := StallService.StallSkus.List(req)
	if err != nil {
		response.FailWithMessage("获取列表错误"+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
