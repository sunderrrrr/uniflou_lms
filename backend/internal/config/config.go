package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	DBSSLMode       string
	JWTSecret       string
	JWTAccessTTL    time.Duration
	JWTRefreshTTL   time.Duration
	LLMAPIURL       string
	LLMAPIKey       string
	LLMModel        string
	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	AppURL          string
	CORSOrigins     []string
	AdminBootstrapEnabled bool
	AdminEmail      string
	AdminPassword   string
	RateLimitAI     int
	RateLimitLogin  int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "465"))
	rateAI, _ := strconv.Atoi(getEnv("RATE_LIMIT_AI_PER_DAY", "10"))
	rateLogin, _ := strconv.Atoi(getEnv("RATE_LIMIT_LOGIN_PER_15MIN", "5"))

	origins := strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	adminEnabled, _ := strconv.ParseBool(getEnv("ADMIN_BOOTSTRAP_ENABLED", "false"))

	return &Config{
		Port:                 getEnv("PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBName:               getEnv("DB_NAME", "uniflou"),
		DBUser:               getEnv("DB_USER", "uniflou"),
		DBPassword:           getEnv("DB_PASSWORD", "change_me"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		JWTSecret:            getEnv("JWT_SECRET", "dev_secret_at_least_32_characters_long"),
		JWTAccessTTL:         accessTTL,
		JWTRefreshTTL:        refreshTTL,
		LLMAPIURL:            getEnv("LLM_API_URL", ""),
		LLMAPIKey:            getEnv("LLM_API_KEY", ""),
		LLMModel:             getEnv("LLM_MODEL", "yandexgpt-lite"),
		SMTPHost:             getEnv("SMTP_HOST", ""),
		SMTPPort:             smtpPort,
		SMTPUser:             getEnv("SMTP_USER", ""),
		SMTPPassword:         getEnv("SMTP_PASSWORD", ""),
		AppURL:               getEnv("APP_URL", "http://localhost:3000"),
		CORSOrigins:          origins,
		AdminBootstrapEnabled: adminEnabled,
		AdminEmail:           getEnv("ADMIN_EMAIL", ""),
		AdminPassword:        getEnv("ADMIN_PASSWORD", ""),
		RateLimitAI:          rateAI,
		RateLimitLogin:       rateLogin,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
