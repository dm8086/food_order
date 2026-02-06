package initialize

import (
	"net/http"
	"order_food/docs"
	"order_food/middleware"
	"order_food/router"

	"order_food/global"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middleware.DefaultLogger())

	Router.HEAD("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	// Router.StaticFS(global.GVA_CONFIG.Local.StorePath, http.Dir(global.GVA_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PublicGroup.GET("/", func(ctx *gin.Context) {
	})

	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	// 路由注册
	router.RT.RouteInit(PrivateGroup, PublicGroup)

	global.GVA_LOG.Info("router register success")
	return Router
}
