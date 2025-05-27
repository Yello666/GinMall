package router

import (
	"User/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	userG := r.Group("/user")
	{
		userG.GET("/get/:id", handler.GetUser)
		userG.POST("/register", handler.AddUser)
		userG.PUT("/update/:id", handler.UpdateUser)
		userG.DELETE("/delete/:id", handler.DeleteUser)
		userG.POST("/updatePsw/:id", handler.UpdatePsw)
	}
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
	return r
}
