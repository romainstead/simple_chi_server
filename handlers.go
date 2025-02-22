package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	db *gorm.DB
}

// GetBalance возвращает баланс кошелька по адресу
// @Summary возврат баланс кошелька по заданному адресу
// @Description возвращает JSON с адресом и балансом кошелька
// @Tags примеры
// @Produce json
// @Param address path string true "Адрес кошелька"
// @Success 200 {object} map[string]string
// @Router /wallet/{address}/balance [get]
func (p *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	wallet, err := GetWallet(address, p.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetNLast возвращает N последних по времени транзакций
// @Summary показ N последних транзакций
// @Description возвращает JSON Array длиной N объектов Transaction
// @Tags примеры
// @Produce json
// @Param count query number true "Количество передаваемых транзакций"
// @Success 200 {array} map[string]string
// @Router /transactions [get]
func (p *Handler) GetNLast(w http.ResponseWriter, r *http.Request) {
	n := r.URL.Query().Get("n")
	count, err := strconv.Atoi(n)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(502)
		if err := json.NewEncoder(w).Encode("'n' must be a positive number"); err != nil {
			log.Println("couldn't write response", err)
			return
		}
		return
	}
	if count <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode("'n' must be a positive number"); err != nil {
			log.Println("couldn't write response", err)
			return
		}
		return
	}
	rows, err := GetNLast(count, p.db)
	if err != nil {
		json.NewEncoder(w).Encode("internal server error")
		return
	}
	json.NewEncoder(w).Encode(rows)
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
func (p *Handler) Send(w http.ResponseWriter, r *http.Request) {}
