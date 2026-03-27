package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App           AppConfig
	Admin         AdminConfig
	Socket        SocketConfig
	DB            DBConfig
	Redis         RedisConfig
	NATS          NATSConfig
	Elasticsearch ElasticsearchConfig
	Log           LogConfig
	JWT           JWTConfig
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type AdminConfig struct {
	Port string
}

type ElasticsearchConfig struct {
	Addresses    []string
	Username     string
	Password     string
	ProductIndex string
}

type LogConfig struct {
	Level        string
	LogstashAddr string
}

type NATSConfig struct {
	URL string
}

type SocketConfig struct {
	Port string
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Admin: AdminConfig{
			Port: getEnv("ADMIN_PORT", "8082"),
		},
		Socket: SocketConfig{
			Port: getEnv("SOCKET_PORT", "8081"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "retail_store"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
		},
		Elasticsearch: ElasticsearchConfig{
			Addresses:    strings.Split(getEnv("ELASTICSEARCH_ADDRESSES", "http://localhost:9200"), ","),
			Username:     getEnv("ELASTICSEARCH_USERNAME", ""),
			Password:     getEnv("ELASTICSEARCH_PASSWORD", ""),
			ProductIndex: getEnv("ELASTICSEARCH_PRODUCT_INDEX", "products"),
		},
		Log: LogConfig{
			Level:        getEnv("LOG_LEVEL", "info"),
			LogstashAddr: getEnv("LOGSTASH_ADDRESS", ""),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: time.Duration(jwtExpHours) * time.Hour,
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
