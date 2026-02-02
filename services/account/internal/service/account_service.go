package service

import (
	"bank_micro/pkg/rabbitmq"
	"bank_micro/services/account/internal/model"
	"bank_micro/services/account/internal/repository"
)

type AccountService struct {
	repo   *repository.AccountRepository
	rabbit *rabbitmq.RabbitMQClient
}

func NewAccountService(repo *repository.AccountRepository, rabbit *rabbitmq.RabbitMQClient) *AccountService {
	return &AccountService{repo: repo, rabbit: rabbit}
}

// 1. Create
func (s *AccountService) Create(currency string, initialBalance int64) (*model.Account, error) {
	acc := &model.Account{
		Currency: currency,
		Balance:  initialBalance,
	}
	return acc, s.repo.Create(acc)
}

// 2. GetByID
func (s *AccountService) GetByID(id string) (*model.Account, error) {
	return s.repo.GetByID(id)
}

// 3. GetAll
func (s *AccountService) GetAll() ([]model.Account, error) {
	return s.repo.GetAll()
}

// 4. Update
func (s *AccountService) Update(id string, balance int64, isLocked bool) (*model.Account, error) {
	acc, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	acc.Balance = balance
	acc.IsLocked = isLocked

	if err := s.repo.Update(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

// 5. Delete (Soft Delete)
func (s *AccountService) Delete(id string) error {
	return s.repo.Delete(id)
}

// 6. ProcessDeposit
func (s *AccountService) ProcessDeposit(id string, amount int64) (string, error) {
	// 1. Account control
	_, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	// 2. Prepare the message content.
	event := map[string]interface{}{
		"account_id": id,
		"amount":     amount,
		"type":       "DEPOSIT",
	}

	// 3. Send to RabbitMQ
	err = s.rabbit.Publish("deposit_queue", event)
	if err != nil {
		return "", err
	}

	return "Your request has been added to the queue!", nil
}
