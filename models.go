package main

import "time"

// Transaction - сущность из базы данных
// Ну что есть в транзакции? Отправитель, получатель, количество денег и время перевода
// Id - суррогатный ключ, в postgre это PRIMARY KEY с IDENTITY
type Transaction struct {
	Id          int       `json:"id" db:"id"`
	FromAddress string    `json:"from_address" db:"from_address"`
	ToAddress   string    `json:"to_address" db:"to_address"`
	Amount      float64   `json:"amount" db:"amount"`
	Time        time.Time `json:"time" db:"time"`
}

// Wallet Тут наверное понятно всё
type Wallet struct {
	Address string  `json:"wallet_address" db:"wallet_address"`
	Balance float64 `json:"balance" db:"balance"`
}
