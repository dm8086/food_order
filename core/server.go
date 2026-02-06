package core

import (
	"fmt"
	"order_food/initialize"
	"time"

	"order_food/global"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		initialize.Redis()
	}

	if global.GVA_CONFIG.System.UseMongo {
		initialize.Mongo()
	}

	if global.GVA_CONFIG.ApiLog.Enabled {
		initialize.ApiLog()
	}

	if global.GVA_CONFIG.System.UseMiniProgram {
		initialize.MiniProgram()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	// 判断是本地的时候不添加到etcd
	if global.GVA_CONFIG.Etcd.Env != "local" {
		uid, _ := uuid.NewV4()
		// 生成实例id
		caId := uid.String()
		endpoints := []string{global.GVA_CONFIG.Etcd.Uri + global.GVA_CONFIG.Etcd.Port}
		serviceKey := "/jhBase/orderServ/" + caId

		serviceVal := global.GVA_CONFIG.Etcd.ServUri + address

		etcdServ := initialize.ServiceRegister{}
		_, err := etcdServ.NewServiceRegister(endpoints, serviceKey, serviceVal, global.GVA_CONFIG.Etcd.Lease)
		if err != nil {
			fmt.Println("微服务注册失败:", err.Error())
			return
		}
		// 注册成功把当前实例id存到global中
		global.ServUuid = caId
	}

	fmt.Printf(`
	欢迎使用 寄海管理系统
	默认自动化文档地址:http://127.0.0.1%s/api/swagger/index.html`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
