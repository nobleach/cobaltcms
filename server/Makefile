ifneq (,$(wildcard ./.env))
	include .env
	export
endif

export PORT=:8080

migrate_up:
	migrate -database ${DATABASE_URL} -path db/migrations up

migrate_down:
	migrate -database ${DATABASE_URL} -path db/migrations down

build:
	sqlc generate
	go build -o bin/server ./cmd/server/...

run: build
	PORT=$(PORT) ./bin/server

watch:
	./bin/air
