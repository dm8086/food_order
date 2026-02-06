package global

import (
	"order_food/config"
	"sync"

	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/socifi/jazz"
	"gopkg.in/mgo.v2"

	"order_food/utils/timer"

	"github.com/songzhibin97/gkit/cache/local_cache"

	"golang.org/x/sync/singleflight"

	"go.uber.org/zap"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	GVA_DB                  *gorm.DB
	GVA_DBList              map[string]*gorm.DB
	GVA_REDIS               *redis.Client
	GVA_MONGO               *mgo.Database
	GVA_CONFIG              config.Server
	GVA_VP                  *viper.Viper
	GVA_LOG                 *zap.Logger
	GVA_APILOG              *mgo.Database
	GVA_RABBITMQ            *jazz.Connection
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
	GVA_WEIXINMP            *miniprogram.MiniProgram
	ServUuid                string
	GVA_EMQX                mqtt.Client
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
