package sea

import (
	"order_food/global"
	"time"
)

type Table struct {
	global.SeaModel
	StoreId       string         `gorm:"column:store_id;comment:门店id" json:"storeId"`           // 门店id
	AreaId        uint           `gorm:"column:area_id;comment:区域id" json:"areaId"`             // 区域id
	Name          string         `gorm:"column:table_name;comment:桌台名称" json:"tableName"`       // 桌台名称
	HumanNum      int            `gorm:"column:human_num;comment:可容纳人数" json:"humanNum"`        // 可容人数
	SortNo        int            `gorm:"column:sort_no;comment:排序号" json:"sortNo"`              // 排序号
	OperationId   uint           `gorm:"column:operation_id;comment:操作人id" json:"operationId"`  // 操作人id
	Status        int            `gorm:"column:status;comment:状态" json:"status"`                // 状态 1:启用 2:禁用
	TableCode     string         `gorm:"column:table_code;comment:桌台编码" json:"tableCode"`       // 桌台编码
	KryTableId    string         `gorm:"column:kry_table_id;comment:客如云桌台id" json:"kryTableId"` // 客如云桌台id
	PosX          int            `gorm:"column:pos_x;comment:x坐标" json:"posX"`                  //
	PosY          int            `gorm:"column:pos_y;comment:x坐标" json:"posY"`                  //
	PosWidth      int            `gorm:"column:pos_width;comment:宽度" json:"posWidth"`           //
	PosHeight     int            `gorm:"column:pos_height;comment:高度" json:"posHeight"`         //
	QrCode        string         `gorm:"column:qrCode;comment:点餐码" json:"qrCode"`               //点餐码
	IsBox         int            `gorm:"column:is_box;comment:是否包厢" json:"isBox"`               //是否包厢
	InternalId    string         `gorm:"column:internal_id;comment:内部桌子id" json:"internalId"`   //内部桌子id
	IsVirtual     int            `gorm:"column:is_virtual" json:"isVirtual"`                    // 是否虚拟桌台
	IsFixed       int            `gorm:"column:is_fixed" json:"isFixed"`                        // 是否固定桌台
	IsLock        int            `gorm:"column:is_lock" json:"isLock"`                          // 是否绑定
	VirtualData   string         `gorm:"column:virtual_data" json:"virtualData"`                // 虚拟桌台json
	UseNum        int            `gorm:"column:use_num" json:"useNum"`                          // 可用的电磁炉数量
	Remark        string         `gorm:"column:remark" json:"remark"`
	Business      *TableBusiness `gorm:"foreignKey:ID;references:table_id;comment:业务信息" json:"business"`
	StoreArea     *StoreArea     `gorm:"foreignKey:AreaId;references:ID;comment:区域" json:"area"`
	VirtualTables any            `gorm:"-" json:"virtualTables"`
}

func (*Table) TableName() string {
	return "table"
}

type TableBusiness struct {
	global.SeaModel
	TableId        uint       `gorm:"column:table_id;comment:桌台id" json:"tableId"`             // 桌台id
	OrderId        string     `gorm:"column:order_id;comment:订单编号" json:"orderId"`             // 订单编号
	CustomerNum    int        `gorm:"column:customer_num;comment:用餐人数" json:"customerNum"`     // 用餐人数
	BusinessStatus int        `gorm:"column:business_status;comment:状态" json:"businessStatus"` // 业务状态：空闲（0） 占位(1) 用餐中（10）待服务（11）服务中（21）待买单（12）买单中（22）待收桌（13）收桌中（23）
	IsUnion        int        `gorm:"column:is_union;comment:是否合桌" json:"isUnion"`             //是否合桌
	IsOccupy       int        `gorm:"column:is_occupy;comment:是否占桌" json:"isOccupy"`           //是否占桌
	OpenTime       *time.Time `gorm:"column:open_time;comment:开台时间" json:"openTime"`           //开台时间
	StoreId        string     `gorm:"column:store_id;comment:门店编号" json:"storeId"`             //门店编号
	QueueCode      string     `gorm:"column:queue_code;comment:排队号"`                           //排队号
	QueueCodeId    int        `gorm:"column:queue_code_id;comment:排队码编号"`                      //排队码编号
	Src            string     `gorm:"column:src;comment:事件来源" json:"src"`                      // 来源
	ExtendStatus   string     `gorm:"column:extend_status" json:"extendStatus"`                // 扩展订单状态
	UnionId        string     `gorm:"column:union_id" json:"unionId"`                          // 联合表的id
	SettleStatus   int        `gorm:"column:settle_status" json:"settleStatus"`
	Mobile         string     `gorm:"column:mobile" json:"mobile"`            // 电话号码
	SettleServ     int        `gorm:"column:settle_serv" json:"settleServ"`   // 结算服务
	ServiceServ    int        `gorm:"column:service_serv" json:"serviceServ"` // 服务
}

func (TableBusiness) TableName() string {
	return "table_business"
}

type TableBusinessLog struct {
	global.SeaModel
	TableId   uint      `gorm:"column:table_id;comment:桌台id" json:"tableId"`           // 桌台id
	OrderId   string    `gorm:"column:order_id" json:"orderId"`                        //订单编号
	LogType   int       `gorm:"column:log_type;comment;日志类型" json:"logType"`           // 业务状态 1:启用 2:禁用
	EventType int       `gorm:"column:event_type;comment:事件类型" json:"eventType"`       // 事件类型
	LogTime   time.Time `gorm:"column:open_time;comment:开台时间" json:"openTime"`         // 日志时间
	LogReq    string    `gorm:"type:json;column:log_req;comment:日志请求内容" json:"logReq"` // 日志请求内容
	LogResp   string    `gorm:"type:json;log_resp;comment:日志响应内容" json:"logResp"`      // 日志响应内容
}

