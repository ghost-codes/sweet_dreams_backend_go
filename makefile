migrateup:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams?sslmode=disable" -verbose down
	
migratedown1:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/sweet_dreams?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

.PHONY: sqlc migrateup migrateup1 migratedown migratedown1