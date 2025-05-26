// 未完成
package Order

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetOrder(c *gin.Context, DB *gorm.DB) {
	//获取userID
	id, ok := c.Params.Get("userID")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't get user ID",
		})
		log.Error("获取用户id失败")
		return
	}
	//保存到order变量的UserID中
	var order Order
	idval, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "转换类型失败",
		})
		log.Error("转换类型失败")
		return
	}
	order.UserID = uint(idval)
	//查询数据库
	if err := DB.Table("orders").Where("user_id=?", order.UserID).Find(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询数据库失败",
		})
		log.Errorf("查询数据库失败:%v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"orderID": order.ID,
		"userID":  order.UserID,
		//"list":    order.List,
	})
	log.Infof("成功查看user%v的订单信息", id)
}
