package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	log "github.com/sirupsen/logrus"

	"Goods/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var ctx = context.Background()

func GetGoodsService(c *gin.Context, db *gorm.DB, rdb *redis.Client) {
	// 1. 获取商品ID
	log.Info("get goods info")
	goodsID := c.Param("id")
	redisKey := fmt.Sprintf("goods:%s", goodsID)

	// 2. 先查 Redis 缓存
	val, err := rdb.Get(ctx, redisKey).Result()
	if err == nil {
		// 缓存命中
		var cachedGoods model.GoodsInfo
		if err = json.Unmarshal([]byte(val), &cachedGoods); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"source": "cache",
				"data":   cachedGoods,
			})
			return
		}
	}

	// 3. 缓存未命中或反序列化失败 → 查数据库
	var goods model.GoodsInfo
	if goods, err = model.GetGoodsByID(db, goodsID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "商品不存在",
		})
		return
	}

	// 4. 写入 Redis 缓存（设置 10 分钟过期）
	jsonVal, _ := json.Marshal(goods)
	rdb.Set(ctx, redisKey, jsonVal, 10*time.Minute)

	// 5. 返回数据
	c.JSON(http.StatusOK, gin.H{
		"source": "database",
		"data":   goods,
	})
}
