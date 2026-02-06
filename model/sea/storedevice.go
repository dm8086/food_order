package sea

import (
	"order_food/global"
	"time"
)

type StoreDevice struct {
	global.SeaModel
	DeviceUuid    string      `gorm:"column:device_uuid;" json:"deviceUuid"`
	TypeId        uint        `gorm:"column:type_id;" json:"typeId"`
	Name          string      `gorm:"column:name;" json:"name"`
	StoreId       string      `gorm:"column:store_id;" json:"storeId"`
	TableId       int         `gorm:"column:table_id;" json:"tableId"`
	Status        int         `gorm:"column:status;" json:"status"`              //设备状态（1：启用中，2：禁用中）
	OnlineStatus  int         `gorm:"column:online_status;" json:"onlineStatus"` //设备在线状态（1：在线，2：离线）
	HeartbeatTime *time.Time  `gorm:"column:heartbeat_time" json:"heartbeatTime"`
	SortNum       int         `gorm:"column:sort_num;" json:"sortNum"`
	WebUrl        string      `gorm:"column:web_url" json:"webUrl"`
	ShowWebUrl    string      `gorm:"column:show_web_url" json:"showWebUrl"`
	WorkingStatus int         `gorm:"column:working_status" json:"workingStatus"` //设备自定义工作状态（如：有人/无人等）
	DeviceType    *DeviceType `gorm:"foreignKey:TypeId;references:ID;comment:设备分类" json:"deviceType"`
	TableLabel    string      `gorm:"-" json:"tableLabel"`  // 桌子显示的名称
	TableQrCode   string      `gorm:"-" json:"tableQrCode"` // 桌子二维码
	Version       string      `gorm:"column:version" json:"version"`
	Remark        string      `gorm:"column:remark" json:"remark"`        // 设备备注
	WorkLabel     string      `gorm:"column:work_label" json:"workLabel"` // 工作状态
	AppVersion    string      `gorm:"column:app_version" json:"appVersion"`
	RingTime      int         `gorm:"-" json:"ringTime"`
}

func (StoreDevice) TableName() string {
	return "store_device"
}

type DeviceType struct {
	global.SeaModel
	TypeName string `gorm:"column:type_name;" json:"typeName"`
	TypeKey  string `gorm:"type_key" json:"typeKey"`
	Remark   string `gorm:"column:remark;" json:"remark"`
	WebUrl   string `gorm:"column:web_url;" json:"webUrl"`
	SortNum  int    `gorm:"column:sort_num;" json:"sortNum"`
	Status   int    `gorm:"column:status;" json:"status"`
}

func (DeviceType) TableName() string {
	return "device_type"
}
