package utils

import (
	"User"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type updatepswReq struct {
	OldPsw string `json:"old_psw" binding:"required,min=6,max=20"`
	NewPsw string `json:"new_psw" binding:"required,min=6,max=20"`
}

// 校验密码，第一个参数是需要校验的password，第二个参数是数据库里面的password
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
func UpdatePassword(c *gin.Context, DB *gorm.DB) {
	//1.绑定用户信息
	var req updatepswReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "密码格式错误",
			"error":   err.Error(),
		})
		return
	}
	//2.校验密码是否正确
	//获得数据库里的密码
	var user User.Userinfo
	var err error
	id, _ := c.Params.Get("id")
	user, err = GetUserByID(DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取用户信息失败",
			"error":   err.Error(),
		})
		return
	}
	//校验密码
	hash := user.PasswordHash
	if checkPasswordHash(req.OldPsw, hash) {
		log.Debug("用户验证成功，可以修改密码")
		//将新密码哈希加密
		hashedPsw, _ := hashPassword(req.NewPsw)
		//存储新密码
		updateFields := map[string]interface{}{}
		updateFields["password_hash"] = hashedPsw

		err = updateUser(user, updateFields, DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "修改出错",
				"error":   err.Error(),
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "密码错误,请重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功修改密码",
	})

}
