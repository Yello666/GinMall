package service

import (
	"User/AuthCasbin"
	"User/AuthJwt" // 你自己的 JWT 工具包路径
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//var Enforcer *casbin.Enforcer // 应该在初始化时设置,一定要引用初始化文件的Enforcer

type AuthCheckRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

func AuthCheckService(c *gin.Context) {
	// 1. 从 Header 获取 token
	token := c.GetHeader("Authorization")
	if token == "" {
		fmt.Println("缺少 Token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少 Token"})
		return
		//return 401, "缺少 Token"
	}

	// 2. 解析 token 获取用户信息
	userClaims, err := AuthJwt.ParseToken(token)
	if err != nil {
		fmt.Println("errorToken 无效", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 无效"})
		return
		//return 401, "Token 无效"
	}
	role := userClaims.Role // "user" / "seller" / "admin"

	// 3. 获取请求体 path 和 method
	var req AuthCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
		//return 400, "参数错误"
	}

	// 4. 使用 Casbin 校验权限
	role = strings.TrimSpace(role)
	path := strings.TrimRight(req.Path, "/") // 只去掉结尾的 /
	path = strings.TrimSpace(path)
	method := strings.ToUpper(strings.TrimSpace(req.Method))

	ok, err := AuthCasbin.Enforcer.Enforce(role, path, method)
	fmt.Printf("role = %q, path = %q, method = %q\n", role, req.Path, req.Method)
	if err != nil {
		fmt.Println("权限引擎错误")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "权限引擎错误"})
		//return 500, "权限引擎错误"
		return
	}

	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
		return
		//return 403, "没有权限"
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	return
	//return 200, "ok"return
}
