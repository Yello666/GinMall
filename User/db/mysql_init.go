package db

import (
	"User/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

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
	if err := DB.AutoMigrate(&model.Userinfo{}); err != nil {
		return err
	}
	// 获取底层的sql.DB连接并测试
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
