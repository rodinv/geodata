build:
	go build -o bin/geodata ./cmd/geodata

mod:
	go mod download

swag:
	swag init --parseDependency -g cmd/geodata/main.go
