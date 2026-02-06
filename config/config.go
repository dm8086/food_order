package config

type Server struct {
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Email    Email    `mapstructure:"email" json:"email" yaml:"email"`
	System   System   `mapstructure:"system" json:"system" yaml:"system"`
	Captcha  Captcha  `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	NodeRed  NodeRed  `mapstructure:"nodered" json:"nodered" yaml:"nodered"`
	MongoDB  MongoDB  `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	ApiLog   ApiLog   `mapstructure:"apilog" json:"apilog" yaml:"apilog"`

	// gorm
	Mysql  Mysql           `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	DBList []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	// oss
	Local Local `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu Qiniu `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`

	Excel Excel `mapstructure:"excel" json:"excel" yaml:"excel"`
	Timer Timer `mapstructure:"timer" json:"timer" yaml:"timer"`

	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`

	// 客如云配置
	Keruyun Keruyun `mapstructure:"keruyun" json:"keruyun" yaml:"keruyun"`

	// 企业微信配置
	Wecom         Wecom          `mapstructure:"wecom" json:"wecom" yaml:"wecom"`
	WecomWebHooks Wecom_Webhooks `mapstructure:"wecom-webhooks" json:"wecom-webhooks" yaml:"wecom-webhooks"`

	// 小程序配置
	MiniProgram MiniProgram `mapstructure:"miniProgram" json:"miniProgram" yaml:"miniProgram"`

	// 不需要记录日志的接口
	NoApilog NoApilog `mapstructure:"noApilog" json:"noApilog" yaml:"noApilog"`

	// 订单关闭变为待收桌时间
	OrderToClearTime int `mapstructure:"orderToClearTime" json:"orderToClearTime" yaml:"orderToClearTime"`

	// 设备刷新
	DeviceRefreshUri string `mapstructure:"deviceRefreshUri" json:"deviceRefreshUri" yaml:"deviceRefreshUri"`

	// etcd
	Etcd Etcd `mapstructure:"etcd" json:"etcd" yaml:"etcd"`

	// emqx
	Emqx Emqx `mapstructure:"emqx" json:"emqx" yaml:"emqx"`
}
