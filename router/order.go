package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Or = OrderRouter{}

type OrderRouter struct{}

func (OrderRouter) OrderRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	orderRouter := publicRouter.Group("")
	api := api.OrderApi{}

	orderRouter.POST("add", api.Add)
	orderRouter.POST("sub", api.Sub)
	orderRouter.POST("pay/callback", api.Update)
	orderRouter.POST("update", api.Update)
	orderRouter.GET("info", api.Info)
	orderRouter.GET("list", api.List)
	orderRouter.POST("sku/remove", api.SkuRemove)

	orderRouter.POST("sku/update", api.OrderSkuUpdate)

	orderRouter.POST("soup/remove", api.SoupRemove)

}
