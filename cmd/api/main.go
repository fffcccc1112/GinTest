package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"test/cmd/api/container"
	"test/config"
	"test/internal/api/handler"
	"test/internal/api/router"
	"test/internal/repository"
	"test/internal/service"
	"test/pkg/db"
	"test/pkg/logger"
	"test/pkg/register"
	"time"
)

type Hello struct {
	Name string `json:"Name"`
}

func main() {
	//db.InitGenerator()
	cfg := config.Load()
	logger.Init(cfg)
	logger.Info("加载配置成功...")
	fmt.Print(cfg.Consul)
	client := register.Register(cfg)
	info := register.ServiceInfo{
		ID:       "demo-1",
		Name:     "demo",
		Address:  "127.0.0.1",
		Port:     cfg.ServerConfig.Port,
		CheckUrl: "http://localhost:8080/healthy",
		Tags:     []string{"primary", "v1"},
	}
	go register.RegisterService(client, info)
	//初始化业务层
	//TODO:修改为依赖注入的方式
	diContainer := container.NewDIContainer()
	mysqlType := reflect.TypeOf((*db.MySqlWrapper)(nil)).Elem()

	mysqlWrapper, err := container.RegisterFactoryAndBased(diContainer, mysqlType, db.NewMysqlWrapper)
	if err != nil {
		fmt.Println(err)
	}
	//注入redis的底层template
	redisType := reflect.TypeOf((*db.RedisTemplate)(nil)).Elem()
	redisTemplate, err := container.RegisterFactoryAndBased(diContainer, redisType, db.NewRedisTemplates)
	if err != nil {
		fmt.Println("redis容器注入失败！")
	}

	redisTemplate.SetJson("hello", Hello{
		Name: "hello——fc",
	}, 0)
	json, err := redisTemplate.GetJson("hello")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", json)
	dataMap, ok := json.(map[string]interface{})
	if !ok {
		fmt.Println("数据不是map类型")
		return
	} else {
		data, _ := dataMap["Name"]
		fmt.Printf("%v", data)
	}

	userType := reflect.TypeOf((*repository.UserRepository)(nil)).Elem()
	userRepository, err2 := container.RegisterFactoryAndInstance(diContainer, userType, mysqlWrapper, repository.NewUserRepository)
	if err2 != nil {
		fmt.Println(err)
	}

	userServiceType := reflect.TypeOf((*service.UserService)(nil)).Elem()

	userServiceIn, err := container.RegisterFactoryAndInstance(diContainer, userServiceType, userRepository, service.NewUserService)
	if err != nil {
		fmt.Println(err)
	}

	userHandler := handler.NewUserHandler(userServiceIn)

	r := router.NewRouter(cfg, userHandler)
	//启动服务
	server := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.ServerConfig.Port),
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
			logger.Error("服务启动失败")
		} else {
			logger.Info("启动成功...")
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	<-sigChan // 阻塞！直到按Ctrl+C才往下走
	// 3. 按了Ctrl+C后，执行退出逻辑（终于到这了）
	diContainer.DestoryInstance(mysqlType)
	fmt.Println("优雅的退出...") // 这行100%会打印
	// 4. 退出程序
	os.Exit(0)

}
