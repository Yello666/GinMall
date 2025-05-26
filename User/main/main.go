/*
用户服务，负责用户的增删改查
*/
package main

import (
	"User"
	"User/Logger"
	"User/utils"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm" //gorm v1必须以空白的方式导入驱动包(前面+ _)，不然会说没有注册驱动,所以这里使用v2
)

var (
	DB *gorm.DB
)

func main() {
	//初始化数据库和日志记录器
	err := InitServer()
	if err != nil {
		fmt.Printf("FATAL!!服务器启动失败:%v", err)
		return
	}
	sqlDB, _ := DB.DB()
	defer sqlDB.Close()
	//注册服务
	//创建consul客户端
	consulClient, err := createConsulClient()
	if err != nil {
		log.Fatalf("创建consul客户端失败：%v", err)
	}
	//注册服务到consul
	serviceID := "user-service-1"
	err = registerService(consulClient, serviceID)
	if err != nil {
		log.Fatalf("注册服务失败:%v", err)
	}
	defer deregisterService(consulClient, serviceID)

	//初始化服务器和路由设置
	r := gin.Default()

	//TODO casbin用户登录验证
	//用户的增删改查
	userG := r.Group("user")
	{
		userG.GET("/get/:id", GetUser)
		userG.POST("/register", AddUser)
		userG.PUT("/update/:id", UpdateUser)
		userG.DELETE("/delete/:id", DeleteUser)
		userG.POST("/updatePsw/:id", UpdatePsw)
	}
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
	go func() {
		err := r.Run(":8081")
		if err != nil {
			log.Fatalf("服务器启动失败:%v", err)
		}
		log.Info("服务器启动成功，监听8081端口")
	}()
	waitForShutdown()
}

func GetUser(c *gin.Context) {
	log.Info("Get User")
	utils.GetUser(c, DB)
}
func AddUser(c *gin.Context) {
	log.Info("add user")
	utils.AddUser(c, DB)
}
func UpdateUser(c *gin.Context) {
	log.Info("update user")
	utils.UpdateUser(c, DB)
}
func DeleteUser(c *gin.Context) {
	log.Info("delete user")
	utils.DeleteUser(c, DB)
}
func UpdatePsw(c *gin.Context) {
	log.Info("update password")
	utils.UpdatePassword(c, DB)
}

func InitServer() error {
	if err := Logger.InitLogger(); err != nil {
		fmt.Printf("日志记录器初始化失败:%v\n", err)
		return err
	} else {
		log.Info("日志记录器初始化成功")
	}
	err := InitMysql()
	if err != nil {
		fmt.Printf("数据库连接失败:%v\n", err)
		return err
	}
	log.Info("数据库连接成功")
	return nil
}
func InitMysql() (e error) {
	//TODO viper+.env获取配置
	dsn := "root:123456@tcp(192.168.64.2:3306)/User?charset=utf8mb4&parseTime=True&loc=Local"
	//DB是全局变量，所以此处不需要冒号
	//如果创建没问题就看能不能ping通
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //上面的返回已经声明了返回变量
	if err != nil {
		//创建失败就返回e
		return err
	}
	// 自动迁移表结构
	if err := DB.AutoMigrate(&User.Userinfo{}); err != nil {
		return err
	}
	// 获取底层的sql.DB连接并测试
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func waitForShutdown() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("收到sigint信号，关闭服务器！！！！")
}
func createConsulClient() (*api.Client, error) {
	config := api.DefaultConfig()
	return api.NewClient(config)
}

func registerService(client *api.Client, serviceID string) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "user-service",
		Port:    8081,
		Address: "localhost",
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://localhost:8081/health",
			Interval:                       "30s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return client.Agent().ServiceRegister(registration)
}
func deregisterService(client *api.Client, serviceID string) {
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Fatalf("注销服务失败")
	}
	log.Info("注销了用户服务")
}
