include .env

LOCAL_BIN:=$(CURDIR)/bin
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc

generate:
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto


#DB
MIGRATIONS_DIR=$(LOCAL_MIGRATIONS)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) sslmode=disable"

MIGRATION?=init_migration

docker-up:
	docker compose -f ./db/docker-compose.yml up -d
docker-down:
	docker compose -f ./db/docker-compose.yml down

l-migration-status:
	${LOCAL_BIN}/goose -dir ${MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

l-migration-create:
	mkdir -p ${MIGRATIONS_DIR} 
	${LOCAL_BIN}/goose -dir ${MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} create ${MIGRATION} sql


l-migration-up:
	${LOCAL_BIN}/goose -dir ${MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

l-migration-down:
	${LOCAL_BIN}/goose -dir ${MIGRATIONS_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v
