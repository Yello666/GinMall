package utils

import (
	"User"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type updateUserReq struct {
	UserName *string `json:"user_name" binding:"omitempty,min=1,max=50"` // 增加长度限制
	Sex      *string `json:"sex" binding:"omitempty,oneof=male female unknown"`
	Email    *string `json:"email" binding:"omitempty,email"`       // 忽略空值
	Age      *int    `json:"age" binding:"omitempty,gte=0,lte=150"` // 数据库级检查
}

//使用指针类型可以区分传入的是空字符串还是没传入值，是0还是没传入值

func UpdateUser(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "UpdateUser").Info("进入UpdateUser")
	var req updateUserReq
	var err error
	//获取更新的参数
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数格式错误",
			"error":   err.Error(),
		})
		return
	}
	//1.查找原来的用户
	str_id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取id失败"})
		log.Error("获取id失败")
		return
	}
	origin := User.Userinfo{}
	origin, err = GetUserByID(DB, str_id)
	if err != nil {
		log.Errorf("找不到原来的user信息:%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "账号未注册",
			"error":   err.Error(),
		})
		return
	}
	//2.更新非空字段
	updateFields := map[string]interface{}{}
	if req.UserName != nil {
		updateFields["user_name"] = *req.UserName
	}
	if req.Sex != nil {
		updateFields["sex"] = *req.Sex
	}
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.Age != nil {
		updateFields["age"] = *req.Age
	}
	err = updateUser(origin, updateFields, DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新出错",
			"error":   err.Error(),
		})
		return
	}
	//更新成功
	updatedUser, err := GetUserByID(DB, str_id)

	c.JSON(http.StatusOK, gin.H{
		"message":   "更新用户信息成功",
		"user_name": updatedUser.UserName,
		"age":       updatedUser.Age,
		"sex":       updatedUser.Sex,
		"email":     updatedUser.Email,
		"id":        updatedUser.ID,
	})

	log.Infof("用户%v更新数据成功", str_id)

}
func updateUser(origin User.Userinfo, updateFields map[string]interface{}, DB *gorm.DB) error {
	log.WithField("func", "updateUser").Info("进入updateUser")

	//只进行修改操作
	if err := DB.Model(&origin).Updates(updateFields).Error; err != nil {
		log.Errorf("更新出错:%v", err)
		return err
	}
	return nil

}
