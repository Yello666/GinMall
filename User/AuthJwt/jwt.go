package AuthJwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 这里声明的值要和密钥一样，不然就会覆盖密钥！！！！
var JwtSecret = []byte("YQ FOREVER")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}
		fmt.Println("authHeader:", authHeader)

		// 验证token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}
		fmt.Println("token:", parts[1])
		// 解析token
		tokenStr := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtSecret, nil
		})
		fmt.Println("JwtSecret:", string(JwtSecret))

		// 签名方法必须是 HMAC 且是 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "意外的签名方法"})
			return
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "无效的token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "无效的token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

//func ParseToken(tokenString string) (*Claims, error) {
//	if len(tokenString) > 6 && tokenString[:7] == "Bearer " {
//		tokenString = tokenString[7:]
//	}
//	fmt.Println(tokenString)
//	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		return JwtSecret, nil
//	})
//	if err != nil {
//		fmt.Println("token无效", err.Error())
//		return nil, err
//	}
//	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
//		return claims, nil
//	}
//	return nil, errors.New("invalid token")
//}

func ParseToken(authHeader string) (*Claims, error) {
	// 0. 复制输入值，防止潜在的内存问题
	tokenString := authHeader

	// 1. 打印原始Token用于调试
	fmt.Printf("原始Token (长度 %d): %q\n", len(tokenString), tokenString)

	// 2. 使用改进的方法提取Token
	//cleanToken, err := extractJWTToken(tokenString)
	//if err != nil {
	//	return nil, fmt.Errorf("提取token失败: %v", err)
	//}
	cleanToken := tokenString[7:]

	fmt.Printf("清理后的Token (长度 %d): %q\n", len(cleanToken), cleanToken)

	// 3. 验证Token格式（基本检查）
	if !isValidJWTFormat(cleanToken) {
		fmt.Println("无效的JWT格式")
		return nil, errors.New("无效的JWT格式")
	}

	// 4. 解析Token
	token, err := jwt.ParseWithClaims(cleanToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return JwtSecret, nil
	})

	if err != nil {
		// 提供更详细的错误信息
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token格式错误: %v", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, fmt.Errorf("token已过期或尚未生效: %v", err)
			} else {
				return nil, fmt.Errorf("token签名验证失败: %v", err)
			}
		}
		return nil, fmt.Errorf("token解析失败: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("Token解析成功: %+v\n", claims)
		return claims, nil
	}

	return nil, errors.New("无效的token")
}

// extractJWTToken 安全地从授权头中提取JWT Token
func extractJWTToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)

	// 检查是否为空
	if authHeader == "" {
		return "", errors.New("授权头为空")
	}

	// 检查是否以"Bearer "开头
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		// 尝试直接解析，可能是没有Bearer前缀的Token
		return authHeader, nil
	}

	// 分割Bearer前缀和Token部分
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("无效的授权头格式")
	}

	// 返回清理后的Token部分
	return strings.TrimSpace(parts[1]), nil
}

// isValidJWTFormat 检查字符串是否符合JWT的基本格式
func isValidJWTFormat(token string) bool {
	// JWT格式: header.payload.signature
	parts := strings.Split(token, ".")
	return len(parts) == 3
}
