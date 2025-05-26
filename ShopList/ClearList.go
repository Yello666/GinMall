package ShopList

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func ClearShopList(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "ClearShopList")

	// 接收用户ID参数
	userIDStr := c.Param("userID")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user ID"})
		log.Error("Missing user ID parameter")
		return
	}

	// 将用户ID字符串转换为uint
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// 找到要清空的购物车
	var shopList Relation
	shopList.UserID = uint(userID)
	if err := DB.Model(&shopList).Where("UserID=?", shopList.UserID).Update("Is_deleted", true).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 如果购物车不存在，我们可以选择创建一个空的购物车或者返回一个错误
			// 这里我们选择创建一个空的购物车
			shopList.UserID = uint(userID)
			shopList.GoodsID = 0
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			log.WithError(err).Error("Failed to fetch user shop list")
			return
		}
	}

	//// 清空购物车商品列表
	//shopList.GoodsID= 0

	//// 更新用户的购物车
	//if err := DB.Save(&shopList).Error; err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update shopping list"})
	//	log.WithError(err).Error("Failed to update user's shopping list")
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	log.WithField("userID", userID).Info("Successfully cleared user's shopping list")
}
