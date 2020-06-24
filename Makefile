GOCMD=go
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

PROTO_DIR=internal/proto
DOC_PORT=19970

.PHONY: all test clean

all:
test:
	$(GOTEST) ./test/e2e/id_test.go

build:
	protoc  --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) $(PROTO_DIR)/*.proto

build-docs:
	$(info Please ensure godoc, protoc, and protoc-gen-doc are installed)

	protoc \
		--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
		--proto_path=$(PROTO_DIR)\
		--doc_out=./docs \
		--doc_opt=markdown,proto.md \
			./$(PROTO_DIR)/*.proto

deps:
	$(GOGET) -d ./...


room_server:
	go run ./cmd/room_server.go
