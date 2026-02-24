package initialize

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"order_food/global"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var etcdCli *clientv3.Client
var cliId string
var smp sync.Map

func EtcdInit() {
	endpoints := []string{global.GVA_CONFIG.Etcd.Uri + global.GVA_CONFIG.Etcd.Port}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic("etcd init err:" + err.Error())
	}

	etcdCli = cli

}

type ServiceRegister struct {
	cli     *clientv3.Client
	leaseID clientv3.LeaseID
	key     string
	val     string
}

func PutKey(key, value string) error {

	_, err := etcdCli.Put(context.Background(), key, value)

	return err
}

func GetKey(key string) (string, error) {

	res, err := etcdCli.Get(context.Background(), key)
	value := ""
	for _, v := range res.Kvs {
		value = v.String()
	}

	return value, err

}

func (ServiceRegister) NewServiceRegister(endpoints []string, key, val string, lease int64) (*ServiceRegister, error) {

	ser := &ServiceRegister{
		cli: etcdCli,
		key: key,
		val: val,
	}

	// 创建租约
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}

	// 保持心跳
	go ser.keepAlive()

	return ser, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Grant(context.Background(), lease)
	if err != nil {
		return err
	}

	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}

	s.leaseID = resp.ID
	return nil
}

func (s *ServiceRegister) keepAlive() {
	ch, err := s.cli.KeepAlive(context.Background(), s.leaseID)
	if err != nil {
		log.Fatal(err)
	}

	for {
		<-ch
	}
}

func (s *ServiceRegister) Close() error {
	if _, err := s.cli.Revoke(context.Background(), s.leaseID); err != nil {
		return err
	}
	return s.cli.Close()
}

// func main() {

// }

type ServiceDiscovery struct {
	cli    *clientv3.Client
	prefix string
	nodes  map[string]string
}

func (ServiceDiscovery) NewServiceDiscovery(endpoints []string, prefix string) (*ServiceDiscovery, error) {

	sd := &ServiceDiscovery{
		cli:    etcdCli,
		prefix: prefix,
		nodes:  make(map[string]string),
	}

	// 初始获取服务列表
	if err := sd.getServices(); err != nil {
		return nil, err
	}

	// 监听服务变化
	if prefix != "" {
		if _, ok := smp.Load(prefix); !ok {
			go sd.watchServices()
			smp.Store(prefix, global.ServUuid)
		}
	}

	return sd, nil
}

func (s *ServiceDiscovery) getServices() error {
	resp, err := s.cli.Get(context.Background(), s.prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		s.nodes[string(kv.Key)] = string(kv.Value)
	}

	log.Printf("Current services: %v", s.nodes)
	return nil
}

func (s *ServiceDiscovery) watchServices() {
	watchChan := s.cli.Watch(context.Background(), s.prefix, clientv3.WithPrefix())

	for resp := range watchChan {
		for _, ev := range resp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				s.nodes[string(ev.Kv.Key)] = string(ev.Kv.Value)
				log.Printf("Service added: %s => %s", ev.Kv.Key, ev.Kv.Value)
			case clientv3.EventTypeDelete:
				delete(s.nodes, string(ev.Kv.Key))
				log.Printf("Service deleted: %s", ev.Kv.Key)
			}
		}
	}
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

func (s *ServiceDiscovery) GetService() (string, error) {
	if len(s.nodes) == 0 {
		return "", fmt.Errorf("no available services")
	}

	// 简单轮询负载均衡
	var keys []string
	for k := range s.nodes {
		keys = append(keys, k)
	}
	selected := keys[rand.Intn(len(keys))]
	return s.nodes[selected], nil
}
