package initialize

import "order_food/global"

func WeCom() {
	cfg := global.GVA_CONFIG.Wecom
	if cfg.AppId == "" || cfg.AesKey == "" || cfg.Token == "" {
		global.GVA_LOG.Error("企业微信配置信息不完整")
	}
}
