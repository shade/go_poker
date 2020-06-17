GOCMD=go
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

PROTO_DIR=internal/proto
DOC_PORT=19970

all:
	$(GORUN) src/main.go -- -wsport=8081

test:
	$(GOTEST) ./pkg/...  -count=1

build:
	protoc  --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) $(PROTO_DIR)/*.proto

build-docs:
	$(info Please ensure godoc, protoc, and protoc-gen-doc is installed)

	protoc \
		--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
		--doc_out=./docs \
		--doc_opt=markdown,proto.md \
			./$(PROTO_DIR)/*.proto

deps:
	$(GOGET) -d ./...