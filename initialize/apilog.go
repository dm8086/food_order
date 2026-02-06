package initialize

import (
	"order_food/global"
	"time"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

func ApiLog() {
	cfg := global.GVA_CONFIG.ApiLog

	if cfg.Uri == "" {
		global.GVA_LOG.Error("mongodb未配置连接信息")
		return
	}

	mongo, err := mgo.Dial(cfg.Uri)
	if err != nil {
		global.GVA_LOG.Error("mongo connect failed, err:", zap.Error(err))
	}

	client := mongo.DB(cfg.DB)
	global.GVA_APILOG = client

	// 开启携程增加检查是否断链,重新恢复 5分钟检查一次
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-ticker.C:
				var err error
				if global.GVA_APILOG.Session != nil {
					err = global.GVA_APILOG.Session.Ping()
					global.GVA_LOG.Info("mongoDb检查结果:", zap.Error(err))
				}

				if global.GVA_APILOG.Session == nil || err != nil {
					mongo, err := mgo.Dial(cfg.Uri)
					if err != nil {
						global.GVA_LOG.Error("mongo connect failed, err:", zap.Error(err))
					}

					clientRetry := mongo.DB(cfg.DB)
					global.GVA_APILOG = clientRetry
				}
			}
		}
	}()
}
