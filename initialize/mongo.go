package initialize

import (
	"order_food/global"
	"time"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

func Mongo() {
	cfg := global.GVA_CONFIG.MongoDB

	if cfg.Uri == "" {
		global.GVA_LOG.Error("mongodb未配置连接信息")
		return
	}

	mongo, err := mgo.Dial(cfg.Uri)
	if err != nil {
		global.GVA_LOG.Error("mongo connect failed, err:", zap.Error(err))
	}

	client := mongo.DB(cfg.DB)
	global.GVA_MONGO = client

	// 开启携程增加检查是否断链,重新恢复 30秒检查一次
	go func() {
		ticker := time.NewTicker(30 * time.Second)
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

					clientRetryM := mongo.DB(cfg.DB)
					global.GVA_APILOG = clientRetryM
				}
			}
		}
	}()
}
