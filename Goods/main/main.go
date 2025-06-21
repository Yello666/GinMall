/*
用户服务，负责用户的增删改查
*/
package main

import (
	"Goods/Logger"
	"Goods/consul"
	"Goods/db"
	"Goods/router"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	//1.初始化casbin,数据库和日志记录器
	fmt.Println("hello")
	if err := Logger.InitLogger(); err != nil {
		fmt.Printf("日志记录器初始化失败:%v\n", err)
		return
	} else {
		log.Info("日志记录器初始化成功")
	}
	err := db.InitMysql()
	if err != nil {
		fmt.Printf("mysql数据库连接失败:%v\n", err)
		return
	}
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	err = db.InitRedis()
	if err != nil {
		fmt.Printf("redis数据库连接失败:%v\n", err)
		return
	}

	log.Info("数据库连接成功")

	//// 初始化Casbin
	//if err := AuthCasbin.InitCasbin(); err != nil {
	//	log.Printf("Casbin初始化失败: %v\n", err)
	//	return
	//}
	//log.Info("Casbin初始化成功")

	log.Info("服务器初始化成功！！！")

	//2.注册服务
	//创建consul客户端
	consulClient, err := consul.CreateConsulClient()
	if err != nil {
		log.Fatalf("创建consul客户端失败：%v", err)
	}
	//注册服务到consul
	serviceID := "goods-service-1"
	err = consul.RegisterService(consulClient, serviceID)
	if err != nil {
		log.Fatalf("注册服务失败:%v", err)
	}
	defer consul.DeregisterService(consulClient, serviceID)

	//3.设置路由器
	r := router.SetupRouter()
	//4.启动服务器
	go func() {
		err := r.Run(":8081")
		if err != nil {
			log.Fatalf("服务器启动失败:%v", err)
		}
		log.Info("商品服务启动成功，监听8081端口")
	}()
	waitForShutdown()
}

func waitForShutdown() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("收到sigint信号，关闭服务器！！！！")
}
