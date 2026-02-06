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

var BS = ButtonService{}

type ButtonService struct{}

func (*ButtonService) Click(req request.ButtonServiceReq) error {

	tableInfo := &sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.TableId).Preload("Business").First(&tableInfo).Error
	if tableInfo.ID == 0 {
		return errors.New("桌台获取错误")
	}
	tableBusiness := &sea.TableBusiness{}
	if tableInfo.Business != nil {
		tableBusiness = tableInfo.Business
	} else {
		tableBusiness.TableId = tableInfo.ID
		tableBusiness.StoreId = tableInfo.StoreId
	}

	businessLog := &sea.TableBusinessLog{}
	businessLog.TableId = tableInfo.ID
	businessLog.LogType = 1
	businessLog.EventType = req.EventType
	businessLog.LogTime = time.Now()

	var err error
	switch req.EventType {
	case 2: // 买单
		err = serviceHandler(tableInfo, tableBusiness, businessLog)
	case 3: // 服务
		err = settleHandler(tableInfo, tableBusiness, businessLog)
	case 4: // 单击
		err = singleHandler(req, tableInfo, tableBusiness, businessLog)
	case 13: // 买单服务取消
		err = cancelHandler(req, tableInfo, tableBusiness, businessLog)
	case 34: // 双击 当前是锅底送达
		err = doubleHandler(req, tableInfo, tableBusiness, businessLog)
	default:
		return errors.New("事件类型错误")
	}

	// 处理错误 直接返回
	if err != nil {
		return err
	}

	return nil

}

// serviceHandler 服务处理
func serviceHandler(table *sea.Table, business *sea.TableBusiness, log *sea.TableBusinessLog) error {
	if business.ServiceServ > 11 {
		return errors.New("服务中不可再点服务...")
	}

	logRes := getButtonLog(*table, *business, "请求服务")
	logRes["businessStatus"] = table.Business
	logRes["eventType"] = 2

	business.ServiceServ = 11
	business.SettleServ = -1

	logResByte, _ := json.Marshal(logRes)

	log.LogResp = string(logResByte)

	tx := global.GVA_DB.Begin()

	err := tx.Save(&table).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// @todo 增加推送到对应的emqx
	buttonServiceSendEmqx(*table, business.BusinessStatus, 2)

	return nil
}

// settleHandler 买单处理
func settleHandler(table *sea.Table, business *sea.TableBusiness, log *sea.TableBusinessLog) error {

	if business.ServiceServ > 12 {
		return errors.New("结算中不可再点结算...")
	}

	logRes := getButtonLog(*table, *business, "请求买单")
	logRes["businessStatus"] = table.Business
	logRes["eventType"] = 3

	business.ServiceServ = -1
	business.SettleServ = 12

	logResByte, _ := json.Marshal(logRes)

	log.LogResp = string(logResByte)

	tx := global.GVA_DB.Begin()

	err := tx.Save(&table).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// @todo 增加推送到对应的emqx
	buttonServiceSendEmqx(*table, business.BusinessStatus, 3)
	return nil
}

// cancelHandler 服务买单取消
func cancelHandler(req request.ButtonServiceReq, table *sea.Table, business *sea.TableBusiness, log *sea.TableBusinessLog) error {
	tableEvent := ""
	if business.ServiceServ == 11 {
		tableEvent = "取消服务"
	}
	if business.SettleServ == 12 {
		tableEvent = "取消买单"
	}
	logRes := getButtonLog(*table, *business, tableEvent)
	logRes["businessStatus"] = table.Business
	logRes["eventType"] = 13

	business.ServiceServ = -1
	business.SettleServ = -1

	logResByte, _ := json.Marshal(logRes)

	log.LogResp = string(logResByte)

	tx := global.GVA_DB.Begin()

	err := tx.Save(&table).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// @todo 增加推送到对应的emqx
	buttonServiceSendEmqx(*table, business.BusinessStatus, 6)
	return nil
}

// singleHandler // 单击处理
func singleHandler(req request.ButtonServiceReq, table *sea.Table, business *sea.TableBusiness, log *sea.TableBusinessLog) error {

	tableEvent := "单击"
	logRes := getButtonLog(*table, *business, tableEvent)
	logRes["eventType"] = req.EventType

	businessStatus := business.BusinessStatus
	if business.ServiceServ > 0 {
		businessStatus = business.ServiceServ
	}
	if business.SettleServ > 0 {
		businessStatus = business.SettleServ
	}

	upBool := false

	switch businessStatus {
	case 0: // 空闲
		business.BusinessStatus = 1
		upBool = true
	case 1: // 占桌
		table.IsLock = 1
	case 11: // 等服务
		business.BusinessStatus = 21
	case 12: // 等结算
		business.BusinessStatus = 22
	case 13: // 收桌
		business.BusinessStatus = 23
	case 21: // 服务中
		business.BusinessStatus = 31
	case 22: // 买单中
		business.BusinessStatus = 32
	case 23: // 收桌中
		business.BusinessStatus = 33
	}
	var err error
	tx := global.GVA_DB.Begin()
	if upBool {
		err = tx.Model(&sea.Table{}).Where("id = ?", table.ID).Updates(map[string]any{
			"is_lock": 0,
		}).Error
	} else {
		err = tx.Save(&table).Error
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&business).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Save(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	buttonServiceSendEmqx(*table, business.BusinessStatus, 6)
	return nil
}

// doubleHandler // 双击处理
func doubleHandler(req request.ButtonServiceReq, table *sea.Table, business *sea.TableBusiness, log *sea.TableBusinessLog) error {

	return nil
}

// buttonServiceSendEmqx 按钮服务丢到emqx
func buttonServiceSendEmqx(tableInfo sea.Table, businessStatus, shopType int) {

	// todo 设置topic 和内容

	topic := fmt.Sprintf("storeId/%s:areaId/%d", tableInfo.StoreId, tableInfo.AreaId)

	msg := response.TopScreenResp{}
	msg.TableId = int(tableInfo.ID)
	msg.TableName = tableInfo.Name
	msg.Status = businessStatus
	msg.TableEvent = "按钮推送天上屏幕"
	msg.ShopList = nil
	msg.ShopType = shopType
	msg.EventName = "topScreenStatus"
	msg.EventTime = time.Now().Local().Format("2006-01-02 15:04:05")
	msg.StoreId = tableInfo.StoreId

	// 设置重试次数
	retry := 3

LoogRetry:
	token := global.GVA_EMQX.Publish(topic, 0, true, msg)
	token.Wait()
	if token.Error() != nil {
		if retry > 0 {
			goto LoogRetry
		} else {
			// @todo 增加日志 推送错误
		}
		retry--
	}
}

func getButtonLog(table sea.Table, business sea.TableBusiness, tableEvent string) map[string]any {
	logs := map[string]any{}

	logs["storeId"] = table.StoreId
	logs["tableId"] = table.ID
	logs["tableName"] = table.Name
	logs["tableEvent"] = tableEvent
	logs["tableBeforeStatus"] = business.BusinessStatus
	logs["orderId"] = business.OrderId
	logs["openTime"] = business.OpenTime
	logs["isUnion"] = business.IsUnion == 1
	logs["eventTime"] = time.Now()
	// logs["eventType"] = "button-event"

	return logs

}
