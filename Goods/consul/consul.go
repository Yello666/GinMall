package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

func CreateConsulClient() (*api.Client, error) {
	config := api.DefaultConfig()
	return api.NewClient(config)
}

func RegisterService(client *api.Client, serviceID string) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "goods-service",
		Port:    8081,
		Address: "localhost",
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://localhost:8081/health",
			Interval:                       "30s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	return client.Agent().ServiceRegister(registration)
}
func DeregisterService(client *api.Client, serviceID string) {
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Fatalf("注销服务失败")
	}
	log.Info("注销了商品服务")
	fmt.Println("注销了商品服务")
}
