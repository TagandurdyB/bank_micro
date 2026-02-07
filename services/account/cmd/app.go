package main

import (
	"bank_micro/pkg/database"
	"bank_micro/pkg/rabbitmq"
	"bank_micro/services/account/internal/account/handler"
	repository "bank_micro/services/account/internal/account/pgx_repository"
	"bank_micro/services/account/internal/account/service"

	"github.com/gin-gonic/gin"
)

type Container struct {
	Gin    *gin.Engine
	Rabbit *rabbitmq.RabbitMQClient
}

func InitDependencies(pgx database.PgxClient, rabbitURL string) (*Container, error) {
	// Launch RabbitMQ
	rabbit, err := rabbitmq.NewRabbitMQClient(rabbitURL, rabbitmq.RabbetQueueDeposit)
	if err != nil {
		return nil, err
	}

	repo := repository.NewPgxAccountRepository(pgx)
	svc := service.NewPgxAccountService(repo)
	hnd := handler.NewPgxAccountHandler(svc)

	r := gin.Default()
	hnd.RegisterAccountRoutes(r)

	return &Container{
		Gin:    r,
		Rabbit: rabbit,
	}, nil
}
