# =========================================================
# Config
# =========================================================
-include .env
export

COMPOSE      := docker compose -f build/docker-compose.local.yml -p shortlink
ENV_FILE     := .env

DATABASE_URL ?= postgres://shortlink:shortlink@localhost:5432/shortlink?sslmode=disable

# =========================================================
# Build / Run
# =========================================================
.PHONY: build run require-env

build:
	go build -o bin/server ./cmd/server

require-env:
	@test -f $(ENV_FILE) || (echo "Missing $(ENV_FILE) — copy .env.example to .env first"; exit 1)

run: require-env
	@set -a; . ./$(ENV_FILE); set +a; go run ./cmd/server

# =========================================================
# Test
# =========================================================
.PHONY: test test-integration

test:
	go test ./... -race -count=1

test-integration: migrate
	DATABASE_URL="$(DATABASE_URL)" go test -tags=integration ./...

# =========================================================
# Stack lifecycle
# =========================================================
.PHONY: up down logs ps db

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

db:
	$(COMPOSE) up -d postgres

# =========================================================
# Migration / Codegen
# =========================================================
.PHONY: migrate migration sqlc

migrate:
	$(COMPOSE) run --rm migrate

migration:
	@if [ -z "$(name)" ]; then echo "Usage: make migration name=<name>"; exit 1; fi
	@last=$$(ls data/migrations/ 2>/dev/null | grep -oE '^[0-9]+' | sort -n | tail -1); \
	next=$$(printf "%06d" $$(( 10#$${last:-0} + 1 ))); \
	touch data/migrations/$${next}_$(name).up.sql data/migrations/$${next}_$(name).down.sql; \
	echo "Created data/migrations/$${next}_$(name).{up,down}.sql"

sqlc:
	sqlc generate
