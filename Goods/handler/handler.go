package handler

import (
	"Goods/db"
	"Goods/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetOneGoods(c *gin.Context) {
	log.Info("Get Goodsinfo")
	service.GetGoodsService(c, db.DB)
}
func AddGoods(c *gin.Context) {
	log.Info("Goods Register")
	service.AddGoodsService(c, db.DB)
}
func UpdateGoods(c *gin.Context) {
	log.Info("Goods Update Goodsinfo")
	service.UpdateGoodsService(c, db.DB)
}
func DeleteGoods(c *gin.Context) {
	log.Info("Goods Delete Goods")
	service.DeleteGoodsService(c, db.DB)
}
