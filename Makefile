SERVER_SRC = bggdb.go server.go
GATEWAY_SRC = bggdb.go proxy.go

all:
	@mkdir -p bin
	@go build -o bin/bggserver $(SERVER_SRC)
	@go build -o bin/bgggateway $(GATEWAY_SRC)

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
