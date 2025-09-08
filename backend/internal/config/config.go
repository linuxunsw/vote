package config

import (
	"strings"
	"time"
)

type Config struct {
	API      APIConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OTP      OTPConfig
	Mailer   MailerConfig
	Logger   LoggerConfig
}

type APIConfig struct {
	Version string
}

type ServerConfig struct {
	Port int
	// used for CrossOriginProtection
	AllowedOrigins []string
	// comma seperated allowlist of IPs to accept "Real IP" headers from
	// i.e authorised servers that handle mutliple people's requests
	RealIPAllowlist []string
	RateLimit       RateLimitConfig
}

type DatabaseConfig struct {
	Address            string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        string
}

type JWTConfig struct {
	Secret     string
	CookieName string
	Duration   time.Duration
	Issuer     string
}

type OTPConfig struct {
	Secret   string
	MaxRetry int
	Duration time.Duration
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
	RequestIDHeader string
}

type RateLimitConfig struct {
	RequestLimit int
	WindowLength time.Duration
}

func Load() Config {
	config := Config{
		API: APIConfig{
			Version: "1.0.0",
		},
		Server: ServerConfig{
			Port:            GetInt("PORT", 8888),
			AllowedOrigins:  SplitAndTrim(GetString("ALLOWED_ORIGINS", "")),
			RealIPAllowlist: SplitAndTrim(GetString("REAL_IP_ALLOWLIST", "")),
			RateLimit: RateLimitConfig{
				RequestLimit: 10,
				WindowLength: 10 * time.Second,
			},
		},
		Database: DatabaseConfig{
			Address:            GetString("DB_ADDR", "postgres://admin:admin@localhost/vote?sslmode=disable"),
			MaxOpenConnections: GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConnections: GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:        GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		JWT: JWTConfig{
			Secret:     GetString("JWT_SECRET", "DONTUSEMEINPRODPLEASE"),
			CookieName: "SESSION",
			Duration:   time.Minute * 30,
			Issuer:     "vote-api",
		},
		OTP: OTPConfig{
			Secret:   GetString("OTP_SECRET", "DONTUSEMEINPRODPLEASE-IMEANIT!"),
			MaxRetry: 10,
			Duration: time.Minute * 10,
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
			RequestIDHeader: GetString("LOGGER_REQUEST_ID_HEADER", "X-Request-ID"),
		},
	}

	return config
}

// takes a comma separated string and returns a slice of
// trimmed, non-empty components. If s is empty it returns nil.
func SplitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
