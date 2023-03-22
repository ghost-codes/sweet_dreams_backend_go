migrateup:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams_go?sslmode=disable" -verbose down 1

swag:
	swag init
sqlc: 
	sqlc generate

server: swag
	go run main.go

redis:
	docker run --name redis -p 6379:6379 -d redis

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)


.PHONY: sqlc migrateup migrateup1 migratedown migratedown1 server redis new_migration swag