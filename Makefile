include .env
export $(shell sed 's/=.*//' .env)
migrate-up:
	migrate -path=$(MIGRATION_PATH) -database=$(MIGRATION_CONNECTION) up
migrate-down:
	migrate -path=$(MIGRATION_PATH) -database=$(MIGRATION_CONNECTION) down
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)
run-app:
	go run ./cmd/api/main.go