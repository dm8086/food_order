package initialize

import (
	"order_food/global"
)

func NodeRed() {
	cfg := global.GVA_CONFIG.NodeRed
	if cfg.Host == "" {
		global.GVA_LOG.Error("nodered主机信息为空")
	}
}
