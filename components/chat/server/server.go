package main

import (
	"log"
	"net"
	"os"

	pb "github.com/ryo29wx/caolila_interfaces/chat"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger *zap.Logger
)

type server struct {
	pb.UnimplementedChatServiceServer
}

func (s *server) Chat(stream pb.ChatService_ChatServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Received message from %s: %s", in.User, in.Message)
		if err := stream.Send(in); err != nil {
			return err
		}
	}
}

func main() {
	// set-up logging environment using zap
	var err error
	environment := os.Getenv("CAOLILA_ENV")

	if environment == "development" || environment == "" {
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = config.Build()
	} else {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		logger, err = config.Build()
	}

	if err != nil {
		log.Println("failed to set-up zap log in recommend component. \n")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("success set-up logging function.")

	defer logger.Sync()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
