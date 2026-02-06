package request

type StallAddReq struct {
	StoreId string `json:"storeId"`
	Name    string `json:"name"`
	Sort    int    `json:"sort"`
	Remark  string `json:"remark"`
	Status  int    `json:"status"`
	IsShow  int    `json:"isShow"`
}

type StallEditReq struct {
	Id      int    `json:"id,omitempty"`
	StoreId string `json:"storeId,omitempty"`
	Name    string `json:"name,omitempty"`
	Status  *int   `json:"status,omitempty"`
	IsShow  *int   `json:"isShow,omitempty"`
	Sort    *int   `json:"sort,omitempty"`
	Remark  string `json:"remark,omitempty"`
}

type StallListReq struct {
	StoreId  string `form:"storeId" json:"storeId,omitempty"`
	Name     string `form:"name" json:"name,omitempty"`
	Remark   string `form:"remark" json:"remark,omitempty"`
	Sort     int    `form:"sort" json:"sort,omitempty"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type StallBingFoodAddReq struct {
	StoreId string `json:"storeId"`
	StallId int    `json:"stallId"`
	SkuId   string `json:"skuId"`
}

type StallBingFoodEditReq struct {
	StoreId string `json:"storeId"`
	StallId int    `json:"stallId"`
	SkuId   int    `json:"skuId"`
}

type StallBingFoodListReq struct {
	StoreId  string `form:"storeId" json:"storeId"`
	StallId  int    `form:"stallId" json:"stallId"`
	Sort     int    `form:"sort" json:"sort"`
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type StallSkuListReq struct {
	StoreId string `form:"storeId" json:"storeId"`
	StallId int    `form:"stallId" json:"stallId"`
	DeskId  int    `form:"deskId" json:"deskId"`
	AreaId  int    `form:"areaId" json:"areaId"`
	OrderId string `form:"orderId" json:"orderId"`
	BatchId string `form:"batchId" json:"batchId"`
}