func (TableBusinessLog) TableName() string {
	return "table_business_log"
}

// StoreInfo 门店信息
type Store struct {
	global.SeaModel
	StoreId        string `gorm:"column:store_id;comment:门店id" json:"storeId"`             // 门店id
	StoreName      string `gorm:"column:store_name;comment:门店名称" json:"storeName"`         // 门店名称
	StoreStatus    int    `gorm:"column:store_status;comment:门店状态" json:"storeStatus"`     // 门店状态 1:启用 2:禁用
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                  // 备注
	SortNo         int    `gorm:"column:sort_no;comment:排序号" json:"sortNo"`                // 排序号
	OperationId    uint   `gorm:"column:operation_id;comment:操作人id" json:"operationId"`    // 操作人id
	KryStoreId     int64  `gorm:"column:kry_store_id;comment:客如云门店id" json:"kryStoreId"`   // 客如云门店id
	Address        string `gorm:"column:address;comment:门店地址" json:"address"`              // 门店地址
	BusinessTime   string `gorm:"column:business_time;comment:营业时间描述" json:"businessTime"` // 营业时间描述
	Latitude       string `gorm:"column:latitude;comment:纬度" json:"latitude"`              // 地图定位
	Longitude      string `gorm:"column:longitude;comment:经度" json:"longitude"`            // 地图定位
	Tel            string `gorm:"column:tel;comment:门店联系电话" json:"tel"`                    // 门店联系电话
	BusinessStatus int    `gorm:"column:business_status" json:"businessStatus"`            // 营业状态
	HostIp         string `gorm:"column:host_ip" json:"hostIp"`                            // nr主机IP
	OpenTime       string `gorm:"column:open_time" json:"openTime"`                        // 开始营业时间
	ClosedTime     string `gorm:"column:closed_time" json:"closedTime"`                    // 结束营业时间
	StoreType      int    `gorm:"column:store_type" json:"storeType"`                      // 门店类型  1门店 2仓库
	ParentId       int    `gorm:"column:parent_id" json:"parentId"`                        // 上级id
	AreaId         int    `gorm:"column:area_id" json:"areaId"`                            // 区域id
}

func (*Store) TableName() string {
	return "sea_store"
}

type StoreArea struct {
	global.SeaModel
	StoreId     string `gorm:"column:store_id;comment:门店id" json:"storeId"`          // 门店id
	AreaName    string `gorm:"column:area_name;comment:区域名称" json:"areaName"`        // 区域名称
	OperationId uint   `gorm:"column:operation_id;comment:操作人id" json:"operationId"` // 操作人id
	SortNo      int    `gorm:"column:sort_no;comment:排序号" json:"sortNo"`             // 排序号
	Status      int    `gorm:"column:status;comment:状态" json:"status"`               // 状态 1:启用 2:禁用
	KryAreaId   string `gorm:"column:kry_area_id;comment:客如云区域id" json:"kryAreaId"`  // 客如云区域id
}

func (StoreArea) TableName() string {
	return "sea_store_area"
}

// 联合桌台
type VirtualTableConnection struct {
	Id        int        `gorm:"column:id" json:"id"`
	VTid      string     `gorm:"column:vt_id" json:"vtid"`                    // 虚拟桌台id
	TableId   uint       `gorm:"column:table_id;comment:桌台id" json:"tableId"` // 桌台id
	Status    int        `gorm:"column:status" json:"status"`                 // 状态
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt"`
	TableInfo *Table     `gorm:"foreignKey:TableId;references:id"`
}

func (VirtualTableConnection) TableName() string {
	return "virtual_table_connection"
}

type LocalTable struct {
	global.SeaModel
	StoreId        string `gorm:"column:store_id" json:"storeId"`
	DeskId         int    `gorm:"column:desk_id" json:"deskId"`
	DeskName       string `gorm:"column:desk_name" json:"deskName"`
	AreaId         int    `gorm:"column:area_id" json:"areaId"`
	AreaName       string `gorm:"column:area_name" json:"areaName"`
	Status         int    `gorm:"column:status" json:"status"`
	SettleStatus   int    `gorm:"column:settle_status" json:"settleStatus"`
	BusinessStatus int    `gorm:"column:business_status" json:"businessStatus"`
	OrderId        string `gorm:"column:order_id" json:"orderId"`
	QrCode         string `gorm:"column:qr_code" json:"qrCode"`
	OrderCode      string `gorm:"column:order_code" json:"orderCode"`
	DefaultCode    string `gorm:"column:default_code" json:"defaultCode"`
	PeopleCount    int    `gorm:"column:people_count" json:"peopleCount"`
	ExtendStatus   string `gorm:"column:extend_status" json:"extendStatus"`
	QueueCodeId    int    `gorm:"column:queue_code_id" json:"queueCodeId"`
	QueueCode      string `gorm:"column:queue_code" json:"queueCode"`
	IsVirtual      int    `gorm:"column:is_virtual" json:"isVirtual"`
	VirtualData    string `gorm:"column:virtual_data" json:"virtualData"`
	IsUnion        int    `gorm:"column:is_union" json:"isUnion"`
	UnionData      string `gorm:"column:union_data" json:"unionData"`
	IsOccupy       int    `gorm:"column:is_occupy" json:"isOccupy"`
	IsLock         int    `gorm:"column:is_lock" json:"isLock"`
	Src            string `gorm:"column:src" json:"src"`
}

func (*LocalTable) TableName() string {
	return "local_tables"
}
