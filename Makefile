PROTO_SRC_DIR=proto
PROTO_TARGET_DIR=$(GOPATH)/src/poker_backend/messages

all:
	go run src/main.go -- -wsport=8081

build:
	protoc $(PROTO_SRC_DIR)/messages.proto --go_out=$(PROTO_TARGET_DIR)

build-docs:
	$(info Please ensure protoc and protoc-gen-doc is installed, you can find it here: https://github.com/pseudomuto/protoc-gen-doc)

	protoc \
		--plugin=protoc-gen-doc=$(GOPATH)/bin/protoc-gen-doc \
		--doc_out=./docs \
		--doc_opt=markdown,proto.md \
			./internal/proto/*.proto