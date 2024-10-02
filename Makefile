SERVICE_NAME=tiny-bank

start:
	go run main.go

start-docker:
	docker-compose up --build

tests:
	go test ./...

coverage:
	go test --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
