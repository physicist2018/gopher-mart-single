package config

import (
	"flag"
	"time"
)

type Config struct {
	ServerAddress        string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecret            string
	DBType               string
	TokenExpiry          time.Duration
}

func LoadConfig() *Config {
	// Определение флагов
	serverAddress := flag.String("a", ":8080", "Address and port to run the server (RUN_ADDRESS)")
	databaseURI := flag.String("d", "postgres://kshmirko:123123@kshmirko1.fvds.ru:5432/kshmirko", "Database connection string (DATABASE_URI)")
	accrualSystemAddress := flag.String("r", "http://accrual-system.local", "Accrual system address (ACCRUAL_SYSTEM_ADDRESS)")
	jwtSecret := flag.String("jwt-secret", "secret", "Secret key for JWT token signing (JWT_SECRET)")
	tokenExpiry := flag.Duration("token-expiry", time.Hour*24, "JWT token expiry duration")
	dbType := flag.String("db-type", "postgres", "Database type (DB_TYPE)")

	// Парсинг флагов
	flag.Parse()

	// Создание конфигурации на основе переданных флагов
	config := &Config{
		ServerAddress:        *serverAddress,
		DatabaseURI:          *databaseURI,
		AccrualSystemAddress: *accrualSystemAddress,
		JWTSecret:            *jwtSecret,
		TokenExpiry:          *tokenExpiry,
		DBType:               *dbType,
	}

	return config
}
