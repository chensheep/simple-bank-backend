DB_URL=vim .mfddfpostgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable

migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose down 1

migratedrop:
	migrate -path ./db/migration -database "$(DB_URL)" -verbose drop

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/chensheep/simple-bank-backend/db/sqlc Store

db_doc:
	dbdocs build doc/db.dbml 

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PONY: migratecreate migrateup migratedown migratedrop migratedown1 migratedrop1 sqlc test server mock db_schema db_doc