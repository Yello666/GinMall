package Goods

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"net/http"
	"sync"
	"time"
)

var sf *sonyflake.Sonyflake
var once sync.Once

//func initSonyflake() {
//	var settings sonyflake.Settings
//	sf = sonyflake.NewSonyflake(settings)
//}
//
//func GenerateID() (uint64, error) {
//	once.Do(initSonyflake)
//	return sf.NextID()
//}
//var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			// 返回一个唯一的机器 ID
			return 1, nil
		},
		StartTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if sf == nil {
		log.Fatal("Failed to initialize Sonyflake")
	}
}

func GenerateID() (uint64, error) {
	if sf == nil {
		return 0, fmt.Errorf("Sonyflake instance is nil")
	}
	return sf.NextID()
}
func AddGoods(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "AddGoods")
	var goods Goodsinfo
	c.ShouldBindJSON(&goods)
	id, err := GenerateID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成商品id失败",
		})
		log.Errorf("生成商品id失败:%v", err)
		return
	}
	//将商品id改为雪花算法生成的id
	goods.ID = int64(id)
	//存进数据库里
	err = DB.Create(&goods).Error
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
	})
	log.Infof("成功添加商品，商品id:%v", goods.ID)
}
