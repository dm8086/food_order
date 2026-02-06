package config

type Keruyun struct {
	Domain    string `mapstructure:"domain" json:"domain" yaml:"domain"`          // 域名
	AppKey    string `mapstructure:"appKey" json:"appKey" yaml:"appKey"`          // 应用标识
	SecretKey string `mapstructure:"secretKey" json:"secretKey" yaml:"secretKey"` // 应用密钥
	Version   string `mapstructure:"version" json:"version" yaml:"version"`       // 版本号, 默认为 2.0
	Debug     bool   `mapstructure:"debug" json:"debug" yaml:"debug"`             // 启用调试日志
}
