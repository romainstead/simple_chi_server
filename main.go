// Вот эта хурма для сваггера прост
// @title Go-Chi-Swagger-pgx project
// @version 1.0
// @description Пример веб-сервера на Chi с использованием Swagger и работой с БД Postgres
// @host localhost:8080
// @BasePath /api
package main

import (
	_ "chi-crud-api/docs"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func initConfig() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	initConfig()
	config := NewConfig()
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=5433 sslmode=disable",
		config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("couldn't connect to database: ", err)
	}
	handler := Handler{db: db}
	//GenerateWallets(handler.db)
	r := createRouter(handler)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Printf("couldn't start server: %v", err)
	}
}

func createRouter(handler Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", Ping)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/api/wallet/{address}/balance", handler.GetBalance)
	r.Get("/api/transactions", handler.GetNLast)
	r.Post("/api/send", handler.Send)
	return r
}

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("couldn't write response: %v", err)
	}
}
