.PHONY: run build docker-up docker-down start stop
run:
	go run cmd/api/main.go

build:
	go build -o bin/main cmd/api/main.go

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

getdocs:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go -o docs

start: docker-up 
stop: docker-down 
