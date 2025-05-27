package AuthCasbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	enforcer *casbin.Enforcer
)

// CasbinMiddleware Casbin权限控制中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到用户角色"})
			c.Abort()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 检查权限
		allowed, err := enforcer.Enforce(role, path, method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 初始化Casbin
func InitCasbin() error {
	//// 使用GORM适配器连接数据库
	//a, err := gormadapter.NewAdapterByDB(db.DB)
	//if err != nil {
	//	return fmt.Errorf("创建Casbin适配器失败: %v", err)
	//}

	// 加载模型配置
	var err error
	enforcer, err = casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		return fmt.Errorf("创建Casbin执行器失败: %v", err)
	}

	// 加载策略
	if err = enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("加载Casbin策略失败: %v", err)
	}

	return nil
}
