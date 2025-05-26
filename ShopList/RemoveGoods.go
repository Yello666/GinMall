package ShopList

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 请求数据的结构体，用于接收 JSON
type RemoveGoodsRequest struct {
	List []int64 `json:"list"`
}

// 错误响应封装
func sendErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"error": errorMessage})
}

// RemoveGoods 处理删除购物车商品的请求
func RemoveGoods(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "RemoveGoods").Info("RemoveGoods function called")

	// 接收请求参数
	var req RemoveGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Invalid request payload")
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.List) == 0 {
		log.Error("Empty list provided")
		sendErrorResponse(c, http.StatusBadRequest, "Empty list")
		return
	}

	// 获取用户 ID
	userIDStr := c.Param("userID")
	if userIDStr == "" {
		log.Error("Missing user ID parameter")
		sendErrorResponse(c, http.StatusBadRequest, "Missing user ID")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		log.WithError(err).Error("Invalid user ID")
		sendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// 查询用户的购物车商品列表
	var existingGoodsList []int64
	if err := DB.Table("relations").
		Where("user_id = ?", userID).
		Pluck("goods_id", &existingGoodsList).Error; err != nil {
		log.WithError(err).Error("Failed to fetch user's shopping list")
		sendErrorResponse(c, http.StatusInternalServerError, "Server error")
		return
	}

	// 将请求的商品 ID 转换为 map，便于快速查找
	goodsToDelete := make(map[int64]struct{})
	for _, goodsID := range req.List {
		goodsToDelete[goodsID] = struct{}{}
	}

	// 找出需要标记删除的商品
	var goodsToMarkDeleted []int64
	for _, goodsID := range existingGoodsList {
		if _, shouldDelete := goodsToDelete[goodsID]; shouldDelete {
			goodsToMarkDeleted = append(goodsToMarkDeleted, goodsID)
		}
	}

	if len(goodsToMarkDeleted) == 0 {
		log.WithField("userID", userID).Info("No goods to delete")
		c.JSON(http.StatusOK, gin.H{"message": "No goods to delete"})
		return
	}

	// 使用事务批量更新商品为已删除状态
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("relations").
			Where("goods_id IN (?)", goodsToMarkDeleted).
			Update("is_deleted", true).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.WithError(err).Error("Failed to delete goods from shopping list")
		sendErrorResponse(c, http.StatusInternalServerError, "Failed to delete goods")
		return
	}

	// 响应成功
	log.WithField("userID", userID).Info("Successfully deleted goods from shopping list")
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
