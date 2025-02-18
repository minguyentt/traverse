include .env

# For local development
LOCAL_DSN=user=$(LOCAL_DB_USER) password=$(LOCAL_DB_PASSWORD) dbname=$(LOCAL_DB_NAME) host=$(LOCAL_DB_HOST) port=$(LOCAL_DB_PORT) sslmode=disable

MIGRATIONS_DIR = $(MIGRATION_DIR)

# Phony targets
.PHONY: migrate-up migrate-down migrate-create migrate-fix migrate-status migrate-reset

run:
	@go run ./cmd/api

migrate-create:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" create $(filter-out $@,$(MAKECMDGOALS)) sql

migrate-up:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" up

migrate-down:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" down

migrate-status:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" status

migrate-reset:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" reset

migrate-fix:
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(LOCAL_DSN)" fix

