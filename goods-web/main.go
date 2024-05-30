package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	_ "os"
	_ "os/signal"
	"shop-api/goods-web/global"
	"shop-api/goods-web/initialize"
	"shop-api/goods-web/utils"
	_ "shop-api/goods-web/utils/register/consul"
	_ "syscall"
)

func main() {
	//	初始化logger
	initialize.InitLogger()
	//	初始化配置文件
	initialize.InitConfig()
	//	初始化routers
	Router := initialize.Routers()
	//	初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 初始化连接
	initialize.InitClient()

	viper.AutomaticEnv()
	debug := viper.GetBool("SHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	zap.S().Debugf("启动web服务，端口为:%d", global.ServerConfig.Port)
	go func() {
		err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
}
