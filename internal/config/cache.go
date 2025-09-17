package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type CacheConfig struct {
	host string
	port int
}

func NewCacheConfig() *CacheConfig {
	godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return &CacheConfig{
		host: viper.GetString("CACHE.HOST"),
		port: viper.GetInt("CACHE.PORT"),
	}
}
