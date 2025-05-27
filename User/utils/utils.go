package utils

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"math/big"
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
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
