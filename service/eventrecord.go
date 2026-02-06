package service

import (
	"errors"
	"fmt"
	"order_food/global"
	"order_food/model/sea"
	"order_food/model/sea/request"
	"time"

	"github.com/gofrs/uuid"
)

var ERS = EventRecordService{}

type EventRecordService struct{}

func (em EventRecordService) EventRecordAdd(req request.EventRecordReq) (string, error) {
	info := sea.EventRecord{}

	reqUuid, _ := uuid.NewV4()
	reqId := reqUuid.String()

	nowTime := time.Now().Local()

	info.StoreId = req.StoreId
	info.StoreName = req.StoreName
	info.TableId = req.TableId
	info.DeskName = req.DeskName
	info.AreaId = req.AreaId
	info.AreaName = req.AreaName
	info.RequestId = reqId
	info.EventType = req.EventType
	info.RequestTime = &nowTime
	info.Date = nowTime.Format("20060102")

	err := global.GVA_DB.Create(&info).Error
	return reqId, err
}

// EventRecordUpdate 修改记录相关
func (em EventRecordService) EventRecordUpdate(req request.EventRecordReq) error {
	info := sea.EventRecord{}
	err := global.GVA_DB.Where("store_id = ? and table_id = ? and event_type = ? and event_status = 0", req.StoreId, req.TableId, req.EventType).Last(&info).Error
	if err != nil {
		return err
	}
	info.AreaId = req.AreaId
	nowTime := time.Now().Local()
	// if isCompletion && info.ResponseTime == nil {
	// 	return errors.New("流程错误响应还没完成")
	// }
	// 直接失效
	if req.IsInvalid {
		info.EventStatus = -1
		err = global.GVA_DB.Save(&info).Error
		return err
	}

	if req.Status == 21 || req.Status == 22 || req.Status == 23 || req.Status == 35 { // 响应时间
		if info.RequestTime == nil {
			global.GVA_LOG.Info(fmt.Sprintf("请求时间为空,请求编号:%v", info.RequestId))
			return errors.New("请求时间为空")
		}
		info.ResponseTime = &nowTime
		info.ResponseDuration = int(nowTime.Sub(*info.RequestTime).Seconds())
	} else if req.Status == 31 || req.Status == 32 || req.Status == 33 || req.Status == 36 { // 完成
		if info.ResponseTime == nil {
			global.GVA_LOG.Info(fmt.Sprintf("响应时间为空,请求编号:%v", info.RequestId))
			return errors.New("响应时间为空")
		}
		if info.EventStatus == 1 {
			global.GVA_LOG.Info(fmt.Sprintf("记录已完成:%v", info.RequestId))
			return errors.New("记录已完成")
		}
		info.EventStatus = 1
		info.CompletionTime = &nowTime
		info.ProcessingDuration = int(nowTime.Sub(*info.ResponseTime).Seconds())
		info.CompletionDuration = int(nowTime.Sub(*info.RequestTime).Seconds())
	} else {
		return errors.New("状态错误")
	}
	err = global.GVA_DB.Save(&info).Error

	return err
}
