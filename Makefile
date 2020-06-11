PROTO_DIR=internal/proto

all:
	go run src/main.go -- -wsport=8081

build:
	protoc  --proto_path=$(PROTO_DIR) --go_out=$(PROTO_DIR) $(PROTO_DIR)/*.proto

build-docs:
	$(info Please ensure protoc and protoc-gen-doc is installed)

	protoc \
		--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
		--doc_out=./docs \
		--doc_opt=markdown,proto.md \
			./$(PROTO_DIR)/*.proto