package main

import (
	"os"
)

// DBConfig Здесь мы просто указываем структуру DBConfig
// Здесь перечислены основные данные для соединения к БД
type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

// Config нужен чтобы хранить в себе подконфиги (по типу конфига баз данных)
type Config struct {
	DBConfig DBConfig
}

// NewConfig возвращает конфиг
// Пока что только конфиг для БД
// os.Getenv ищет в переменных окружения (env variables) данные о подключении к БД
// Ключи должны быть такие же, как в .env
func NewConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
