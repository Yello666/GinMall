package handler

import (
	"User/db"
	"User/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetUser(c *gin.Context) {
	log.Info("User Get userinfo")
	service.GetUserService(c, db.DB)
}
func AddUser(c *gin.Context) {
	log.Info("User Register")
	service.AddUserService(c, db.DB)
}
func UpdateUser(c *gin.Context) {
	log.Info("User Update userinfo")
	service.UpdateUserService(c, db.DB)
}
func DeleteUser(c *gin.Context) {
	log.Info("User Delete user")
	service.DeleteUserService(c, db.DB)
}
func UpdatePsw(c *gin.Context) {
	log.Info("User Update password")
	service.UpdatePassword(c, db.DB)
}

func Login(c *gin.Context) {
	log.Info("Login")
	service.UserLogin(c, db.DB)
}

func ListUsers(c *gin.Context) {
	log.Info("Admin List users")
	service.ListUsers(c, db.DB)
}
func DeleteAnyUser(c *gin.Context) {
	log.Info("Admin Delete user")

}
