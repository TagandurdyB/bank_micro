package service

import (
	"bank_micro/services/account/internal/account/model"
	repository "bank_micro/services/account/internal/account/pgx_repository"
	"fmt"
)

type PgxAccountService struct {
	repo *repository.PgxAccountRepository
}

func NewPgxAccountService(repo *repository.PgxAccountRepository) *PgxAccountService {
	return &PgxAccountService{repo: repo}
}

func (s *PgxAccountService) CreateAccount(acc *model.Account) (*model.Account, error) {
	err := s.repo.Create(acc)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create account: %w", err)
	}
	return acc, nil
}

func (s *PgxAccountService) GetAccountByID(id string) (*model.Account, error) {
	acc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get account: %w", err)
	}
	return acc, nil
}

func (s *PgxAccountService) GetAllAccounts(balance *int64) ([]model.Account, error) {
	accounts, err := s.repo.GetAll(balance)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get all accounts: %w", err)
	}
	return accounts, nil
}

func (s *PgxAccountService) UpdateAccount(acc *model.Account) (*model.Account, error) {
	updatedAcc, err := s.repo.Update(acc)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update account: %w", err)
	}
	return updatedAcc, nil
}

func (s *PgxAccountService) SoftDeleteAccount(id string) (*model.Account, error) {
	deletedAcc, err := s.repo.Delete(id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to soft delete account: %w", err)
	}
	return deletedAcc, nil
}

func (s *PgxAccountService) UpdateAccountBalance(id string, amount int64) (int64, error) {
	newBalance, err := s.repo.UpdateBalance(id, amount)
	if err != nil {
		return 0, fmt.Errorf("service: failed to update balance: %w", err)
	}
	return newBalance, nil
}
