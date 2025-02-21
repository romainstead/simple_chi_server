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

// GetBalance возвращает баланс кошелька по адресу
// @Summary возврат баланс кошелька по заданному адресу
// @Description возвращает JSON с адресом и балансом кошелька
// @Tags примеры
// @Produce json
// @Param address path string true "Адрес кошелька"
// @Success 200 {object} map[string]string
// @Router /wallet/{address}/balance [get]
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

// GetNLast возвращает N последних по времени транзакций
// @Summary показ N последних транзакций
// @Description возвращает JSON Array длиной N объектов Transaction
// @Tags примеры
// @Produce json
// @Param N query number true "Количество передаваемых транзакций"
// @Success 200 {array} map[string]string
// @Router /transactions/{N} [get]
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

// Send переводит средства с кошелька на кошелек
// @Summary снимает с баланса кошелька отправителя заданную сумму и прибавляет её к балансу кошелька получателя
// @Description взаимодействует с БД через UPDATE, INSERT, BEGIN TRANSACTION
// @Tags примеры
// @Produce json
// @Param to query string true "Адрес получателя"
// @Param from query string true "Адрес отправителя"
// @Param amount query number true "Сумма для отправки"
// @Success 200
// @Failure 400
// @Router /send [post]
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
