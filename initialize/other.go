package initialize

import (
	"context"
	"fmt"
	"order_food/global"
	"order_food/utils"

	"github.com/loebfly/keruyun-sdk-go/keruyun"
	"github.com/loebfly/keruyun-sdk-go/kry_model"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.GVA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}

	kryConfig := kry_model.SdkConfig{
		Domain:    global.GVA_CONFIG.Keruyun.Domain,
		AppKey:    global.GVA_CONFIG.Keruyun.AppKey,
		SecretKey: global.GVA_CONFIG.Keruyun.SecretKey,
		Version:   global.GVA_CONFIG.Keruyun.Version,
		SetTokenForShopIdHandle: func(shopId int64, token string) {
			global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("keruyun_token_%d", shopId), token, 0)
		},
		GetTokenForShopIdHandle: func(shopId int64) string {
			token, tokenErr := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("keruyun_token_%d", shopId)).Result()
			if tokenErr != nil {
				global.GVA_LOG.Error("GetTokenForShopIdHandle", zap.Error(tokenErr))
				return ""
			}
			return token
		},
		PrintApiLogHandle: func(ctx kry_model.ReqCtx) {
			if global.GVA_CONFIG.Keruyun.Debug {
				global.GVA_LOG.Debug("keruyunRequest", zap.String("data", fmt.Sprintf("%+v", ctx)))
			}
		},
	}
	keruyun.RegisterSdk(kryConfig)

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
