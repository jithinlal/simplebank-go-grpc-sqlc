postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: mock
	go test -v -cover ./...

postgresdown:
	docker stop postgres12

postgresup:
	docker start postgres12

server:
	go run main.go

lint:
	golangci-lint run --fix

setup:
	make postgres
	make createdb
	make migrateup
	make sqlc

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/jithinlal/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test postgresdown postgresup server mock lint setup
