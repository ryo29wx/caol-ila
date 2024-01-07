package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/ryo29wx/caolila_interfaces/recommend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type User struct {
	Likes []string `bson:"likes"`
}

type Feedback struct {
	UserId string `json:"UserId"`
	// ItemId string `json:"ItemId"`
	Event string `json:"Event"`
}

/* This function is called from frontend service via gRPC.
 * if you want to see difinition of protbuf interface, see: (https://github.com/ryo29wx/caolila_interfaces/tree/main/recommend)
 *
 */
func (s *server) Recommend(ctx context.Context, req *pb.RecommnedRequest) (*pb.RecommendResponseList, error) {
	// get data from mongoDB
	userId := pb.user_id

	var result User
	err = collection.FindOne(context.TODO(), bson.M{
		"user_id": userId,
	}).Decode(&result)
	if err != nil {
		logger.Error("user id is null.")
	}

	feedback := Feedback{
		UserId: result.Likes[0],
		Event:  "like",
	}

	go func() {
		err := insertFeedback(feedback)
		if err != nil {
			logger.Error("failed to insert data to gorse.")
		}
	}()

	// Call Gorse's recommend API
	resp, err := http.Get("http://localhost:8087/recommend/" + userId + "?n=" + strconv.Itoa(n))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var recommendations []string
	json.Unmarshal(body, &recommendations)
	return recommendations, nil

}

func insertFeedback(feedback Feedback) error {
	jsonValue, _ := json.Marshal(feedback)
	_, err := http.Post("http://localhost:8087/insert-feedback", "application/json", bytes.NewBuffer(jsonValue))
	return err
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

	// Connect MongoDB
	mongoSvc := os.Getenv("MONGO_SVC")
	if mongoSvc == "" {
		logger.Error("failed to get mongo dns name from os env.")
		mongoSvc = "mongodb://mongo-svc.dev_caolila.svc.cluster.local:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoSvc))
	if err != nil {
		logger.Error("failed to connect mongoDB.")
		panic()
	}
	col := client.Database("caolila").Collection("user")

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
