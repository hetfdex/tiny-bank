SERVICE_NAME=tiny-bank

start:
	docker-compose up --build

tests:
	go test ./...

coverage:
	go test --coverprofile=coverage.out ./...
	go tool cover -func=coverage.out