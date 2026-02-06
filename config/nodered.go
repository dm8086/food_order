package config

type NodeRed struct {
	Host  string `mapstructure:"host" json:"host" yaml:"host"`    // 服务器地址:端口
	Token string `mapstructure:"token" json:"token" yaml:"token"` //token
}
