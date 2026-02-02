package service

import (
	"bank_micro/pkg/rabbitmq"
	proto "bank_micro/proto/gen"
	"bank_micro/services/transaction/internal/client"
	"bank_micro/services/transaction/internal/model"
	"bank_micro/services/transaction/internal/repository"
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

type TransactionWorker struct {
	rabbit    *rabbitmq.RabbitMQClient
	repo      *repository.TransactionRepository
	accClient *client.AccountClient
}

func NewTransactionWorker(r *rabbitmq.RabbitMQClient, repo *repository.TransactionRepository, acc *client.AccountClient) *TransactionWorker {
	return &TransactionWorker{rabbit: r, repo: repo, accClient: acc}
}

func (w *TransactionWorker) Start() {
	msgs, err := w.rabbit.Channel.Consume(
		"deposit_queue", "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			var event struct {
				AccountID string `json:"account_id"`
				Amount    int64  `json:"amount"`
			}
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling deposit event: %v", err)
				continue
			}

			w.handleDeposit(event.AccountID, event.Amount)
		}
	}()
}

func (w *TransactionWorker) handleDeposit(accID string, amount int64) {
	ctx := context.Background()

	// 1. Get current account via gRPC
	acc, err := w.accClient.Client.GetAccount(ctx, &proto.GetAccountRequest{Id: accID})
	if err != nil {
		log.Printf("Error getting account: %v", err)
		return
	}

	// 2. Update Balance via gRPC
	newBalance := acc.Balance + amount
	_, err = w.accClient.Client.UpdateAccount(ctx, &proto.UpdateAccountRequest{
		Id:       accID,
		Balance:  newBalance,
		IsLocked: acc.IsLocked,
	})
	if err != nil {
		log.Printf("Error updating balance: %v", err)
		return
	}

	// 3. Save Transaction to DB
	uID, _ := uuid.Parse(accID)

	err = w.repo.Create(&model.Transaction{
		AccountID:       uID,
		Amount:          amount,
		TransactionType: model.TransactionTypeDeposit,
	})
	if err != nil {
		log.Printf("Critical: Failed to save transaction to DB: %v", err)
		// Note: There should be a 'rollback' mechanism or log alert here
		// because the balance was updated but the record could not be kept!
		return
	}

	log.Printf("Successfully processed deposit for %s", accID)
}
