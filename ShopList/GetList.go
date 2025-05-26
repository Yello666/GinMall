package ShopList

import (
	"codeup.aliyun.com/codeup/go-micro/Goods"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// GetList 获取用户购物车列表
func GetList(c *gin.Context, DB *gorm.DB) {
	// 获取用户ID
	userIDStr, exists := c.Params.Get("userID")
	if !exists {
		respondWithError(c, http.StatusInternalServerError, "can't get user ID")
		log.Error("获取用户ID失败")
		return
	}

	// 转换用户ID为uint
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid user ID")
		log.Errorf("用户ID转换失败: %v", err)
		return
	}

	// 查询数据库，使用 JOIN 优化查询
	var goodsList []Goods.Goodsinfo

	err = DB.Table("relations").
		Select("goodsinfos.*").
		Joins("JOIN goodsinfos ON relations.goods_id = goodsinfos.id").
		Where("relations.user_id = ? AND relations.is_deleted = ? AND goodsinfos.is_deleted = ?", uint(userID), false, false).
		Find(&goodsList).Error
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "failed to query database")
		log.Errorf("查询数据库失败: %v", err)
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"userID":  userID,
		"list":    goodsList,
	})
	log.Infof("成功查看user %d 的购物车信息", userID)
}

// respondWithError 返回统一的错误响应
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}
