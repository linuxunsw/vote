package config

import (
	"time"
)

type Config struct {
	API      APIConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Mailer   MailerConfig
	Logger   LoggerConfig
}

type APIConfig struct {
	Version string
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Address            string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        string
}

type JWTConfig struct {
	Secret   string
	Duration time.Duration
	Issuer   string
}

type MailerConfig struct {
	ResendAPIKey string
	FromEmail    string
}

type LoggerConfig struct {
	Level    string
	Encoding string
}

func Load() Config {
	config := Config{
		API: APIConfig{
			Version: "1.0.0",
		},
		Server: ServerConfig{
			Port: GetInt("PORT", 8888),
		},
		Database: DatabaseConfig{
			Address:            GetString("DB_ADDR", "postgres://admin:admin@localhost/vote?sslmode=disable"),
			MaxOpenConnections: GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConnections: GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:        GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		JWT: JWTConfig{
			Secret:   GetString("JWT_SECRET", "DONTUSEMEINPRODPLEASE"),
			Duration: time.Minute * 30,
			Issuer:   "vote-api",
		},
		Mailer: MailerConfig{
			ResendAPIKey: GetString("RESEND_API_KEY", ""),
			FromEmail:    GetString("MAILER_FROM_EMAIL", ""),
		},
		Logger: LoggerConfig{
			Level:    GetString("LOGGER_LEVEL", "debug"),
			Encoding: GetString("LOGGER_ENCODING", "console"),
		},
	}

	return config
}
