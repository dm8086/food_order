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
	resp "order_food/model/sea/response"
	"order_food/utils"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/mgo.v2/bson"
)

var ES = EmqxService{}

type EmqxService struct{}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	go addLog(msg.Topic(), string(msg.Payload()), "defaultSub", msg.Qos())
	// 会接受emqx的信息处理 更具不同的topic做不同的处理
	go msgHandler(msg.Topic(), msg.Payload())
}

// 发布
func pub(topics map[string]byte, msg string) error {
	if len(topics) == 0 {
		return errors.New("topic不能为空")
	}

	for k, v := range topics {
		token := global.GVA_EMQX.Publish(k, v, true, msg)
		token.Wait()
		if token.Error() != nil {
			go addLog(k, msg, "Pub", v)
		}
	}

	return nil
}

// 发布
func PubToEmqx(topics map[string]byte, msg string) error {
	if len(topics) == 0 {
		return errors.New("topic不能为空")
	}

	for k, v := range topics {
		token := global.GVA_EMQX.Publish(k, v, true, msg)
		token.Wait()
		if token.Error() != nil {
			go addLog(k, msg, "Pub", v)
		}
	}

	return nil
}

func (e *EmqxService) SendToTopScreen(topic string, info response.TopScreenResp) error {

	msg, _ := json.Marshal(info)
	// 构建成和之前一样的数据就可以实现通用
	token := global.GVA_EMQX.Publish(topic, 1, true, msg)
	token.Wait()
	if token.Error() != nil {
		go addLog(topic, string(msg), "Pub", 1)
		return token.Error()
	}

	return nil
}

func (e *EmqxService) SendToTopScreenBase(tableId int, req ServiceSettle) error {

	tableInfo := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.TableId).Preload("Business").Preload("StoreArea").First(&tableInfo).Error
	if tableInfo.ID == 0 {
		return errors.New("桌台id错误")
	}

	topic := fmt.Sprintf("storeId/%s:areaId/%d", tableInfo.StoreId, tableInfo.AreaId)

	info := response.TopScreenResp{}

	nowTime := time.Now()

	shopType := req.ServType
	status := req.Status
	eventName := tableInfo.Name
	// 1服务 2买单 3锅底 4 加工类 5 饮料
	switch req.ServType {
	case 1:
		eventName = eventName + "请求服务"
		if req.Status == 2 {
			eventName = eventName + "取消服务"
		}
	case 2:
		eventName = eventName + "请求买单"
		if req.Status == 2 {
			eventName = eventName + "取消买单"
		}
	case 3:
		eventName = eventName + "锅底服务"
		// @todo  构建锅底
	case 4:
		eventName = eventName + "加工类服务"
		// @todo  构建加工类
	case 5:
		eventName = eventName + "饮料服务"
		// @todo  构建饮料
	}

	info.TableId = int(tableInfo.ID)
	info.TableName = tableInfo.Name
	info.Status = status
	info.TableEvent = eventName
	info.EventName = "topScreenStatus"
	info.EventTime = nowTime.Local().Format("2006-01-02 15:04:05")
	info.ShopList = nil
	info.ShopType = shopType
	info.StoreId = tableInfo.StoreId

	msg, _ := json.Marshal(info)
	// 构建成和之前一样的数据就可以实现通用
	token := global.GVA_EMQX.Publish(topic, 1, true, msg)
	token.Wait()
	if token.Error() != nil {
		go addLog(topic, string(msg), "Pub", 1)
	}

	return nil
}

func (e *EmqxService) OrderSendSatll(orderId string, ids []int) error {

	stallSkus := []sea.StallSkus{}
	_ = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ?", orderId, ids).Find(&stallSkus).Error
	for _, info := range stallSkus {
		e.SendSatll(info)
		SendToScreen(info)
	}

	return nil
}

// ActionEvent  针对事件的处理  下发数据
func (e *EmqxService) SendSatll(info sea.StallSkus) error {
	topic := fmt.Sprintf("storeId/%s:areaId/%d/stallId/%d", info.StoreId, info.AreaId, info.StallId)

	msg, _ := json.Marshal(info)
	// 构建成和之前一样的数据就可以实现通用
	token := global.GVA_EMQX.Publish(topic, 1, true, msg)
	token.Wait()
	if token.Error() != nil {
		go addLog(topic, string(msg), "Pub", 1)
	}

	return nil
}

