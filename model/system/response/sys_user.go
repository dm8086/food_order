package response

import (
	"order_food/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User               system.SysUser `json:"user"`
	LastLoginStoreId   string         `json:"lastLoginStoreId"`
	LastLoginStoreName string         `json:"lastLoginStoreName"`
	Token              string         `json:"token"`
	ExpiresAt          int64          `json:"expiresAt"`
}
