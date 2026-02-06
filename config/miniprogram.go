package config

type MiniProgram struct {
	AppId        string     `mapstructure:"appid" json:"appid" yaml:"appid"`
	SecretKey    string     `mapstructure:"secretkey" json:"secretkey" yaml:"secretkey"`
	Env          string     `mapstructure:"env" json:"env" yaml:"env"`
	MsgTemplates []Template `mapstructure:"msgTemplate" json:"msgTemplate" yaml:"msgTemplate"`
}

type Template struct {
	Name       string `mapstructure:"name" json:"name" yaml:"name"`
	TemplateId string `mapstructure:"templateId" json:"templateId" yaml:"templateId"`
	Page       string `mapstructure:"page" json:"page" yaml:"page"`
}
