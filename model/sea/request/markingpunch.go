package request

type MarkingConfAddReq struct {
	Name        string `json:"name"`
	DeviceType  int    `json:"deviceType"`
	DeviceConf  int    `json:"deviceConf"`
	Devices     string `json:"devices"`
	DeskConf    int    `json:"deskConf"`
	Desks       string `json:"desks"`
	SkuShowSort int    `json:"skuShowSort"`
	SkuMerge    int    `json:"skuMerge"`
	OrderFrom   int    `json:"orderFrom"`
	OrderType   int    `json:"orderType"`
	SkuStall    int    `json:"skuStall"`
	SkuStallIds string `json:"skuStallIds"`
}

type MarkingConfUpdateReq struct {
	Id          int    `json:"-"`
	Name        string `json:"name,omitempty"`
	DeviceType  int    `json:"deviceType,omitempty"`
	DeviceConf  int    `json:"deviceConf,omitempty"`
	Devices     string `json:"devices,omitempty"`
	DeskConf    int    `json:"deskConf,omitempty"`
	Desks       string `json:"desks,omitempty"`
	SkuShowSort int    `json:"skuShowSort,omitempty"`
	SkuMerge    int    `json:"skuMerge,omitempty"`
	OrderFrom   int    `json:"orderFrom,omitempty"`
	OrderType   int    `json:"orderType,omitempty"`
	SkuStall    int    `json:"skuStall,omitempty"`
	SkuStallIds string `json:"skuStallIds,omitempty"`
}

type MarkingPunchListReq struct {
	MarkTypeId int `json:"markTypeId"`
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
}

type MarkingConfListReq struct {
	DeviceType  *int `json:"deviceType" form:"deviceType"`
	DeviceConf  *int `json:"deviceConf" form:"deviceConf"`
	DeskConf    *int `json:"deskConf" form:"deskConf"`
	SkuShowSort *int `json:"skuShowSort" form:"skuShowSort"`
	SkuMerge    *int `json:"skuMerge" form:"skuMerge"`
	OrderFrom   *int `json:"orderFrom" form:"orderFrom"`
	OrderType   *int `json:"orderType" form:"orderType"`
	SkuStall    *int `json:"skuStall" form:"skuStall"`
	Page        int  `json:"page"`
	PageSize    int  `json:"pageSize"`
}

type MarkingPunchAddReq struct {
	Name       string `json:"name"`
	MarkTypeId uint   `json:"markTypeId"`
	TmpIds     string `json:"tmpIds"`
	Ip         string `json:"ip"`
}

type MarkingPunchUpdateReq struct {
	Id         int    `json:"-"`
	Name       string `json:"name"`
	MarkTypeId *uint  `json:"markTypeId"`
	TmpIds     string `json:"tmpIds"`
	Ip         string `json:"ip"`
}

type MarkingTmpAddReq struct {
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

type MarkingTmpUpdateReq struct {
	Id     int    `json:"-"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

type MarkingTmpListReq struct {
	MarkTypeId int `json:"markTypeId" form:"markTypeId"`
	Page       int `json:"page" form:"page"`
	PageSize   int `json:"pageSize" form:"pageSize"`
}

type MarkingBillAddReq struct {
	MarkingPunchId int `json:"markingPunchId"`
	BillType       int `json:"billType"`
}

type MarkingBillUpdateReq struct {
	Id             int `json:"-"`
	MarkingPunchId int `json:"markingPunchId"`
	BillType       int `json:"billType"`
}

type MarkingBillListReq struct {
	MarkingPunchId int `json:"markingPunchId" form:"markingPunchId"`
	BillType       int `json:"billType" form:"billType"`
	Page           int `json:"page" form:"page"`
	PageSize       int `json:"pageSize" form:"pageSize"`
}
