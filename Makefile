include .env

run:
	go run ./cmd/api

sql-gen:
	sqlc generate

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

migrate-up:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down 1

migrate-down-all:
	migrate -path ./migrations/${module} -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

pre-commit:
	pre-commit run --all-files
