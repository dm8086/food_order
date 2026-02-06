package config

type Etcd struct {
	Uri     string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	Lease   int64  `mapstructure:"lease" json:"lease" yaml:"lease"`
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	ServUri string `mapstructure:"servUri" json:"servUri" yaml:"servUri"`
}
