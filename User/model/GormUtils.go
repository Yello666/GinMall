package model

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

func AddUser(DB *gorm.DB, user *Userinfo) error {
	return DB.Create(&user).Error
}

func DeleteUserByID(DB *gorm.DB, str_id string, user *Userinfo) error {
	id_64, _ := strconv.ParseUint(str_id, 10, 64)
	id := uint(id_64)
	return DB.Where("id=?", id).Delete(&user).Error
}
func GetUserByID(DB *gorm.DB, str_id string) (Userinfo, error) {
	log.WithField("func", "GetUserByID")
	var userInfo Userinfo
	int64ID, _ := strconv.ParseInt(str_id, 10, 64)
	id := uint(int64ID)
	if err := DB.Where("id=?", id).Find(&userInfo).Error; err != nil {
		return userInfo, err
	} else {
		return userInfo, nil
	}
}

// UpdateUser 此处的origin必须存在主键id的值
func UpdateUser(origin *Userinfo, updateFields map[string]interface{}, DB *gorm.DB) error {
	log.WithField("func", "updateUser").Info("进入updateUser")

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
