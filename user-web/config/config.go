package config

type UserConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
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

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	JWTInfo    JWTInfo      `mapstructure:"jwt" json:"jwt"`
	Port       int          `mapstructure:"port" json:"port"`
	UserConfig UserConfig   `mapstructure:"user" json:"user"`
	AliSMS     AliSMSConfig `mapstructure:"ali-sms" json:"ali-sms"`
	Redis      RedisConfig  `mapstructure:"redis" json:"redis"`
}
