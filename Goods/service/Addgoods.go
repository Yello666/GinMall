package service

import (
	"Goods/model"
	"Goods/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func AddGoodsService(c *gin.Context, DB *gorm.DB, rdb *redis.Client) {
	log.WithField("func", "AddGoods").Info("add goods")

	// 1. 解析表单字段
	goodsName := c.PostForm("goodsName")
	priceStr := c.PostForm("price")
	stockStr := c.PostForm("stock")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "价格格式错误"})
		log.Errorf("价格格式错误: %v", err)
		return
	}

	stock, err := strconv.ParseInt(stockStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "库存格式错误"})
		log.Errorf("库存格式错误: %v", err)
		return
	}

	// 2. 获取上传的文件（图片）
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片上传失败"})
		log.Errorf("图片上传失败: %v", err)
		return
	}

	// 3. 保存图片文件
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	imagePath := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		log.Errorf("保存图片失败: %v", err)
		return
	}

	// 4. 生成雪花 ID
	id, err := utils.GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成商品id失败"})
		log.Errorf("生成商品id失败: %v", err)
		return
	}

	// 5. 构建商品模型
	goods := model.GoodsInfo{
		ID:         int64(id),
		Is_deleted: false,
		GoodsName:  goodsName,
		Price:      price,
		Stock:      stock,
		ImagePath:  "/" + imagePath, // 前端可以直接通过静态路由访问
	}

	// 6. 存入数据库
	err = model.AddGoods(DB, &goods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储失败"})
		log.Errorf("存储失败: %v", err)
		return
	}
	//7.保存进缓存里面
	// 手动序列化为 JSON
	goodsJSON, err := json.Marshal(goods)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize user"})
		return
	}
	//保存进缓存里面
	redisKey := fmt.Sprintf("goods:%d", goods.ID)
	err = rdb.Set(ctx, redisKey, goodsJSON, time.Hour).Err()

	// 8. 响应成功
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"id":        goods.ID,
		"goodsName": goods.GoodsName,
		"price":     goods.Price,
		"stock":     goods.Stock,
		"image":     goods.ImagePath,
	})
	log.Infof("成功添加商品，商品id:%v", goods.ID)
}
