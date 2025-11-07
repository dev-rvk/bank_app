postgres:
	docker run --name bankapp -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 -d postgres

postgres_start:
	docker start bankapp

createdb:
	docker exec -it bankapp createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it bankapp dropdb --username=postgres simple_bank --force

migrateup:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/mock.go github.com/devrvk/simplebank/db/sqlc Store 

docker_build:
	docker build -t simplebank:latest .

docker_run:
	docker run --network bank-network --name simplebank -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE="postgres://postgres:postgres@bankapp:5432/simple_bank?sslmode=disable"  simplebank:latest

.PHONY: postgres postgres_start createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock