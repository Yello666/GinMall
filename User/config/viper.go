package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
)

// 加载环境变量和配置文件
func Load() error {
	// 加载 .env 文件（开发环境）
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("加载 .env 文件失败: %v", err)
	}

	// 初始化 Viper
	viper.SetConfigName("config")   // 配置文件名（config.yaml）
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath(".")        // 项目根目录
	viper.AddConfigPath("./config") // 配置文件目录
	viper.AutomaticEnv()            // 自动读取环境变量（优先级最高）

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("未找到配置文件，使用环境变量或默认值")
		} else {
			return fmt.Errorf("解析配置文件失败: %v", err)
		}
	}
	return nil
}

// 设置默认配置
func setDefaults() {
	viper.SetDefault("db.host", "192.168.64.2")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.user", "root")
	viper.SetDefault("db.name", "User")
	viper.SetDefault("jwt.secret", "YQ FOREVER")
	viper.SetDefault("jwt.expire_millis", 3600000)
	viper.SetDefault("jwt.issuer", "gin-mall")
}
