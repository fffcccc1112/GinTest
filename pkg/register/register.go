package register

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"test/config"
	"test/pkg/retry"
	"time"
)

type ServiceInfo struct {
	ID       string
	Name     string
	Address  string
	Port     int
	CheckUrl string
	Tags     []string
}

func Register(c *config.Config) *api.Client {
	//初始化客户端
	config := api.DefaultConfig()
	address := fmt.Sprintf("%s:%d", c.Consul.Addr, c.Consul.Port)
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic("consul连接失败")
	}
	// 输出解密后的 token
	return client
}
func RegisterService(client *api.Client, info ServiceInfo) {
	//定义服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      info.ID,
		Name:    info.Name,
		Address: info.Address,
		Port:    info.Port,
		Tags:    info.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           info.CheckUrl,
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	//注册
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println("fail to register your service...")
		//无法注册
		//TODO：定时重试
		ReConn(client, info)

	} else {
		fmt.Println("register service success...")
	}

}
func Search(client *api.Client, serviceName string) {
	//查询服务
	services, _, err := client.Catalog().Service(serviceName, "", nil)
	if err != nil {
		log.Fatal("failed to search service")
	}
	for _, service := range services {
		fmt.Println(service)
	}
}
func deregisterService(client *api.Client, serviceID string) {
	err := client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Fatalf("Failed to deregister service: %v", err)
	}

	fmt.Println("Service deregistered successfully")
}
func ReConn(client *api.Client, info ServiceInfo) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*60*60)
	defer cancelFunc()
	conf := retry.RetryConfig{
		Interval:     time.Second * 10,
		BaseInterval: time.Second * 10,
		MaxInterval:  time.Second * 10,
		MaxRetries:   10,
		RetryMode:    "common",
	}
	retryTask := func() bool {
		// 调用你的带参数函数，传入已有变量
		return re(client, info)
	}
	retry.RetryTask(ctx, conf, retryTask, "reCONN_consul")
}
func re(client *api.Client, info ServiceInfo) bool {
	registration := &api.AgentServiceRegistration{
		ID:      info.ID,
		Name:    info.Name,
		Address: info.Address,
		Port:    info.Port,
		Tags:    info.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:                           info.CheckUrl,
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	//注册
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		return false
	}
	return true
}
