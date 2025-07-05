package configs

import (
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	IPLimitRate  int    `mapstructure:"IP_LIMIT_RATE"`
	IPLimitBurst int    `mapstructure:"IP_LIMIT_BURST"`
	IPLimitTTL   string `mapstructure:"IP_LIMIT_TTL"`
	TokenRate    int    `mapstructure:"TOKEN_RATE"`
	TokenBurst   int    `mapstructure:"TOKEN_BURST"`
	TokenTTL     string `mapstructure:"TOKEN_TTL"`
}

// LoadConfig() => recebe um arquivo de configuração  - init() => função executada antes do método main
func LoadConfig(path string) (*conf, error) {
	viper.SetConfigFile("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
