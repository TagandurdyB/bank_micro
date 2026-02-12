package service

import (
	accModel "bank_micro/services/account/internal/account/model"
	accRepo "bank_micro/services/account/internal/account/pgx_repository"
	"bank_micro/services/account/internal/user/model"
	repository "bank_micro/services/account/internal/user/repository"

	"fmt"
)

type PgxUserService struct {
	repo    *repository.PgxUserRepository
	accRepo *accRepo.PgxAccountRepository
}

func NewPgxUserService(repo *repository.PgxUserRepository, accRepo *accRepo.PgxAccountRepository) *PgxUserService {
	return &PgxUserService{repo: repo, accRepo: accRepo}
}

func (s *PgxUserService) CreateUser(user *model.User) (*model.User, error) {
	if user.Account == nil {
		user.Account = &accModel.Account{
			Balance:  0,
			Currency: "TMT",
			IsLocked: false,
		}
	}

	if err := s.accRepo.Create(user.Account); err != nil {
		return nil, fmt.Errorf("failed to create account for user: %w", err)
	}
	user.AccountID = user.Account.ID

	err := s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("user_service: failed to create user: %w", err)
	}

	return user, nil
}

func (s *PgxUserService) GetUserByID(id string, loadAccount bool) (*model.User, error) {
	user, err := s.repo.GetByID(id, loadAccount)
	if err != nil {
		return nil, fmt.Errorf("user_service: failed to get user: %w", err)
	}
	return user, nil
}

func (s *PgxUserService) GetAllUsers(loadAccount bool) ([]model.User, error) {
	users, err := s.repo.GetAll(loadAccount)
	if err != nil {
		return nil, fmt.Errorf("user_service: failed to get all users: %w", err)
	}
	return users, nil
}

func (s *PgxUserService) UpdateUser(user *model.User) (*model.User, error) {
	if user.Account != nil {
		if _, err := s.accRepo.Update(user.Account); err != nil {
			return nil, fmt.Errorf("failed to update account for user: %w", err)
		}
		user.AccountID = user.Account.ID
	}
	updatedUser, err := s.repo.Update(user)
	if err != nil {
		return nil, fmt.Errorf("user_service: failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (s *PgxUserService) SoftDeleteUser(id string) (*model.User, error) {
	deletedUser, err := s.repo.Delete(id)
	if err != nil {
		return nil, fmt.Errorf("user_service: failed to soft delete user: %w", err)
	}
	return deletedUser, nil
}

func (s *PgxUserService) UpdateUserReference(userID string, refID *string) error {
	if userID == "" {
		return fmt.Errorf("user_service: user id cannot be empty")
	}

	if refID != nil && userID == *refID {
		return fmt.Errorf("user_service: a user cannot be their own referrer")
	}

	err := s.repo.UpdateRef(userID, refID)
	if err != nil {
		return fmt.Errorf("user_service: %w", err)
	}

	return nil
}

func (s *PgxUserService) DepositToUser(userID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("user_service: amount must be greater than zero")
	}

	currentUser, err := s.repo.GetByID(userID, true)
	if err != nil {
		return err
	}
	bonusAmount := amount

	for range 3 {
		if _, err := s.accRepo.UpdateBalance(currentUser.AccountID.String(), bonusAmount); err != nil {
			return fmt.Errorf("failed to create account for user: %w", err)
		}
		if currentUser.RefUserID == nil {
			break
		}

		bonusAmount = bonusAmount / 10
		currentUser, err = s.repo.GetByID(currentUser.RefUserID.String(), true)
		if err != nil {
			return fmt.Errorf("ref user not found: %w", err)
		}
		if currentUser == nil {
			break
		}
	}
	return nil

}
