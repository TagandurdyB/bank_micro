package app

import (
	"bank_micro/pkg/rabbitmq"
	"bank_micro/services/transaction/internal/client"
	"bank_micro/services/transaction/internal/handler"
	"bank_micro/services/transaction/internal/repository"
	"bank_micro/services/transaction/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	Worker  *service.TransactionWorker
	Handler *handler.TransactionHandler
}

func InitDependencies(db *gorm.DB, rabbitURL string, accountServiceAddr string) (*Container, error) {
	rabbit, err := rabbitmq.NewRabbitMQClient(rabbitURL, rabbitmq.RabbetQueueDeposit)
	if err != nil {
		return nil, err
	}

	accClient, err := client.NewAccountClient(accountServiceAddr)
	if err != nil {
		return nil, err
	}

	repo := repository.NewTransactionRepository(db)

	hnd := handler.NewTransactionHandler(repo)
	worker := service.NewTransactionWorker(rabbit, repo, accClient)

	return &Container{
		Worker:  worker,
		Handler: hnd,
	}, nil
}
