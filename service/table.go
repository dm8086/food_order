package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order_food/global"
	"order_food/model/sea"
	"order_food/model/sea/request"
	"order_food/model/sea/response"
	"order_food/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

var TS = TableService{}

type TableService struct{}

var ers = EventRecordService{}    // 事件记录服务
var orderserv = OrderService{}    // 订单服务
var emqxService = EmqxService{}   // emqx服务
var stallService = StallService{} // 明档服务

func (t *TableService) DeskList(req request.TableListReq) ([]sea.Table, error) {
	list := []sea.Table{}
	tableDb := global.GVA_DB.Model(&sea.Table{})
	if req.StoreId != "" {
		tableDb = tableDb.Where("store_id = ?", req.StoreId)
	}
	if len(req.AreaIds) > 0 {
		tableDb = tableDb.Where("area_id IN ?", req.AreaIds)
	}

	var err error

	if req.BusinessStatus != nil && len(req.BusinessStatus) > 0 {
		err = tableDb.Preload("StoreArea").Preload("Business", " business_status in ?", req.BusinessStatus).Order("status asc, sort_no asc, id desc").Find(&list).Error
	} else {
		err = tableDb.Preload("StoreArea").Order("status asc, sort_no asc, id desc").Find(&list).Error
	}

	return list, err
}

func (t *TableService) DeskOrder(tableId int) (string, int, error) {
	desk := sea.TableBusiness{}
	_ = global.GVA_DB.Where("table_id = ?", tableId).First(&desk).Error

	return desk.OrderId, desk.BusinessStatus, nil

}

func (t *TableService) ClosedDesk(req request.OpenDeskReq) error {
	// 增加开台逻辑 必点商品加入
	tableInfo := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.TableId).Preload("Business").Last(&tableInfo).Error
	if tableInfo.ID == 0 {
		return errors.New("桌台详情查询错误")
	}
	if tableInfo.Business.BusinessStatus == 0 || tableInfo.Business.BusinessStatus == 13 {
		return errors.New("桌台已经关闭")
	}

	tableInfo.IsLock = 1

	tableBusiness := tableInfo.Business

	orderId := tableBusiness.OrderId

	tableBusiness.BusinessStatus = 13
	tableBusiness.OpenTime = nil
	tableBusiness.OrderId = ""
	tableBusiness.SettleStatus = 0

	tx := global.GVA_DB.Begin()
	err := tx.Save(&tableInfo).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Save(&tableBusiness).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	nowTime := time.Now()
	err = tx.Model(&sea.Order{}).Where("order_id = ?", orderId).Updates(map[string]any{
		"open_time":     nil,
		"close_time":    nowTime,
		"status":        6,
		"settle_status": 0,
	}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// 增加更新一单一码
	updateTableLocalCode(int(tableInfo.ID))
	// 增加推送到emqx
	emqxService.SendOrder(orderId, tableInfo.StoreId, "closed", 6)
	// 增加发送到指定明档
	stallService.StallSkus.DelOrder(orderId)
	// 下发桌台状态
	emqxService.SendDesk(int(tableInfo.ID), "businessStatus")

	return nil
}

