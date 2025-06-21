package service

import (
	"User/model"
	"User/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

func GetUserService(c *gin.Context, DB *gorm.DB) {
	log.WithField("func", "GetUser")
	var userInfo model.Userinfo
	//从上下文中获取userID
	str_id, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息过期，请重新登录"})
		return
	}
	//从数据库中获取user信息
	userInfo, err := model.GetUserByID(DB, str_id)
	if err != nil {
		//找不到id
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未注册登录",
		})
		log.Errorf("id:%v未注册登录:%v", str_id, err)
		return
	} else {
		//返回id,用户名,性别,
		if userInfo.Age == nil {
			age := 0
			userInfo.Age = &age
		}
		c.JSON(http.StatusOK, gin.H{
			"message":   "success",
			"id":        userInfo.ID,
			"user_name": userInfo.UserName,
			"sex":       userInfo.Sex,
			"age":       *userInfo.Age,
			"email":     userInfo.Email,
		})
		log.Infof("获取用户%v信息成功", str_id)
	}
}

// ListUsers 获取全部用户信息（带分页和权限控制）
func ListUsers(c *gin.Context, DB *gorm.DB) {
	//登录校验
	_, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户信息过期，请重新登录"})
		return
	}
	// 权限校验：仅管理员可访问（JWT中角色为"admin"）
	role, ok := utils.GetRole(c)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足，仅管理员可操作"})
		log.Warn("非管理员尝试获取用户列表")
		return
	}

	// 获取分页参数
	page, pageSize, err := utils.GetPaginationParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 构建查询条件（示例：排除敏感字段，按ID降序排序）
	var users []model.Userinfo
	// 使用Select指定返回字段，避免暴露敏感信息（如Password）
	result := DB.
		Select("id, user_name, sex, age, email, created_at"). // 筛选字段
		Order("id DESC"). // 按ID降序排列
		Scopes(model.Paginate(page, pageSize)). // 分页插件（见下方）
		Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		log.Errorf("获取用户列表失败: %v", err)
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"total":    result.RowsAffected,
		"page":     page,
		"pageSize": pageSize,
		"data":     users,
	})
	log.Infof("管理员获取用户列表，共%d条记录", result.RowsAffected)
}
