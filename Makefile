ifneq (,$(wildcard ./.env))
include .env
export 
ENV_FILE_PARAM = --env-file .env

endif

build:
	docker-compose up --build -d --remove-orphans

up:
	docker-compose up -d

down:
	docker-compose down

show-logs:
	docker-compose logs

mig:
	go generate ./ent

test:
	go test ./tests -v -count=1

ureqm:
	go mod tidy

swag:
	swag init --md .