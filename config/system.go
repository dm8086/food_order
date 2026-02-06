package config

type System struct {
	Env            string `mapstructure:"env" json:"env" yaml:"env"`                                     // 环境值
	Addr           int    `mapstructure:"addr" json:"addr" yaml:"addr"`                                  // 端口值
	DbType         string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                         // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType        string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                      // Oss类型
	UseMultipoint  bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`    // 多点登录拦截
	UseRedis       bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                   // 使用redis
	UseRabbitMQ    bool   `mapstructure:"use-rabbitmq" json:"use-rabbitmq" yaml:"use-rabbitmq"`          // 使用rabbitmq
	UseNodeRed     bool   `mapstructure:"use-nodered" json:"use-nodered" yaml:"use-nodered"`             // 使用nodered
	UseWecom       bool   `mapstructure:"use-wecom" json:"use-wecom" yaml:"use-wecom"`                   // 使用企业微信
	UseMiniProgram bool   `mapstructure:"use-miniprogram" json:"use-miniprogram" yaml:"use-miniprogram"` // 使用微信小程序
	UseMongo       bool   `mapstructure:"use-mongo" json:"use-mongo" yaml:"use-mongo"`                   // 使用mongoDB
	LimitCountIP   int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP    int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	RouterPrefix   string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}
