package config

type Wecom struct {
	Token  string `mapstructure:"token" json:"token" yaml:"token"`
	AppId  string `mapstructure:"appid" json:"appid" yaml:"appid"`
	AesKey string `mapstructure:"aeskey" json:"aeskey" yaml:"aeskey"`
}
