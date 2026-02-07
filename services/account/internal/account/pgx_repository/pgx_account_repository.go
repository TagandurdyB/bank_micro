package repository

import (
	"bank_micro/pkg/database"
	"bank_micro/services/account/internal/account/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const accountInsert = "balance,currency,is_locked"
const accountReturning = "id,created_at"
const accountSellect = "id,balance,currency,is_locked,created_at,deleted_at"

type PgxAccountRepository struct {
	pgx database.PgxClient
}

func NewPgxAccountRepository(client database.PgxClient) *PgxAccountRepository {
	return &PgxAccountRepository{pgx: client}
}

// 1. Create
func (r *PgxAccountRepository) Create(acc *model.Account) error {
	query := "insert into accounts (" + accountInsert + ") values ($1,$2,$3) returning " + accountReturning
	err := r.pgx.QueryRow(query,
		acc.Balance,
		acc.Currency,
		acc.IsLocked,
	).Scan(&acc.ID, &acc.CreatedAt)

	if err != nil {
		return fmt.Errorf("pgx repo create account: %w", err)
	}

	return nil
}

// 2. GetByID
func (r *PgxAccountRepository) GetByID(id string) (*model.Account, error) {
	var acc model.Account
	query := "select " + accountSellect + " from accounts where id=$1"
	err := r.pgx.QueryRow(query, id).Scan(&acc.ID, &acc.Balance, &acc.Currency, &acc.IsLocked, &acc.CreatedAt, &acc.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	return &acc, nil
}

// 3. GetAll
func (r *PgxAccountRepository) GetAll(balance *int64) ([]model.Account, error) {
	query := "select " + accountSellect + " from accounts where deleted_at is NULL"

	var rows pgx.Rows
	var err error

	if balance != nil {
		query += " AND balance = $1"
		rows, err = r.pgx.Query(query, *balance)
	} else {
		rows, err = r.pgx.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []model.Account

	for rows.Next() {
		var a model.Account
		if err := rows.Scan(
			&a.ID,
			&a.Balance,
			&a.Currency,
			&a.IsLocked,
			&a.CreatedAt,
			&a.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return accounts, nil
}

// 4. Update
func (r *PgxAccountRepository) Update(acc *model.Account) (*model.Account, error) {
	query := "update accounts set balance=$1, currency=$2, is_locked=$3, deleted_at=$4 where id = $5 returning created_at"

	err := r.pgx.QueryRow(
		query,
		acc.Balance,
		acc.Currency,
		acc.IsLocked,
		acc.DeletedAt,
		acc.ID,
	).Scan(&acc.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return acc, nil
}

// 5. Delete (Soft Delete)
func (r *PgxAccountRepository) Delete(id string) (*model.Account, error) {
	query := "update accounts set deleted_at=CURRENT_TIMESTAMP where id=$1 and deleted_at is NULL returning " + accountSellect

	var acc model.Account
	err := r.pgx.QueryRow(query, id).Scan(
		&acc.ID,
		&acc.Balance,
		&acc.Currency,
		&acc.IsLocked,
		&acc.CreatedAt,
		&acc.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found or already deleted")
		}
		return nil, fmt.Errorf("failed to soft delete account: %w", err)
	}

	return &acc, nil
}

// 6. UpdateBalance
func (r *PgxAccountRepository) UpdateBalance(id string, amount int64) (int64, error) {
	query := "update accounts set balance=balance+$1 where id=$2 and deleted_at is NULL returning balance"

	var updatedBalance int64
	err := r.pgx.QueryRow(query, amount, id).Scan(&updatedBalance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("account not found or deleted")
		}
		return 0, fmt.Errorf("failed to update balance: %w", err)
	}

	return updatedBalance, nil
}
