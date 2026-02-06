package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Dr = DeskRouter{}

type DeskRouter struct{}

func (DeskRouter) DeskRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	apiTable := api.DeskApi{}

	deskRouter := publicRouter.Group("desk")
	deskRouter.POST("open", apiTable.Open)
	deskRouter.POST("closed", apiTable.Closed)
	deskRouter.GET("order", apiTable.DeskOrder)
	deskRouter.GET("list", apiTable.DeskList)
	deskRouter.POST("change", apiTable.DeskChange)
	deskRouter.POST("joint", apiTable.DeskJoint)
}
