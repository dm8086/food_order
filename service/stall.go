package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"order_food/global"
	"order_food/model/sea"
	"order_food/model/sea/request"
	"order_food/model/sea/response"
	"time"
)

var SS = StallService{}

type StallService struct {
	StoreStall
	StoreBindFood
	StallSkus
}

type StoreStall struct{}

func (*StoreStall) Add(req []request.StallAddReq) error {
	if len(req) <= 0 {
		return errors.New("请求参数为空")
	}
	list := []sea.StoreStall{}

	for _, v := range req {
		list = append(list, sea.StoreStall{
			StoreId: v.StoreId,
			Name:    v.Name,
			Sort:    v.Sort,
			IsShow:  v.IsShow,
			Status:  v.Status,
		})
	}
	if len(list) <= 0 {
		return errors.New("数据为空")
	}
	err := global.GVA_DB.Model(&sea.StoreStall{}).Save(&list).Error

	return err
}

func (*StoreStall) Edit(req []request.StallEditReq) error {
	if len(req) <= 0 {
		return errors.New("请求参数为空")
	}
	list := []sea.StoreStall{}
	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &list)

	err := global.GVA_DB.Model(&sea.StoreStall{}).Save(&list).Error
	return err
}

func (*StoreStall) List(req request.StallListReq) ([]sea.StoreStall, int64, error) {
	list := []sea.StoreStall{}
	var count int64

	stallDb := global.GVA_DB.Model(&sea.StoreStall{})
	if req.StoreId != "" {
		stallDb = stallDb.Where("store_id = ?", req.StoreId)
	}
	if req.Name != "" {
		stallDb = stallDb.Where("name like ?", req.Name+"%")
	}
	if req.Remark != "" {
		stallDb = stallDb.Where("remark like ?", req.Remark+"%")
	}
	if req.Sort != 0 {
		stallDb = stallDb.Order("id asc")
	} else {
		stallDb = stallDb.Order("id desc")
	}
	err := stallDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error

	return list, count, err
}

func (*StoreStall) ListNopage(stroreId string) ([]sea.StoreStall, error) {
	list := []sea.StoreStall{}

	err := global.GVA_DB.Model(&sea.StoreStall{}).Where("store_id = ?", stroreId).Find(&list).Error

	return list, err
}

type StoreBindFood struct{}

func (*StoreBindFood) Add(req []request.StallBingFoodAddReq) error {
	list := []sea.StallBingFood{}
	for _, v := range req {
		list = append(list, sea.StallBingFood{
			StoreId: v.StoreId,
			StallId: v.StallId,
			SkuId:   v.SkuId,
		})
	}
	if len(list) <= 0 {
		return errors.New("数据为空")
	}

	err := global.GVA_DB.Model(&sea.StallBingFood{}).Save(&list).Error

	return err
}

func (*StoreBindFood) Edit(req []request.StallBingFoodEditReq) error {
	list := []sea.StallBingFood{}

	reqByte, _ := json.Marshal(req)
	json.Unmarshal(reqByte, &list)

	err := global.GVA_DB.Model(&sea.StallBingFood{}).Save(&list).Error

	return err
}

func (*StoreBindFood) Info(id int) (sea.StallBingFood, error) {
	info := sea.StallBingFood{}

	err := global.GVA_DB.Where("id = ?", id).First(&info).Error

	return info, err
}

func (*StoreBindFood) List(req request.StallBingFoodListReq) ([]sea.StallBingFood, int64, error) {
	list := []sea.StallBingFood{}
	var count int64

	stallbDb := global.GVA_DB.Model(&sea.StallBingFood{})
	if req.StoreId != "" {
		stallbDb = stallbDb.Where("store_id = ?", req.StoreId)
	}
	if req.StallId != 0 {
		stallbDb = stallbDb.Where("stall_id = ?", req.StallId)
	}
	if req.Sort == 0 {
		stallbDb = stallbDb.Order("id desc")
	} else {
		stallbDb = stallbDb.Order("id asc")
	}
	err := stallbDb.Count(&count).Offset((req.Page - 1) * req.PageSize).Offset(req.PageSize).Find(&list).Error

	return list, count, err
}

type StallSkus struct{}

// Add 添加到商品绑定明档
func (*StallSkus) Add(orderId string) error {
	info := sea.FoodOrder{}
	_ = global.GVA_DB.Where("order_id = ?", orderId).Preload("OrderBatchs.OrderBatchDishs").First(&info).Error
	if info.ID == 0 {
		return errors.New("订单不存在,请重新输入")
	}

	stallFoods := []sea.StallBingFood{}
	_ = global.GVA_DB.Where("store_id = ?", info.StoreId).Find(&stallFoods).Error

	StallBingFoodMap := map[string][]sea.StallBingFood{}
	for _, v := range stallFoods {
		StallBingFoodMap[v.SkuId] = append(StallBingFoodMap[v.SkuId], v)
	}

	skuLists := []sea.StallSkus{}
	for _, v := range info.OrderBatchs {
		for _, v2 := range v.OrderBatchDishs {
			if len(StallBingFoodMap[v2.SkuId]) > 0 && v2.Num > 0 {
				for _, v3 := range StallBingFoodMap[v2.SkuId] {
					if v3.IsSplit == 1 {
						for i := 0; i < v2.Num; i++ {
							skuLists = append(skuLists, sea.StallSkus{
								StoreId:     v2.StoreId,
								AreaId:      v2.AreaId,
								DeskId:      v2.DeskId,
								DeskName:    info.DeskName,
								OrderId:     v2.OrderId,
								BatchId:     v2.BatchId,
								StallId:     v3.StallId,
								SkuId:       v2.SkuId,
								Num:         1,
								BatchDishId: int(v2.ID),
								Status:      v2.Status,
								DishName:    v2.DishName,
							})
						}
					} else {
						skuLists = append(skuLists, sea.StallSkus{
							StoreId:     v2.StoreId,
							AreaId:      v2.AreaId,
							DeskId:      v2.DeskId,
							DeskName:    info.DeskName,
							OrderId:     v2.OrderId,
							BatchId:     v2.BatchId,
							StallId:     v3.StallId,
							SkuId:       v2.SkuId,
							Num:         v2.Num,
							BatchDishId: int(v2.ID),
							Status:      v2.Status,
							DishName:    v2.DishName,
						})
					}
				}
			}
		}
	}

	var err error
	if len(skuLists) > 0 {
		err = global.GVA_DB.Omit("status").Save(&skuLists).Error
		if err == nil {
			// 发送  增加 明档和天上屏幕
			for _, v := range skuLists {
				SendSatll(v)
				SendToScreen(v)
			}
		}
	}
	return err
}

