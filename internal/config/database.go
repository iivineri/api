package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	host string
	port uint16

	database string
	username string
	password string

	minConns     int
	maxConns     int
	minIdleConns int

	logLevel string
}

func NewDatabaseConfig() *DatabaseConfig {
	godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("DB.HOST", "localhost")
	viper.SetDefault("DB.PORT", "5432")
	viper.SetDefault("DB.NAME", "dev")
	viper.SetDefault("DB.USERNAME", "develop")
	viper.SetDefault("DB.PASSWORD", "develop")
	viper.SetDefault("DB.MIN.CONNS", "5")
	viper.SetDefault("DB.MAX.CONNS", "25")
	viper.SetDefault("DB.MIN.IDLE.CONNS", "5")
	viper.SetDefault("DB.LOG.LEVEL", "error")

	return &DatabaseConfig{
		host:         viper.GetString("DB.HOST"),
		port:         viper.GetUint16("DB.PORT"),
		database:     viper.GetString("DB.NAME"),
		username:     viper.GetString("DB.USERNAME"),
		password:     viper.GetString("DB.PASSWORD"),
		minConns:     viper.GetInt("DB.MIN.CONNS"),
		maxConns:     viper.GetInt("DB.MAX.CONNS"),
		minIdleConns: viper.GetInt("DB.MIN.IDLE.CONNS"),
		logLevel:     viper.GetString("DB.LOG.LEVEL"),
	}
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.username, c.password, c.host, c.port, c.database)
}

func (c *DatabaseConfig) DBMinConns() int {
	return c.minConns
}

func (c *DatabaseConfig) DBMaxConns() int {
	return c.maxConns
}

func (c *DatabaseConfig) DBMinIdleConns() int {
	return c.minIdleConns
}

func (c *DatabaseConfig) DBMaxIdleConns() int {
	return c.minIdleConns
}

func (c *DatabaseConfig) LogLevel() string {
	return c.logLevel
}
