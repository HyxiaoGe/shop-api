package config

type UserConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTInfo struct {
	SigningKey string `mapstructure:"signing-key" json:"signing-key"`
}

type ServerConfig struct {
	Name       string     `mapstructure:"name" json:"name"`
	JWTInfo    JWTInfo    `mapstructure:"jwt" json:"jwt"`
	Port       int        `mapstructure:"port" json:"port"`
	UserConfig UserConfig `mapstructure:"user" json:"user"`
}
