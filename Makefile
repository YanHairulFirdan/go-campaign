include .env
export $(shell sed 's/=.*//' .env)
N ?= 1
migrate-up:
	migrate -path=$(MIGRATION_PATH) -database=$(MIGRATION_CONNECTION) up
migrate-down:
	migrate -path=$(MIGRATION_PATH) -database=$(MIGRATION_CONNECTION) down $(N)
migrate-create:
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(NAME)
migrate-force:
	migrate -database=$(MIGRATION_CONNECTION) -path=$(MIGRATION_PATH) force $(VERSION)
run-app:
	go run ./main.go