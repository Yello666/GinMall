package service

import (
	"Goods/model"
	"Goods/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"

	"net/http"
)

type AddGoodsModel struct {
	GoodsName string  `json:"goodsName" binding:"required,min=6,max=50"`
	Price     float64 `json:"price" binding:"required"`
	Stock     int64   `json:"stock" binding:"required"`
} //如何添加商品图片？

func AddGoods(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "AddGoods").Info("add goods")
	var goodsModel AddGoodsModel
	err := c.ShouldBindJSON(&goodsModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "商品数据格式错误",
			"error":   err.Error(),
		})
		log.Errorf("商品数据格式错误:%v", err)
		return
	}
	id, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成商品id失败",
		})
		log.Errorf("生成商品id失败:%v", err)
		return
	}
	//3.构建商品模型
	//将商品id改为雪花算法生成的id
	goods := model.GoodsInfo{
		ID:         int64(id),
		Is_deleted: false,
		GoodsName:  goodsModel.GoodsName,
		Price:      goodsModel.Price,
		Stock:      goodsModel.Stock,
	}

	//存进数据库里
	err = model.AddGoods(DB, &goods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "存储失败",
		})
		log.Errorf("存储失败:%v", err)
		return
	}
	idstr := fmt.Sprintf("%d", goods.ID)
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"id":        idstr,
		"goodsName": goods.GoodsName,
		"price":     goods.Price,
		"stock":     goods.Stock,
	})
	log.Infof("成功添加商品，商品id:%v", goods.ID)
}
