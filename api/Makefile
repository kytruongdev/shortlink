# =========================================================
# Config
# =========================================================
-include .env
export

ENV          ?= local
COMPOSE      := docker compose -f build/docker-compose.$(ENV).yml -p shortlink $(if $(filter prod,$(ENV)),--env-file .env.prod,)
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
.PHONY: test

test: migrate
	PG_URL="$(DATABASE_URL)" go test ./... -race -count=1

# =========================================================
# Stack lifecycle
# =========================================================
# Local: make up | make down | make reset (wipes data)
# Prod : make up ENV=prod | make logs ENV=prod | make down ENV=prod (keeps data)
.PHONY: up down reset logs ps db

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

reset:
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
.PHONY: migrate migration sqlc swag

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

swag:
	swag init -g cmd/server/main.go -o internal/docs
