package response

import "time"

type DishList struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	SrvTime string `json:"srvTime"`
	Data    struct {
		List     []Dish `json:"list"`
		Total    int    `json:"total"`
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
	} `json:"data"`
}

type Dish struct {
	ID           string  `json:"id"`
	StoreID      string  `json:"storeId"`
	DishID       string  `json:"dishId"`
	DishName     string  `json:"dishName"`
	DishType     string  `json:"dishType"`
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	UnitName     string  `json:"unitName"`
	Price        int     `json:"price"` // 单位：分
	Required     bool    `json:"required"`
	Pinyin       string  `json:"pinyin"`
	Sort         int     `json:"sort"`
	ComboID      string  `json:"comboId"`
	Combo        *Combo  `json:"combo"`
	Skus         []Sku   `json:"skus"`
	Specs        *Specs  `json:"specs"`
	Images       []Image `json:"images"`
	Tags         *Tags   `json:"tags"`
	Status       string  `json:"status"`
	Config       Config  `json:"config"`
}

type Combo struct{} // 空对象保留
type Specs struct{} // 空对象保留
type Tags struct{}  // 空对象保留

type Sku struct {
	ID        int       `json:"id"`
	DishID    string    `json:"dishId"`
	SkuID     string    `json:"skuId"`
	Name      string    `json:"name"`
	StoreID   string    `json:"storeId"`
	SpecIds   string    `json:"specIds"`
	Price     int       `json:"price"`
	Status    string    `json:"status"`
	Stock     *int      `json:"stock"` // null 时指针为 nil
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Image struct {
	ID           int    `json:"id"`
	DishID       string `json:"dishId"`
	ImageURL     string `json:"imageUrl"`
	Sort         int    `json:"sort"`
	DefaultImage int    `json:"default"` // 避开关键字
}

type Config struct {
	HideOnMiniProgram bool    `json:"hideOnMiniProgram"`
	OnSaleCron        *string `json:"onSaleCron"` // null 用指针
}
