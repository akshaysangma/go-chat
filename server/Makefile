postgresinit:
	docker run --name postgresql15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgresql15 psql

createdb:
	docker exec -it postgresql15 createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgresql15 dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/go-chat?sslmode=disable" -verbose down

build:
	go mod tidy
	go build -o ./bin/server.exe ./cmd

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown build

