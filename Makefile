include .env
export

run:
	go run cmd/main.go

build:
	go build -o bin/clfd cmd/main.go

pg-up:
	podman run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d docker.io/postgres:14

dev:
	air

down:
	migrate -path ./schema -database $${DATABASE_URL} down

up:
	migrate -path ./schema -database $${DATABASE_URL} up

fix-version:
	migrate -path schema/ -database $${DATABASE_URL} force 1