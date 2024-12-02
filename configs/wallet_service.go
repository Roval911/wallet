package services

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/google/uuid"
	"wallet-service/database"
)

var mu sync.Mutex // Глобальная блокировка для работы с кошельками

func GetWalletBalance(walletId uuid.UUID) (int64, error) {
	var balance int64
	query := "SELECT balance FROM wallets WHERE id = $1"

	err := database.DB.QueryRow(query, walletId).Scan(&balance)
	if err == sql.ErrNoRows {
		return 0, errors.New("wallet not found")
	} else if err != nil {
		return 0, err
	}

	return balance, nil
}

func UpdateWallet(walletId uuid.UUID, operationType string, amount int64) error {
	mu.Lock()
	defer mu.Unlock()

	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	var currentBalance int64
	query := "SELECT balance FROM wallets WHERE id = $1 FOR UPDATE"
	err = tx.QueryRow(query, walletId).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return errors.New("wallet not found")
	} else if err != nil {
		tx.Rollback()
		return err
	}

	var newBalance int64
	if operationType == "DEPOSIT" {
		newBalance = currentBalance + amount
	} else if operationType == "WITHDRAW" {
		if currentBalance < amount {
			tx.Rollback()
			return errors.New("insufficient funds")
		}
		newBalance = currentBalance - amount
	} else {
		tx.Rollback()
		return errors.New("invalid operation type")
	}

	updateQuery := "UPDATE wallets SET balance = $1 WHERE id = $2"
	_, err = tx.Exec(updateQuery, newBalance, walletId)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}