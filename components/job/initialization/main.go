package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.uber.org/zap"
)

type User struct {
	Name string
	Age  int
	City string
}

const (
	port           = ":50051"
	jsonMountPoint = "/data/index/users.json"
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
		log.Println("failed to set-up zap log in search component. \n")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("success set-up logging function.")

	defer logger.Sync()

	// set-up MongoDB client
	mongoHost := os.Getenv("MONGO_SVC_SERVICE_HOST")
	mongoPort := os.Getenv("MONGO_SVC_SERVICE_PORT")
	mongoPass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")

	if mongoHost == "" {
		logger.Error("does not exist remote mongo host.")
		mongoHost = "localhost"
	}
	if mongoPort == "" {
		logger.Error("does not exist remote mongo port.")
		mongoPort = "27017"
	}
	if mongoPass == "" {
		logger.Error("does not exist remote mongo password.")
		mongoPass = "bar"
	}
	if mongoUser == "" {
		logger.Error("does not exist remote mongo username.")
		mongoUser = "bar"
	}

	remoteMongoHost := "mongodb://" + mongoUser + ":" + mongoPass + "@" + mongoHost + ":" + mongoPort
	client, err := mongo.NewClient(options.Client().ApplyURI(remoteMongoHost))
	if err != nil {
		logger.Error("does not exist remote mongo port.")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logger.Error("unexpected error occur when connect to mongo.")
		panic(err)
	}
	defer client.Disconnect(ctx)

	mongoUserDb := os.Getenv("MONGO_USER_DB_NAME")
	mongoUserCollection := os.Getenv("MONGO_USER_COLLECTION_NAME")
	if mongoHost == "" {
		logger.Error("does not exist MONGO_USER_DB_NAME.")
		mongoUserDb = "user_info"
	}
	if mongoPort == "" {
		logger.Error("does not exist MONGO_USER_COLLECTION_NAME.")
		mongoUserCollection = "users"
	}

	logger.Debug("mongo user db: " + mongoUserDb + "mongo user collection: " + mongoUserCollection)
	collection := client.Database(mongoUserDb).Collection(mongoUserCollection)

	// データを挿入する
	user := User{"John Doe", 30, "New York"}
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Success to insert data.")

}
