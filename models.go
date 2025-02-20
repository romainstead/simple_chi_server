package main

import "time"

type Transaction struct {
	Id          int       `json:"id" db:"id"`
	FromAddress string    `json:"from_address" db:"from_address"`
	ToAddress   string    `json:"to_address" db:"to_address"`
	Amount      float64   `json:"amount" db:"amount"`
	Time        time.Time `json:"time" db:"time"`
}

type Wallet struct {
	Address string  `json:"wallet_address" db:"wallet_address"`
	Balance float64 `json:"balance" db:"balance"`
}
