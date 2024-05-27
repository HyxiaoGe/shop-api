package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop-api/user-web/global"
	"shop-api/user-web/proto"
)

func InitClient() {
	// 从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	consulInfo := global.ServerConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userHost := ""
	userPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Errorf("consul服务连接失败: %v", err)
		return
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserConfig.Name))
	if err != nil {
		zap.S().Errorf("consul服务查询失败: %v", err)
		return
	}
	for _, value := range data {
		userHost = value.Address
		userPort = value.Port
		break
	}
	if userHost == "" {
		zap.S().Fatal("用户服务不可用")
		return
	}

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userHost, userPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorf("连接用户服务失败: %v", err)
		return
	}
	global.UserClient = proto.NewUserClient(userConn)
}
