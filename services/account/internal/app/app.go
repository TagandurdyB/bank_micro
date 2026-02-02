package app

import (
	"bank_micro/pkg/rabbitmq"
	"bank_micro/services/account/internal/handler"
	"bank_micro/services/account/internal/repository"
	"bank_micro/services/account/internal/service"
	"gorm.io/gorm"
)

type Container struct {
	Handler *handler.AccountHandler
	Rabbit  *rabbitmq.RabbitMQClient
}

func InitDependencies(db *gorm.DB, rabbitURL string) (*Container, error) {
	// Launch RabbitMQ
	rabbit, err := rabbitmq.NewRabbitMQClient(rabbitURL)
	if err != nil {
		return nil, err
	}

	repo := repository.NewAccountRepository(db)
	svc := service.NewAccountService(repo, rabbit) 
	hnd := handler.NewAccountHandler(svc)

	return &Container{
		Handler: hnd,
		Rabbit:  rabbit,
	}, nil
}