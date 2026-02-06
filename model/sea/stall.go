package sea

import "order_food/global"

type StoreStall struct {
	global.SeaModel
	StoreId string `gorm:"column:store_id" json:"storeId"`
	Name    string `gorm:"column:name" json:"name"`
	Status  int    `gorm:"column:status" json:"status"`
	IsShow  int    `gorm:"column:is_show" json:"isShow"`
	Sort    int    `gorm:"column:sort" json:"sort"`
	Remark  string `gorm:"column:remark" json:"remark"`
}

func (*StoreStall) TableName() string {
	return "store_stall"
}

type StallBingFood struct {
	global.SeaModel
	StoreId string `gorm:"column:store_id" json:"storeId"`
	StallId int    `gorm:"column:stall_id" json:"stallId"`
	SkuId   string `gorm:"column:sku_id" json:"skuId"`
	IsSplit int    `gorm:"column:is_split" json:"isSplit"`
}

func (*StallBingFood) TableName() string {
	return "stall_bing_food"
}

type StallSkus struct {
	global.SeaModel
	StoreId     string         `gorm:"column:store_id" json:"storeId"`
	AreaId      int            `gorm:"column:area_id" json:"areaId"`
	DeskId      int            `gorm:"column:desk_id" json:"deskId"`
	DeskName    string         `gorm:"column:desk_name" json:"deskName"`
	OrderId     string         `gorm:"column:order_id" json:"orderId"`
	BatchId     string         `gorm:"column:batch_id" json:"batchId"`
	StallId     int            `gorm:"column:stall_id" json:"stallId"`
	SkuId       string         `gorm:"column:sku_id" json:"skuId"`
	Num         int            `gorm:"column:num" json:"num"`
	BatchDishId int            `gorm:"column:batch_dish_id" json:"batchDishId"`
	Status      int            `gorm:"column:status" json:"status"` // 未加工 1送出 2完成 3退单
	DishName    string         `gorm:"column:dish_name" json:"dishName"`
	Detail      OrderBatchDish `gorm:"foreignKey:batch_dish_id;references:id;" json:"detail"`
}

func (*StallSkus) TableName() string {
	return "stall_skus"
}
