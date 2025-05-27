package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"strconv"
)

// 生成随机的数字
func GetRandomIndex(length int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(length))) //返回一个【0~length-1】范围的随机数
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

// 哈希加密
// 使用哈希加密存储密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 校验密码，第一个参数是需要校验的password，第二个参数是数据库里面的password
func CheckPasswordHash(passwordToCheck, hashedPsw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPsw), []byte(passwordToCheck)) == nil
}

// 获取user的ID（将interface转换为string）
func GetUserID(c *gin.Context) (string, bool) {
	if userIDRaw, exists := c.Get("user_id"); exists {
		userID, ok := userIDRaw.(string)
		return userID, ok
	}
	return "", false
}

// GetRole （将interface转换为string）
func GetRole(c *gin.Context) (string, bool) {
	if roleRaw, exists := c.Get("role"); exists {
		role, ok := roleRaw.(string)
		return role, ok
	}
	return "", false
}

// GetPaginationParams 解析分页参数（page从1开始）
func GetPaginationParams(c *gin.Context) (int, int, error) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, 0, fmt.Errorf("page参数必须为正整数")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		return 0, 0, fmt.Errorf("pageSize必须为1-100的整数")
	}

	return page, pageSize, nil
}
