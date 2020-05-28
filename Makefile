SERVER_SRC = bggdb.go server.go
GATEWAY_SRC = bggdb.go gateway.go
CLIENT_SRC = bggdb.go client.go

all: generate
	@mkdir -p bin
	@go build -o bin/bggserver $(SERVER_SRC)
	@go build -o bin/bgggateway $(GATEWAY_SRC)
	@go build -o bin/bggclient $(CLIENT_SRC)

run_client:
	@clear
	@go run $(CLIENT_SRC)

run_gateway:
	@clear
	@go run $(GATEWAY_SRC)

run_server:
	@clear
	@go run $(SERVER_SRC)

generate:
	@protoc messages/bgg.proto \
		-I. -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:.
