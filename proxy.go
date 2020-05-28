package main

import (
	ctx "context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"

	gw "github.com/chukmunnlee/boardgamegrpc/messages"
)

func main() {
	// Parse command line
	// --port 3000 default 8080
	port := flag.Uint("port", 8080, "HTTP port")
	flag.Parse()

	// Create a server mux
	mux := runtime.NewServeMux()

	// Options for conneting the the service
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the proxy
	log.Printf("Registering the proxy\n")
	checkError(gw.RegisterBoardgameServiceHandlerFromEndpoint(ctx.Background(), mux, "localhost:50051", opts))

	// Open the 8080 port to allow HTTP traffic
	log.Printf("Opening HTTP connection on port %d\n", *port)
	connStr := fmt.Sprintf("0.0.0.0:%d", *port)
	checkError(http.ListenAndServe(connStr, mux))
}
