package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Br = ButtonRouter{}

type ButtonRouter struct{}

func (ButtonRouter) ButtonRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {

	bApi := api.BA
	button := publicRouter.Group("button")

	button.POST("/click", bApi.Click)

}
