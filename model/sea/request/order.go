package request

import "order_food/model/sea/response"

type OrderAddReq struct {
	OrderId     string         `json:"orderId"`
	OrderNo     string         `json:"orderNo"`
	RequestId   string         `json:"requestId"`
	MemberId    int            `json:"memberId"`
	TableId     int            `json:"tableId"`
	StoreId     string         `json:"storeId"`
	PeopleCount int            `json:"peopleCount"`
	Remark      string         `json:"remark"`
	Src         string         `json:"src"`
	GoodsList   []response.Sku `json:"goodsList"`
	Dishes      []DishInfo     `json:"dishes"`
	AddType     string         `json:"addType"`
}

type OrderSubReq struct {
	OrderId   string         `json:"orderId"`
	RequestId string         `json:"requestId"`
	TableId   int            `json:"tableId"`
	StoreId   string         `json:"storeId"`
	GoodsList map[string]int `json:"goodsList"`
}

type DishInfo struct {
	DishId         string         `json:"dishId"`
	SkuId          string         `json:"skuId"`
	DishName       string         `json:"dishName"`
	DishType       string         `json:"dishType"`
	Quantity       int            `json:"quantity"`
	Price          int            `json:"price"`
	Amount         int            `json:"amount"`
	ComboSubDishes []comboSubDish `json:"comboSubDishes"`
	BatchId        string         `json:"batchId"`
}

type comboSubDish struct {
	DishId   string `json:"dishId"`
	SkuId    string `json:"skuId"`
	DishName string `json:"dishName"`
	SkuName  string `json:"skuName"`
	Quantity int    `json:"quantity"`
}

type OrderAddBase struct {
	DishId         string          `json:"dishId"`
	DishType       string          `json:"dishType"` // 单盘还是套餐
	SkuId          string          `json:"skuId"`
	Quantity       int             `json:"quantity"`
	ComboSubDishes []ComboSubDishe `json:"comboSubDishes"`
	CookingWay     CookingWay      `json:"cookingWay"`
}
type CookingWay struct {
	AddPrice       int    `json:"addPrice"`
	CookingwayId   string `json:"cookingwayId"`
	CookingwayName string `json:"cookingwayName"`
}

type ComboSubDishe struct {
	DishId   string `json:"dishId"`
	SkuId    string `json:"skuId"`
	Quantity int    `json:"quantity"`
	Amount   int    `json:"amount"`
}

type OrderUpdateReq struct {
	OrderId             string `json:"orderId"`
	RequestId           string `json:"requestId"`
	Status              *int   `json:"status"`
	SettleStatus        *int   `json:"settleStatus"`
	OrderAmount         *int   `json:"orderAmount"`
	OrderReceivedAmount *int   `json:"orderReceivedAmount"`
	PromoAmount         *int   `json:"promoAmount"`
	Coupons             string `json:"coupons"`
}

type OrderAdd struct {
	Src       string         `json:"src"`
	DeskId    int            `json:"deskId"`
	RequestId string         `json:"requestId"`
	StoreId   string         `json:"storeId"`
	Skus      []response.Sku `json:"skus"`
}
type Sku struct {
	SkuId   string `json:"skuId"`
	Num     int    `json:"num"`
	SkuType string `json:"skuType"`
	Skus    []Sku  `json:"skus"`
}

type OrderSkuRemoveReq struct {
	OrderId          string `json:"orderId"`
	BatchId          string `json:"batchId"`
	StallId          int    `json:"stallId"`
	OrderBatchDishId int    `json:"orderBatchDishId"`
}

