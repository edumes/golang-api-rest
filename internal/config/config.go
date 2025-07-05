package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	URL      string `mapstructure:"url"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

func LoadConfig() (*Config, error) {
	logger := logrus.New()

	setDefaults()

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Failed to read .env file, using environment variables")
	}

	viper.AutomaticEnv()

	if err := bindEnvVars(); err != nil {
		return nil, fmt.Errorf("failed to bind env vars: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger.Info("Configuration loaded successfully")
	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "golang_api_rest")
	viper.SetDefault("database.sslmode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration", "24h")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"*"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.exposed_headers", []string{"Content-Length"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 86400)
}

func bindEnvVars() error {
	// Server
	if err := viper.BindEnv("server.port", "APP_PORT"); err != nil {
		return fmt.Errorf("bindEnv server.port: %w", err)
	}
	if err := viper.BindEnv("server.read_timeout", "APP_READ_TIMEOUT"); err != nil {
		return fmt.Errorf("bindEnv server.read_timeout: %w", err)
	}
	if err := viper.BindEnv("server.write_timeout", "APP_WRITE_TIMEOUT"); err != nil {
		return fmt.Errorf("bindEnv server.write_timeout: %w", err)
	}
	if err := viper.BindEnv("server.idle_timeout", "APP_IDLE_TIMEOUT"); err != nil {
		return fmt.Errorf("bindEnv server.idle_timeout: %w", err)
	}

	// Database
	if err := viper.BindEnv("database.url", "DATABASE_URL"); err != nil {
		return fmt.Errorf("bindEnv database.url: %w", err)
	}
	if err := viper.BindEnv("database.host", "DB_HOST"); err != nil {
		return fmt.Errorf("bindEnv database.host: %w", err)
	}
	if err := viper.BindEnv("database.port", "DB_PORT"); err != nil {
		return fmt.Errorf("bindEnv database.port: %w", err)
	}
	if err := viper.BindEnv("database.user", "DB_USER"); err != nil {
		return fmt.Errorf("bindEnv database.user: %w", err)
	}
	if err := viper.BindEnv("database.password", "DB_PASSWORD"); err != nil {
		return fmt.Errorf("bindEnv database.password: %w", err)
	}
	if err := viper.BindEnv("database.dbname", "DB_NAME"); err != nil {
		return fmt.Errorf("bindEnv database.dbname: %w", err)
	}
	if err := viper.BindEnv("database.sslmode", "DB_SSL_MODE"); err != nil {
		return fmt.Errorf("bindEnv database.sslmode: %w", err)
	}

	// JWT
	if err := viper.BindEnv("jwt.secret", "APP_JWT_SECRET"); err != nil {
		return fmt.Errorf("bindEnv jwt.secret: %w", err)
	}
	if err := viper.BindEnv("jwt.expiration", "APP_JWT_EXPIRATION"); err != nil {
		return fmt.Errorf("bindEnv jwt.expiration: %w", err)
	}

	// Logging
	if err := viper.BindEnv("logging.level", "LOG_LEVEL"); err != nil {
		return fmt.Errorf("bindEnv logging.level: %w", err)
	}
	if err := viper.BindEnv("logging.format", "LOG_FORMAT"); err != nil {
		return fmt.Errorf("bindEnv logging.format: %w", err)
	}

	return nil
}

func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if config.Database.URL == "" {
		if config.Database.Host == "" || config.Database.Port == "" || config.Database.User == "" || config.Database.DBName == "" {
			return fmt.Errorf("database configuration is incomplete")
		}
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

func (c *Config) GetDatabaseURL() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.DBName, c.Database.SSLMode)
}
