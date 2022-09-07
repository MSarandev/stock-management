-include .env

.PHONY: init
init:
	@cp .env.dist .env

.PHONY: start
start:
	@docker-compose up -d

.PHONY: down
down:
	@docker-compose down --remove-orphans
	@docker-compose ps

.PHONY: down-v
down-v:
	@docker-compose down -v --remove-orphans
	@docker-compose ps

.PHONY: migrate
migrate:
	@go run cmd/migrator/main.go migrate

.PHONY: rollback
rollback:
	@go run cmd/migrator/main.go rollback

.PHONY: serve
serve:
	@go run cmd/server/main.go