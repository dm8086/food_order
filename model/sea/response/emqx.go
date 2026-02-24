package response

import (
	"order_food/model/sea"
	"time"
)

// 所有的桌台事件

type TableEventResp struct {
	TableInfo   TableInfo `json:"tableInfo"`
	Message     string    `json:"message"`
	IsAvailable int       `json:"isAvailable"`
	Src         string    `json:"src"`
}

// TableInfo 桌台事件
type TableInfo struct {
	TableId            uint       `json:"tableId"`
	TableName          string     `json:"tableName"`
	TableEvent         string     `json:"tableEvent"`
	BusinessStatus     int        `json:"businessStatus"`
	OrderId            string     `json:"orderId"`
	OpenTime           *time.Time `json:"openTime"`
	StoreId            string     `json:"storeId"`
	AreaName           string     `json:"areaName"`
	LastServiceRequest time.Time  `json:"lastServiceRequest"`
	IsUnion            bool       `json:"isUnion"`
	IsOccupy           bool       `json:"isOccupy"`
	EventTime          time.Time  `json:"eventTime"`
	TransferTableId    uint       `json:"transferTableId"` //是否转桌
	QueueCode          string     `json:"queueCode"`
	IsLock             bool       `json:"isLock"`
	UnionId            string     `json:"unionId"`
	SettleStatus       int        `json:"settleStatus"`
	MemberId           int        `json:"memberId"`
	MobileLabel        string     `json:"mobileLabel"`
	PeopleCount        int        `json:"peopleCount"`
	PeopleCountLabel   string     `json:"peopleCountLabel"`
}

// Extend 扩展需要处理的字段
type Extend struct {
	Soup              int `json:"soup"`
	Rice              int `json:"rice"`
	ProcessedProducts int `json:"processedProducts"`
}

// DispatchInfo 调度屏事件数据结构
type DispatchInfo struct {
	TableId        uint             `json:"tableId"`
	TableName      string           `json:"tableName"`
	TableEvent     string           `json:"tableEvent"`
	BusinessStatus int              `json:"businessStatus"`
	StoreId        string           `json:"storeId"`
	AreaName       string           `json:"areaName"`
	OrderId        string           `json:"orderId"`
	SubOrderId     string           `json:"subOrderId"` // 子订单id
	SubOrder       []sea.FoodOrderBatch `json:"subOrder"`   // 子订单列表
}

type RespTableItem struct {
	TableId            uint       `json:"tableId"`
	BusinessStatus     int        `json:"businessStatus"`
	LastServiceRequest time.Time  `json:"lastServiceRequest"`
	TableName          string     `json:"tableName"`
	Area               string     `json:"area"`
	AreaId             uint       `json:"areaId"`
	OrderId            string     `json:"orderId"`
	OpenTime           *time.Time `json:"openTime"`
	IsOccupy           bool       `json:"isOccupy"`
	IsUnion            bool       `json:"isUnion"`
	PeopleNum          int        `json:"peopleNum"`
	QrCode             string     `json:"qrCode"`
	QueueCode          string     `json:"queueCode"`
	InternalId         string     `json:"internalId"` //内部桌子id
	Status             int        `json:"status"`     // 桌台禁用状态
	IsVirtual          int        `json:"isVirtual"`
	IsLock             int        `json:"isLock"`
	VirtualTables      any        `json:"virtualTables"`
	EventData          any        `json:"eventData"`
	UnionTables        any        `json:"unionTables"`
	UnionId            string     `json:"unionId"`
	Table_id           uint       `json:"table_id,omitempty"`
	Table_name         string     `json:"table_name,omitempty"`
	UseNum             int        `json:"useNum"`
	Remark             string     `json:"remark"`
	TableCode          string     `json:"tableCode"`
	SettleStatus       int        `json:"settleStatus"`
	MobileLabel        string     `json:"mobileLabel"`
	IsBooking          bool       `json:"isBooking"`
	PeopleCount        int        `json:"peopleCount"`
	PeopleCountLabel   string     `json:"peopleCountLabel"`
	OrderCode          string     `json:"orderCode"`
	RingTime           int        `json:"ringTime"`
}

// DeskButtonResp 桌台按钮推送 买单和服务
type DeskButtonResp struct {
	TableId        int    `json:"tableId"`
	TableName      string `json:"tableName"`
	BusinessStatus int    `json:"businessStatus"`
	ButtonStatus   int    `json:"buttonStatus"`
	TableEvent     string `json:"tableEvent"`
	EventName      string `json:"eventName"` // buttonStatus
	EventTime      string `json:"eventTime"`
}

// DeskStatusResp 桌台状态推送
type DeskStatusResp struct {
	TableId          int    `json:"tableId"`
	TableName        string `json:"tableName"`
	BusinessStatus   int    `json:"businessStatus"`
	TableEvent       string `json:"tableEvent"`
	EventName        string `json:"eventName"` // businessStatus
	EventTime        string `json:"eventTime"`
	IsLock           bool   `json:"isLock"`
	IsOccupy         bool   `json:"isOccupy"`
	MemberId         int    `json:"memberId"`
	MemberUsedDy     bool   `json:"memberUsedDy"`
	MobileLabel      string `json:"mobileLabel"`
	OrderCode        string `json:"orderCode"`
	OrderId          string `json:"orderId"`
	PeopleCountLabel string `json:"peopleCountLabel"`
	PeopleCount      int    `json:"peopleCount"`
	QueueCode        string `json:"queueCode"`
	SettleStatus     int    `json:"settleStatus"`
	StoreId          string `json:"storeId"`
}

// DeskSoupResp 锅底状态推送
type DeskSoupResp struct {
	TableId    int    `json:"tableId"`
	SoupStatus int    `json:"soupStatus"`
	SoupList   any    `json:"soupList"`
	TableEvent string `json:"tableEvent"`
	EventName  string `json:"eventName"` // soupStatus
	EventTime  string `json:"eventTime"`
	StoreId    string `json:"storeId"`
}

// TopScreenResp 天上屏幕的内容
type TopScreenResp struct {
	Id         uint   `json:"id"`
	TableId    int    `json:"tableId"`
	TableName  string `json:"tableName"`
	Status     int    `json:"status"`     // 0 1 2 3
	TableEvent string `json:"tableEvent"` // 买单还是服务  取消
	ShopList   any    `json:"shopList"`   // ShopList  商品列表
	ShopType   int    `json:"shopType"`   // ShopType  商品类型  买单 服务 锅底 加工类  饮料  1 2 3 4 5 @todo
	EventType  int    `json:"eventType"`  // eventType  事件类型  买单 服务 取消
	EventName  string `json:"eventName"`  // topScreenStatus
	EventTime  string `json:"eventTime"`  // 事件时间
	StoreId    string `json:"storeId"`    // 门店id
	Num        int    `json:"num"`        // 数量
}
