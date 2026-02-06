package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Mr = MarkingpunchRouter{}

type MarkingpunchRouter struct{}

func (MarkingpunchRouter) MarkingpunchRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	mkpa := privateRouter.Group("marking")
	api := api.MarkingpunchApi{}

	mkpa.POST("add", api.Add)
	mkpa.POST("edit", api.Edit)
	mkpa.GET("info", api.Info)
	mkpa.GET("list", api.List)

	mkpa.POST("conf/add", api.ConfAdd)
	mkpa.POST("conf/edit", api.ConfEdit)
	mkpa.GET("conf/info", api.ConfInfo)
	mkpa.GET("conf/list", api.ConfList)

	mkpa.POST("tmp/add", api.TmpAdd)
	mkpa.POST("tmp/edit", api.TmpEdit)
	mkpa.GET("tmp/info", api.TmpInfo)
	mkpa.GET("tmp/list", api.TmpList)

	mkpa.POST("bill/add", api.BillAdd)
	mkpa.POST("bill/edit", api.BillEdit)
	mkpa.GET("bill/info", api.BillInfo)
	mkpa.GET("bill/list", api.BillList)
}
