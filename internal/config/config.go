package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDb        string
	PostgresPort      string
	PostgresUser      string
	PostgresPassword  string
	RedisPort         string
	JwtSecretKey      string
	AuthorizationPort string
	ProfilePort       string
	GatewayPort       string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		PostgresDb:        os.Getenv("POSTGRES_DB"),
		PostgresPort:      os.Getenv("POSTGRES_PORT"),
		PostgresUser:      os.Getenv("POSTGRES_USER"),
		PostgresPassword:  os.Getenv("POSTGRES_PASSWORD"),
		RedisPort:         os.Getenv("REDIS_PORT"),
		JwtSecretKey:      os.Getenv("JWT_SECRET_KEY"),
		AuthorizationPort: os.Getenv("AUTHORIZATION_PORT"),
		ProfilePort:       os.Getenv("PROFILE_PORT"),
		GatewayPort:       os.Getenv("GATEWAY_PORT"),
	}, nil
}
