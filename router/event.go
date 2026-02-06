package router

import (
	"order_food/api"

	"github.com/gin-gonic/gin"
)

var Er = EventRouter{}

type EventRouter struct{}

func (EventRouter) EventRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup) {
	eventRouter := privateRouter.Group("event")
	evenapi := api.EventApi{}

	eventRouter.POST("action", evenapi.Action)
}
