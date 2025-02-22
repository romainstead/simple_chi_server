package main

import (
	"time"
)

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	FromAddress string    `json:"from_address" gorm:"column:from_address"`
	ToAddress   string    `json:"to_address" gorm:"column:to_address"`
	Amount      float64   `json:"amount" gorm:"column:amount"`
	Time        time.Time `json:"time" gorm:"column:time"`
}

type Wallet struct {
	Address string  `json:"wallet_address" gorm:"column:wallet_address"`
	Balance float64 `json:"balance" gorm:"column:balance"`
}
