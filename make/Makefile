
migrateUp:
	migrate -path ../migrations -database "postgres://ts:pass@localhost:5432/bank?sslmode=disable" up

migrateDown:
	migrate -path ../migrations -database "postgres://ts:pass@localhost:5432/bank?sslmode=disable" down

postgresUp:
	docker-compose up -d

postgresDown:
	docker-compose down
