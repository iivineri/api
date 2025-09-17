package config

type Config struct {
	App      *AppConfig
	Database *DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		App:      NewAppConfig(),
		Database: NewDatabaseConfig(),
	}
}
