include .env
export

run:
	go run cmd/main.go

dev:
	air

build:
	go build -o bin/clfd cmd/main.go

pg-up:
	podman run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d docker.io/postgres:14

pg-up-docker:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres

down:
	migrate -path ./schema -database $${DATABASE_URL} down

up:
	migrate -path ./schema -database $${DATABASE_URL} up

keygen:
	ssh-keygen -t rsa -P "" -b 2048 -m PEM -f secrets/refresh.key