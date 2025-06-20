package router

import (
	"Goods/handler"
	"github.com/gin-gonic/gin"
)

// JWT和casbin的区别：
/*
JWT是用于验证这个人是否为登录的用户，如果用户登录了，那么就获得一个JWT，使用这个JWT可以访问登陆的用户才能访问的东西（修改用户信息，等）
没有登录的时候只能访问注册页面，登录页面

casbin是验证这个人是否为管理员，登录之后，如果是管理员，就可以查看所有用户信息，注销某一个用户，如果是普通人，就只能注销自己的用户
*/

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//需要JWT验证的接口
	goodsG := r.Group("/goods")
	//goodsG.Use(AuthJwt.JWTAuthMiddleware())
	{
		goodsG.GET("", handler.GetOneGoods)
		goodsG.PUT("", handler.UpdateGoods)
		//注销自己的用户
		goodsG.DELETE("", handler.DeleteGoods)
	}

	//consul健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.Status(200)
	})
	return r
}