func (t *TableService) OpenDesk(req request.OpenDeskReq) (string, error) {
	// 增加开台逻辑 必点商品加入
	tableInfo := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.TableId).Preload("Business").First(&tableInfo).Error
	if tableInfo.ID == 0 {
		return "", errors.New("桌台详情查询错误")
	}
	tableBusiness := &sea.TableBusiness{}
	if tableInfo.Business == nil {
		tableBusiness.TableId = tableInfo.ID
		tableBusiness.StoreId = tableInfo.StoreId
		tableInfo.Business = tableBusiness
	} else {
		tableBusiness = tableInfo.Business
	}

	// 判断状态
	if tableBusiness.BusinessStatus != 1 && tableBusiness.BusinessStatus != 13 {
		return "", errors.New("已经开台不需要重复开台")
	}

	storeInfo := sea.Store{}
	_ = global.GVA_DB.Where("store_id = ?", tableInfo.StoreId).First(&storeInfo).Error
	if storeInfo.ID == 0 {
		return "", errors.New("门店详情错误")
	}

	// 获取必点商品
	res := utils.Get("https://dev-api.jihaihotpot.com/api/jhris/dish/v1/admin/dish/detail/list", map[string]string{"storeId": storeInfo.StoreId, "required": "1"}, nil)
	dishSkus := response.DishList{}
	json.Unmarshal([]byte(res), &dishSkus)

	if dishSkus.Code != 1 {
		return "", errors.New("获取必点商品失败")
	}
	// 构建商品列表
	skus := []response.Sku{}
	for _, v := range dishSkus.Data.List {
		if len(v.Skus) > 0 {
			for _, v2 := range v.Skus {
				skus = append(skus, v2)
			}
		}
	}

	// 订单id生成
	orderId, orderNo := getOrderId(storeInfo)

	// 订单入库
	_, _, err := orderserv.orderOpen(request.OrderAddReq{
		OrderId:     orderId,
		OrderNo:     orderNo,
		PeopleCount: req.Num,
		GoodsList:   skus,
		Src:         req.Src,
	}, &storeInfo, &tableInfo)
	if err != nil {
		return "", err
	}

	// 增加推送订单到emqx
	emqxService.SendOrder(orderId, tableInfo.StoreId, "open", 3)
	// 增加发送到指定明档
	stallService.StallSkus.Add(orderId)
	// 下发桌台状态
	emqxService.SendDesk(int(tableInfo.ID), "businessStatus")

	return orderId, nil
}
func getOrderId(storeInfo sea.Store) (string, string) {

	nowTime := time.Now().Local()
	prefix := nowTime.Format("20060102")
	noPrefix := "1125"
	suffix := ""
	if storeInfo.ID <= 9 {
		suffix = "0000" + strconv.Itoa(int(storeInfo.ID))
	} else if storeInfo.ID <= 99 {
		suffix = "000" + strconv.Itoa(int(storeInfo.ID))
	} else if storeInfo.ID <= 999 {
		suffix = "00" + strconv.Itoa(int(storeInfo.ID))
	} else if storeInfo.ID <= 9999 {
		suffix = "0" + strconv.Itoa(int(storeInfo.ID))
	} else if storeInfo.ID <= 99999 {
		suffix = strconv.Itoa(int(storeInfo.ID))
	}
	timespe := strconv.FormatInt(nowTime.UnixNano(), 10)

	// orderId = "20251208 061316000166243565250 189"
	orderId := prefix + timespe + suffix
	orderNo := noPrefix + timespe + suffix
	return orderId, orderNo
}

func updateTableLocalCode(tableId int) (string, error) {

	redisIdKey := fmt.Sprintf("tablecode:id:%d", tableId)
	res, _ := global.GVA_REDIS.Get(context.Background(), redisIdKey).Result()

	redisCodeKey := fmt.Sprintf("tablecode:code:%s", res)
	global.GVA_REDIS.Del(context.Background(), redisCodeKey)

	// 更新新的key
	uid, _ := uuid.NewV4()
	code := strings.Replace(uid.String(), "-", "", -1)
	codes := code[:10]

	redisCodeKeyTo := fmt.Sprintf("tablecode:code:%s", codes)
	global.GVA_REDIS.Set(context.Background(), redisCodeKeyTo, tableId, 0)

	global.GVA_REDIS.Set(context.Background(), redisIdKey, codes, 0)

	return codes, nil
}

