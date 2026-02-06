package config

type ApiLog struct {
	Uri        string `mapstructure:"uri" json:"uri" yaml:"uri"`                      // 服务器地址
	DB         string `mapstructure:"db" json:"db" yaml:"db"`                         // 数据库
	Collection string `mapstructure:"collection" json:"collection" yaml:"collection"` // 表名
	Enabled    bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`          // 日志开启状态
}
