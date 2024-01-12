DB_PATH=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

start:
	docker start postgres12
	docker start redis
	make server

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migrations -database "$(DB_PATH)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_PATH)" -verbose down

sqlc:
	sqlc generate

test: mock
	go test -v -cover -short ./...

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

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
        proto/*.proto
				statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test postgresdown postgresup server mock lint setup proto evans redis