func (t *TableService) DeskHandler(req request.TableEventReq) (any, error) {
	if req.EventType == 0 {
		return 0, errors.New("事件类型错误")
	}

	tx := global.GVA_DB.Begin()

	tableInfo := sea.LocalTable{}
	_ = tx.Where("desk_id = ?", req.TableId).First(&tableInfo).Error
	if tableInfo.ID == 0 {
		tx.Rollback()
		return 0, errors.New("桌台不存在")
	}
	log := sea.TableBusinessLog{}
	var err error
	eventRecordType := 0
	isUnion := false // 是否有合桌 有合桌需要处理下发

	switch req.EventType {
	case sea.EventTypeServiceRequest: // 服务
		eventRecordType, err = serviceRequestHandler(&tableInfo, &log)
	case sea.EventTypeSettleRequest: // 买单
		eventRecordType, err = settleRequestHandler(&tableInfo, &log)
	case sea.EventTypeCancelRequest: // 取消服务/买单
		eventRecordType, err = cancelRequestHandler(&tableInfo, &log)
	case sea.EventTypeSingle: // 单击
		eventRecordType, isUnion, err = singleRequestHandler(&tableInfo, &log)
	case sea.EventTypeDoubleClick: // 双击
		eventRecordType, err = doubleRequestHandler(&tableInfo, &log)

	}
	if err != nil { // 等待处理完成 继续执行提交
		return 0, errors.New("处理异常...")
	}

	// 增加记录
	if eventRecordType == 1 || eventRecordType == 2 {
		go ers.EventRecordAdd(request.EventRecordReq{ // 新增记录
			StoreId:   tableInfo.StoreId,
			StoreName: "", // 渲染的时候获取
			TableId:   tableInfo.ID,
			DeskName:  tableInfo.DeskName,
			AreaId:    uint(tableInfo.AreaId),
			AreaName:  "",              // 渲染的时候获取
			EventType: eventRecordType, // 服务
		})
	} else if eventRecordType == 1 || eventRecordType == 2 {
		tmpEventType := 0
		if eventRecordType == 21 || eventRecordType == 31 {
			tmpEventType = 1
		}
		if eventRecordType == 22 || eventRecordType == 32 {
			tmpEventType = 2
		}
		if eventRecordType == 23 || eventRecordType == 33 {
			tmpEventType = 3
		}
		if eventRecordType == 35 || eventRecordType == 36 {
			tmpEventType = 4
		}

		go ers.EventRecordUpdate(request.EventRecordReq{
			StoreId:   tableInfo.StoreId,
			TableId:   tableInfo.ID,
			EventType: tmpEventType,
			Status:    eventRecordType,
			AreaId:    uint(tableInfo.AreaId),
		})
	} else if eventRecordType == -1 {
		go ers.EventRecordUpdate(request.EventRecordReq{
			StoreId:   tableInfo.StoreId,
			TableId:   tableInfo.ID,
			IsInvalid: true, // 记录无效化
			AreaId:    uint(tableInfo.AreaId),
		})
	}

	err = tx.Save(&tableInfo).Error
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 增加联合桌台枷锁处理
	if isUnion {
		vtList1 := make([]sea.VirtualTableConnection, 0)
		_ = global.GVA_DB.Where("vt_id = ?", tableInfo.UnionData).Find(&vtList1).Error
		vtIds := make([]uint, 0)
		for _, v := range vtList1 {
			vtIds = append(vtIds, v.TableId)
		}
		if len(vtIds) > 0 {
			err = tx.Model(&sea.Table{}).Where("id IN ?", vtIds).Updates(map[string]any{
				"is_lock": tableInfo.IsLock,
			}).Error
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	_ = tx.Save(&log).Error
	tx.Commit()
	// 推送逻辑 推送到天上屏幕和远程
	return formatAndSendToQmqx(tableInfo.DeskId)
}

func formatAndSendToQmqx(deskId int) (*response.TableInfo, error) {
	info := sea.LocalTable{}

	_ = global.GVA_DB.Where("desk_id = ?", deskId).First(&info).Error
	if info.ID == 0 {
		return nil, errors.New("桌台不存在")
	}

	return nil, nil
}

func (t *TableService) sendToQmqx(res response.TableInfo, tids []int) error {
	list := []sea.LocalTable{}

	_ = global.GVA_DB.Where("desk_id IN ?", tids).Find(&list).Error

	topics := map[string]byte{}
	for _, v := range list {
		if v.IsUnion == 1 {
			uids := strings.Split(v.UnionData, ",")
			for _, v := range uids {
				vdeskIdInt, _ := strconv.Atoi(v)
				topics[fmt.Sprintf(sea.TABLEEVENT, vdeskIdInt)] = 0
			}
		} else {
			topics[fmt.Sprintf(sea.TABLEEVENT, v.DeskId)] = 0
		}
	}

	if len(topics) <= 0 {
		return errors.New("需要推送的转台id错误")
	}
	resByte, _ := json.Marshal(res)
	err := PubToEmqx(topics, string(resByte))

	return err
}

func (t *TableService) formatSendToQmqx(id uint) (*response.TableInfo, []int, error) {
	info := sea.LocalTable{}

	_ = global.GVA_DB.Where("desk_id = ?", id).First(&info).Error
	if info.ID == 0 {
		return nil, nil, errors.New("桌台不存在")
	}

	// 处理是否需要下发多个 关联
	tableIds := []int{}
	if info.IsVirtual == 1 { // 是虚拟桌台
		if info.VirtualData != "" {
			vtIds := strings.Split(strings.Trim(info.VirtualData, ","), ",")
			for _, v := range vtIds {
				vint, _ := strconv.Atoi(v)
				tableIds = append(tableIds, vint)
			}
		}
	} else {
		tableIds = append(tableIds, int(info.ID))
	}

	nowTime := time.Now().Local()

	res := response.TableInfo{}
	res.TableId = info.ID
	res.TableName = info.DeskName
	// res.TableEvent = ""
	res.BusinessStatus = info.BusinessStatus
	res.OrderId = info.OrderId
	// res.OpenTime = info.OpenTime
	res.StoreId = info.StoreId
	res.AreaName = info.AreaName
	res.LastServiceRequest = nowTime
	res.IsUnion = info.IsUnion == 1
	res.IsOccupy = info.IsOccupy == 1
	res.EventTime = nowTime
	res.TransferTableId = 0 // 转桌
	res.QueueCode = ""      // 排队码
	res.IsLock = info.IsLock == 1
	res.UnionId = info.UnionData
	res.SettleStatus = info.SettleStatus

	res.PeopleCount = info.PeopleCount
	res.PeopleCountLabel = ""
	if info.PeopleCount > 0 {
		res.PeopleCountLabel = strconv.Itoa(info.PeopleCount) + ",人"
	}

	return &res, tableIds, nil
}

func (t *TableService) Info(id int) (*sea.Table, error) {

	return nil, nil
}

func (t *TableService) DeskInfo(id int) (*response.TableInfo, []int, error) {

	info := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", id).Preload("Business").First(&info).Error
	if info.ID == 0 {
		return nil, nil, errors.New("桌台不存在")
	}

	// 处理是否需要下发多个 关联
	tableIds := []int{}
	if info.IsVirtual == 1 { // 是虚拟桌台
		if info.VirtualData != "" {
			vtIds := strings.Split(strings.Trim(info.VirtualData, ","), ",")
			for _, v := range vtIds {
				vint, _ := strconv.Atoi(v)
				tableIds = append(tableIds, vint)
			}
		}
	} else {
		tableIds = append(tableIds, int(info.ID))
	}

	tmpBusinessStatus := sea.TableBusiness{}
	if info.Business != nil {
		tmpBusinessStatus = *info.Business
	}

	nowTime := time.Now().Local()

	res := response.TableInfo{}
	res.TableId = info.ID
	res.TableName = info.Name
	// res.TableEvent = ""
	res.BusinessStatus = tmpBusinessStatus.BusinessStatus
	res.OrderId = tmpBusinessStatus.OrderId
	res.OpenTime = tmpBusinessStatus.OpenTime
	res.StoreId = info.StoreId
	res.AreaName = info.StoreArea.AreaName
	res.LastServiceRequest = nowTime
	res.IsUnion = tmpBusinessStatus.IsUnion == 1
	res.IsOccupy = tmpBusinessStatus.IsOccupy == 1
	res.EventTime = nowTime
	res.TransferTableId = 0 // 转桌
	res.QueueCode = ""      // 排队码
	res.IsLock = info.IsLock == 1
	res.UnionId = tmpBusinessStatus.UnionId
	res.SettleStatus = tmpBusinessStatus.SettleStatus

	res.PeopleCount = tmpBusinessStatus.CustomerNum
	res.PeopleCountLabel = ""
	if tmpBusinessStatus.CustomerNum > 0 {
		res.PeopleCountLabel = strconv.Itoa(tmpBusinessStatus.CustomerNum) + ",人"
	}
	// 用户  电话 @todo
	// res.MemberId = 0
	// res.MobileLabel = ""

	return &res, tableIds, nil
}

// ChangeDesk 换桌
func (t *TableService) ChangeDesk(req request.ChangeDeskReq) error {

	fromTable := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.FromTableId).Preload("Business").First(&fromTable).Error
	if fromTable.ID == 0 {
		return errors.New("原桌台不存在")
	}
	fromTableBusiness := fromTable.Business
	orderId := fromTable.Business.OrderId
	openTime := fromTable.Business.OpenTime

	orderInfo := sea.Order{}
	_ = global.GVA_DB.Where("order_id = ?", orderId).First(&orderInfo).Error
	if orderInfo.ID == 0 {
		return errors.New("订单信息不存在，不能换桌")
	}

	toTable := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.ToTableId).Preload("Business").First(&toTable).Error
	if toTable.ID == 0 {
		return errors.New("目标桌台不存在")
	}
	if toTable.Business.BusinessStatus != 0 {
		return errors.New("目标桌台非空闲状态，不能换桌")
	}
	toTableBusiness := toTable.Business

	// 开始处理事务
	tx := global.GVA_DB.Begin()
	nowTime := time.Now()

	businessStatus := 0
	if openTime.Sub(nowTime).Minutes() < -10 {
		businessStatus = 13
	}
	fromTableBusiness.OpenTime = nil
	fromTableBusiness.BusinessStatus = businessStatus
	fromTableBusiness.OrderId = ""
	fromTableBusiness.SettleStatus = 0

	toTableBusiness.OpenTime = openTime
	toTableBusiness.BusinessStatus = 10
	toTableBusiness.OrderId = orderId

	err := tx.Save(&fromTableBusiness).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&toTableBusiness).Error
	if err != nil {
		tx.Rollback()
	}

	orderInfo.FromDeskId = int(fromTable.ID)
	orderInfo.FromDeskName = fromTable.Name
	orderInfo.DeskId = int(toTable.ID)
	orderInfo.DeskName = toTable.Name
	err = tx.Save(&orderInfo).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// 增加推送到emqx 回调
	emqxService.SendOrder(orderId, orderInfo.StoreId, "chenge", orderInfo.Status)

	// 下发来源桌台状态
	emqxService.SendDesk(int(fromTable.ID), "businessStatus")
	// 下发来源桌台状态
	emqxService.SendDesk(int(toTable.ID), "businessStatus")
	return nil
}

