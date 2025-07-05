APP_NAME=golang-api-rest

run:
	go run cmd/api/main.go

build:
	go build -o $(APP_NAME).exe cmd/api/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" up

migrate-down:
	migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSLMODE" down

seeds:
	go run cmd/seeds/main.go

seeds-users:
	go run cmd/seeds/main.go -type=users

seeds-projects:
	go run cmd/seeds/main.go -type=projects

seeds-project-items:
	go run cmd/seeds/main.go -type=project-items

seeds-all:
	go run cmd/seeds/main.go -type=all

swag:
	swag init -g cmd/api/main.go 