package config

type Emqx struct {
	Broker        string `mapstructure:"broker" json:"broker" yaml:"broker"`
	ClientID      string `mapstructure:"clientID" json:"clientID" yaml:"clientID"`
	Username      string `mapstructure:"username" json:"username" yaml:"username"`
	Password      string `password:"password" json:"password" yaml:"password"`
	CleanSession  bool   `mapstructure:"cleanSession" json:"cleanSession" yaml:"cleanSession"`
	KeepAlive     int    `mapstructure:"keepAlive" json:"keepAlive" yaml:"keepAlive"`
	AutoReconnect bool   `mapstructure:"autoReconnect" json:"autoReconnect" yaml:"autoReconnect"`
}
