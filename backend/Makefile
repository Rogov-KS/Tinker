.PHONY: run build migrate migrate-down seed test clean swagger

run:
	go run ./cmd/server/main.go

build:
	go build -o backend ./cmd/server/main.go

migrate:
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5444/tnews?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgresql://postgres:postgres@localhost:5444/tnews?sslmode=disable" down

seed:
	go run ./seed/seed.go

test:
	go test ./...

swagger:
	swag init -g cmd/server/main.go -o docs

clean:
	rm -f backend seed
