package sea

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Events struct {
	TableEvent
	OrderEvent
	DeviceEvent
}

const (
	DEVICEHEARTBEAT         = "devices/heartbeat"                 // 心跳topic
	TABLEEVENTLABEL         = "/table/event"                      // 桌台事件
	TABLEEVENT              = "storeId:%s:tableId:%d:event"       // 桌台事件
	TABLEEVENTLOCAL         = "storeId:%s:tableId:%d:local:event" // 本地桌台事件
	STOREEVENTLABEL         = "store/event"                       // 门店事件
	STOREEVENT              = "storeId:%s:event"                  // 门店事件
	STOREEVENTLOCAL         = "storeId:%s:local:event"            // 门店事件
	DEVICEREFRESHLABEL      = "device/refresh"                    // 设备页面刷新
	DEVICEREFRESH           = "storeId:%s:device:refresh:%s"      // 设备页面刷新
	TABLEORDERLABEL         = "table/order"                       // 订单事件
	TABLEORDER              = "storeId:%s:tableId:%s:order"       // 订单事件
	STOREORDERLABEL         = "store/stall"                       // 明档使用
	STOREORDER              = "storeId:%s:stallId%d"              // 明档使用
	TABLESCANLABEL          = "table/scan"                        // 买单扫描
	TABLESCAN               = "storeId:%s:scan:%d"                // 买单扫描
	STOREBROADCASTLABEL     = "store/broadcast"                   // 平板主题
	STOREBROADCAST          = "storeId:%s:broadcast:%s"           // 平板主题
	STOREQUEUELABEL         = "store/queue"                       // 排队叫号
	STOREQUEUE              = "storeId:%s:queue:%d"               // 排队叫号
	STORESERVICESETTLELABEL = "storeId:service:settle"            // 服务结算
	STORESERVICE            = "storeId:%s:service:settle"         // 服务结算
)

const (
	BusinessStatusEmpty       = 0  //空闲
	BusinessStatusOccupy      = 1  //占桌
	BusinessStatusDinnering   = 10 //用餐中
	BusinessStatusWaitService = 11 //待服务
	BusinessStatusWaitSettle  = 12 //待结算
	BusinessStatusWaitClean   = 13 //待收桌
	BusinessStatusServicing   = 21 //服务中
	BusinessStatusSettling    = 22 //结算中
	BusinessStatusCleaning    = 23 //收桌中
)

type EmqxLogs struct {
	Id      bson.ObjectId `bson:"_id" json:"_id"`
	Topic   string        `bson:"topic" json:"topic"`
	Msg     string        `bson:"msg" json:"msg"`
	Qos     int           `bson:"qos" json:"qos"`
	LogType string        `bson:"logType" json:"logType"`
	Created string        `bson:"created" json:"created"`
}

type TableInfo struct {
	TableId            uint       `json:"tableId"`
	TableName          string     `json:"tableName"`
	TableEvent         string     `json:"tableEvent"`
	BusinessStatus     int        `json:"businessStatus"`
	TableBeforeStatus  int        `json:"tableBeforeStatus"`
	OrderId            string     `json:"orderId"`
	OpenTime           *time.Time `json:"openTime"`
	LastServiceRequest time.Time  `json:"lastServiceRequest"`
	IsUnion            bool       `json:"isUnion"`
	IsOccupy           bool       `json:"isOccupy"`
	EventTime          time.Time  `json:"eventTime"`
	TransferTableId    uint       `json:"transferTableId"` //是否转桌
	QueueCode          string     `json:"queueCode"`
	StoreId            string     `json:"storeId"`
	EventType          string     `json:"eventType"`
	IsDelayed          int        `json:"isDelayed"`
	TypeId             int        `json:"typeId"`
	AreaName           string     `json:"areaName"`
	IsLock             bool       `json:"isLock"`
	UnionId            string     `json:"unionId"`
	SettleStatus       int        `json:"settleStatus"`
	MemberId           int        `json:"memberId"`
	MobileLabel        string     `json:"mobileLabel"`
	KryMemberId        string     `json:"kryMemberId"`
	IsBooking          bool       `json:"isBooking"`
	PeopleCount        int        `json:"peopleCount"`
	PeopleCountLabel   string     `json:"peopleCountLabel"`
}

// TableEvent 桌台事件
type TableEvent struct{}

const (
	EventTypeServiceRequest   = 2
	EventTypeSettleRequest    = 3
	EventTypeSingle           = 4
	EventTypeHold             = 5
	EventTypeOpen             = 6
	EventTypeClose            = 7
	EventTypeCancel           = 8
	EventTypeOccupy           = 9
	EventTypeCancelOccupy     = 10
	EventTypeUnion            = 11
	EventTypeCancelUnion      = 12
	EventTypeCancelRequest    = 13
	EventTypeStartService     = 21
	EventTypeStartSettle      = 22
	EventTypeStartClean       = 23
	EventTypeEndService       = 31
	EventTypeEndSettle        = 32
	EventTypeEndClean         = 33
	EventTypeDoubleClick      = 34
	EventTypeSoupStart        = 35
	EventTypeSoupEnd          = 36
	EventTypeDishChange       = 37
	EventTypeLock             = 38
	EventTypeUnlock           = 39
	EventTypeOccupyWithCode   = 40
	EventTypeOccupyUpdate     = 41
	EventTypeTableUsedSeat    = 42
	EventTypeTableStartSettle = 43
	EventTypeTableEndSettle   = 44
	EventTypeTableRefresh     = 45
	EventTypeTableChange      = 46
)

