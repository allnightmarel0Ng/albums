package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDb          string
	PostgresPort        string
	PostgresUser        string
	PostgresPassword    string
	KafkaPort           string
	RedisPort           string
	JwtSecretKey        string
	AuthorizationPort   string
	ProfilePort         string
	OrderManagementPort string
	SearchEnginePort    string
	GatewayPort         string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		PostgresDb:          os.Getenv("POSTGRES_DB"),
		PostgresPort:        os.Getenv("POSTGRES_PORT"),
		PostgresUser:        os.Getenv("POSTGRES_USER"),
		PostgresPassword:    os.Getenv("POSTGRES_PASSWORD"),
		KafkaPort:           os.Getenv("KAFKA_PORT"),
		RedisPort:           os.Getenv("REDIS_PORT"),
		JwtSecretKey:        os.Getenv("JWT_SECRET_KEY"),
		AuthorizationPort:   os.Getenv("AUTHORIZATION_PORT"),
		ProfilePort:         os.Getenv("PROFILE_PORT"),
		OrderManagementPort: os.Getenv("ORDER_MANAGEMENT_PORT"),
		SearchEnginePort:    os.Getenv("SEARCH_ENGINE_PORT"),
		GatewayPort:         os.Getenv("GATEWAY_PORT"),
	}, nil
}
