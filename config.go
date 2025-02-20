package main

import "os"

type DBConfig struct {
	Username  string
	Password  string
	Host      string
	Port      string
	TableName string
}
type Config struct {
	DBConfig DBConfig
}

func NewConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			Username:  os.Getenv("DB_USERNAME"),
			Password:  os.Getenv("DB_PASSWORD"),
			Host:      os.Getenv("DB_HOST"),
			Port:      os.Getenv("DB_PORT"),
			TableName: os.Getenv("DB_TABLE_NAME"),
		},
	}
}
