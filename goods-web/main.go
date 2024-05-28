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

	//register_client := consul.NewRegisterClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	//e := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, []string{global.ServerConfig.Tags}, "goods_test_id")
	//if e != nil {
	//	zap.S().Panic("注册服务失败:", e.Error())
	//}

	zap.S().Debugf("启动web服务，端口为:%d", global.ServerConfig.Port)
	go func() {
		err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
		if err != nil {
			zap.S().Panic("启动失败:", err.Error())
		}
	}()
	//	接收终止信号
	//quit := make(chan os.Signal)
	//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//<-quit
	//if e = register_client.DeRegister("test_id"); e != nil {
	//	zap.S().Info("注销服务失败:", e.Error())
	//} else {
	//	zap.S().Info("注销服务成功")
	//}
}
