package main

import (
	chess "github.com/cbotte21/chess-go/pb"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/enviroment"
	"github.com/cbotte21/queue-go/internal"
	"github.com/cbotte21/queue-go/pb"
	"github.com/cbotte21/queue-go/schema"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Verify environment variables exist
	enviroment.VerifyEnvVariable("port")
	enviroment.VerifyEnvVariable("chess_addr")

	port := enviroment.GetEnvVariable("port")

	// Setup tcp listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port: %s", port)
	}
	grpcServer := grpc.NewServer()

	// Register handlers to attach
	redisClient := datastore.RedisClient[schema.Queue]{}
	redisClient.Init()

	chessClient := chess.NewChessServiceClient(getChessConn())

	// Initialize hive
	queue := internal.NewQueue(&chessClient, &redisClient)
	pb.RegisterQueueServiceServer(grpcServer, &queue)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}

func getChessConn() *grpc.ClientConn {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(enviroment.GetEnvVariable("chess_addr"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf(err.Error())
	}
	return conn
}
