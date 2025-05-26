package utils

import (
	"User"
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/big"
	"net/http"
)

// 接收请求时候的结构体与数据库中定义的结构体要求不一样
type AddUserReq struct {
	UserName string `json:"user_name" binding:"omitempty,min=1,max=50"` //omitempty：如果为空则忽略后面的限制
	Sex      string `json:"sex" binding:"omitempty,oneof=male female unknown"`
	Password string `json:"password" binding:"required,min=6,max=20" `
	Email    string `json:"email" binding:"omitempty,email"`
	Age      *int   `json:"age" binding:"omitempty"`
}

func AddUser(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "AddUser")
	//1.参数绑定与验证
	var req AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		log.Errorf("参数格式错误：%v", err)
		return
	}

	//2.密码哈希处理
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "密码加密失败",
			"error":   err.Error(),
		})
		log.Errorf("密码加密失败，error:%v", err)
		return
	}
	//3.构建用户模型
	user := User.Userinfo{
		UserName:     req.UserName,
		Sex:          req.Sex,
		Email:        req.Email,
		Age:          req.Age,
		PasswordHash: hashedPassword,
	}
	//4.自动生成用户名（如果没提供）
	if user.UserName == "" {
		user.UserName = GenerateRandomUsername()
	}
	// 5. 数据库存储
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "用户创建失败",
			"error":   err.Error(),
		})
		log.Errorf("数据库写入失败: %v", err)
		return
	}
	//6.成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":   "用户注册成功",
		"id":        user.ID,
		"user_name": user.UserName,
		"sex":       user.Sex,
		"age":       user.Age,
		"email":     user.Email,
	})
	log.Infof("新用户注册成功 ID:%d Username:%s", user.ID, user.UserName)
}

var namePrefixes = []string{
	"happy", "fast", "cool", "brave", "sunny", "clever", "lazy", "kind",
}

var nameSuffixes = []string{
	"tiger", "panda", "cat", "dog", "wolf", "fox", "lion", "bear",
}

// 生成随机的数字
func getRandomIndex(length int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(length))) //返回一个【0~length-1】范围的随机数
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

// 生成随机的名字
func GenerateRandomUsername() string {
	prefix := namePrefixes[getRandomIndex(len(namePrefixes))] //从namePrefixes切片中随机选择一个
	suffix := nameSuffixes[getRandomIndex(len(nameSuffixes))]
	number := getRandomIndex(10000)
	return fmt.Sprintf("%s_%s_%04d", prefix, suffix, number)
}

// 使用哈希加密存储密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
