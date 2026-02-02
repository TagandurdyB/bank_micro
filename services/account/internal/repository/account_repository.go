package repository

import (
	"bank_micro/services/account/internal/model"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// 1. Create
func (r *AccountRepository) Create(acc *model.Account) error {
	return r.db.Create(acc).Error
}

// 2. GetByID
func (r *AccountRepository) GetByID(id string) (*model.Account, error) {
	var acc model.Account
	err := r.db.First(&acc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// 3. GetAll
func (r *AccountRepository) GetAll() ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Find(&accounts).Error
	return accounts, err
}

// 4. Update
func (r *AccountRepository) Update(acc *model.Account) error {
	return r.db.Save(acc).Error
}

// 5. Delete (Soft Delete)
func (r *AccountRepository) Delete(id string) error {
	return r.db.Model(&model.Account{}).Delete("id = ?", id).Error
}

// 6. UpdateBalance
func (r *AccountRepository) UpdateBalance(id string, newBalance int64) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("balance", newBalance).Error
}