// EventIsValid 事件是否有效
func (TableEvent) EventIsValid(eventId int) bool {
	switch eventId {
	case EventTypeServiceRequest, EventTypeSettleRequest, EventTypeSingle, EventTypeHold, EventTypeOpen, EventTypeClose, EventTypeCancel,
		EventTypeOccupy, EventTypeCancelOccupy, EventTypeUnion, EventTypeCancelUnion, EventTypeCancelRequest,
		EventTypeStartService, EventTypeStartSettle, EventTypeStartClean, EventTypeDoubleClick,
		EventTypeEndService, EventTypeEndSettle, EventTypeEndClean, EventTypeTableRefresh, EventTypeTableChange:
		return true
	}
	return false
}

// GetEventLabel 获取事件的解释
func (TableEvent) GetEventLabel(eventId int) string {
	switch eventId {
	case EventTypeServiceRequest:
		return "顾客请求服务"
	case EventTypeSettleRequest:
		return "顾客请求买单"
	case EventTypeSingle:
		return "顾客请求服务"
	case EventTypeCancelRequest:
		return "顾客取消服务/买单"
	case EventTypeOpen: // 开台
		return "开始用餐。"
	case EventTypeClose: // 关闭
		return "桌台买单完成。"
	case EventTypeCancel: // 取消
		return "桌台订单取消。"
	case EventTypeOccupy: // 占桌
		return "准备用餐。"
	case EventTypeCancelOccupy: // 取消占桌
		return "取消用餐。"
	case EventTypeUnion: // 合桌
		return "开始用餐。"
	case EventTypeCancelUnion: // 取消合桌
		return "订单取消。"
	case EventTypeStartService: // 开始服务
		return "开始服务。"
	case EventTypeStartSettle: // 开始结算
		return "开始买单。"
	case EventTypeStartClean: // 开始收桌
		return "开始收桌。"
	case EventTypeEndService: // 服务完成
		return "服务完成。"
	case EventTypeEndSettle: // 买单完成
		return "买单服务完成。"
	case EventTypeEndClean: // 收桌完成
		return "收桌完成。"
	case EventTypeSoupStart: // 开始送锅底
		return "送锅底开始。"
	case EventTypeSoupEnd: // 完成送锅底
		return "送锅底完成。"
	case EventTypeDishChange: // 菜品变更
		return "菜品变更。"
	case EventTypeLock: // 锁定桌台
		return "桌台绑定。"
	case EventTypeUnlock: // 解锁桌台
		return "桌台解绑。"
	case EventTypeOccupyWithCode: //使用排队码占桌
		return "准备用餐。"
	case EventTypeOccupyUpdate: // 占桌更新
		return "桌台信息更新。"
	case EventTypeTableUsedSeat: // 修改可用桌台
		return "修改桌台可用座位。"
	case EventTypeTableStartSettle: //桌台开始结算
		return "桌台开始结算。"
	case EventTypeTableEndSettle: // 桌台买单完成
		return "桌台结算完成。"
	case EventTypeTableRefresh: // 桌台刷新
		return "桌台刷新。"
	case EventTypeTableChange: // 转台
		return "转台。"
	}
	return ""
}

// OrderEvent 订单事件(订单状态变更)
type OrderEvent struct{}

const (
	OrderStatusDEFAULT       = 0  // 默认值
	OrderStatusWAITPROCESSED = 1  // 待处理
	OrderStatusSUCCESS       = 2  // 已完成
	OrderStatusWAITSETTLED   = 3  // 待结账
	OrderStatusSELLTED       = 4  // 已结账
	OrderStatusREFUND        = 5  // 已退单
	OrderStatusCLOSED        = 6  // 已关闭
	OrderStatusINVALID       = 7  // 已作废
	OrderStatusCANCELLED     = 8  // 已取消
	OrderStatusREJECTED      = 9  // 已拒绝
	OrderStatusANTISETTLED   = 10 // 已反结账
)

func (OrderEvent) GetStatus(orderStatus int) string {
	switch orderStatus {
	case OrderStatusWAITPROCESSED:
		return "待处理"
	case OrderStatusSUCCESS:
		return "已完成"
	case OrderStatusWAITSETTLED:
		return "待结账"
	case OrderStatusSELLTED:
		return "已结账"
	case OrderStatusREFUND:
		return "已退单"
	case OrderStatusCLOSED:
		return "已关闭"
	case OrderStatusINVALID:
		return "已作废"
	case OrderStatusCANCELLED:
		return "已取消"
	case OrderStatusREJECTED:
		return "已拒绝"
	case OrderStatusANTISETTLED:
		return "已反结账"
	default:
		return ""
	}
}

// DeviceEvent 订单事件(订单状态变更)
type DeviceEvent struct{}

const (
	DevicesRefresh = 1 // 设备刷新
)

func (DeviceEvent) GetLabel(eventId int) string {
	switch eventId {
	case DevicesRefresh:
		return "设备刷新。"
	}
	return ""
}

type StallEvent struct{}