// DeskJoint 合桌/合单
func (t *TableService) DeskJoint(req request.DeskJointReq) error {
	// 把两个桌台的订单合并到一个桌台 需要保持数据完整

	fromTable := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.FromTableId).Preload("Business").First(&fromTable).Error

	if fromTable.ID == 0 {
		return errors.New("原桌台不存在")
	}

	if fromTable.Business.OrderId == "" {
		return errors.New("原桌台无订单，不能合桌")
	}

	fromOrder := sea.Order{}
	_ = global.GVA_DB.Where("order_id = ?", fromTable.Business.OrderId).Preload("OrderBatchs.OrderBatchDishs").First(&fromOrder).Error
	if fromOrder.ID == 0 {
		return errors.New("原桌台订单不存在")
	}

	toTable := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.ToTableId).Preload("Business").First(&toTable).Error
	if toTable.ID == 0 {
		return errors.New("目标桌台不存在")
	}
	if toTable.Business.OrderId == "" {
		return errors.New("目标桌台无订单，不能合桌")
	}
	toOrder := sea.Order{}
	_ = global.GVA_DB.Where("order_id = ?", toTable.Business.OrderId).First(&toOrder).Error
	if toOrder.ID == 0 {
		return errors.New("目标桌台订单不存在")
	}
	// 开始处理事务
	tx := global.GVA_DB.Begin()

	nowTime := time.Now() // 处理桌台状态
	fromTable.Business.OpenTime = nil
	fromTable.Business.BusinessStatus = 0
	fromTable.Business.OrderId = ""
	fromTable.Business.SettleStatus = 0

	fromOrder.Status = 7 // 作废
	fromOrder.CloseTime = &nowTime

	// 增加订单批次处理
	addBatch := []sea.OrderBatch{}
	addBatchDish := []sea.OrderBatchDish{}
	for _, v := range fromOrder.OrderBatchs {
		if v.BatchType == 1 { // 普通商品批次
			if v.OrderBatchDishs != nil && len(v.OrderBatchDishs) > 0 {
				for _, v2 := range v.OrderBatchDishs {
					tmpDish := v2
					tmpDish.ID = 0
					tmpDish.OrderId = toOrder.OrderId
					addBatchDish = append(addBatchDish, tmpDish)
				}
			}
			tmpBatch := v
			tmpBatch.ID = 0
			tmpBatch.OrderId = toOrder.OrderId
			addBatch = append(addBatch, tmpBatch)
		}
	}

	toOrder.OrderAmount += fromOrder.OrderAmount

	err := tx.Save(&fromTable.Business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&fromOrder).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Save(&toTable.Business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&toOrder).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&addBatch).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&addBatchDish).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// 增加推送到emqx 回调
	// 关闭订单
	emqxService.SendOrder(fromOrder.OrderId, fromOrder.StoreId, "closed", fromOrder.Status)
	// 换订单
	emqxService.SendOrder(toOrder.OrderId, toOrder.StoreId, "chenge", toOrder.Status)

	// 下发来源桌台状态
	emqxService.SendDesk(int(fromTable.ID), "businessStatus")
	// 下发来源桌台状态
	emqxService.SendDesk(int(toTable.ID), "businessStatus")

	return nil
}
