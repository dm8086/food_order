package response

import "order_food/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
