include .env.example
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)

docker:
	docker-compose up
.PHONY: docker

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 ### swag run
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run -tags migrate ./cmd/app
.PHONY: run

migrate-create:  ### create new migration
	migrate create -ext sql -dir migrations 'migrate_name'
.PHONY: migrate-create

migrate-up: ### migration up
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
.PHONY: migrate-up

bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
