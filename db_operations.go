package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func ConnectToDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:12345@localhost:5433/mydb?sslmode=disable")
	return conn, err
}

func GenerateWallets(conn *pgx.Conn) error {
	rows := []Wallet{}
	for i := 0; i < 10; i++ {
		rows = append(rows, Wallet{uniuri.NewLen(30), 100})
	}
	_, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"wallets"},
		[]string{"wallet_address", "balance"},
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return []any{rows[i].Address, rows[i].Balance}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("couldn't generate wallets: %v", err)
	}
	return nil
}

func GetBalance(address string, conn *pgx.Conn) (Wallet, error) {
	// TODO: CHECK IF ADDRESS IS PRESENT IN DB
	// IF YES, RETURN JSON WITH THAT ADDRESS
	// IF NO, RETURN ERROR 404
	row := conn.QueryRow(context.Background(), "SELECT * FROM wallets w WHERE w.wallet_address = $1", address)
	var wallet Wallet
	err := row.Scan(&wallet.Address, &wallet.Balance)
	if err != nil {
		fmt.Printf("couldn't fetch data from db: %v", err)
	}
	return wallet, nil
}

func GetNLast(n int, conn *pgx.Conn) ([]Transaction, error) {
	rows, err := conn.Query(context.Background(), "SELECT * FROM transactions ORDER BY time DESC FETCH FIRST $1 ROWS ONLY", n)
	if err != nil {
		log.Fatalf("couldn't fetch data from db: %v", err)
	}
	defer rows.Close()
	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err = rows.Scan(&t.ToAddress, &t.FromAddress, &t.Amount, &t.Id, &t.Time)
		if err != nil {
			log.Fatalf("couldn't parse data: %v", err)
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func Send(FromAddress string, ToAddress string, Amount float64, conn *pgx.Conn) error {
	fromWallet, err := GetBalance(FromAddress, conn)
	toWallet, err := GetBalance(ToAddress, conn)
	fmt.Println(fromWallet.Address+": ", fromWallet.Balance)
	fmt.Println(toWallet.Address+": ", toWallet.Balance)
	if err != nil {
		log.Fatalf("couldn't read balance from sender wallet: %v", err)
		return err
	}
	if fromWallet.Balance == 0 || fromWallet.Balance < 0 {
		return errors.New("balance is zero or negative")
	}

	if fromWallet.Balance-Amount < 0 {
		return errors.New("sender wallet has insufficient funds")
	}

	if fromWallet.Address == "" {
		return errors.New("sender wallet address does not exist")
	}

	if toWallet.Address == "" {
		return errors.New("recipient wallet address does not exist")
	}
	tx, err := conn.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return errors.New("database server error")
	}
	tag, err := tx.Exec(context.Background(),
		"UPDATE wallets SET balance = balance - $1 WHERE wallet_address = $2", Amount, fromWallet.Address)
	if tag.RowsAffected() == 0 {
		return errors.New("cant find sender wallet")
	}
	tag, err = tx.Exec(context.Background(),
		"UPDATE wallets SET balance = balance + $1 WHERE wallet_address = $2", Amount, toWallet.Address)
	if err != nil {
		return errors.New("cant find recipient wallet")
	}
	tag, err = tx.Exec(context.Background(),
		"INSERT INTO transactions(to_address, from_address, amount, time) VALUES ($1, $2, $3, $4)", ToAddress, FromAddress, Amount, time.Now())

	if err != nil {
		return errors.New("couldn't insert into 'transactions'")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return errors.New("database server error")
	}
	return err
}
