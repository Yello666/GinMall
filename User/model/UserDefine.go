package model

import (
	"gorm.io/gorm"
)

// 新增：参数校验，默认值设置
type Userinfo struct {
	gorm.Model
	UserName     string `json:"user_name" gorm:"not null;size:50"` // 增加长度限制
	Sex          string `json:"sex" gorm:"size:10;default:'unknown'"`
	Email        string `json:"email,omitempty" gorm:"size:100"` // 忽略空值
	Age          *int   `json:"age" `                            // 数据库级检查
	PasswordHash string `json:"-" gorm:"not null"`               // 禁止JSON序列化，存储bcrypt哈希值
}
