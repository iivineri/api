package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	environment     string
	port            uint16
	logLevel        string
	prefork         bool
	jwtSecret       string
	swaggerEnabled  bool
	swaggerHost     string
	swaggerBasePath string
	swaggerSchemes  []string
}

func NewAppConfig() *AppConfig {
	godotenv.Load()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("LOG.LEVEL", "error")
	viper.SetDefault("PREFORK", "false")
	viper.SetDefault("JWT.SECRET", "your-secret-key-change-this-in-production")
	viper.SetDefault("SWAGGER.ENABLED", "false")
	viper.SetDefault("SWAGGER.HOST", "localhost:8080")
	viper.SetDefault("SWAGGER.BASE_PATH", "/api/v1")
	viper.SetDefault("SWAGGER.SCHEMES", "http,https")

	envValue := viper.GetString("ENV")
	switch envValue {
	case "dev", "develop":
		envValue = "development"
	case "prod", "production":
		envValue = "production"
	default:
		envValue = "development"
	}

	swaggerSchemes := strings.Split(viper.GetString("SWAGGER.SCHEMES"), ",")
	for i, scheme := range swaggerSchemes {
		swaggerSchemes[i] = strings.TrimSpace(scheme)
	}

	return &AppConfig{
		environment:     envValue,
		port:            viper.GetUint16("PORT"),
		logLevel:        viper.GetString("LOG.LEVEL"),
		prefork:         viper.GetBool("PREFORK"),
		jwtSecret:       viper.GetString("JWT.SECRET"),
		swaggerEnabled:  viper.GetBool("SWAGGER.ENABLED"),
		swaggerHost:     viper.GetString("SWAGGER.HOST"),
		swaggerBasePath: viper.GetString("SWAGGER.BASE_PATH"),
		swaggerSchemes:  swaggerSchemes,
	}
}

func (c *AppConfig) Environment() string {
	return c.environment
}
func (c *AppConfig) Port() uint16 {
	return c.port
}

func (c *AppConfig) IsDevelopment() bool {
	return c.environment == "development"
}

func (c *AppConfig) IsProduction() bool {
	return c.environment == "production"
}

func (c *AppConfig) Prefork() bool {
	return c.prefork
}

func (c *AppConfig) LogLevel() string {
	return c.logLevel
}

func (c *AppConfig) JWTSecret() string {
	return c.jwtSecret
}

func (c *AppConfig) SwaggerEnabled() bool {
	return c.swaggerEnabled
}

func (c *AppConfig) SwaggerHost() string {
	return c.swaggerHost
}

func (c *AppConfig) SwaggerBasePath() string {
	return c.swaggerBasePath
}

func (c *AppConfig) SwaggerSchemes() []string {
	return c.swaggerSchemes
}
