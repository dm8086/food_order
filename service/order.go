package service

import (
	"encoding/json"
	"errors"
	"order_food/global"
	"order_food/model/sea"
	"order_food/model/sea/request"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

var OrderS = OrderService{}

type OrderService struct{}

// OrderAdd  订单增加商品
func (o *OrderService) OrderSub(req request.OrderSubReq) (*sea.FoodOrder, *sea.FoodOrderBatch, error) {
	info := sea.FoodOrder{}

	nowTime := time.Now()
	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).Preload("OrderBatchs.OrderBatchDishs").First(&info).Error
	if info.ID == 0 {
		return nil, nil, errors.New("订单不存在")
	}
	if info.Status != 3 {
		return nil, nil, errors.New("订单已关闭")
	}

	// 获取订单里所有的商品
	orderBatchDishMap := map[int]bool{}
	for _, v := range info.OrderBatchs {
		for _, v1 := range v.OrderBatchDishs {
			orderBatchDishMap[int(v1.ID)] = true
		}
	}

	subList := []sea.FoodOrderBatchDish{}
	subIdMaps := map[int]bool{}
	subids := []int{}
	subNumMap := map[int]int{}
	subNum := 0
	subAmount := 0

	for k, v := range req.GoodsList {
		kInt, _ := strconv.Atoi(k)
		if !orderBatchDishMap[kInt] {
			return nil, nil, errors.New("商品不属于该订单")
		}
		subIdMaps[kInt] = true
		subids = append(subids, kInt)
		subNumMap[kInt] = v
	}

	_ = global.GVA_DB.Where("order_id = ? and id IN ?", req.OrderId, subids).Find(&subList).Error

	orderBatch := sea.FoodOrderBatch{}
	batchUuid, _ := uuid.NewV4()
	bathId := batchUuid.String() //
	// 退单的批次
	orderBatch.AddTime = &nowTime
	orderBatch.OrderId = info.OrderId
	orderBatch.BatchId = bathId
	orderBatch.StoreId = info.StoreId
	orderBatch.DeskId = info.DeskId
	orderBatch.AreaId = info.AreaId
	orderBatch.DeskId = info.DeskId

	// 需要退单的批次
	subOrderBatch := []string{}

	type NumAmount struct {
		Num    int `json:"num"`
		Amount int `json:"amount"`
	}
	batchChangeMap := map[string]NumAmount{}

	// orderBatch dish
	subAddList := []sea.FoodOrderBatchDish{}
	for k1, v1 := range subList {
		subOrderBatch = append(subOrderBatch, v1.BatchId)
		tmpNum := subNumMap[int(v1.ID)]
		tmpAmount := v1.Price * tmpNum
		subAddList = append(subAddList, sea.FoodOrderBatchDish{
			AddTime:   &nowTime,
			OrderId:   v1.OrderId,
			BatchId:   bathId,
			StoreId:   v1.StoreId,
			DeskId:    v1.DeskId,
			AreaId:    v1.AreaId,
			DishId:    v1.DishId,
			DishName:  v1.DishName,
			SkuId:     v1.SkuId,
			Price:     v1.Price,
			Num:       tmpNum,
			Amount:    tmpAmount,
			Detail:    v1.Detail,
			BatchType: 2,
		})
		subNum += tmpNum
		subAmount += tmpAmount

		tmpbatchChangeMap := batchChangeMap[v1.BatchId]
		tmpbatchChangeMap.Num += tmpNum
		tmpbatchChangeMap.Amount += tmpAmount
		batchChangeMap[v1.BatchId] = tmpbatchChangeMap

		v1.Num = v1.Num - tmpNum
		v1.Amount = v1.Amount - tmpAmount
		subList[k1] = v1
	}

	// orderBatch
	orderBatch.Num = subNum
	orderBatch.BatchType = 2 // 退单
	orderBatch.Amount = subAmount

	// 修改订单金额
	info.OrderAmount = info.OrderAmount - subAmount

	// 需要修改的批次
	subOriginBatch := []sea.FoodOrderBatch{}
	_ = global.GVA_DB.Where("order_id = ? and batch_id IN ?", req.OrderId, subOrderBatch).Find(&subOriginBatch).Error

	for k2, v2 := range subOriginBatch {
		if v2.Amount < 0 || v2.Num < 0 {
			return nil, nil, errors.New("当前批次不可退")
		}
		v2.Num = v2.Num - batchChangeMap[v2.BatchId].Num
		v2.Amount = v2.Amount - batchChangeMap[v2.BatchId].Amount
		subOriginBatch[k2] = v2
	}

	tx := global.GVA_DB.Begin()

	err := tx.Save(&info).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Save(&orderBatch).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if len(subAddList) > 0 {
		err = tx.Save(&subAddList).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if len(subOriginBatch) > 0 {
		err = tx.Save(&subOriginBatch).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if len(subList) > 0 {
		err = tx.Save(&subList).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if len(subids) > 0 {
		stallSkuList := []sea.StallSkus{}
		_ = global.GVA_DB.Where("batch_dish_id IN ?", subids).Find(&stallSkuList).Error
		for _, v := range stallSkuList {
			if subNumMap[v.BatchDishId] > 0 {
				status := v.Status
				num := 0
				if v.Num >= subNumMap[v.BatchDishId] {
					num = v.Num - subNumMap[v.BatchDishId]

				} else {
					subNumMap[v.BatchDishId] = subNumMap[v.BatchDishId] - v.Num
					status = 3
					num = 0
				}
				if num == 0 {
					status = 3
				}
				_ = global.GVA_DB.Model(&sea.StallSkus{}).Where("id = ?", v.ID).Updates(map[string]any{
					"status": status,
					"num":    num,
				}).Error
			}
		}
	}

	tx.Commit()

	// 商品减少下发回调
	emqxService.SendOrder(info.OrderId, info.StoreId, "subSku", info.Status)
	//  增加发送到指定明档  消除商品
	emqxService.OrderSendSatll(info.OrderId, subids)
	//  开台下发锅底状态 有减少锅底就推送减少锅底
	emqxService.SendDesk(info.DeskId, "soupStatus")
	return nil, nil, nil
}

// OrderAdd  订单增加商品
func (o *OrderService) OrderAdd(req request.OrderAddReq) (*sea.FoodOrder, *sea.FoodOrderBatch, error) {

	info := sea.FoodOrder{}
	orderBatch := sea.FoodOrderBatch{}
	bathDIshList := []sea.FoodOrderBatchDish{}

	nowTime := time.Now()
	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&info).Error
	if info.ID == 0 {
		return nil, nil, errors.New("订单不存在")
	}
	if info.Status != 3 {
		return nil, nil, errors.New("订单已经关闭")
	}

	orderAmount := 0
	num := 0
	if len(req.Dishes) > 0 {
		batchUuid, _ := uuid.NewV4()
		bathId := batchUuid.String() // 批次id

		orderBatch.AddTime = &nowTime
		orderBatch.OrderId = info.OrderId
		orderBatch.BatchId = bathId
		orderBatch.StoreId = info.StoreId
		orderBatch.DeskId = info.DeskId
		orderBatch.AreaId = info.AreaId
		orderBatch.DeskId = info.DeskId
		orderBatch.Num = req.PeopleCount
		orderBatch.BatchType = 1

		for _, v := range req.Dishes {
			detail := ""
			if len(v.ComboSubDishes) > 0 {
				detailByte, _ := json.Marshal(v.ComboSubDishes)
				detail = string(detailByte)
			}
			tmp2 := sea.FoodOrderBatchDish{}
			tmp2.AddTime = &nowTime
			tmp2.OrderId = info.OrderId
			tmp2.BatchId = bathId
			tmp2.StoreId = info.StoreId
			tmp2.DeskId = info.DeskId
			tmp2.AreaId = info.AreaId
			tmp2.DeskId = info.DeskId
			tmp2.DishId = v.DishId
			tmp2.DishName = v.DishName
			tmp2.SkuId = v.SkuId
			tmp2.Price = v.Price
			tmp2.Amount = v.Price * v.Quantity
			tmp2.Num = v.Quantity
			tmp2.Detail = detail
			tmp2.BatchType = 1
			bathDIshList = append(bathDIshList, tmp2)
			num += v.Quantity
			orderAmount += tmp2.Amount
		}
		orderBatch.Amount = orderAmount
		orderBatch.Num = num
		info.OrderAmount += orderAmount
	}

	tx := global.GVA_DB.Begin()
	err := tx.Save(&info).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Save(&orderBatch).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if len(bathDIshList) > 0 {
		err = tx.Save(&bathDIshList).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}
	tx.Commit()
	// 增加推送到emqx
	emqxService.SendOrder(info.OrderId, info.StoreId, "addSku", info.Status)
	// 增加发送到指定明档
	stallService.StallSkus.Add(info.OrderId)
	// 开台下发锅底状态
	emqxService.SendDesk(info.DeskId, "soupStatus")

	return &info, &orderBatch, nil
}

func (o *OrderService) orderOpen(req request.OrderAddReq, storeInfo *sea.Store, tableInfo *sea.Table) (*sea.FoodOrder, *sea.FoodOrderBatch, error) {

	info := sea.FoodOrder{}
	orderBatch := sea.FoodOrderBatch{}
	bathDIshList := []sea.FoodOrderBatchDish{}

	nowTime := time.Now()
	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&info).Error
	if info.ID == 0 {
		info.AreaId = int(tableInfo.AreaId)
		info.OpenTime = &nowTime
		info.OrderId = req.OrderId
		info.OrderNo = req.OrderNo
		info.DeskId = int(tableInfo.ID)
		info.DeskName = tableInfo.Name
		info.StoreId = req.StoreId
		info.StoreName = storeInfo.StoreName
		info.PeopleCount = req.PeopleCount
		info.Remark = req.Remark
		info.Status = 3
		info.StoreId = storeInfo.StoreId
		info.PeopleCount = req.PeopleCount
		info.Src = req.Src
	}

	orderAmount := 0
	nums := 0
	if len(req.GoodsList) > 0 {
		batchUuid, _ := uuid.NewV4()
		bathId := batchUuid.String() // 批次id

		orderBatch.AddTime = &nowTime
		orderBatch.OrderId = info.OrderId
		orderBatch.BatchId = bathId
		orderBatch.StoreId = storeInfo.StoreId
		orderBatch.DeskId = info.DeskId
		orderBatch.AreaId = int(tableInfo.AreaId)
		orderBatch.DeskId = int(tableInfo.ID)
		orderBatch.BatchType = 1

		for _, v := range req.GoodsList {
			tmp2 := sea.FoodOrderBatchDish{}
			tmp2.AddTime = &nowTime
			tmp2.OrderId = info.OrderId
			tmp2.BatchId = bathId
			tmp2.StoreId = storeInfo.StoreId
			tmp2.DeskId = info.DeskId
			tmp2.AreaId = int(tableInfo.AreaId)
			tmp2.DeskId = int(tableInfo.ID)
			tmp2.DishId = v.DishID
			tmp2.DishName = "酱料"
			tmp2.SkuId = v.SkuID
			tmp2.Price = v.Price
			tmp2.Amount = v.Price * req.PeopleCount
			tmp2.Num = req.PeopleCount
			tmp2.BatchType = 1
			bathDIshList = append(bathDIshList, tmp2)
			orderAmount += v.Price * req.PeopleCount
			nums += req.PeopleCount
		}
		orderBatch.Amount = orderAmount
		orderBatch.Num = nums
		info.OrderAmount += orderAmount
	}

	tx := global.GVA_DB.Begin()
	err := tx.Save(&info).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Save(&orderBatch).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if len(bathDIshList) > 0 {
		err = tx.Save(&bathDIshList).Error
		if err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}
	tableInfo.IsLock = 0
	err = tx.Model(&tableInfo).Where("id = ?", tableInfo.ID).Updates(map[string]any{
		"is_lock": 0,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tableInfo.Business.BusinessStatus = 10
	tableInfo.Business.CustomerNum = req.PeopleCount
	tableInfo.Business.OpenTime = &nowTime
	tableInfo.Business.OrderId = req.OrderId

	err = tx.Save(&tableInfo.Business).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tx.Commit()

	return &info, &orderBatch, nil
}

// orderId, batchId string, stallId, orderBatchDishId int
func (o *OrderService) SkuRemove(req request.OrderSkuRemoveReq) error {
	orderInfo := sea.FoodOrder{}
	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&orderInfo).Error
	if orderInfo.ID == 0 {
		return errors.New("订单查询失败")
	}
	orderBatchDish := sea.FoodOrderBatchDish{}
	_ = global.GVA_DB.Where("order_id = ? and batch_id = ? and id = ?", req.OrderId, req.BatchId, req, orderBatchDish).First(&orderBatchDish).Error
	if orderBatchDish.ID == 0 {
		return errors.New("订单查询失败")
	}
	orderBatchDish.Status = 1 // 完成
	err := global.GVA_DB.Save(&orderBatchDish).Error

	stallSkus := sea.StallSkus{}
	_ = global.GVA_DB.Where("order_id = ? and batch_id = ? and stall_id = ? and batch_dish_id = ?", req.OrderId, req.BatchId, req.StallId, req.OrderBatchDishId).First(&stallSkus).Error

	stallSkus.Status = 1
	err = global.GVA_DB.Save(&stallSkus).Error
	// @todo   增加 下发
	return err
}

func (o *OrderService) SoupRemove(req request.SoupRemoveReq) error {
	orderInfo := sea.FoodOrder{}
	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&orderInfo).Error
	if orderInfo.ID == 0 {
		return errors.New("订单获取失败")
	}
	if orderInfo.Status != 3 {
		return errors.New("桌台已关闭")
	}

	status := *req.Status
	status = status - 1
	if status < 0 {
		status = 0
	}

	stallSkus := []sea.StallSkus{}
	_ = global.GVA_DB.Where("order_id = ? and stall_id = 1 and status = ?", orderInfo.OrderId, status).Find(&stallSkus).Error

	ids := []int{}
	for _, v := range stallSkus {
		ids = append(ids, v.BatchDishId)
	}

	if len(ids) > 0 {

		err := global.GVA_DB.Model(&sea.StallSkus{}).Where("order_id = ? and  batch_dish_id IN ?", req.OrderId, ids).Updates(map[string]any{
			"status": &req.Status,
		}).Error
		if err != nil {
			return err
		}
		// 开台下发锅底状态 更新条数大于0才下发
		emqxService.SendDesk(orderInfo.DeskId, "soupStatus")
		// 增加明档下发数据
		emqxService.OrderSendSatll(req.OrderId, ids)

	}

	return nil
}

func (o *OrderService) List(req request.OrderListReq) ([]sea.FoodOrder, int64, error) {
	list := []sea.FoodOrder{}

	orderDb := global.GVA_DB.Model(&sea.FoodOrder{})
	if req.StoreId != "" {
		orderDb = orderDb.Where("order_id = ?", req.StoreId)
	}
	if req.AreaId != 0 {
		orderDb = orderDb.Where("area_id = ?", req.AreaId)
	}
	if req.Status != nil {
		orderDb = orderDb.Where("status = ?", req.Status)
	}

	var count int64
	orderBy := "id  desc "
	if req.Sort > 0 {
		orderBy = "id asc "
	}
	err := orderDb.Count(&count).Order(orderBy).Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&list).Error

	return list, count, err
}

