package main

import (
	"go.uber.org/zap"

	"order_food/core"
	"order_food/global"
	"order_food/initialize"
)

func main() {
	global.GVA_VP = core.Viper() // 初始化Viper
	initialize.Mongo()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.EmqxInit()
	initialize.EtcdInit()
	core.RunWindowsServer()
}

