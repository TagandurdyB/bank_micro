package repository

import (
	"bank_micro/pkg/database"
	accModel "bank_micro/services/account/internal/account/model"
	"bank_micro/services/account/internal/user/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

const userInsert = "account_id,ref_user_id,is_locked"
const userReturning = "id,created_at"
const userSelect = "u.id,u.account_id,u.ref_user_id,u.is_locked,u.created_at,u.deleted_at"
const accountSelect = "a.id, a.balance, a.currency, a.is_locked, a.created_at, a.deleted_at"

type PgxUserRepository struct {
	pgx     database.PgxClient
}

func NewPgxUserRepository(client database.PgxClient) *PgxUserRepository {
	return &PgxUserRepository{pgx: client}
}

func (r *PgxUserRepository) Create(user *model.User) error {
	query := "insert into users (" + userInsert + ") values ($1,$2,$3) returning " + userReturning
	err := r.pgx.QueryRow(query, user.AccountID, user.RefUserID, user.IsLocked).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return fmt.Errorf("pgx repo create user: %w", err)
	}
	return nil
}

func (r *PgxUserRepository) GetByID(id string, loadAccount bool) (*model.User, error) {
	var user model.User

	query := "select " + userSelect
	if loadAccount {
		query += ", " + accountSelect + " from users u left join accounts a on u.account_id = a.id where u.id = $1"
		user.Account = &accModel.Account{}
		err := r.pgx.QueryRow(query, id).Scan(
			&user.ID, &user.AccountID, &user.RefUserID, &user.IsLocked, &user.CreatedAt, &user.DeletedAt,
			&user.Account.ID, &user.Account.Balance, &user.Account.Currency, &user.Account.IsLocked, &user.Account.CreatedAt, &user.Account.DeletedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("user not found: %w", err)
			}
			return nil, fmt.Errorf("failed to query user with account: %w", err)
		}
		return &user, nil
	} else {
		query += " from users u where u.id = $1"
		if err := r.pgx.QueryRow(query, id).Scan(&user.ID, &user.AccountID, &user.RefUserID, &user.IsLocked, &user.CreatedAt, &user.DeletedAt); err != nil {
			if err == pgx.ErrNoRows {
				return nil, fmt.Errorf("user not found: %w", err)
			}
			return nil, fmt.Errorf("failed to query user: %w", err)
		}

		return &user, nil
	}

}

func (r *PgxUserRepository) GetAll(loadAccount bool) ([]model.User, error) {
	var users []model.User
	query := "select " + userSelect

	if loadAccount {
		query += ", " + accountSelect + " from users u left join accounts a on u.account_id = a.id"
	} else {
		query += " from users u"
	}
	query += " where u.deleted_at is NULL"

	rows, err := r.pgx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u model.User

		if loadAccount {
			u.Account = &accModel.Account{}
			err := rows.Scan(
				&u.ID, &u.AccountID, &u.RefUserID, &u.IsLocked, &u.CreatedAt, &u.DeletedAt,
				&u.Account.ID, &u.Account.Balance, &u.Account.Currency, &u.Account.IsLocked, &u.Account.CreatedAt, &u.Account.DeletedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan user and account: %w", err)
			}
		} else {
			err := rows.Scan(&u.ID, &u.AccountID, &u.RefUserID, &u.IsLocked, &u.CreatedAt, &u.DeletedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to scan user: %w", err)
			}
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *PgxUserRepository) Update(u *model.User) (*model.User, error) {
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

func (r *PgxUserRepository) Delete(id string) (*model.User, error) {
	query := "update users set deleted_at=CURRENT_TIMESTAMP where id=$1 and deleted_at is NULL returning id, account_id, ref_user_id, is_locked, created_at, deleted_at"
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

func (r *PgxUserRepository) UpdateRef(userID string, refID *string) error {
	if refID != nil && userID == *refID {
		return fmt.Errorf("user cannot reference themselves!")
	}

	query := "update users set ref_user_id = $1 where id = $2 and deleted_at is NULL"

	result, err := r.pgx.Exec(query, refID, userID)
	if err != nil {
		return fmt.Errorf("failed to update user reference: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found or already deleted")
	}

	return nil
}

// func (r *PgxUserRepository) Deposit(userID string, amount int64) error {
// 	currentUser, err := r.GetByID(userID, true)
// 	if err != nil {
// 		return err
// 	}

// 	bonusAmount := amount
// 	for range 2 {
// 		if _, err := r.accRepo.UpdateBalance(currentUser.AccountID.String(), amount); err != nil {
// 			return fmt.Errorf("failed to create account for user: %w", err)
// 		}

// 		bonusAmount = bonusAmount / 10
// 		if currentUser.RefUserID == nil {
// 			break
// 		}

// 		currentUser, err = r.GetByID(currentUser.RefUserID.String(), true)
// 		if err != nil {
// 			return fmt.Errorf("ref user not found: %w", err)
// 		}
// 		if currentUser == nil {
// 			break
// 		}

// 	}

// 	return nil
// }
