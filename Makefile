include .env

LOCAL_BIN:=$(CURDIR)/bin
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.51.1
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u google.golang.org/grpc

generate:
	make generate-user-api
	make generate-auth-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/user_v1/user.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 --proto_path vendor.protogen \
	--go_out pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	--validate_out lang=go:pkg/auth_v1 --validate_opt=paths=source_relative \
	api/auth_v1/auth.proto



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


test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/um3ra/auth-microservice/internal/service/...,github.com/um3ra/auth-microservice/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "\n/coverage.out" >> .gitignore



TEST_DB_DSN=host=localhost port=55433 dbname=test user=test password=test sslmode=disable
integration-tests:
	docker compose -f ./db/docker-compose.test.yml build --no-cache
	docker compose -f ./db/docker-compose.test.yml up -d --wait
	sleep 3
	TEST_DB_DSN="$(TEST_DB_DSN)" GRPC_TEST_ADDR="localhost:8889" go test -v ./tests; \
	EXIT_CODE=$$?; \
	docker compose -f ./db/docker-compose.test.yml down --volumes; \
	exit $$EXIT_CODE

lint:
	bin/golangci-lint -c .golangci.workflow.yaml run ./...


vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi