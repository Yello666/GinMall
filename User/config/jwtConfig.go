package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// JWTConfig JWT配置结构体
type JWTConfig struct {
	Secret       string        `mapstructure:"secret"`
	ExpireMillis int           `mapstructure:"expire_millis"`
	Issuer       string        `mapstructure:"issuer"`
	ExpireTime   time.Duration `mapstructure:"-"` // 不反序列化，通过计算生成
}

// 计算过期时间
func (c *JWTConfig) CalculateExpireTime() {
	c.ExpireTime = time.Duration(c.ExpireMillis) * time.Millisecond
	if c.ExpireTime <= 0 {
		c.ExpireTime = 24 * time.Hour // 默认24小时
	}
}

// 获取JWT配置
func GetJWTConfig() (*JWTConfig, error) {
	if err := Load(); err != nil {
		return nil, err
	}

	var cfg JWTConfig
	if err := viper.UnmarshalKey("jwt", &cfg); err != nil {
		return nil, fmt.Errorf("解析JWT配置失败: %v", err)
	}
	cfg.CalculateExpireTime() // 计算过期时间
	return &cfg, nil
}