func (e *EmqxService) SendDesk(tableId int, eventName string) error {
	tableInfo := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", tableId).Preload("Business").First(&tableInfo).Error

	tableBusiness := &sea.TableBusiness{}
	if tableInfo.Business != nil {
		tableBusiness = tableInfo.Business
	} else {
		// 增加没有状态的时候
	}

	nowTime := time.Now().Local().Format("2006-01-02 15:04:05")

	// storeId/ae48ee12-f901-4cbb-ba0d-4c7497a26c23/tableId/553
	topic := fmt.Sprintf("storeId/%s/tableId/%d", tableInfo.StoreId, tableId)

	var sendRes any

	switch eventName {
	case "businessStatus":

		orderId := ""
		memberId := 0
		mobileLabel := ""
		queueCode := ""
		memberUsedDy := false
		peopleCountLabel := ""
		peopleCount := 0
		if tableBusiness.OrderId != "" {
			orderId = tableBusiness.OrderId
			orderInfo := sea.Order{}
			_ = global.GVA_DB.Where("order_id = ?", orderId).First(&orderInfo).Error
			memberId = orderInfo.MemberId
			if orderInfo.Mobile != "" && len(orderInfo.Mobile) == 11 {
				mobileLabel = orderInfo.Mobile[:3] + "****" + orderInfo.Mobile[7:]
			}
			if orderInfo.MemberId != 0 {
				memberUsedDy, _ = global.GVA_REDIS.SIsMember(context.Background(), "dy:members", orderInfo.MemberId).Result()
			}
			if tableBusiness.CustomerNum > 0 {
				peopleCount = tableBusiness.CustomerNum
				peopleCountLabel = strconv.Itoa(tableBusiness.CustomerNum) + ",人"
			}
		}

		orderCode, _ := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("tablecode:id:%d", tableInfo.ID)).Result()

		res := resp.DeskStatusResp{}
		res.TableId = int(tableInfo.ID)
		res.TableName = tableInfo.Name
		res.BusinessStatus = tableInfo.Business.BusinessStatus
		res.TableEvent = tableInfo.Name + "桌台状态"
		res.EventName = "businessStatus"
		res.EventTime = nowTime
		res.IsLock = tableInfo.IsLock == 1
		res.IsOccupy = tableInfo.Business.IsOccupy == 1
		res.MemberId = memberId                       // 用户id
		res.MemberUsedDy = memberUsedDy               // 用户是否使用过抖音
		res.MobileLabel = mobileLabel                 // 电话号码
		res.OrderCode = orderCode                     // 一单一码
		res.OrderId = orderId                         // 订单号
		res.PeopleCountLabel = peopleCountLabel       // 用户人数
		res.PeopleCount = peopleCount                 // 用户人数
		res.QueueCode = queueCode                     // 排队码
		res.SettleStatus = tableBusiness.SettleStatus // 结算状态
		res.StoreId = tableInfo.StoreId               // 门店id

		sendRes = res

	case "buttonStatus":
		buttonStatus := 0
		// 服务和买单是互斥的  不会有影响
		if tableBusiness.ServiceServ > 0 {
			buttonStatus = tableBusiness.ServiceServ
		}
		if tableBusiness.SettleServ > 0 {
			buttonStatus = tableBusiness.SettleServ
		}

		res := resp.DeskButtonResp{}
		res.TableId = int(tableInfo.ID)
		res.TableName = tableInfo.Name
		res.BusinessStatus = tableBusiness.BusinessStatus
		res.ButtonStatus = buttonStatus
		res.TableEvent = tableInfo.Name + "桌台按钮事件"
		res.EventName = "buttonStatus"
		res.EventTime = nowTime
		sendRes = res
	case "soupStatus":
		soupStatus := 0

		orderId := tableBusiness.OrderId
		if orderId == "" {
			return nil
		}

		list := []sea.StallSkus{}
		_ = global.GVA_DB.Where("order_id = ? and stall_id = 1 ", orderId).Find(&list).Error
		if len(list) == 0 {
			return nil
		}

		batchDishIds := []int{}
		for _, v := range list {
			batchDishIds = append(batchDishIds, v.BatchDishId)
		}

		res := resp.DeskSoupResp{}
		res.TableId = int(tableInfo.ID)
		res.TableEvent = tableInfo.Name + "桌台锅底事件"
		res.EventName = "soupStatus"
		res.EventTime = nowTime

		batchDishs := []sea.StallSkus{}
		_ = global.GVA_DB.Where("order_id = ? and batch_dish_id IN ? and status < 2", orderId, batchDishIds).Find(&batchDishs).Error
		if len(batchDishs) == 0 {
			res.SoupList = nil
			res.SoupStatus = 2 // 全部送达
			sendRes = res
			break
		}

		type soupStr struct {
			Id     int    `json:"id"`
			Name   string `json:"name"`
			Qty    int    `json:"qyt"`
			Status int    `json:"status"`
		}

		soupStrs := []soupStr{}

		for _, v := range batchDishs {
			if v.Status == 1 {
				soupStatus = 1
			}

			soupStrs = append(soupStrs, soupStr{
				Id:     int(v.ID),
				Name:   v.DishName,
				Qty:    v.Num,
				Status: v.Status,
			})
		}
		res.SoupList = soupStrs
		res.SoupStatus = soupStatus
		sendRes = res
	}

	msg, _ := json.Marshal(sendRes)
	// 构建成和之前一样的数据就可以实现通用
	token := global.GVA_EMQX.Publish(topic, 1, true, msg)
	token.Wait()
	if token.Error() != nil {
		go addLog(topic, string(msg), "Pub", 1)
	}

	return nil
}

