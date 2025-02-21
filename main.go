// Вот эта хурма для сваггера прост
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

// initConfig Инициализация конфига, должен быть текстовый файл с именем .env
// В нём должны быть строки по типу DB_USERNAME=postgres и т.д.
func initConfig() {
	// godotenv это просто библиотека чтобы не париться с ручной настройкой парсинга данных из .env
	// godotenv возвращает error, поэтому написано так
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// Загружаем переменные
	initConfig()

	// Про это инфа в config.go
	conf := NewConfig()
	// Connection to Postgres DB, generating wallets with random addresses
	//conn, err := ConnectToDB()

	// Тут происходит подключение к локальной БД
	// в переменной conf записаны мои локальные данные из файла .env благодаря NewConfig()
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.DBConfig.Username, conf.DBConfig.Password, conf.DBConfig.Host, conf.DBConfig.Port, conf.DBConfig.Name))

	// После завершения main() закрываем соединение с БД
	defer conn.Close(context.Background())

	// Новый раутер (или роутер, пох)
	r := chi.NewRouter()

	// Используем логгер чтобы он выводил инфу о запросах и ответах в консоль
	r.Use(middleware.Logger)

	// В раутер добавляем пустой эндпоинт ("/"), а когда мы туда попадаем, вызывается хэндлер Ping
	r.Get("/", Ping)

	// Добавляем эндпоинт для сваггера
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Присобачиваем к основному роутеру роутер, который возвращается из функции psRoutes
	// Основные функции будут лежать в http://localhost:8080/api
	// conn - это соединение с БД (выше)
	r.Mount("/api", psRoutes(conn))

	// Запуск сервера
	err = http.ListenAndServe(":8080", r)

	// Вот эта хрень почему-то не работает
	fmt.Println("server is running and listening on port 8080")
	if err != nil {
		fmt.Printf("couldn't start server: %v", err)
	}
}

// psRoutes Return a new router with pgx connection
// Создаём новый роутер с соединением к БД postgre
func psRoutes(conn *pgx.Conn) chi.Router {
	// Новый экземпляр роутера
	r := chi.NewRouter()

	// Смотри pshandler.go
	handler := PsHandler{conn: conn}

	// в эндпоинт http://localhost:8080/api/wallet/{address}/balance добавляем хэндлер GetBalance
	// address - это переменная ПУТИ
	r.Get("/wallet/{address}/balance", handler.GetBalance)
	// И т.д.
	r.Get("/transactions", handler.GetNLast)
	r.Post("/send", handler.Send)

	// Настроили наш новый роутер, можно его отдавать в main()
	return r
}

// Ping Функция для проверки что всё работает нормально
func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("couldn't write response: %v", err)
	}
}
