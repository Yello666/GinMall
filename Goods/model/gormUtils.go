package model

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

func AddGoods(DB *gorm.DB, Goods *GoodsInfo) error {
	return DB.Create(&Goods).Error
}

func DeleteGoodsByID(DB *gorm.DB, str_id string, Goods *GoodsInfo) error {
	id_64, _ := strconv.ParseUint(str_id, 10, 64)
	id := uint(id_64)
	return DB.Where("id=? and is_deleted=?", id, false).Delete(&Goods).Error
}
func GetGoodsByID(DB *gorm.DB, str_id string) (GoodsInfo, error) {
	log.WithField("func", "GetGoodsByID")
	var GoodsInfo GoodsInfo
	int64ID, _ := strconv.ParseInt(str_id, 10, 64)
	id := uint(int64ID)
	if err := DB.Where("id=? and is_deleted=?", id, false).Find(&GoodsInfo).Error; err != nil {
		return GoodsInfo, err
	} else {
		return GoodsInfo, nil
	}
}

// UpdateGoods 此处的origin必须存在主键id的值
func UpdateGoods(origin *GoodsInfo, updateFields map[string]interface{}, DB *gorm.DB) error {
	log.WithField("func", "updateGoods").Info("进入updateGoods")

	//只进行修改操作
	if err := DB.Model(origin).Updates(updateFields).Error; err != nil {
		log.Errorf("更新出错:%v", err)
		return err
	}
	return nil

}

// Paginate 分页查询插件
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