// ActionEvent  针对事件的处理  下发数据
func (e *EmqxService) ActionEvent(req request.TableEventReq) (any, error) {

	return nil, nil
}

// 简单的自动订阅订阅
func (e *EmqxService) AutoSub(topics map[string]byte) {

	token := global.GVA_EMQX.SubscribeMultiple(topics, messageHandler)
	token.Wait()
	if token.Error() != nil {
		label := "emqx 默认监听失败" + token.Error().Error()
		panic(label)
	}
}

// 往sea-admin 推送订单消息
func (e *EmqxService) SendOrder(orderId, storeId, orderType string, orderStatus int) error {
	topic := "storeId/order"
	data := map[string]any{}
	data["orderId"] = orderId
	data["storeId"] = storeId
	data["status"] = orderStatus
	data["type"] = "orderType"
	msg, _ := json.Marshal(data)
	global.GVA_EMQX.Publish(topic, 1, true, msg) // qos 1 保证送达
	return nil
}

func addLog(topic, msg, logType string, qos byte) {
	info := sea.EmqxLogs{}
	info.Id = bson.NewObjectId()
	info.Created = time.Now().Local().Format("2006-01-02 15:04:05")
	info.Topic = topic
	info.Qos = int(qos)
	info.LogType = logType
	info.Msg = msg
	global.GVA_MONGO.C("emqx_logs").Insert(&info)
}

func msgHandler(topic string, msg []byte) {
	switch topic {
	case sea.DEVICEHEARTBEAT: // 处理心跳
		deviceHeartBeatHandler(msg)
	// case sea.TABLEEVENTLABEL: // 桌台事件
	// 	tableEventHandler(msg)
	// case sea.STOREEVENTLABEL: // 桌台事件
	// 	storeEventHandler(msg)
	// case sea.DEVICEREFRESHLABEL: // 设备页面刷新
	// 	deviceRefreshHandler(msg)
	// case sea.TABLEORDERLABEL: // 订单事件
	// 	tableOrderHandler(msg)
	// case sea.STOREORDERLABEL: // 桌台事件
	// 	storeOrderHandler(msg)
	// case sea.TABLESCANLABEL: // 明档使用
	// 	tableScanHandler(msg)
	// case sea.STOREBROADCASTLABEL: // 页面刷新
	// 	storeBroadcastHandler(msg)
	case sea.STORESERVICESETTLELABEL: // 买单服务
		storeServiceSettleHandler(msg)
	default:
		fmt.Println("topic错误...")
	}
}

