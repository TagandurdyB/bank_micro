package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TransactionTypeDeposit  = "DEPOSIT"
	TransactionTypeTransfer = "TRANSFER"
	TransactionTypeWithdraw = "WITHDRAW"
)

type Transaction struct {
	ID              uuid.UUID `json:"id" db:"id"`
	AccountID       uuid.UUID `json:"account_id" db:"account_id"`
	ToAccountID     *uuid.UUID `json:"to_account_id,omitempty" db:"to_account_id"` 
	Amount          int64     `json:"amount" db:"amount"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}