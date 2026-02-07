package model

import (
	"time"
	"github.com/google/uuid"
)


type Account struct {
	ID        uuid.UUID  `json:"id"`
	Balance   int64      `json:"balance"`
	Currency  string     `json:"currency"`
	IsLocked  bool       `json:"is_locked"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}