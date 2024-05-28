package config

type GoodsConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTInfo struct {
	SigningKey string `mapstructure:"signing-key" json:"signing-key"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string       `mapstructure:"name" json:"name"`
	JWTInfo     JWTInfo      `mapstructure:"jwt" json:"jwt"`
	Port        int          `mapstructure:"port" json:"port"`
	GoodsConfig GoodsConfig  `mapstructure:"goods" json:"goods"`
	ConsulInfo  ConsulConfig `mapstructure:"consul" json:"consul"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        uint64 `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespace-id" json:"namespace-id"`
	Group       string `mapstructure:"group" json:"group"`
	DataId      string `mapstructure:"data-id" json:"data-id"`
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
}
