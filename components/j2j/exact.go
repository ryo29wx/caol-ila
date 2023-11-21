package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	port           = ":50051"
	jsonMountPoint = "/opt/mnt/output.json"
)

var (
	exactStartCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "exact_start",
		Help: "Total number of start that exact/load data from mongodb to local search component",
	})

	exactDoneCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "exact_done",
		Help: "Total number of done that exact/load data from mongodb to local search component",
	})

	logger *zap.Logger
)

func main() {
	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()
	time.Sleep(5 * time.Second)
	exactStartCount.Inc()

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
		log.Println("failed to set-up zap log in j2j component. \n")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("success set-up logging function.")

	defer logger.Sync()

	// setup mongodb client
	mongoSvcHost := os.Getenv("MONGODB_SVC_SERVICE_HOST")
	mongoSvcPort := os.Getenv("MONGODB_SVC_SERVICE_PORT")

	if mongoSvcHost == "" {
		logger.Error("mongodb host doesn't exist.")
		mongoSvcHost = "localhost"
	}
	if mongoSvcPort == "" {
		logger.Error("mongodb host doesn't exist.")
		mongoSvcHost = "27017"
	}

	mongoURI := "mongodb://" + mongoSvcHost + ":" + mongoSvcPort

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Error("mongodb host doesn't exist.")
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("caolila").Collection("users")

	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		logger.Error("failed to read json object.")
		panic(err)
	}
	defer cur.Close(context.Background())

	var results []bson.M
	if err = cur.All(context.Background(), &results); err != nil {
		logger.Error("failed to transform data from json.")
		panic(err)
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Error("failed to write data to search component.")
		panic(err)
	}

	ioutil.WriteFile(jsonMountPoint, jsonData, 0644)
	// increment counter
	exactDoneCount.Inc()
	time.Sleep(30 * time.Second)
}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
