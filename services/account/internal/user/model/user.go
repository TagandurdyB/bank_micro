package model

import (
	"time"

	acc "bank_micro/services/account/internal/account/model"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	RefUserID *uuid.UUID   `json:"ref_user_id,omitempty" db:"ref_user_id"`
	AccountID uuid.UUID    `json:"account_id" db:"account_id"`
	Account   *acc.Account `json:"account,omitempty" db:"-"`
	IsLocked  bool         `json:"is_locked" db:"is_locked"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	DeletedAt *time.Time   `json:"deleted_at,omitempty" db:"deleted_at"`
}
