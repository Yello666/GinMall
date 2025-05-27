/*
用户服务，负责用户的增删改查
*/
package main

import (
	"User/AuthCasbin"
	"User/Logger"
	"User/consul"
	"User/db"
	"User/router"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// TODO JWT权限认证
// TODO casbin用户登录验证
func main() {
	//1.初始化casbin,数据库和日志记录器
	if err := Logger.InitLogger(); err != nil {
		fmt.Printf("日志记录器初始化失败:%v\n", err)
		return
	} else {
		log.Info("日志记录器初始化成功")
	}
	err := db.InitMysql()
	if err != nil {
		fmt.Printf("数据库连接失败:%v\n", err)
		return
	}
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()
	log.Info("数据库连接成功")

	// 初始化Casbin
	if err := AuthCasbin.InitCasbin(); err != nil {
		log.Printf("Casbin初始化失败: %v\n", err)
		return
	}
	log.Info("Casbin初始化成功")

	log.Info("服务器初始化成功！！！")

	//2.注册服务
	//创建consul客户端
	consulClient, err := consul.CreateConsulClient()
	if err != nil {
		log.Fatalf("创建consul客户端失败：%v", err)
	}
	//注册服务到consul
	serviceID := "user-service-1"
	err = consul.RegisterService(consulClient, serviceID)
	if err != nil {
		log.Fatalf("注册服务失败:%v", err)
	}
	defer consul.DeregisterService(consulClient, serviceID)

	//初始化服务器和路由设置
	//r := gin.Default()
	//userG := r.Group("user")
	//{
	//	userG.GET("/get/:id", handler.GetUser)
	//	userG.POST("/register", handler.AddUser)
	//	userG.PUT("/update/:id", handler.UpdateUser)
	//	userG.DELETE("/delete/:id", handler.DeleteUser)
	//	userG.POST("/updatePsw/:id", handler.UpdatePsw)
	//}
	//r.GET("/health", func(c *gin.Context) {
	//	c.Status(200)
	//})
	//3.设置路由器
	r := router.SetupRouter()
	//4.启动服务器
	go func() {
		err := r.Run(":8081")
		if err != nil {
			log.Fatalf("服务器启动失败:%v", err)
		}
		log.Info("服务器启动成功，监听8081端口")
	}()
	waitForShutdown()
}

//func GetUser(c *gin.Context) {
//	log.Info("Get User")
//	service.GetUser(c, DB)
//}
//func AddUser(c *gin.Context) {
//	log.Info("add user")
//	service.AddUser(c, DB)
//}
//func UpdateUser(c *gin.Context) {
//	log.Info("update user")
//	service.UpdateUser(c, DB)
//}
//func DeleteUser(c *gin.Context) {
//	log.Info("delete user")
//	service.DeleteUser(c, DB)
//}
//func UpdatePsw(c *gin.Context) {
//	log.Info("update password")
//	service.UpdatePassword(c, DB)
//}

//	func InitServer() error {
//		if err := Logger.InitLogger(); err != nil {
//			fmt.Printf("日志记录器初始化失败:%v\n", err)
//			return err
//		} else {
//			log.Info("日志记录器初始化成功")
//		}
//		err := InitMysql()
//		if err != nil {
//			fmt.Printf("数据库连接失败:%v\n", err)
//			return err
//		}
//		log.Info("数据库连接成功")
//		return nil
//	}

func waitForShutdown() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("收到sigint信号，关闭服务器！！！！")
}

//func createConsulClient() (*api.Client, error) {
//	config := api.DefaultConfig()
//	return api.NewClient(config)
//}
//
//func registerService(client *api.Client, serviceID string) error {
//	registration := &api.AgentServiceRegistration{
//		ID:      serviceID,
//		Name:    "user-service",
//		Port:    8081,
//		Address: "localhost",
//		Check: &api.AgentServiceCheck{
//			HTTP:                           "http://localhost:8081/health",
//			Interval:                       "30s",
//			Timeout:                        "5s",
//			DeregisterCriticalServiceAfter: "1m",
//		},
//	}
//	return client.Agent().ServiceRegister(registration)
//}
//func deregisterService(client *api.Client, serviceID string) {
//	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
//		log.Fatalf("注销服务失败")
//	}
//	log.Info("注销了用户服务")
//}
