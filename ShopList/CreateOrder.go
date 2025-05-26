package ShopList

import (
	"codeup.aliyun.com/codeup/go-micro/Goods"
	"codeup.aliyun.com/codeup/go-micro/Order"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type OrderRequest struct {
	UserID      uint
	GoodsIDList []int64
}

func CreateOrder(c *gin.Context, DB *gorm.DB) {
	//接收订单商品参数
	var orderR OrderRequest
	err := c.ShouldBindJSON(&orderR)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	Orderid, err := Goods.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成订单id失败",
		})
		log.Errorf("生成订单id失败:%v", err)
		return
	}
	//将商品id改为雪花算法生成的id
	var order Order.Order
	order.ID = int64(Orderid)
	order.UserID = orderR.UserID
	idStr := fmt.Sprintf("%d", order.ID)

	for _, goodsid := range orderR.GoodsIDList {
		order.GoodsID = goodsid
		if err := DB.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "server error",
			})
			log.Errorf("数据库访问错误:%v", err)
			return
		}

	}
	//if err := DB.Create(&order).Error; err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": "server error",
	//	})
	//	log.Errorf("数据库访问错误:%v", err)
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"orderID": idStr,
	})
	log.Infof("用户%v创建了订单%v", order.UserID, order.ID)
}
