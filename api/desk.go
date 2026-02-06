package api

import (
	"context"
	"encoding/json"
	"order_food/global"
	"order_food/model/common/response"
	"order_food/model/sea/request"
	"order_food/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var TableService = service.TableService{}

type DeskApi struct{}

// Open 开台
// @Tags 桌台管理
// @Summary   开台
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     req body   request.OpenDeskReq true "请求体"
// @Success   200   {object}  response.Response "成功"
// @Router    /desk/open [post]
func (*DeskApi) Open(c *gin.Context) {
	req := request.OpenDeskReq{}
	_ = c.ShouldBindJSON(&req)
	if req.TableId == 0 {
		response.FailWithMessage("桌台id不能为空", c)
		return
	}
	if req.Num == 0 {
		response.FailWithMessage("用餐人数不能为空", c)
		return
	}

	// 检测请求id是否存在
	if req.RequestId == "" {
		response.FailWithMessage("请求Id不能为空", c)
		return
	}

	reqByte, _ := json.Marshal(req)
	reqRes, _ := global.GVA_REDIS.SetNX(context.Background(), req.RequestId, string(reqByte), time.Hour).Result()

	if !reqRes {
		response.FailWithMessage("请求已失效", c)
		return
	}

	res, err := TableService.OpenDesk(req)

	if err != nil {
		response.FailWithMessage("开台失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]any{
		"orderId":     res,
		"tableId":     req.TableId,
		"peopleCount": req.Num,
		"status":      3,
	}, "开台成功", c)
}

// Closed 关台
// @Tags 桌台管理
// @Summary   关台
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     req body      request.OpenDeskReq true "请求体"
// @Success   200   {object}  response.Response "成功"
// @Router    /desk/closed [post]
func (*DeskApi) Closed(c *gin.Context) {
	req := request.OpenDeskReq{}
	_ = c.ShouldBindJSON(&req)
	if req.TableId == 0 {
		response.FailWithMessage("桌台id不能为空", c)
		return
	}
	err := TableService.ClosedDesk(req)

	if err != nil {
		response.FailWithMessage("关闭桌台失败:"+err.Error(), c)
		return
	}
	response.OkWithData("桌台关闭成功", c)
}

// DeskOrder 获取桌台订单
// @Tags 桌台管理
// @Summary 获取桌台订单
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   deskId 	query  int   true "桌台id"
// @Success   200   {object} response.Response  "成功"
// @Router    /desk/order [get]
func (*DeskApi) DeskOrder(c *gin.Context) {
	deskIdStr := c.Query("deskId")
	deskId, _ := strconv.Atoi(deskIdStr)
	if deskId == 0 {
		response.FailWithMessage("桌台id不能为空", c)
		return
	}

	orderId, status, err := TableService.DeskOrder(deskId)
	if err != nil {
		response.FailWithMessage("获取桌台详情:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]any{
		"orderId": orderId,
		"status":  status,
	}, "获取桌台信息成功", c)
}

// DeskList 桌台列表
// @Tags 桌台管理
// @Summary 桌台列表
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   storeId 	query  string   true "门店id"
// @Param   humanNum 	query  int   true "用餐人数"
// @Param   areaIds 	query  []int   true "区域id列表"
// @Param   businessStatus 	query  []int   true "状态列表"
// @Success   200   {object} response.Response  "成功"
// @Router    /desk/list [get]
func (*DeskApi) DeskList(c *gin.Context) {
	req := request.TableListReq{}
	_ = c.ShouldBindQuery(&req)
	list, err := TableService.DeskList(req)
	if err != nil {
		response.FailWithMessage("获取桌列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取桌台信息成功", c)
}

// DeskChange 换桌
// @Tags 桌台管理
// @Summary 换桌
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.ChangeDeskReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /desk/change [post]
func (*DeskApi) DeskChange(c *gin.Context) {
	req := request.ChangeDeskReq{}
	_ = c.ShouldBindJSON(&req)
	if req.FromTableId == 0 {
		response.FailWithMessage("原桌台id不能为空", c)
		return
	}
	if req.ToTableId == 0 {
		response.FailWithMessage("目标桌台id不能为空", c)
		return
	}
	err := TableService.ChangeDesk(req)
	if err != nil {
		response.FailWithMessage("获取桌列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(req, "换桌成功", c)
}

// DeskJoint 合桌/合单
// @Tags 桌台管理
// @Summary 合桌/合单
// @Security  ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param     req body      request.DeskJointReq true "请求体"
// @Success   200   {object} response.Response  "成功"
// @Router    /desk/joint [post]
func (*DeskApi) DeskJoint(c *gin.Context) {
	req := request.DeskJointReq{}
	_ = c.ShouldBindJSON(&req)
	if req.FromTableId == 0 {
		response.FailWithMessage("原桌台id不能为空", c)
		return
	}
	if req.ToTableId == 0 {
		response.FailWithMessage("目标桌台id不能为空", c)
		return
	}
	err := TableService.DeskJoint(req)
	if err != nil {
		response.FailWithMessage("获取桌列表失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(req, "合桌成功", c)
}
