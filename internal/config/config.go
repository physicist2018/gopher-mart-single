package config

import (
	"flag"
	"os"
	"strings"
)

type Configuration struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

func NewConfiguration() *Configuration {
	cfg := &Configuration{}
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "параметры подключения к базе данных")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "адрес системы расчёта начислений")
	return cfg
}

func LoadConfig() (*Configuration, error) {
	cfg := NewConfiguration()
	cfg.Parse()
	return cfg, nil
}

func (c *Configuration) Parse() {
	flag.Parse()
	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		c.RunAddress = envRunAddress
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		c.DatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		c.AccrualSystemAddress = envAccrualSystemAddress
	}
}

func (c *Configuration) String() string {
	return strings.Join([]string{
		"RunAddress: " + c.RunAddress,
		"DatabaseURI: " + c.DatabaseURI,
		"AccrualSystemAddress: " + c.AccrualSystemAddress,
	}, ", \n")
}
