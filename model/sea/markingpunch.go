package sea

import "order_food/global"

// Markings 打印机管理
type MarkingPunch struct {
	global.SeaModel
	Name       string `gorm:"column:name" json:"name"`
	MarkTypeId uint   `gorm:"column:mark_type_id" json:"markTypeId"`
	TmpIds     string `gorm:"column:tmp_ids" json:"tmpIds"` // tmpIds 模板信息
	Ip         string `gorm:"column:ip" json:"ip"`
}

func (*MarkingPunch) TableName() string {
	return "marking_punch"
}

// MarkingPunch 打印机类型
type MarkingConf struct {
	global.SeaModel
	Name        string `gorm:"column:name" json:"name"`
	DeviceType  int    `gorm:"column:device_type" json:"deviceType"`    // deviceType 0 1收银出票 2后厨出票 3标签出票
	DeviceConf  int    `gorm:"column:device_conf" json:"deviceConf"`    // deviceConf 是否关联其他设备 0全部设备 1 部分设备
	Devices     string `gorm:"column:devices" json:"devices"`           // devices 关联设备集合
	DeskConf    int    `gorm:"column:desk_conf" json:"deskConf"`        // deskConf 桌台配置 0全部桌台 1部分桌台
	Desks       string `gorm:"column:desks" json:"desks"`               // desks  关联桌台集合
	SkuShowSort int    `gorm:"column:sku_show_sort" json:"skuShowSort"` // SkuShowSort 0下单时间 1菜品分类
	SkuMerge    int    `gorm:"column:sku_merge" json:"skuMerge"`        // skuMerge 0合并 1拆开 2 购物车
	OrderFrom   int    `gorm:"column:order_from" json:"orderFrom"`      //  orderFrom 订单来源 1点餐app下单 2扫码下单
	OrderType   int    `gorm:"column:order_type" json:"orderType"`      //  OrderType 订单类型 0堂食
	SkuStall    int    `gorm:"column:sku_stall" json:"skuStall"`        // SkuStall 菜品分档出品 0 不分档 1按菜品 2按分类 3按做法
	SkuStallIds string `gorm:"column:sku_stall_ids" json:"skuStallIds"` // SkuStallIds  菜品分档出品 集合
}

func (*MarkingConf) TableName() string {
	return "marking_conf"
}

type MarkingTmp struct {
	global.SeaModel
	Name   string `gorm:"column:name" json:"name"`     // 模板名称
	Detail string `gorm:"column:detail" json:"detail"` // 模板详情
}

func (*MarkingTmp) TableName() string {
	return "marking_tmp"
}

type Bills struct {
	global.SeaModel
	MarkingPunchId uint `gorm:"column:marking_punch_id" jorn:"markingPunchId"` // 打印机id
	BillType       int  `gorm:"column:bill_type" json:"billType"`
}

func (*Bills) TableName() string {
	return "bill"
}
