package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ThumborConfig struct {
	host   string
	secret string
}

func NewThumborConfig() *ThumborConfig {
	godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return &ThumborConfig{
		host:   viper.GetString("THUMBOR.HOST"),
		secret: viper.GetString("THUMBOR.SECRET"),
	}
}

func (c *ThumborConfig) Url() string {
	return c.host
}

func (c *ThumborConfig) Secret() string {
	return c.secret
}
