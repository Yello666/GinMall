package Goods

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetSomeGoods(c *gin.Context, DB *gorm.DB) {
	// 从查询字符串中获取分页参数
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	// 计算偏移量
	offset := (page - 1) * limit
	// 查询数据库，应用分页逻辑
	var goodss []Goodsinfo
	var total int64
	if err := DB.Model(&Goodsinfo{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Error("查询错误")
		return
	}
	if err := DB.Offset(offset).Limit(limit).Find(&goodss).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Error("批量查询错误")
		return
	}
	var response []GoodsResponse
	for k, v := range goodss {
		response[k].Name = v.GoodsName
		response[k].ID = v.ID
		response[k].Price = v.Price
	}
	// 计算总页数（可选，但通常很有用）
	//totalPages := int(math.Ceil(float64(total) / float64(limit)))
	// 返回分页结果和元数据
	c.JSON(http.StatusOK, gin.H{
		"goods": response,
	})
}
