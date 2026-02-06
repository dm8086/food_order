package api

import (
	"order_food/model/common/response"
	"order_food/model/sea/request"
	"order_food/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Markingpunch = service.MarkingpunchService{}

type MarkingpunchApi struct{}

func (*MarkingpunchApi) Add(c *gin.Context) {
	req := request.MarkingPunchAddReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.Markingpunch.Add(req)
	if err != nil {
		response.FailWithMessage("添加失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

func (*MarkingpunchApi) Edit(c *gin.Context) {
	req := request.MarkingPunchUpdateReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.Markingpunch.Update(req)
	if err != nil {
		response.FailWithMessage("修改失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

func (*MarkingpunchApi) Info(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	info, err := Markingpunch.Markingpunch.Info(id)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(info, "获取成功", c)
}

func (*MarkingpunchApi) List(c *gin.Context) {
	req := request.MarkingPunchListReq{}
	_ = c.ShouldBindQuery(&req)
	list, count, err := Markingpunch.Markingpunch.List(req)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

func (*MarkingpunchApi) ConfAdd(c *gin.Context) {
	req := request.MarkingConfAddReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingConf.Add(req)
	if err != nil {
		response.FailWithMessage("添加失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

func (*MarkingpunchApi) ConfEdit(c *gin.Context) {
	req := request.MarkingConfUpdateReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingConf.Update(req)
	if err != nil {
		response.FailWithMessage("修改失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

func (*MarkingpunchApi) ConfInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	info, err := Markingpunch.MarkingConf.Info(id)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(info, "获取成功", c)
}

func (*MarkingpunchApi) ConfList(c *gin.Context) {
	req := request.MarkingConfListReq{}
	_ = c.ShouldBindQuery(&req)
	list, count, err := Markingpunch.MarkingConf.List(req)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

func (*MarkingpunchApi) TmpAdd(c *gin.Context) {
	req := request.MarkingTmpAddReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingTmp.Add(req)
	if err != nil {
		response.FailWithMessage("添加失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

func (*MarkingpunchApi) TmpEdit(c *gin.Context) {
	req := request.MarkingTmpUpdateReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingTmp.Update(req)
	if err != nil {
		response.FailWithMessage("修改失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

func (*MarkingpunchApi) TmpInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	info, err := Markingpunch.MarkingTmp.Info(id)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(info, "获取成功", c)
}

func (*MarkingpunchApi) TmpList(c *gin.Context) {
	req := request.MarkingTmpListReq{}
	_ = c.ShouldBindQuery(&req)
	list, count, err := Markingpunch.MarkingTmp.List(req)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

func (*MarkingpunchApi) BillAdd(c *gin.Context) {
	req := request.MarkingBillAddReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingBill.Add(req)
	if err != nil {
		response.FailWithMessage("添加失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

func (*MarkingpunchApi) BillEdit(c *gin.Context) {
	req := request.MarkingBillUpdateReq{}
	_ = c.ShouldBindJSON(&req)
	err := Markingpunch.MarkingBill.Update(req)
	if err != nil {
		response.FailWithMessage("修改失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

func (*MarkingpunchApi) BillInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.Atoi(idStr)
	info, err := Markingpunch.MarkingBill.Info(id)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(info, "获取成功", c)
}

func (*MarkingpunchApi) BillList(c *gin.Context) {
	req := request.MarkingBillListReq{}
	_ = c.ShouldBindQuery(&req)
	list, count, err := Markingpunch.MarkingBill.List(req)
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}
