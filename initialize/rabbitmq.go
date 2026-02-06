package initialize

import (
	"order_food/global"

	"github.com/socifi/jazz"
)

func RabbitMQ() {
	cfg := global.GVA_CONFIG.RabbitMQ
	if cfg.Addr == "" {
		global.GVA_LOG.Error("RabbitMQ连接信息为空")

	}
	client, err := jazz.Connect(cfg.Addr)
	if err != nil {
		global.GVA_LOG.Error("RabbitMQ连接错误:" + err.Error())
	} else {
		global.GVA_RABBITMQ = client
	}
}
