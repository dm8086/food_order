package system

import "order_food/global"

type SystemConst struct {
	global.SeaModel
	Key       string `gorm:"column:key" json:"key"`
	Title     string `gorm:"column:title" json:"title"`
	Value     string `gorm:"column:value" json:"value"`
	ConstType string `gorm:"column:const_type" json:"constType"`
	Remark    string `gorm:"column:remark" json:"remark"`
	Upload    int    `gorm:"column:upload" json:"upload"`
}

func (SystemConst) TableName() string {
	return "system_const"
}
