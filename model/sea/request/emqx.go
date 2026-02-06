package request

type DevicesHeartbeatReq struct {
	DeviceUUID    string `json:"deviceUuid"`
	OnlineStatus  *int   `json:"onlineStatus"`
	WorkingStatus *int   `json:"workingStatus"`
	Version       string `json:"version"`
	WorkLabel     string `json:"workLabel"`
}

type TableEventReq struct {
	EventType   int    `json:"eventType"`   //事件类型
	StoreId     string `json:"storeId"`     //门店id
	TableId     int    `json:"tableId"`     //桌台id
	QueueCodeId int    `json:"queueCodeId"` //排队码编号
	QueueCode   string `json:"queueCode"`   //排队码
	Src         string `json:"src"`         //桌台来源
	RingId      string `json:"ringId"`      // 手环id
}