func (o *OrderService) Info(orderId string) (*sea.FoodOrder, error) {

	info := sea.FoodOrder{}
	_ = global.GVA_DB.Where("order_id = ?", orderId).Preload("OrderBatchs.OrderBatchDishs").First(&info).Error
	if info.ID == 0 {
		return nil, errors.New("订单不存在,请重新输入")
	}

	return &info, nil
}

func (o *OrderService) Update(req request.OrderUpdateReq) error {
	info := sea.FoodOrder{}

	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&info).Error
	if info.ID == 0 {
		return errors.New("订单不存在,请重新输入")
	}
	if req.OrderAmount != nil {
		info.OrderAmount = *req.OrderAmount
	}
	if req.Status != nil {
		info.Status = *req.Status
	}
	if req.SettleStatus != nil {
		info.SettleStatus = *req.SettleStatus
	}
	if req.OrderReceivedAmount != nil {
		info.OrderReceivedAmount = *req.OrderReceivedAmount
	}
	if req.PromoAmount != nil {
		info.PromoAmount = *req.PromoAmount
	}
	err := global.GVA_DB.Where("order_id = ?", req.OrderId).Save(&info).Error
	return err
}

// DishSkuUpdate  菜品券扣减
func (o *OrderService) DishSkuUpdate(req request.CombWriteoffReq) error {
	info := sea.FoodOrder{}

	_ = global.GVA_DB.Where("order_id = ?", req.OrderId).First(&info).Error
	if info.ID == 0 {
		return errors.New("订单不存在,请重新输入")
	}

	// @todo

	bathDishs := []sea.FoodOrderBatchDish{}
	_ = global.GVA_DB.Where("order_id = ? and sku_id = ? and num > 0", req.OrderId, "").Find(&bathDishs).Error

	err := global.GVA_DB.Where("order_id = ?", req.OrderId).Save(&info).Error
	return err
}

// 增加调用写日志记录
func addOrderLog(orderId, detail string, logEvent int) error {
	info := sea.OrderLog{}
	info.OrderId = orderId
	info.OrderEvent = logEvent
	info.Detail = detail

	return global.GVA_DB.Save(&info).Error
}
