package request

type EventActionReq struct {
	ActionType  int    `json:"actionType"`
	StoreId     string `json:"storeId"`
	TableId     int    `json:"tableId"`
	QueueCodeId int    `json:"queueCodeId"`
	QueueCode   string `json:"queueCode"`
	Src         string `json:"src"`
	RingId      string `json:"ringId"`
}
