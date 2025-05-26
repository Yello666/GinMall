package utils

import (
	"User"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetUser(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "GetUser")
	var userInfo User.Userinfo
	str_id, _ := c.Params.Get("id")
	userInfo, err := GetUserByID(DB, str_id)
	if err != nil {
		//找不到id
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未注册登录",
		})
		log.Errorf("id:%v未注册登录:%v", str_id, err)
		return
	} else {
		//返回id,用户名,性别,
		c.JSON(http.StatusOK, gin.H{
			"message":   "success",
			"id":        userInfo.ID,
			"user_name": userInfo.UserName,
			"sex":       userInfo.Sex,
			"age":       *userInfo.Age,
			"email":     userInfo.Email,
		})
		log.Infof("获取用户%v信息成功", str_id)
	}
}
func GetUserByID(DB *gorm.DB, str_id string) (User.Userinfo, error) {
	log.WithField("func", "GetUserByID")
	var userInfo User.Userinfo
	int64ID, _ := strconv.ParseInt(str_id, 10, 64)
	id := uint(int64ID)
	if err := DB.Where("id=?", id).Find(&userInfo).Error; err != nil {
		return userInfo, err
	} else {
		return userInfo, nil
	}
}
