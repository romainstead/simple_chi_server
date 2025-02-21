// @title Go-Chi-Swagger-pgx project
// @version 1.0
// @description Пример веб-сервера на Chi с использованием Swagger и работой с БД Postgres
// @host localhost:8080
// @BasePath /api
package main

import (
	_ "chi-crud-api/docs"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
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
	conf := NewConfig()
	// Connection to Postgres DB, generating wallets with random addresses
	//conn, err := ConnectToDB()
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.DBConfig.Username, conf.DBConfig.Password, conf.DBConfig.Host, conf.DBConfig.Port, conf.DBConfig.TableName))
	defer conn.Close(context.Background())
	//err = GenerateWallets(conn)
	if err != nil {
		fmt.Printf("couldn't generate wallets: %v", err)
	}
	// Setting router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", Ping)
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Mount new subrouter
	r.Mount("/api", psRoutes(conn))

	// Launch the server
	err = http.ListenAndServe(":8080", r)
	fmt.Println("server is running and listening on port 8080")
	if err != nil {
		fmt.Printf("couldn't start server: %v", err)
	}
}

// Return a new router with pgx connection
func psRoutes(conn *pgx.Conn) chi.Router {
	r := chi.NewRouter()
	handler := PsHandler{conn: conn}
	r.Get("/wallet/{address}/balance", handler.GetBalance)
	r.Get("/transactions", handler.GetNLast)
	r.Post("/send", handler.Send)
	return r
}
func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("couldn't write response: %v", err)
	}
}
