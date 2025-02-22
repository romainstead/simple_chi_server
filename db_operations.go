package main

import (
	"github.com/dchest/uniuri"
	"gorm.io/gorm"
)

type DbConnection struct {
	database *gorm.DB
}

func GenerateWallets(db *gorm.DB) error {
	var rows []*Wallet
	for i := 0; i < 10; i++ {
		rows = append(rows, &Wallet{uniuri.NewLen(30), 1000})
	}
	result := db.Create(rows)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetWallet(address string, db *gorm.DB) (Wallet, error) {
	WalToFind := Wallet{address, 0}
	wal := db.First(&WalToFind, "wallet_address = ?", address)
	if wal.Error != nil {
		return Wallet{"", 0}, wal.Error
	}
	return WalToFind, nil
}

func GetNLast(n int, db *gorm.DB) ([]*Transaction, error) {
	var rows []*Transaction
	db.Order("time desc").Limit(n).Find(&rows)
	if db.Error != nil {
		return nil, db.Error
	}
	return rows, nil
}

func Send(FromAddress string, ToAddress string, Amount float64) {}
