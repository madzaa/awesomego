.PHONY: up down restart migrate build up-only

# Docker commands
build:
	docker-compose build

up-only: build
	docker-compose up -d
	$(MAKE) check-db

up: build
	docker-compose up -d
	$(MAKE) check-db
	$(MAKE) migrate

down:
	docker-compose down

restart: down up

# Database migrations
migrate:
	@echo "Running migrations..."
	docker-compose exec db psql -U user -d core -f /docker-entrypoint-initdb.d/0001_init.sql || \
	PGPASSWORD=password psql -h localhost -U user -d core -f migrations/0001_init.sql

# Helper to check database connection
check-db:
	@echo "Checking database connection..."
	@while ! docker-compose exec db pg_isready -U user -d core >/dev/null 2>&1; do \
		echo "Waiting for PostgreSQL to start..."; \
		sleep 1; \
	done

# Initialize database and run migrations
init-db: up
