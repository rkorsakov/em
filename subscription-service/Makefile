.PHONY: build up down gen-swagger

build:
	@docker-compose up --build

up:
	@docker-compose up

down:
	@docker-compose down

gen-swagger:
	@swag init -g cmd/server/main.go