type OrderListReq struct {
	StoreId  string `json:"storeId" form:"storeId"`
	AreaId   int    `json:"areaId" form:"areaId"`
	Status   *int   `json:"status" form:"status"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Sort     int    `json:"sort" form:"sort"`
}

type SoupRemoveReq struct {
	OrderId string `json:"orderId"`
	Ids     []int  `json:"ids"`
	Status  *int   `json:"status"`
}

type CombWriteoffReq struct {
	RequestId string `json:"requestId"`
	OrderId   string `json:"orderId"`
	CombId    string `json:"combId"`
}

// ----------------
type KposAddDishesToOrderReqBody struct {
	OrderID               string     `json:"orderId"`
	BusinessFormat        string     `json:"businessFormat"`
	SceneCode             string     `json:"sceneCode"`
	GoodList              []GoodItem `json:"goodList"`
	OperateItemSourceType string     `json:"operateItemSourceType,omitempty"`
}

type GoodItem struct {
	SceneCode               string                `json:"sceneCode,omitempty" bson:"sceneCode,omitempty"`
	SkuID                   string                `json:"skuId,omitempty" bson:"skuId,omitempty"`
	SpuId                   string                `json:"spuId,omitempty" bson:"spuId,omitempty"`
	TemporaryDishFlag       string                `json:"temporaryDishFlag,omitempty" bson:"temporaryDishFlag,omitempty"`
	MustOrderFlag           string                `json:"mustOrderFlag,omitempty" bson:"mustOrderFlag,omitempty"`
	ItemName                string                `json:"itemName,omitempty" bson:"itemName,omitempty"`
	ItemType                string                `json:"itemType,omitempty" bson:"itemType,omitempty"`
	ItemExtType             string                `json:"itemExtType,omitempty" bson:"itemExtType,omitempty"`
	TicketExitIdList        []string              `json:"ticketExitIdList,omitempty" bson:"ticketExitIdList,omitempty"`
	ItemPriceStr            string                `json:"itemPriceStr,omitempty" bson:"itemPriceStr,omitempty"`
	ItemNumStr              string                `json:"itemNumStr,omitempty" bson:"itemNumStr,omitempty"`
	WeighDishFlag           string                `json:"weighDishFlag,omitempty" bson:"weighDishFlag,omitempty"`
	UnitID                  string                `json:"unitId,omitempty" bson:"unitId,omitempty"`
	UnitName                string                `json:"unitName,omitempty" bson:"unitName,omitempty"`
	SecondUnitID            string                `json:"secondUnitId,omitempty" bson:"secondUnitId,omitempty"`
	SecondUnitName          string                `json:"secondUnitName,omitempty" bson:"secondUnitName,omitempty"`
	RequisiteCookMethodFlag string                `json:"requisiteCookMethodFlag,omitempty" bson:"requisiteCookMethodFlag,omitempty"`
	RequisiteSideDishFlag   string                `json:"requisiteSideDishFlag,omitempty" bson:"requisiteSideDishFlag,omitempty"`
	CookAttachGroupList     []CookAttachGroup     `json:"cookAttachGroupList,omitempty" bson:"cookAttachGroupList,omitempty"`
	AddSpiceAttachGroupList []AddSpiceAttachGroup `json:"addSpiceAttachGroupList,omitempty" bson:"addSpiceAttachGroupList,omitempty"`
	ComboGroupList          []ComboGroup          `json:"comboGroupList,omitempty" bson:"comboGroupList,omitempty"`
	JoinManualPromoFlag     string                `json:"joinManualPromoFlag,omitempty" bson:"joinManualPromoFlag,omitempty"`
	OrderUnitID             string                `json:"orderUnitId,omitempty" bson:"orderUnitId,omitempty"`
	OrderUnitName           string                `json:"orderUnitName,omitempty" bson:"orderUnitName,omitempty"`
	OrderUnitNum            string                `json:"orderUnitNum,omitempty" bson:"orderUnitNum,omitempty"`
	DoubleUnitWeighDishFlag string                `json:"doubleUnitWeighDishFlag,omitempty" bson:"doubleUnitWeighDishFlag,omitempty"`
	CartItemNoteList        []string              `json:"cartItemNoteList,omitempty" bson:"cartItemNoteList,omitempty"`
	CartItemID              string                `json:"cartItemId,omitempty" bson:"cartItemId,omitempty"`
	ChangePriceFlag         string                `json:"changePriceFlag,omitempty" bson:"changePriceFlag,omitempty"`
	WaitCallFlag            string                `json:"waitCallFlag,omitempty" bson:"waitCallFlag,omitempty"`
	PriceType               string                `json:"priceType,omitempty" bson:"priceType,omitempty"`
	CookbookID              string                `json:"cookbookId,omitempty" bson:"cookbookId,omitempty"`
	CookbookType            string                `json:"cookbookType,omitempty" bson:"cookbookType,omitempty"`
	MealNum                 string                `json:"mealNum,omitempty" bson:"mealNum,omitempty"`
	SinglePackFlag          string                `json:"singlePackFlag,omitempty" bson:"singlePackFlag,omitempty"`
	ItemNoCard              string                `json:"itemNoCard,omitempty" bson:"itemNoCard,omitempty"`
}

type CookAttachGroup struct {
	CookAttachGroupID string         `json:"cookAttachGroupId,omitempty"`
	CookAttachIdList  []CookAttachId `json:"cookAttachIdList,omitempty"`
}

type CookAttachId struct {
	CookAttachID  string `json:"cookAttachId,omitempty"`
	TemporaryFlag string `json:"temporaryFlag,omitempty"`
}

type AddSpiceAttachGroup struct {
	AddSpiceGroupID    string           `json:"addSpiceGroupId"`
	AddSpiceAttachList []AddSpiceAttach `json:"addSpiceAttachList"`
}

type AddSpiceAttach struct {
	ChildSkuID  string `json:"childSkuId"`
	ChildNumStr string `json:"childNumStr"`
}

type ComboGroup struct {
	ComboGroupID string     `json:"comboGroupId,omitempty"`
	ChildSkuList []ChildSku `json:"childSkuList,omitempty"`
}

type ChildSku struct {
	CartItemID                   string                `json:"cartItemId,omitempty"`
	ChildWeighFlag               string                `json:"childWeighFlag,omitempty"`
	ChildWeighItemPrice          string                `json:"childWeighItemPrice,omitempty"`
	RequisiteCookMethodFlag      string                `json:"requisiteCookMethodFlag,omitempty"`
	ChildSkuID                   string                `json:"childSkuId"`
	ChildNumStr                  string                `json:"childNumStr"`
	CartItemNoteList             []string              `json:"cartItemNoteList,omitempty"`
	CookAttachGroupList          []CookAttachGroup     `json:"cookAttachGroupList,omitempty"`
	AddSpiceAttachGroups         []AddSpiceAttachGroup `json:"addSpiceAttachGroups,omitempty"`
	AlternativeSkuID             string                `json:"alternativeSkuId,omitempty"`
	ChildOrderUnitNumStr         string                `json:"childOrderUnitNumStr,omitempty"`
	ChildDoubleUnitWeighDishFlag string                `json:"childDoubleUnitWeighDishFlag,omitempty"`
	RequisiteSideDishFlag        string                `json:"requisiteSideDishFlag,omitempty"`
	SinglePackFlag               string                `json:"singlePackFlag,omitempty"`
	WaitCallFlag                 string                `json:"waitCallFlag,omitempty"`
	TemporaryDishFlag            string                `json:"temporaryDishFlag,omitempty"`
	ItemName                     string                `json:"itemName,omitempty"`
	TicketExitIdList             []string              `json:"ticketExitIdList,omitempty"`
	ItemPriceStr                 string                `json:"itemPriceStr,omitempty"`
	ItemNumStr                   string                `json:"itemNumStr,omitempty"`
	WeighDishFlag                string                `json:"weighDishFlag,omitempty"`
	UnitID                       string                `json:"unitId,omitempty"`
	UnitName                     string                `json:"unitName,omitempty"`
	SecondUnitID                 string                `json:"secondUnitId,omitempty"`
	SecondUnitName               string                `json:"secondUnitName,omitempty"`
}
