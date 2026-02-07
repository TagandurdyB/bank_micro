package repository

import (
	"bank_micro/pkg/database"
	// accModel "bank_micro/services/account/internal/account/model"
	accRepo "bank_micro/services/account/internal/account/pgx_repository"
	"bank_micro/services/account/internal/user/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const userInsert = "account_id,ref_user_id,is_locked"
const userReturning = "id,created_at"
const userSelect = "id,account_id,ref_user_id,is_locked,created_at,deleted_at"

type PgxUserRepository struct {
	pgx     database.PgxClient
	accRepo *accRepo.PgxAccountRepository
}

func NewPgxUserRepository(client database.PgxClient, accRepo *accRepo.PgxAccountRepository) *PgxUserRepository {
	return &PgxUserRepository{pgx: client, accRepo: accRepo}
}

func (r *PgxUserRepository) Create(user *model.User) error {
	if user.Account != nil {
		if err := r.accRepo.Create(user.Account); err != nil {
			return fmt.Errorf("failed to create account for user: %w", err)
		}
		user.AccountID = user.Account.ID
	}

	query := "insert into users (" + userInsert + ") values ($1,$2,$3) returning " + userReturning
	err := r.pgx.QueryRow(query, user.AccountID, user.RefUserID, user.IsLocked).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("pgx repo create user: %w", err)
	}
	return nil
}

func (r *PgxUserRepository) GetByID(id string, loadAccount bool) (*model.User, error) {
	var user model.User
	query := "select " + userSelect + " from users where id=$1"
	if err := r.pgx.QueryRow(query, id).Scan(&user.ID, &user.AccountID, &user.RefUserID, &user.IsLocked, &user.CreatedAt, &user.DeletedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if loadAccount {
		acc, err := r.accRepo.GetByID(user.AccountID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to load account: %w", err)
		}
		user.Account = acc
	}

	return &user, nil
}

// GetAll optional join ile account
func (r *PgxUserRepository) GetAll(loadAccount bool) ([]model.User, error) {
	query := "select " + userSelect + " from users where deleted_at is NULL"
	rows, err := r.pgx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.AccountID, &u.RefUserID, &u.IsLocked, &u.CreatedAt, &u.DeletedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if loadAccount {
			acc, err := r.accRepo.GetByID(u.AccountID.String())
			if err != nil {
				return nil, fmt.Errorf("failed to load account for user: %w", err)
			}
			u.Account = acc
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return users, nil
}

// Update user + optional account update
func (r *PgxUserRepository) Update(u *model.User) (*model.User, error) {
	// Eğer JSON payload içinden account geldi ise önce update et
	if u.Account != nil {
		if _, err := r.accRepo.Update(u.Account); err != nil {
			return nil, fmt.Errorf("failed to update account for user: %w", err)
		}
		u.AccountID = u.Account.ID
	}

	query := "update users set account_id=$1, ref_user_id=$2, is_locked=$3, deleted_at=$4 where id=$5 returning id,created_at"
	err := r.pgx.QueryRow(query, u.AccountID, u.RefUserID, u.IsLocked, u.DeletedAt, u.ID).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u, nil
}

// Soft Delete
func (r *PgxUserRepository) Delete(id string) (*model.User, error) {
	query := "update users set deleted_at=CURRENT_TIMESTAMP where id=$1 and deleted_at is NULL returning " + userSelect
	var u model.User
	err := r.pgx.QueryRow(query, id).Scan(&u.ID, &u.AccountID, &u.RefUserID, &u.IsLocked, &u.CreatedAt, &u.DeletedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found or already deleted")
		}
		return nil, fmt.Errorf("failed to soft delete user: %w", err)
	}
	return &u, nil
}

// Deposit + bonus
func (r *PgxUserRepository) Deposit(userID string, amount int64) error {
	currentUser, err := r.GetByID(userID, false)
	if err != nil {
		return err
	}

	bonusAmount := amount
	for i := 0; i < 3; i++ { // max 3 level
		if currentUser == nil {
			break
		}

		var updatedBalance int64
		query := "update accounts set balance = balance + $1 where id = $2 returning balance"
		err := r.pgx.QueryRow(query, bonusAmount, currentUser.AccountID).Scan(&updatedBalance)
		if err != nil {
			return fmt.Errorf("failed to deposit to account: %w", err)
		}

		bonusAmount = bonusAmount / 10
		if currentUser.RefUserID == nil {
			break
		}

		currentUser, err = r.GetByID(currentUser.RefUserID.String(), false)
		if err != nil {
			return fmt.Errorf("ref user not found: %w", err)
		}
	}

	return nil
}
