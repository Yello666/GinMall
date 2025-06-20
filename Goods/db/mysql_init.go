package db

import (
	"Goods/model"
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

// 数据库连接参数（直接硬编码，生产环境不建议这样做）
const (
	dbUser     = "root"         // 数据库用户名
	dbPassword = "123456"       // 数据库密码
	dbHost     = "192.168.64.2" // 数据库主机
	dbPort     = 3306           // 数据库端口
	dbName     = "goods_db"     // 数据库名称
)

// 构建数据库DSN（连接字符串）
func buildDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
}
func InitMysql() error {
	dsn := buildDSN()
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}
	// 自动迁移表结构
	if err := DB.AutoMigrate(&model.GoodsInfo{}); err != nil {
		return fmt.Errorf("表结构迁移失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return sqlDB.Ping() // 测试连接
}
