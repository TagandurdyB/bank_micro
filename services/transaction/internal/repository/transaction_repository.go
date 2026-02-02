package repository

import (
	"bank_micro/services/transaction/internal/model"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *TransactionRepository) GetAll(accountID string, toAccountID string) ([]model.Transaction, error) {
	var txs []model.Transaction
	query := r.db

	// 1. Case: If there is only one account ID (both sender and receiver)
	if accountID != "" && toAccountID == "" {
		query = query.Where("account_id = ? OR to_account_id = ?", accountID, accountID)
	}

	// 2. Situation: Money specifically allocated to one person
	if accountID == "" && toAccountID != "" {
		query = query.Where("to_account_id = ?", toAccountID)
	}

	// Case 3: Specific transfers from X to Y
	if accountID != "" && toAccountID != "" {
		query = query.Where("account_id = ? AND to_account_id = ?", accountID, toAccountID)
	}

	err := query.Order("created_at desc").Find(&txs).Error
	return txs, err
}
