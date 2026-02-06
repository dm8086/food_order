package router

import (
	"github.com/gin-gonic/gin"
)

var RT = RouteEnt{}

type RouteEnt struct{}

func (*RouteEnt) RouteInit(prv *gin.RouterGroup, pub *gin.RouterGroup) {
	Br.ButtonRouter(prv, pub)
	Dr.DeskRouter(prv, pub)
	Er.EventRouter(prv, pub)
	Mr.MarkingpunchRouter(prv, pub)
	Or.OrderRouter(prv, pub)
	Sr.StallRouter(prv, pub)
}
