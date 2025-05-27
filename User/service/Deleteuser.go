package service

import (
	"User/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func DeleteUserService(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "DeleteUser").Info("进入DeleteUser")
	str_id, _ := c.Params.Get("id")
	//看是否能找到原数据
	if _, err := model.GetUserByID(DB, str_id); err != nil {
		//找不到id
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户未注册登录",
		})
		log.Errorf("用户未注册登录:%v", err)
	} else {
		//执行删除
		if err := model.DeleteUserByID(DB, str_id, &model.Userinfo{}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "删除失败",
				"error":   err.Error(),
			})
			log.Errorf("删除失败:%v，用户id:%v", err, str_id)

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "注销用户成功",
			"id":      str_id,
		})
		log.Infof("删除用户%v成功", str_id)

	}
}
