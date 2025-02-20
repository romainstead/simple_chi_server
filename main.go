package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func main() {
	// Connection to Postgres DB, generating wallets with random addresses
	conn, err := ConnectToDB()
	defer conn.Close(context.Background())
	//err = GenerateWallets(conn)
	if err != nil {
		fmt.Printf("couldn't generate wallets: %v", err)
	}
	// Setting router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", Ping)

	// Mount new subrouter
	r.Mount("/api", psRoutes(conn))

	// Launch the server
	err = http.ListenAndServe(":8080", r)
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
