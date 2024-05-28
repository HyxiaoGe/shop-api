package config

type UserConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type JWTInfo struct {
	SigningKey string `mapstructure:"signing-key" json:"signing-key"`
}

type AliSMSConfig struct {
	ApiKey    string `mapstructure:"api-key" json:"api-key"`
	ApiSecret string `mapstructure:"api-secret" json:"api-secret"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Host       string       `mapstructure:"host" json:"host"`
	Name       string       `mapstructure:"name" json:"name"`
	Tags       string       `mapstructure:"tags" json:"tags"`
	JWTInfo    JWTInfo      `mapstructure:"jwt" json:"jwt"`
	Port       int          `mapstructure:"port" json:"port"`
	UserConfig UserConfig   `mapstructure:"user" json:"user"`
	AliSMS     AliSMSConfig `mapstructure:"ali-sms" json:"ali-sms"`
	Redis      RedisConfig  `mapstructure:"redis" json:"redis"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
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
