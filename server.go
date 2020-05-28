package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	bgg "github.com/chukmunnlee/boardgamegrpc/messages"
)

type BGGServer struct{ Db BoardgameDB }

func (s *BGGServer) FindBoardgamesByName(req *bgg.FindBoardgamesByNameRequest,
	stream bgg.BoardgameService_FindBoardgamesByNameServer) error {

	query := req.GetQuery()
	limit := req.GetLimit()
	offset := req.GetOffset()

	if 0 == limit {
		limit = 10
	}

	log.Printf("FindBoardgamesByName: query=%s, limit=%d, offset=%d", query, limit, offset)

	// read channel
	c := s.Db.FindBoardgamesByName(query, limit, offset)
	count := uint32(0)
	for res := range c {
		if nil != res.Error {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("FindBoardgamesByName: %v", res.Error),
			)
		}

		rec := res.Result.(Boardgame)
		count++
		resp := bgg.FindBoardgamesByNameResponse{
			Count: count,
			Game: &bgg.Boardgame{
				Gid:     rec.Gid,
				Name:    rec.Name,
				Ranking: rec.Ranking,
				Url:     rec.Url,
			},
		}
		stream.Send(&resp)
	}

	return nil
}

func main() {

	// Open a connection to the Boardgame database
	db := BoardgameDB{Username: "fred", Password: "fred"}
	log.Printf("Connecting to the Boardgame database\n")
	checkError(db.Open())
	// Close the database when the program ends
	defer db.Close()

	// Create an instance of the gRPC server
	s := grpc.NewServer()

	// Enable reflection
	log.Printf("Enable reflection\n")
	reflection.Register(s)

	// Create a Boardgame Service
	bggServ := BGGServer{Db: db}

	// Register bggServ as a gRPC service
	bgg.RegisterBoardgameServiceServer(s, &bggServ)

	// Open a port for the service - 50051
	log.Printf("Opening port 50051\n")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	checkError(err)

	log.Printf("Start the Boardgame Service\n")
	checkError(s.Serve(lis))
}
