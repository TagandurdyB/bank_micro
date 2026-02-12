package config

import (
	"bank_micro/pkg/database"
	"bank_micro/pkg/rabbitmq"
	accH "bank_micro/services/account/internal/account/handler"
	accR "bank_micro/services/account/internal/account/pgx_repository"
	accS "bank_micro/services/account/internal/account/service"
	userH "bank_micro/services/account/internal/user/handler"
	userR "bank_micro/services/account/internal/user/repository"
	userS "bank_micro/services/account/internal/user/service"

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

	accRepo := accR.NewPgxAccountRepository(pgx)
	accSvc := accS.NewPgxAccountService(accRepo)
	accHnd := accH.NewPgxAccountHandler(accSvc)

	userRepo := userR.NewPgxUserRepository(pgx)
	userSvc := userS.NewPgxUserService(userRepo, accRepo)
	userHnd := userH.NewPgxUserHandler(userSvc)

	r := gin.Default()
	accHnd.RegisterAccountRoutes(r)
	userHnd.RegisterUserRoutes(r)

	return &Container{
		Gin:    r,
		Rabbit: rabbit,
	}, nil
}
