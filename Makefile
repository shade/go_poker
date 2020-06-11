PROTO_SRC_DIR=proto
PROTO_TARGET_DIR=$(GOPATH)/src/poker_backend/messages

all:
    go run src/main.go -- -wsport=8081

build:
    protoc $(PROTO_SRC_DIR)/messages.proto --go_out=$(PROTO_TARGET_DIR)