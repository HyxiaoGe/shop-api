package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"shop-api/user-web/global"
	"shop-api/user-web/initialize"
	"shop-api/user-web/utils"
	customValidator "shop-api/user-web/validator"
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

	//	注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", customValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 手机号码格式不合法!", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	zap.S().Debugf("启动web服务，端口为:%d", global.ServerConfig.Port)

	err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
}