// SendSatll 推送明档
func SendSatll(info sea.StallSkus) {
	ES.SendSatll(info)
}

// SendToScreen 推送天上屏幕
func SendToScreen(info sea.StallSkus) {
	var stallMap = map[string]int{}
	// 商品类型  买单 服务 锅底 加工类  饮料  1 2 3 4 5
	stallMap["买单"] = 1
	stallMap["服务"] = 2
	stallMap["锅底"] = 3
	stallMap["加工类"] = 4
	stallMap["饮料"] = 5

	// shopType := 0
	storeStallList := []sea.StoreStall{}
	_ = global.GVA_DB.Where("store_id = ?", info.StoreId).Find(&storeStallList).Error
	storeStallMap := map[int]string{}
	for _, v := range storeStallList {
		storeStallMap[int(v.ID)] = v.Name
	}

	shopList := []any{}

	tmp := map[string]any{}
	tmp["name"] = info.DishName
	tmp["qty"] = info.Num
	tmp["status"] = info.Status

	shopList = append(shopList, tmp)

	topic := fmt.Sprintf("storeId/%s:areaId/%d", info.StoreId, info.AreaId)
	msg := response.TopScreenResp{}
	msg.Id = info.ID
	msg.TableId = info.DeskId
	msg.TableName = info.DeskName
	msg.Status = info.Status
	msg.TableEvent = "推送天上屏幕"
	msg.ShopList = shopList
	msg.ShopType = stallMap[storeStallMap[info.StallId]]
	msg.EventName = "topScreenStatus"
	msg.EventTime = time.Now().Local().Format("2006-01-02 15:04:05")
	msg.StoreId = info.StoreId
	msg.Num = info.Num
	ES.SendToTopScreen(topic, msg)
}

func (*StallSkus) Del(orderId string, ids []int) error {
	var err error
	list := []sea.StallSkus{}
	_ = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ?", orderId, ids).Find(&list).Error
	if len(list) > 0 {
		err = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ?", orderId, ids).Updates(map[string]any{
			"status": 2,
		}).Error

		if err == nil {
			// 发送  删除
			for _, v := range list {
				v.Status = 2
				ES.SendSatll(v)
				SendToScreen(v)
			}
		}
	}

	return err
}

func (*StallSkus) DelOrder(orderId string) error {

	orderInfo := sea.FoodOrder{}
	_ = global.GVA_DB.Where("order_id = ?", orderId).Preload("OrderBatchs.OrderBatchDishs").First(&orderInfo).Error
	ids := []int{}
	for _, v := range orderInfo.OrderBatchs {
		for _, v2 := range v.OrderBatchDishs {
			ids = append(ids, int(v2.ID))
		}
	}

	var err error
	list := []sea.StallSkus{}
	_ = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ?", orderId, ids).Find(&list).Error
	if len(list) > 0 {
		err = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ?", orderId, ids).Updates(map[string]any{
			"status": 2,
		}).Error

		if err == nil {
			// 发送  删除
			for _, v := range list {
				v.Status = 2
				SendSatll(v)
				SendToScreen(v)
			}
		}
	}

	return err
}

func (*StallSkus) Info(id int) (*sea.StallSkus, error) {
	info := sea.StallSkus{}
	err := global.GVA_DB.Where("id = ?", id).First(&info).Error

	return &info, err
}

func (*StallSkus) List(req request.StallSkuListReq) ([]sea.StallSkus, error) {
	stallSkuLists := []sea.StallSkus{}

	stallSkuDb := global.GVA_DB.Model(&sea.StallSkus{}).Where("status = 0")
	if req.StoreId != "" {
		stallSkuDb = stallSkuDb.Where("store_id = ?", req.StoreId)
	}
	if req.AreaId != 0 {
		stallSkuDb = stallSkuDb.Where("area_id = ?", req.AreaId)
	}
	if req.DeskId != 0 {
		stallSkuDb = stallSkuDb.Where("desk_id = ?", req.DeskId)
	}
	if req.OrderId != "" {
		stallSkuDb = stallSkuDb.Where("order_id = ?", req.OrderId)
	}
	if req.BatchId != "" {
		stallSkuDb = stallSkuDb.Where("batch_id = ?", req.BatchId)
	}

	err := stallSkuDb.Preload("Detail").Find(&stallSkuLists).Error

	return stallSkuLists, err
}
