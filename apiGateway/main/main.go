/*
实现了api网关，相当于一个微服务的nginx，自动发现可以用的微服务示例，并将请求路由到微服务上
还可以实现负载均衡的功能
需要先启动consul服务器
consul使用：
brew install consul
service start consul
consul agent -dev
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	// 创建Consul客户端
	consulClient, err := createConsulClient()
	if err != nil {
		log.Fatalf("初始化Consul客户端失败: %v", err)
	}

	// 初始化Gin引擎
	r := gin.Default()

	// 创建反向代理
	//这里的serviceName要和当时注册服务时的一样，返回一个RPS，可以将请求转发到对应的服务器上

	//路由配置：
	//不可以和服务的前缀相同不然会重新导到api网关
	serviceMap := map[string]string{
		"/user":   "user-service",
		"/admin":  "user-service",
		"/seller": "user-service",
		"/goods":  "goods-service",
	}
	//创建RPS
	for prefix, serviceName := range serviceMap {
		proxy := createServiceProxy(consulClient, serviceName)

		r.Any(prefix+"/*proxyPath", func(proxy *httputil.ReverseProxy, serviceName string, prefix string) gin.HandlerFunc {
			return func(c *gin.Context) {
				//排除/auth/check路径的鉴权
				if strings.HasPrefix(c.Request.URL.Path, "/auth/check") {
					proxy.ServeHTTP(c.Writer, c.Request)
					fmt.Println("go out")
					return
				}
				//fmt.Println("check")
				ok, code := checkPermission(c, consulClient)
				if !ok {
					c.Abort()
					fmt.Println("访问商品服务失败，权限不够或者token失效")
					c.JSON(code, gin.H{
						"message": "权限不够",
					})
					return
				}
				proxy.ServeHTTP(c.Writer, c.Request)
			}
		}(proxy, serviceName, prefix))
	}

	// 启动服务
	go func() {
		log.Println("API网关已启动，监听端口8000")
		if err := r.Run(":8000"); err != nil {
			log.Fatalf("启动API网关失败: %v", err)
		}
	}()

	// 优雅退出
	waitForShutdown()
	log.Println("API网关已关闭")
}

//	func checkPermission(c *gin.Context, consulClient *api.Client) bool {
//		// 从Header中获取token
//		token := c.GetHeader("Authorization")
//		if token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization token"})
//			return false
//		}
//
//		// 构造鉴权请求体
//		authPayload := map[string]string{
//			"path":   c.Request.URL.Path,
//			"method": c.Request.Method,
//		}
//		data, _ := json.Marshal(authPayload)
//
//		// 获取 user-service 实例
//		instances, _, err := consulClient.Health().Service("user-service", "", true, nil)
//		if err != nil || len(instances) == 0 {
//			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "用户服务不可用"})
//			return false
//		}
//		// 调用 user-service 的 /auth/check 接口，检查是否为seller角色
//		addr := fmt.Sprintf("http://%s:%d/auth/check", instances[0].Service.Address, instances[0].Service.Port)
//		//req, _ := http.NewRequest("POST", addr, nil)
//		//req.Header.Set("Authorization", token)
//		//
//		//resp, err := http.DefaultClient.Do(req)
//		//if err != nil || resp.StatusCode != 200 {
//		//	c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问：仅限seller角色"})
//		//	return false
//		//}
//		//
//		//// 4. 确认角色是seller
//		//var result struct {
//		//	Role string `json:"role"`
//		//}
//		//if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.Role != "seller" {
//		//	c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问：仅限seller角色"})
//		//	return false
//		//}
//		req, _ := http.NewRequest("POST", addr, bytes.NewReader(data))
//		req.Header.Set("Authorization", token)
//		req.Header.Set("Content-Type", "application/json")
//
//		resp, err := http.DefaultClient.Do(req)
//		if err != nil || resp.StatusCode != 200 {
//			c.JSON(http.StatusForbidden, gin.H{"error": "无权限访问"})
//			return false
//		}
//
//		return true
//	}
func checkPermission(c *gin.Context, consulClient *api.Client) (bool, int) {
	// 1. 检查是否是特定路由和方法
	// 排除/auth/check路径
	if strings.HasPrefix(c.Request.URL.Path, "/auth/check") {
		fmt.Println("true")
		return true, 400
	}
	if c.Request.URL.Path == "/goods/add" && c.Request.Method == "POST" {
		// 2. 从Header中获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization token"})
			return false, 401
		}

		// 3. 构造鉴权请求体
		authPayload := map[string]interface{}{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			// 可以添加更多需要的字段
		}
		data, _ := json.Marshal(authPayload)

		// 4. 获取 user-service 实例
		instances, _, err := consulClient.Health().Service("user-service", "", true, nil)
		if err != nil || len(instances) == 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "用户服务不可用"})
			return false, 500
		}

		// 5. 调用 user-service 的 /auth/check 接口
		addr := fmt.Sprintf("http://%s:%d/auth/check", instances[0].Service.Address, instances[0].Service.Port)
		req, err := http.NewRequest("POST", addr, bytes.NewReader(data))
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "内部服务错误"})
			fmt.Println("内部服务错误")
			return false, 500
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")

		// 6. 发送请求并处理响应
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "鉴权服务调用失败"})
			return false, 500
		}
		defer resp.Body.Close()

		// 7. 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			// 避免重定向循环的关键：只返回false，不设置错误响应
			// 因为可能已经在其他中间件中处理过错误
			return false, resp.StatusCode
		}

		return true, 200
	}

	// 如果不是/goods POST请求，默认允许通过
	return true, 200
}

func createConsulClient() (*api.Client, error) {
	config := api.DefaultConfig()
	return api.NewClient(config)
}

func createServiceProxy(client *api.Client, serviceName string) *httputil.ReverseProxy {
	//获取已经注册过的可用的服务实例
	director := func(req *http.Request) {
		// 如果是 /auth/check，直接访问 user-service，不走网关代理
		if strings.HasPrefix(req.URL.Path, "/auth/check") {
			instances, _, err := client.Health().Service("user-service", "", true, nil)
			if err != nil || len(instances) == 0 {
				req.URL = &url.URL{} // 触发 ErrorHandler
				return
			}
			instance := instances[0]
			req.URL.Scheme = "http"
			req.URL.Host = instance.Service.Address + ":" + strconv.Itoa(instance.Service.Port)
			return
		}
		// 从Consul获取健康的服务实例，叫serverName的服务实例可以有很多个，返回到一个切片里面
		serviceInstances, _, err := client.Health().Service(serviceName, "", true, nil)
		//参数解析：serviceName：从consul获取名为serviceName且健康的服务器，true表示只返回健康的实例
		if err != nil || len(serviceInstances) == 0 {
			// 设置一个非法地址，让 Transport 失败，从而触发 ErrorHandler
			req.URL = &url.URL{}
			return
		}

		// 简单的负载均衡：总是选择第一个实例
		instance := serviceInstances[0] // 实际生产中应该实现更复杂的负载均衡算法
		//修改并添加请求头信息
		target := url.URL{
			Scheme: "http",
			//修改成目标服务器的地址//不可以使用string（）来转换成字符串类型，它会将数字按照ascii码变成字符串
			//而不是“8083”
			Host: instance.Service.Address + ":" + strconv.Itoa(instance.Service.Port),
		}

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = target.Host
	}
	// 添加错误处理函数
	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("代理请求失败: %v", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}

	return &httputil.ReverseProxy{
		Director:     director,
		ErrorHandler: errorHandler,
	}
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
