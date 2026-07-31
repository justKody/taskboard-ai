.PHONY: build sqlc schema-dump migrate-up migrate-down migrate-force migration clean run

# Prefer Homebrew libpq's pg_dump (often newer than PATH's postgresql@N; Neon may be ahead).
LIBPQ_PG_DUMP := $(firstword $(wildcard /opt/homebrew/opt/libpq/bin/pg_dump /usr/local/opt/libpq/bin/pg_dump))
PG_DUMP ?= $(if $(LIBPQ_PG_DUMP),$(LIBPQ_PG_DUMP),pg_dump)

sqlc:
	sqlc generate

build:
	go build -o build/taskboard-go-api cmd/main.go

run: build
	./build/taskboard-go-api

clean:
	rm -rf build


migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

# Dump the live Postgres schema to db/schema.sql (committed source of truth).
# Read DATABASE_URL via grep (not `source .env`) so unquoted `&` in the URL is safe.
schema-dump:
	@DATABASE_URL=$$(grep -E '^DATABASE_URL=' .env | cut -d= -f2-) && \
		$(PG_DUMP) --schema-only --no-owner --no-privileges --no-comments -x "$$DATABASE_URL" > db/schema.sql
	@echo "wrote db/schema.sql"

migrate-up:
	@go run cmd/migrate/main.go up
	@$(MAKE) schema-dump

migrate-down:
	@go run cmd/migrate/main.go down
	@$(MAKE) schema-dump

migrate-force:
	@go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))
