package config

type MongoDB struct {
	Uri string `mapstructure:"uri" json:"uri" yaml:"uri"` // 服务器地址
	DB  string `mapstructure:"db" json:"db" yaml:"db"`    // 数据库
}
