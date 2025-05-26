package ShopList

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type listRequest struct {
	UserID  uint  `json:"userID"`
	GoodsID int64 `json:"goodsID"`
}

func AddGoodsToList(c *gin.Context, DB *gorm.DB) {
	//接收商品参数
	var getList listRequest
	err := c.ShouldBindJSON(&getList)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少参数（userid，goodsid）",
		})
		log.Error("缺少参数")
	}
	////接收用户id参数
	//id, ok := c.Params.Get("userID")
	//if !ok {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": "miss id",
	//	})
	//	log.Error("缺少id参数")
	//}
	//
	//// 将id字符串解析为int64
	//intVal, err := strconv.ParseInt(id, 10, 64)
	//if err != nil {
	//	fmt.Println("转换错误:", err)
	//	return
	//}
	//保存用户购物车信息
	var shopList Relation
	shopList.GoodsID = getList.GoodsID
	shopList.UserID = getList.UserID

	if err := DB.Create(&shopList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "server error",
		})
		log.Errorf("数据库访问错误:%v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
	log.Infof("用户%v添加了商品到购物车", shopList.UserID)
}
