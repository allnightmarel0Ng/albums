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
	JwtSecretKey      string
	AuthorizationPort string
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
		JwtSecretKey:      os.Getenv("JWT_SECRET_KEY"),
		AuthorizationPort: os.Getenv("AUTHORIZATION_PORT"),
	}, nil
}
