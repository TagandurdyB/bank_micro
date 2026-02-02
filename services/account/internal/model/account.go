package model

import (
	"time"
	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Balance   int64     `gorm:"not null;default:0"`
	Currency  string    `gorm:"size:3;not null"`
	IsLocked  bool      `gorm:"default:false"`
	CreatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}