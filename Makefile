.PHONY: build sqlc

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

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))
