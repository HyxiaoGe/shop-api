package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop-api/user-web/config"
)

var (
	Trans        ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
