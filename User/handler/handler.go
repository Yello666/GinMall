package handler

import (
	"User/db"
	"User/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetUser(c *gin.Context) {
	log.Info("Get User")
	service.GetUserService(c, db.DB)
}
func AddUser(c *gin.Context) {
	log.Info("add user")
	service.AddUserService(c, db.DB)
}
func UpdateUser(c *gin.Context) {
	log.Info("update user")
	service.UpdateUserService(c, db.DB)
}
func DeleteUser(c *gin.Context) {
	log.Info("delete user")
	service.DeleteUserService(c, db.DB)
}
func UpdatePsw(c *gin.Context) {
	log.Info("update password")
	service.UpdatePassword(c, db.DB)
}
