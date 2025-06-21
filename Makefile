include .env

# Phony targets
.PHONY: help migrate-up migrate-down migrate-create migrate-fix migrate-status migrate-reset

build:
	@go build -o traverse ./cmd/api

run-api: build
	@./traverse

help:
	@echo "Goose Makefile Targets:"
	@echo "  migrate-create <migration_name>  	- Create a new SQL migration file."
	@echo "  migrate-up                     	- Apply all pending migrations."
	@echo "  migrate-down                   	- Roll back the last applied migration."
	@echo "  migrate-reset                 	- Roll back all migrations (use with caution!)."
	@echo "  migrate-status                 	- Show the status of migrations."

migrate-create:
	@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" create $(filter-out $@,$(MAKECMDGOALS)) sql

migrate-up:
	@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" up

migrate-down:
	@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" down

migrate-status:
	@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" status

migrate-reset:
	@echo "WARNING: This will roll back ALL migrations. Are you sure? (y/N)"
	@read -r CONFIRM; \
		if [ "$$CONFIRM" = "y" ]; then \
		@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" reset
	else \
		echo "Operation to reset migration cancelled."; \
		fi

migrate-fix:
	@goose -dir $(MIGRATION_DIR) $(DB_DRIVER) "$(DSN)" fix
