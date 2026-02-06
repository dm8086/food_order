package initialize

import (
	"context"
	"order_food/global"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

func MiniProgram() {
	cfg := global.GVA_CONFIG.MiniProgram
	if cfg.AppId == "" || cfg.SecretKey == "" {
		global.GVA_LOG.Error("小程序配置信息不完整")
	}
	wc := wechat.NewWechat()

	redisCfg := global.GVA_CONFIG.Redis
	memory := cache.NewRedis(context.Background(), &cache.RedisOpts{
		Host:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		Database: redisCfg.DB,       // use default DB
	})
	config := &miniConfig.Config{
		AppID:     cfg.AppId,
		AppSecret: cfg.SecretKey,
		Cache:     memory,
	}

	global.GVA_WEIXINMP = wc.GetMiniProgram(config)
}