type ServiceSettle struct {
	ServType int    `json:"servType"` // 买单服务类型  1服务 2买单 3锅底  4 加工类  5 饮料
	Status   int    `json:"status"`   // 状态值     1请求服务/买单  2  取消服务/买单
	StoreId  string `json:"storeId"`  // 门店id
	TableId  int    `json:"tableId"`  // 桌台id
}

// storeServiceHandler  服务
func storeServiceSettleHandler(msg []byte) error {
	req := ServiceSettle{}
	_ = json.Unmarshal(msg, &req)

	tableInfo := sea.Table{}
	_ = global.GVA_DB.Where("id = ?", req.TableId).Preload("StoreArea").First(&tableInfo).Error
	if tableInfo.ID == 0 {
		return errors.New("获取桌台错误")
	}

	topic := fmt.Sprintf("storeId/%s:areaId/%d", tableInfo.StoreId, tableInfo.AreaId)
	nowTime := time.Now()
	shopType := req.ServType
	status := req.Status

	tableEvent := ""
	if req.ServType == 1 {
		if req.Status == 1 {
			tableEvent = tableInfo.Name + "请求服务"
		}
		if req.Status == 2 {
			tableEvent = tableInfo.Name + "取消服务"
		}
	}
	if req.ServType == 2 {
		if req.Status == 1 {
			tableEvent = tableInfo.Name + "请求买单"
		}
		if req.Status == 2 {
			tableEvent = tableInfo.Name + "取消买单"
		}
	}

	info := resp.TopScreenResp{}
	info.TableId = int(tableInfo.ID)
	info.TableName = tableInfo.Name
	info.Status = status
	info.TableEvent = tableEvent
	info.ShopList = nil
	info.ShopType = shopType
	info.EventName = "topScreenStatus"
	info.EventTime = nowTime.Local().Format("2006-01-02 15:04:05")
	info.StoreId = tableInfo.StoreId

	emqxService.SendToTopScreen(topic, info)
	return nil
}

func storeQueueHandler(msg []byte) error {
	return nil
}

func storeBroadcastHandler(msg []byte) error {
	return nil
}

func tableScanHandler(msg []byte) error {
	return nil
}

func storeOrderHandler(msg []byte) error {
	return nil
}

func tableOrderHandler(msg []byte) error {
	return nil
}

func deviceRefreshHandler(msg []byte) error {
	return nil
}

func storeEventHandler(msg []byte) error {
	return nil
}

func tableEventHandler(msg []byte) error {
	req := request.TableEventReq{}
	_ = json.Unmarshal(msg, &req)

	// 处理事件

	// 构建需要发送的topic 和数据内容
	topics, res := tableEventFormat(req)
	if len(topics) > 0 {
		return errors.New("主题不能为空")
	}
	pub(topics, res)
	return nil
}

// 构建topci 和发送内容
func tableEventFormat(req request.TableEventReq) (map[string]byte, string) {
	// 增加对事件的处理方法

	info, tids, _ := TS.DeskInfo(req.TableId)

	topics := map[string]byte{}
	for _, v := range tids {
		topics[strconv.Itoa(v)] = 0
	}
	infoStr, _ := json.Marshal(info)

	return topics, string(infoStr)
}

func deviceHeartBeatHandler(msg []byte) error {
	req := request.DevicesHeartbeatReq{}
	_ = json.Unmarshal(msg, &req)

	device := sea.StoreDevice{}
	_ = global.GVA_DB.Where("device_uuid = ?", req.DeviceUUID).First(&device).Error
	if device.ID == 0 {
		return errors.New("获取设备失败")
	}
	device.HeartbeatTime = utils.GetTime()

	if req.WorkingStatus != nil {
		device.WorkingStatus = *req.WorkingStatus
	}
	if req.OnlineStatus != nil {
		device.OnlineStatus = *req.OnlineStatus
	} else {
		device.OnlineStatus = 1
	}
	if req.Version != "" {
		device.Version = req.Version
	}
	if req.WorkLabel != "" {
		device.WorkLabel = req.WorkLabel
	}

	err := global.GVA_DB.Save(&device).Error
	if err != nil {
		return errors.New("设备更新失败" + err.Error())
	}

	return nil
}
