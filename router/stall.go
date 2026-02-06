package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Sr = StallRouter{}

type StallRouter struct{}

func (StallRouter) StallRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	stallRouter := privateRouter.Group("stall")
	stallBingRouter := privateRouter.Group("stall/bind")
	api := api.StallApi{}

	stallRouter.POST("add", api.Add)
	stallRouter.POST("edit", api.Edit)
	stallRouter.GET("list", api.StallList)
	stallRouter.GET("list/nopage", api.StallListNoPage)

	stallBingRouter.POST("add", api.StallBingAdd)
	stallBingRouter.POST("edit", api.StallBingEdit)
	stallBingRouter.GET("info", api.StallBingInfo)
	stallBingRouter.GET("list", api.StallBingList)

}
