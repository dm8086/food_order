package request

type OpenDeskReq struct {
	TableId   int    `json:"tableId"`
	Num       int    `json:"num"`
	Src       string `json:"src"`
	QueueCode string `json:"queueCode"`
	CodeId    int    `json:"codeId"`
	CartId    string `json:"cartId"`
	RequestId string `json:"requestId"`
}

type EventRecordReq struct {
	StoreId   string `json:"storeId"`
	StoreName string `json:"storeName"`
	TableId   uint   `json:"tableId"`
	DeskName  string `json:"deskName"`
	AreaId    uint   `json:"areaId,omitempty"`
	AreaName  string `json:"areaName"`
	EventType int    `json:"eventType"`
	IsInvalid bool   `json:"isInvalid"`
	Status    int    `json:"status"`
}

type TableListReq struct {
	StoreId        string `json:"storeId"`
	HumanNum       *int   `json:"humanNum"`
	AreaIds        []int  `json:"areaIds"`
	BusinessStatus []int  `json:"businessStatus"`
}

type ChangeDeskReq struct {
	StoreId     string `json:"storeId"`
	FromTableId int    `json:"fromTableId"`
	ToTableId   int    `json:"toTableId"`
}

type DeskJointReq struct {
	StoreId     string `json:"storeId"`
	FromTableId int    `json:"fromTableId"`
	ToTableId   int    `json:"toTableId"`
}
