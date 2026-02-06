package config

type NoApilog struct {
	Apis string `json:"apis"` // 需要忽略的api用','隔开
}
