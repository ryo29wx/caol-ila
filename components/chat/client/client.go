package main

import (
	"bufio"
	"context"
	"log"
	"os"

	pb "github.com/ryo29wx/caolila_interfaces/chat"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger *zap.Logger
)

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

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Fatalf("failed to receive a message : %v", err)
			}
			log.Printf("Got message %s: %s", in.User, in.Message)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := &pb.ChatMessage{
			User:    "YourUserName", // ユーザー名をセット
			Message: scanner.Text(),
		}
		if err := stream.Send(msg); err != nil {
			log.Fatalf("failed to send a message: %v", err)
		}
	}
}
