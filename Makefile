
migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path ./db/migration -database "postgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable" -verbose down

migratedrop:
	migrate -path ./db/migration -database "postgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable" -verbose drop

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PONY: migratecreate migrateup migratedown migratedrop sqlc test server