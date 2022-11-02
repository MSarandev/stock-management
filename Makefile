-include .env

.PHONY: init
init:
	@cp .env.dist .env && \
	cp .env.test.dist .env.test

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

.PHONY: migrations-generate
migrations-generate:
	@go run cmd/migrator/main.go generate --name=${name}

.PHONY: migrations-apply
migrations-apply:
	@go run cmd/migrator/main.go migrate

.PHONY: migrations-rollback
migrations-rollback:
	@go run cmd/migrator/main.go rollback

.PHONY: pb-generate
pb-generate:
	@protoc --proto_path=protob "protob/stocks.proto"\
		 --go_out=genprotos --go_opt=paths=source_relative \
         --go-grpc_out=genprotos --go-grpc_opt=paths=source_relative
