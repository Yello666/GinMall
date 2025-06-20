package service

import (
	"Goods/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetOneGoods(c *gin.Context, DB *gorm.DB) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取商品id失败",
		})
		log.Error("获取商品id失败")
		return
	}
	var goods model.Goodsinfo
	err := DB.Where("id=? AND is_deleted=?", id, false).Find(&goods).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取失败",
		})
		log.Errorf("获取失败:%v", err)
		return
	}
	idstr := fmt.Sprintf("%d", goods.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"id":        idstr,
		"goodsName": goods.GoodsName,
		"price":     goods.Price,
	})
	log.Infof("成功获取商品%v的信息", goods.ID)
}
