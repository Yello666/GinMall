package service

import (
	"User/AuthJwt"
	"User/model"
	"User/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

const EXPIRETIME = 7 //7天过期，登录有效期为7天，7天后需要重新登录
type login struct {
	UserID   string `json:"user_id" binding:"required,min=1,max=50"`
	Password string `json:"-" binding:"required,min=6,max=20"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

func UserLogin(c *gin.Context, DB *gorm.DB) {
	//1.得到id和密码
	var loginInfo login
	if err := c.BindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "缺少参数或者参数格式错误",
			"error":   err.Error(),
		})
		return
	}
	//得到用户原来的密码
	var user model.Userinfo
	var err error
	user, err = model.GetUserByID(DB, loginInfo.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "获取用户信息失败",
			"error":   err.Error(),
		})
		return
	}
	//校验密码
	hash := user.PasswordHash
	if utils.CheckPasswordHash(loginInfo.Password, hash) {
		//登录成功
		// 登录成功，生成JWT
		claims := AuthJwt.Claims{
			UserID: loginInfo.UserID, // 从数据库获取的用户ID
			Role:   loginInfo.Role,   // 从数据库获取的角色
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(24 * EXPIRETIME * time.Hour).Unix(),
				Issuer:    "gin-mall",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(AuthJwt.JwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":   tokenString,
			"message": "登录成功",
		})
		//TODO 可以进行websocket登录
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "密码错误,登录失败",
		})
		return
	}

}
