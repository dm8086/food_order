package sea

import (
	"order_food/global"
	"time"
)

type Order struct {
	global.SeaModel
	OrderId             string       `gorm:"column:order_id" json:"orderId"`
	AreaId              int          `gorm:"column:area_id" json:"areaId"`
	OrderNo             string       `gorm:"column:order_no" json:"orderNo"`
	DeskId              int          `gorm:"column:desk_id" json:"deskId"`
	DeskName            string       `gorm:"column:deak_name" json:"deakName"`
	OpenTime            *time.Time   `gorm:"column:open_time" json:"openTime"`
	SettleTime          *time.Time   `gorm:"column:settle_time" json:"settleTime"`
	CloseTime           *time.Time   `gorm:"column:close_time" json:"closeTime"`
	StoreId             string       `gorm:"column:store_id" json:"storeId"`
	StoreName           string       `gorm:"column:store_name" json:"storeName"`
	PeopleCount         int          `gorm:"column:people_count" json:"peopleCount"`
	Status              int          `gorm:"column:status" json:"status"`
	SettleStatus        int          `gorm:"column:settle_status" json:"settleStatus"`
	FromDeskName        string       `gorm:"column:from_desk_name" json:"fromDeskName"`
	FromDeskId          int          `gorm:"column:from_desk_id" json:"fromDeskId"`
	MemberId            int          `gorm:"column:member_id" json:"memberId"`
	Mobile              string       `gorm:"column:mobile" json:"mobile"`
	OrderAmount         int          `gorm:"column:order_amount" json:"orderAmount"`
	OrderReceivedAmount int          `gorm:"column:order_received_amount" json:"orderReceivedAmount"`
	PromoAmount         int          `gorm:"column:promo_amount" json:"promoAmount"`
	Remark              string       `gorm:"column:remark" json:"remark"`
	Src                 string       `gorm:"column:src" json:"src"`
	OrderBatchs         []OrderBatch `gorm:"foreignKey:order_id;references:order_id;" json:"orderBatchs"`
}

func (*Order) TableName() string {
	return "order"
}

type OrderBatch struct {
	global.SeaModel
	StoreId         string           `gorm:"column:store_id" json:"storeId"`
	AreaId          int              `gorm:"column:area_id" json:"areaId"`
	DeskId          int              `gorm:"column:desk_id" json:"deskId"`
	OrderId         string           `gorm:"column:order_id" json:"orderId"`
	BatchId         string           `gorm:"column:batch_id" json:"batchId"`
	Amount          int              `gorm:"column:amount" json:"amount"`
	Num             int              `gorm:"column:num" json:"num"`
	ExtendStatus    string           `gorm:"column:extend_status" json:"extendStatus"`
	Src             string           `gorm:"column:src" json:"src"`
	AddTime         *time.Time       `gorm:"column:add_time" json:"addTime"`
	OperatorType    string           `gorm:"column:operator_type" json:"operatorType"`
	OperatorId      int              `gorm:"column:operator_id" json:"operatorId"`
	OperatorName    string           `gorm:"column:operator_name" json:"operatorName"`
	Avatar          string           `gorm:"column:avatar" json:"avatar"`
	BatchType       int              `gorm:"column:batch_type" json:"batchType"`
	OrderBatchDishs []OrderBatchDish `gorm:"foreignKey:batch_id;references:batch_id;" json:"orderBatchDishs"`
}

func (*OrderBatch) TableName() string {
	return "order_batch"
}

type OrderBatchDish struct {
	global.SeaModel
	StoreId      string     `gorm:"column:store_id" json:"storeId"`
	AreaId       int        `gorm:"column:area_id" json:"areaId"`
	DeskId       int        `gorm:"column:desk_id" json:"deskId"`
	OrderId      string     `gorm:"column:order_id" json:"orderId"`
	BatchId      string     `gorm:"column:batch_id" json:"batchId"`
	Price        int        `gorm:"column:price" json:"price"`
	Amount       int        `gorm:"column:amount" json:"amount"`
	Num          int        `gorm:"column:num" json:"num"`
	Status       int        `gorm:"column:status" json:"status"`
	DishId       string     `gorm:"column:dish_id" json:"dishId"`
	DishName     string     `gorm:"column:dish_name" json:"dishName"`
	SkuId        string     `gorm:"column:sku_id" json:"skuId"`
	SkuName      string     `gorm:"column:sku_name" json:"skuName"`
	BatchType    int        `gorm:"column:batch_type" json:"batchType"`
	Src          string     `gorm:"column:src" json:"src"`
	Detail       string     `gorm:"column:detail" json:"detail"`
	AddTime      *time.Time `gorm:"column:add_time" json:"addTime"`
	OperatorType string     `gorm:"column:operator_type" json:"operatorType"`
	OperatorId   int        `gorm:"column:operator_id" json:"operatorId"`
	OperatorName string     `gorm:"column:operator_name" json:"operatorName"`
	Avatar       string     `gorm:"column:avatar" json:"avatar"`
	Deduct       int        `gorm:"column:deduct" json:"deduct"`
}

func (*OrderBatchDish) TableName() string {
	return "order_batch_dish"
}

type OrderLog struct {
	global.SeaModel
	OrderId    string `gorm:"column:order_id" json:"orderId"`
	OrderEvent int    `gorm:"column:order_event" json:"orderEvent"` // 日志事件
	Detail     string `gorm:"column:detail" json:"detail"`
}

func (*OrderLog) TableName() string {
	return "order_log"
}
