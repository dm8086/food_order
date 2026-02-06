package request

type ButtonServiceReq struct {
	StoreId   string `json:"storeId"`
	TableId   int    `json:"tableId"`
	EventType int    `json:"eventType"` // 2服务 3买单 4单击 13取消 34双击
	RindId    string `json:"rindId"`    // 手环id
	RindName  string `json:"rindName"`  // 手环名称
}
