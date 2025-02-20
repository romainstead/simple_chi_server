package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
)

type PsHandler struct {
	conn *pgx.Conn
}

func newPsHandler(conn *pgx.Conn) *PsHandler {
	return &PsHandler{conn: conn}
}

func (p *PsHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// IF ADDRESS IS PRESENT IN DB
	address := chi.URLParam(r, "address")
	wallet, _ := GetBalance(address, p.conn)

	// ENCODE TO JSON AND SEND
	if wallet.Address != "" {
		err := json.NewEncoder(w).Encode(wallet)
		if err != nil {
			log.Fatal(err)
		}
	}
	if wallet.Address == "" {
		err := json.NewEncoder(w).Encode(http.StatusNotFound)
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Header().Set("Content-Type", "application/json")
}
func (p *PsHandler) GetNLast(w http.ResponseWriter, r *http.Request) {
	count := r.URL.Query().Get("count")
	n, err := strconv.Atoi(count)
	if n <= 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Fatalf("couldn't read argument: %v", err)
	}
	transactions, err := GetNLast(n, p.conn)
	if len(transactions) == 0 {
		json.NewEncoder(w).Encode(http.StatusInternalServerError)
		log.Fatalf("transactions not found")
	}
	if err != nil {
		log.Fatalf("couldn't fetch data from db: %v", err)
	}
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		log.Fatalf("couldn't encode data: %v", err)
	}
}
func (p *PsHandler) Send(w http.ResponseWriter, r *http.Request) {
	FromAddress := r.URL.Query().Get("from")
	ToAddress := r.URL.Query().Get("to")
	Amount := r.URL.Query().Get("amount")
	if FromAddress != "" && ToAddress != "" && Amount != "" {
		amount, _ := strconv.ParseFloat(Amount, 64)
		if amount <= 0 {
			w.WriteHeader(400)
			err := json.NewEncoder(w).Encode("amount value must be positive")
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		err := Send(FromAddress, ToAddress, amount, p.conn)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}
		err = json.NewEncoder(w).Encode(http.StatusOK)
		if err != nil {
			log.Fatal(err)
		}
	}
}
