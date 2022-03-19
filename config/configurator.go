package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type configuration struct {
	Auth     *AuthConfiguration
	Database *DatabaseConfiguration
}

type AuthConfiguration struct {
	SigningKey string
	TokenTTL   int64

	Email2FA    string
	Password2FA string
	SmtpHost    string
	SmtpPort    int64
}

type DatabaseConfiguration struct {
	Address  string
	Port     int64
	Name     string
	User     string
	Password string
}

var config *configuration

func GetConfig() *configuration {
	if config == nil {
		config = readConfig()
	}

	return config
}

func readConfig() *configuration {
	err := godotenv.Load("./config/dev.env")
	if err != nil {
		log.Println("Can't load config file")
	}

	value, ok := os.LookupEnv("PG_PORT")
	pgPort, err := strconv.ParseInt(value, 10, 64)
	if !ok || err != nil {
		log.Println("No postgres port passed. Using default 5432 PostgreSQL port")
		pgPort = 5432
	}

	value, ok = os.LookupEnv("TOKEN_TTL")
	tokenTTL, err := strconv.ParseInt(value, 10, 64)
	if !ok || err != nil {
		log.Println("No token ttl passed. Using default 86400 ttl")
		tokenTTL = 86400
	}

	value, ok = os.LookupEnv("SMTP_PORT")
	smtpPort, err := strconv.ParseInt(value, 10, 64)
	if !ok || err != nil {
		log.Println("No smtp port passed")
	}

	return &configuration{
		Auth: &AuthConfiguration{
			SigningKey:  os.Getenv("SIGNING_KEY"),
			TokenTTL:    tokenTTL,
			Email2FA:    os.Getenv("EMAIL_2FA"),
			Password2FA: os.Getenv("PASSWORD_2FA"),
			SmtpHost:    os.Getenv("SMTP_HOST"),
			SmtpPort:    smtpPort,
		},
		Database: &DatabaseConfiguration{
			Address:  os.Getenv("PG_IP"),
			Port:     pgPort,
			Name:     os.Getenv("PG_DATABASE"),
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
		},
	}
}
