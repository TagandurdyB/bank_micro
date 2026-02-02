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
RabbitMQ admin panel:
http://localhost:15672

#If needed=======
#5 sudo docker-compose down
