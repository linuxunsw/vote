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
	Port       int
	Production bool
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
	Level       string
	PrettyPrint bool
	AddSource   bool
	// sets httplog formatting flag, and adds top-level attrs which are useful for production
	Concise bool
	// enables loggging of request/response bodies
	LogHTTPBody bool
	Debug       struct {
		Header      string
		HeaderValue string
	}
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
			Level:       GetString("LOGGER_LEVEL", "debug"),
			PrettyPrint: GetBool("LOGGER_PRETTY_PRINT", true),
			AddSource:   GetBool("LOGGER_ADD_SOURCES", false),
			Concise:     GetBool("LOGGER_CONCISE", true),
			LogHTTPBody: GetBool("LOGGER_LOG_HTTP_BODY", true),
			Debug: struct {
				Header      string
				HeaderValue string
			}{
				Header:      GetString("LOGGER_DEBUG_HEADER", "X-Vote-Debug"),
				HeaderValue: GetString("LOGGER_DEBUG_HEADER_VALUE", "show-body"),
			},
		},
	}

	return config
}
