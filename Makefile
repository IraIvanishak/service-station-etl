DB_URL ?= postgres://postgres:postgres@localhost:5452/service_station?sslmode=disable
# Maintenance database on the same server; integration tests clone a template from it.
TEST_DB_URL ?= postgres://postgres:postgres@localhost:5452/postgres?sslmode=disable

.PHONY: db-up db-reset db-psql api test test-integration sqlc

db-up: ## Start Postgres; the first start creates schemas, seed data and ETL procedures
	docker compose up -d --wait

db-reset: ## Drop the volume and rebuild the database from db/init
	docker compose down -v
	docker compose up -d --wait

db-psql: ## Open psql in the database container
	docker compose exec db psql -U postgres -d service_station

api: ## Run the API on :3000
	DATABASE_URL="$(DB_URL)" go run .

test: ## Unit tests
	go test ./...

test-integration: ## Unit + integration tests (needs `make db-up`)
	TEST_DATABASE_URL="$(TEST_DB_URL)" go test -count=1 ./...

sqlc: ## Regenerate storage/dbs from storage/query
	docker run --rm -v "$(CURDIR)":/src -w /src sqlc/sqlc:1.26.0 generate
