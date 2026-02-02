File System:
tbm@ubuntu24  /media/tbm/Storage/PRJS/golang/bank_micro   master ±  tree ./  
./
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
├── docker
│ ├── account.Dockerfile
│ ├── api-gateway.Dockerfile
│ ├── builder.Dockerfile
│ ├── proto.Dockerfile
│ └── transaction.Dockerfile
├── docker-compose.yml
├── gateway
│ ├── cmd
│ │ └── main.go
│ ├── config
│ │ └── config.go
│ └── internal
│ └── app
│ └── app.go
├── go.mod
├── go.sum
├── init-db.sql
├── pkg
│ ├── config
│ │ └── config.go
│ ├── database
│ │ └── postgres.go
│ ├── rabbitmq
│ │ └── client.go
│ └── redis
│ └── client.go
├── proto
│ ├── account_service.proto
│ ├── gen
│ │ ├── account_service_grpc.pb.go
│ │ ├── account_service.pb.go
│ │ ├── account_service.pb.gw.go
│ │ ├── transaction_service_grpc.pb.go
│ │ ├── transaction_service.pb.go
│ │ └── transaction_service.pb.gw.go
│ └── transaction_service.proto
├── README.md
├── screenshots
│ ├── account_create.png
│ ├── account_deposit.png
│ ├── account_read_all.png
│ ├── account_read.png
│ └── transaction_read.png
└── services
├── account
│ ├── cmd
│ │ └── main.go
│ ├── config
│ │ └── config.go
│ └── internal
│ ├── app
│ │ ├── app.go
│ │ └── grpc.go
│ ├── handler
│ │ └── account_service_handler.go
│ ├── model
│ │ └── account.go
│ ├── repository
│ │ └── account_repository.go
│ └── service
│ └── account_service.go
└── transaction
├── cmd
│ └── main.go
├── config
│ └── config.go
└── internal
├── app
│ ├── app.go
│ └── grpc.go
├── client
│ └── account_client.go
├── handler
│ └── transaction_handler.go
├── model
│ └── transaction.go
├── repository
│ └── transaction_repository.go
└── service
└── transaction_worker.go


#1. change this part in .env "LOCAL_VOLUME_PATH" related full path fore your machine!:
/home/tbm/docker-composes

#2 sudo mkdir -p ~/docker-composes/bank_micro/{postgres_data}
#3 sudo cp ./init-db.sql ~/docker-composes/bank_micro/init-db.sql

👉===============================
//Build images
//RUN only one time to create "golang-bufbuild" docker image:
sudo docker compose --profile tools build
//Create "go_micro_builder" docker image fore micro services's build environment:
sudo docker compose --profile builder build

//Tools
//When you need generate golang proto files:
sudo docker compose run --rm proto-gen dep update
sudo docker compose run --rm proto-gen

//Up services=============
//Init all core services (postgres, redis, rabbetMQ, ....):
sudo docker compose --profile infra up -d
//Build and up all micro services at once:
sudo docker compose --profile runtime up -d

Now you good to go!
👉===============================

REST api port :9080
you can see endpoints, requests & responses inside of ./proto/account_service.proto and ./proto/transaction_service.proto files!

If you use "sudo docker compose --profile runtime up -d" and run all services in docker containers you cannot acces gRPC endpoints.
Because services use docker's internal network inside of containers!
but you can run run services in your local if you have all go packages and go version>=1.25.6:

#Gateway service
go run gateway/cmd/main.go
#Account service
go run services/account/cmd/main.go
#Transaction service
go run services/transaction/cmd/main.go

you need run all 3 of then in different terminals if you want test in local!

you can see request tests result in .png fonmat at ./screenshots:
./screenshots
├── account_create.png
├── account_deposit.png
├── account_read_all.png
├── account_read.png
└── transaction_read.png

RabbitMQ admin panel:
http://localhost:15672

#If needed=======
#5 sudo docker-compose down

fore down spesific prfile:
sudo docker compose --profile runtime down
#Delete all created image from 'runtime' profile
sudo docker compose --profile runtime down --rmi all
