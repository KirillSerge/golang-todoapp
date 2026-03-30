SHELL := C:/Program Files/Git/bin/bash.exe

include .env
export

export PROJECT_ROOT=$(CURDIR)

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Delete all volume files? This will erase data. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres && \
		rm -rf out/pgdata && \
		echo "Volume files deleted"; \
	else \
		echo "Cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down -d port-forwarder

wait-for-db:
	@echo "Waiting for postgres to be ready..."; \
	until docker exec todoapp-env-postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB} > /dev/null 2>&1; do \
		echo "Postgres is not ready yet, retrying in 2s..."; \
		sleep 2; \
	done; \
	echo "Postgres is ready!"

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing required parameter 'seq'. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make wait-for-db && make migrate-action action=up

migrate-down:
	@make wait-for-db && make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing required parameter action. Example: make migrate-action action=up ,down "; \
		exit 1; \
	fi; \
	MSYS_NO_PATHCONV=1 docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		"$(action)"

