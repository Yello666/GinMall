package db

import (
	"User/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	DB *gorm.DB
)

//配置加载顺序（优先级从高到低）：
//环境变量（如通过os.Setenv()或系统环境变量设置）（实际应用时使用）
//.env 文件中的配置 （开发时候可以进行简单的配置）
//config.yaml 文件中的配置(复杂配置）
//Viper 设置的默认值

func InitMysql() error {
	//cfg, err := config.GetDatabaseConfig()
	//if err != nil {
	//	fmt.Printf("加载配置失败: %v\n", err)
	var err error
	DB, err = gorm.Open(mysql.Open("root:123456@tcp(192.168.64.2:3306)/User?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s"), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败1", err)
		return fmt.Errorf("数据库连接失败1: %v", err)
	}
	//} else {
	//	// 连接数据库（使用配置中的 DSN）
	//	DB, err = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	//	if err != nil {
	//		return fmt.Errorf("数据库连接失败2: %v", err)
	//	}
	//}

	// 自动迁移表结构
	if err = DB.AutoMigrate(&model.Userinfo{}); err != nil {
		return fmt.Errorf("表结构迁移失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	//defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return sqlDB.Ping() // 测试连接
}
