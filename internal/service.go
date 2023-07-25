package internal

import (
	"context"
	"errors"
	chess "github.com/cbotte21/chess-go/pb"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	pb "github.com/cbotte21/queue-go/pb"
	"github.com/cbotte21/queue-go/schema"
)

var SearchingNotice = "searching"
var ListenChannel = "game_start"
var FinalError = errors.New("lost connection to server")

type Queue struct {
	ChessClient *chess.ChessServiceClient
	RedisClient *datastore.RedisClient[schema.Queue]
	pb.UnimplementedQueueServiceServer
}

func NewQueue(chessClient *chess.ChessServiceClient, redisClient *datastore.RedisClient[schema.Queue]) Queue {
	return Queue{ChessClient: chessClient, RedisClient: redisClient}
}

func (queue *Queue) Join(joinRequest *pb.JoinRequest, stream pb.QueueService_JoinServer) error {
	// Check if a player is in queue
	res, found := queue.RedisClient.Find(schema.Queue{Status: SearchingNotice})
	if found == nil { // Player found, start game for both players  TODO: Specify error
		_ = queue.RedisClient.Publish(ListenChannel, "")
		_, _ = (*queue.ChessClient).Create(context.Background(), &chess.CreateRequest{
			Player1: &chess.Player{XId: res.Player},
			Player2: &chess.Player{XId: joinRequest.GetJwt()},
			Ranked:  false,
		})
		_ = stream.Send(&pb.JoinResponse{Status: 2})
		return FinalError
	}

	// Player not found, join queue...

	// Listen for when game is started.
	started := 0
	go func(started *int) {
		sub := queue.RedisClient.Subscribe(ListenChannel)
		ch := sub.Channel()
		for range ch {
			*started = 1
		}
	}(&started)

	cachedRequest := schema.Queue{Status: SearchingNotice, Player: joinRequest.GetJwt()}
	_ = queue.RedisClient.Create(cachedRequest)

	// While waiting for opponent
	for stream.Send(&pb.JoinResponse{Status: 1}) == nil {
		if started == 1 {
			_ = stream.Send(&pb.JoinResponse{Status: 2})
		}
	}

	_ = queue.RedisClient.Delete(cachedRequest)
	_ = stream.Send(&pb.JoinResponse{Status: 0})
	// User disconnected
	return FinalError
}
