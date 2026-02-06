package initialize

import (
	"fmt"
	"order_food/global"
	"order_food/model/sea"
	"order_food/service"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type healthCheckReq struct {
	DeviceId   string `json:"deviceId"`
	DeviceType string `json:"deviceType"`
	StoreId    string `json:"storeId"`
	TableId    int    `json:"tableId"`
}

var onlineDevicesMap = map[string]string{}

// var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
// 	// 接收信息的具体业务逻辑
// 	if msg.Topic() == "healthCheck" {
// 		info := healthCheckReq{}
// 		json.Unmarshal(msg.Payload(), &info)
// 		fmt.Println(info, "----")
// 		// fmt.Printf("健康检测: Received: %s \n", msg.Payload())
// 	}

// 	if msg.Topic() == "data" {
// 		fmt.Printf("数据: Received: %s \n", msg.Payload())
// 	}
// }

var connectLost mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
	fmt.Println("链接丢失:" + err.Error())
	exitSignal <- true
}

var exitSignal = make(chan bool)

var serv1 = service.EmqxService{}

// var serv = service.ServiceGroupApp.SeaServiceGroup.EmqxService

func EmqxInit() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(global.GVA_CONFIG.Emqx.Broker)
	opts.SetClientID(global.GVA_CONFIG.Emqx.ClientID)
	opts.SetCleanSession(global.GVA_CONFIG.Emqx.CleanSession)
	opts.SetUsername(global.GVA_CONFIG.Emqx.Username)
	opts.SetPassword(global.GVA_CONFIG.Emqx.Password)
	opts.SetConnectionLostHandler(connectLost)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		msg := "emqx初始化错误:" + token.Error().Error()
		panic(msg)
	}
	global.GVA_EMQX = client
	go checkEmqx()

	go serv1.AutoSub(map[string]byte{
		sea.DEVICEHEARTBEAT:         0,
		sea.STORESERVICESETTLELABEL: 1, // 服务买单
	})

}

func checkEmqx() {
	for {
		select {
		case <-exitSignal:
			// @todo 增加推送到企业微信   服务断了要重启服务
			panic("emqx异常中断.....")
		}
	}
}

// func SubData() {

// 	topics := map[string]byte{}
// 	topics["healthCheck"] = 0
// 	topics["data"] = 2 // 1 2

// 	token := global.GVA_EMQX.SubscribeMultiple(topics, messageHandler)
// 	token.Wait()
// 	if token.Error() != nil {
// 		label := "emqx 默认监听失败" + token.Error().Error()
// 		panic(label)
// 	}
// }
