package main

import (
	"fmt"
	"go.uber.org/zap"
	"shop-api/user-web/initialize"
)

func main() {
	port := 8021
	initialize.InitLogger()
	Router := initialize.Routers()

	zap.S().Debugf("启动web服务，端口为:%d", port)

	err := Router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
