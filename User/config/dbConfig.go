package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	DSN      string `mapstructure:"-"` // 不反序列化，通过计算生成
}

// BuildDSN 生成数据库连接字符串
func (c *DatabaseConfig) BuildDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		c.User, c.Password, c.Host, c.Port, c.Name,
	)
}

// 获取数据库配置
func GetDatabaseConfig() (*DatabaseConfig, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	var cfg DatabaseConfig
	if err := viper.UnmarshalKey("db", &cfg); err != nil {
		return nil, fmt.Errorf("解析数据库配置失败: %v", err)
	}
	cfg.DSN = cfg.BuildDSN() // 生成DSN
	return &cfg, nil
}
