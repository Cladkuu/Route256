.PHONY: run
run:
	go run cmd/bot/main.go

build:
	go build -o bin/bot cmd/bot/main.go

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
    GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
    GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose


MIGRATIONS_DIR:=$(CURDIR)/migrations
.PHONY: .migration
migration:
	goose -dir=$(MIGRATIONS_DIR) create $(NAME) sql


.PHONY: up
up:
	docker-compose up -d --build

.PHONY: .test
.test:
	$(info Running tests...)
	go test ./.../storage/...

.PHONY: cover
cover:
	go test ./.../storage/... -covermode=count -coverprofile=/tmp/c.out
	go tool cover -html=/tmp/c.out