package service

import (
	"context"
	"encoding/json"
	"errors"
	"order_food/global"
	"order_food/model/sea"
	"strconv"
	"time"
)

var oe = sea.OrderEvent{}

func serviceRequestHandler(table *sea.LocalTable, log *sea.TableBusinessLog) (int, error) {
	if table.BusinessStatus >= sea.BusinessStatusWaitService {
		return 0, nil
	}
	redisKey := "table:business:service:" + strconv.Itoa(int(table.ID))
	res, _ := global.GVA_REDIS.SetNX(context.Background(), redisKey, table.BusinessStatus, 15*time.Second).Result()

	nowTime := time.Now().Local()
	log.TableId = table.ID
	log.OrderId = table.OrderId
	log.LogType = 1
	log.EventType = sea.EventTypeServiceRequest
	log.LogTime = nowTime

	if !res {
		logRes := "间隔时间小于15秒，不操作"
		logResByte, _ := json.Marshal(logRes)
		log.LogResp = string(logResByte)
		return 0, nil
	}
	// @todo 处理
	table.BusinessStatus = sea.BusinessStatusWaitService

	logRes := "用户请求服务"
	logResByte, _ := json.Marshal(logRes)
	log.LogResp = string(logResByte)

	return 1, nil
}

func settleRequestHandler(table *sea.LocalTable, log *sea.TableBusinessLog) (int, error) {
	if table.BusinessStatus == sea.BusinessStatusWaitClean || table.BusinessStatus == sea.BusinessStatusCleaning || table.BusinessStatus == sea.BusinessStatusSettling || table.BusinessStatus == sea.BusinessStatusWaitSettle {
		return 0, nil
	}
	redisKey := "table:business:settle" + strconv.Itoa(int(table.ID))
	res, _ := global.GVA_REDIS.SetNX(context.Background(), redisKey, table.BusinessStatus, 15*time.Second).Result()

	nowTime := time.Now().Local()
	log.TableId = table.ID
	log.OrderId = table.OrderId
	log.LogType = 1
	log.EventType = sea.EventTypeSettleRequest
	log.LogTime = nowTime

	if !res {
		logRes := "间隔时间小于15秒，不操作"
		logResByte, _ := json.Marshal(logRes)
		log.LogResp = string(logResByte)
		return 0, nil
	}
	// @todo 处理
	table.BusinessStatus = sea.BusinessStatusWaitSettle

	logRes := "用户请求买单"
	logResByte, _ := json.Marshal(logRes)
	log.LogResp = string(logResByte)

	return 2, nil
}

func cancelRequestHandler(table *sea.LocalTable, log *sea.TableBusinessLog) (int, error) {
	befaultStatus := table.BusinessStatus
	if befaultStatus != sea.BusinessStatusWaitService && befaultStatus != sea.BusinessStatusWaitSettle { // 待服务
		return 0, nil
	}
	nowTime := time.Now().Local()
	log.TableId = table.ID
	log.OrderId = table.OrderId
	log.LogType = 1
	log.EventType = sea.EventTypeCancelRequest
	log.LogTime = nowTime

	logRes := "用户取消服务"
	if befaultStatus == sea.BusinessStatusWaitSettle {
		logRes = "用户取消买单"
	}

	if table.OrderId != "" {
		table.BusinessStatus = int(sea.BusinessStatusDinnering)
	} else if table.IsOccupy == 1 {
		table.BusinessStatus = int(sea.BusinessStatusOccupy)
	} else {
		table.BusinessStatus = int(sea.BusinessStatusEmpty)
	}
	logResByte, _ := json.Marshal(logRes)
	log.LogResp = string(logResByte)

	return -1, nil
}

func singleRequestHandler(table *sea.LocalTable, log *sea.TableBusinessLog) (int, bool, error) {
	businessStatus := table.BusinessStatus

	nowTime := time.Now().Local()
	log.TableId = table.ID
	log.OrderId = table.OrderId
	log.LogType = 1
	// log.EventType = sea.EventTypeCancelRequest// 事件额外处理
	log.LogTime = nowTime

	isUnion := false
	eventRecordType := 0

	switch businessStatus {
	case sea.BusinessStatusEmpty:
		if table.Src == "local" {
			return eventRecordType, isUnion, errors.New("当前桌台状态为：" + oe.GetStatus(businessStatus) + "，本地不支持该操作。")
		}
		table.BusinessStatus = sea.BusinessStatusOccupy
		table.IsLock = 1
		// 合桌是否一起解锁 需要增加逻辑
		if table.IsUnion == 1 {
			// isUnion = true
		}

	case sea.BusinessStatusOccupy:
		if table.Src == "local" {
			return eventRecordType, isUnion, errors.New("当前桌台状态为：" + oe.GetStatus(businessStatus) + "，本地不支持该操作。")
		}
		if table.IsLock != 1 {
			table.IsLock = 1
		}
		// 联合桌台处理额外返回处理
		if table.IsUnion == 1 {
			isUnion = true
		}

	case sea.BusinessStatusWaitService: // 等待服务
		table.BusinessStatus = sea.BusinessStatusServicing
		log.EventType = sea.EventTypeStartService
		eventRecordType = sea.EventTypeStartService

	case sea.BusinessStatusWaitSettle:
		table.BusinessStatus = sea.BusinessStatusSettling
		log.EventType = sea.EventTypeStartSettle
		eventRecordType = sea.EventTypeStartSettle

	case sea.BusinessStatusWaitClean:
		table.BusinessStatus = sea.BusinessStatusCleaning
		log.EventType = sea.EventTypeStartClean
		eventRecordType = sea.EventTypeStartClean

	case sea.BusinessStatusServicing:
		if table.OrderId == "" {
			table.BusinessStatus = sea.BusinessStatusDinnering
		} else if table.IsOccupy == 1 {
			table.BusinessStatus = sea.BusinessStatusOccupy
		} else {
			table.BusinessStatus = sea.BusinessStatusEmpty
		}
		log.EventType = sea.EventTypeEndService
		eventRecordType = sea.EventTypeEndService

	case sea.BusinessStatusSettling:
		if table.OrderId != "" {
			table.BusinessStatus = sea.BusinessStatusDinnering
		} else if table.IsOccupy == 1 {
			table.BusinessStatus = sea.BusinessStatusOccupy
		} else {
			table.BusinessStatus = sea.BusinessStatusEmpty
		}
		log.EventType = sea.EventTypeEndSettle
		eventRecordType = sea.EventTypeEndService

	case sea.BusinessStatusCleaning:
		table.BusinessStatus = sea.BusinessStatusEmpty
		log.EventType = sea.EventTypeEndClean
		eventRecordType = sea.EventTypeEndClean

	}

	return eventRecordType, isUnion, nil
}

// 13799989701

func doubleRequestHandler(table *sea.LocalTable, log *sea.TableBusinessLog) (int, error) {
	// todo  这边实现子订单消除
	soupLogsStatus := 0

	table.ExtendStatus = "{}"

	return soupLogsStatus, nil
}
