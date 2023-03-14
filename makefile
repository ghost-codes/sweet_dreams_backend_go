migrateup:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

redis:
	docker run --name redis -p 6379:6379 -d redis

.PHONY: sqlc migrateup migrateup1 migratedown migratedown1 server redis