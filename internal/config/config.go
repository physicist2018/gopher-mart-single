package config

import (
	"flag"
	"os"
	"strings"
)

// Configuration представляет конфигурацию приложения.
// Содержит параметры для запуска сервера, подключения к базе данных и адреса системы начислений.
type Configuration struct {
	RunAddress           string // Адрес интерфейса, на котором запускается сервер
	DatabaseURI          string // Параметры подключения к базе данных
	AccrualSystemAddress string // Адрес системы расчёта начислений
}

// NewConfiguration создает новый экземпляр Configuration с значениями по умолчанию.
// Возвращает указатель на Configuration.
func NewConfiguration() *Configuration {
	cfg := &Configuration{}
	flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "адрес интерфейса, на котором запускать сервер")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "параметры подключения к базе данных")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "адрес системы расчёта начислений")
	return cfg
}

// LoadConfig загружает конфигурацию из флагов командной строки и переменных окружения.
// Возвращает указатель на Configuration и ошибку (в данном случае всегда nil).
func LoadConfig() (*Configuration, error) {
	cfg := NewConfiguration()
	cfg.Parse()
	return cfg, nil
}

// Parse парсит флаги командной строки и переменные окружения, обновляя значения в Configuration.
// Переменные окружения имеют приоритет над флагами командной строки.
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

// String возвращает строковое представление конфигурации.
// Возвращает строку, содержащую все поля Configuration.
func (c *Configuration) String() string {
	return strings.Join([]string{
		"RunAddress: " + c.RunAddress,
		"DatabaseURI: " + c.DatabaseURI,
		"AccrualSystemAddress: " + c.AccrualSystemAddress,
	}, ", \n")
}
