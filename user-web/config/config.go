package config

type UserConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name       string     `mapstructure:"name" json:"name"`
	Port       int        `mapstructure:"port" json:"port"`
	UserConfig UserConfig `mapstructure:"user" json:"user"`
}
