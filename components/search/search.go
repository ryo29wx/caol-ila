package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	meilisearch "github.com/meilisearch/meilisearch-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	port           = ":50051"
	jsonMountPoint = "/data/index/users.json"
)

var (
	searchReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_request",
		Help: "Total number of requests that have come to search query",
	})

	searchResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_response",
		Help: "Total number of response that send from serch query",
	})

	logger      *zap.Logger
	meiliclient *meilisearch.Client
)

type User struct {
	UserID   int    `json:"user_id"`
	NickName string `json:"nick_name"`
	Sex      string `json:"sex"`
	Title    string `json:"title"`
	Company  string `json:"company"`
	Like     int    `json:"like"`
	ImageURL string `json:"image_url"`
}

type UserData struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

//	message SearchRequest {
//		string query = 1;
//		int32 sort = 2;
//		string token = 3;
//	 }
//
// Implement SearchItemServer using protocol buffer
// func (s *server) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
func searchHandler(c *gin.Context) (UserData, error) {

	query := c.Query("q")
	page := c.Query("p")
	token := c.Query("t")

	var userData UserData
	if query == "" || token == "" || page == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query and Token parameters are required"})
		return userData, errors.New("missing argument of search query")
	}

	// logging request log
	logger.Debug("Request log", zap.String("query", query), zap.Int32("page", page), zap.String("token", token))

	// increment counter
	searchReqCount.Inc()

	searchRes, err := meiliclient.Index("users").Search(query,
		&meilisearch.SearchRequest{
			Limit: 25,
		})
	if err != nil {
		logger.Error("failed to search in some reasons.", zap.Error(err))
		return userData, err
	}

	var users []User
	for _, val := range searchRes.Hits {
		if s, ok := val.(*User); ok {
			users = append(users, s)
		} else {
			logger.Error("Value is not of type pb.ResponseResult")
		}
	}

	userData.Users = users
	userData.Total = 1

	// increment counter
	searchResCount.Inc()

	return userData, nil
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

	ticker := time.NewTicker(300 * time.Second)

	// TODO
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		logger.Error("0: unexpected error occur when find data from mongodb.", zap.Error(err))
	}
	defer cur.Close(context.Background())
	var results []bson.M

	if err = cur.All(context.Background(), &results); err != nil {
		logger.Error("0: failed to get data from mongo.", zap.Error(err))
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		logger.Error("0: failed to write data to search component.", zap.Error(err))
	}

	if _, err := os.Stat("/data/index"); os.IsNotExist(err) {
		if err := os.MkdirAll("/data/index", 0755); err != nil {
			logger.Error("0: failed to create data path ", zap.Error(err))
			return
		}
	}

	err = os.WriteFile(jsonMountPoint, jsonData, 0644)
	if err != nil {
		logger.Error("0: failed to write json data to ", zap.Error(err))
	}
	// TODO END

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				cur, err := collection.Find(context.Background(), bson.D{{}})
				if err != nil {
					logger.Error("unexpected error occur when find data from mongodb.")
				}
				defer cur.Close(context.Background())
				var results []bson.M

				if err = cur.All(context.Background(), &results); err != nil {
					logger.Error("failed to get data from mongo.")
				}

				jsonData, err := json.Marshal(results)
				if err != nil {
					logger.Error("failed to write data to search component.")
				}

				err = os.WriteFile(jsonMountPoint, jsonData, 0644)
				if err != nil {
					logger.Error("failed to write json data to ")
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// set-up meilisearch to register products json(documents) to index.
	meiliBackend := os.Getenv("MEILI_SVC")
	if mongoHost == "" {
		logger.Error("does not exist MEILI_SVC.")
		meiliBackend = "127.0.0.1"
	}

	meiliclient = meilisearch.NewClient(meilisearch.ClientConfig{
		// expect meilisearch sidecar container
		Host:   "http://" + meiliBackend + ":7700",
		APIKey: os.Getenv("MEILISEARCH_MASTERKEY"),
	})

	index := meiliclient.Index("users")

	// If the index 'usres' does not exist, Meilisearch creates it when you first add the documents.
	byteValue, err := os.ReadFile(jsonMountPoint)
	if err != nil {
		logger.Error("failed to load search json file", zap.Error(err))
		panic(err)
	}

	// decode JSON to struct which is defeined this file.
	var users []*pb.ResponseResult
	json.Unmarshal(byteValue, &users)

	task, err := index.AddDocuments(users)
	if err != nil {
		logger.Error("failed to add document to meilesearch task.", zap.Error(err))
		panic(err)
	}

	logger.Info("success to execute meilisearch", zap.Int64("taskid", task.TaskUID))

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// start application
	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	logger.Error("failed to set-up port listen with gRPC.")
	// }
	// grpcserver := grpc.NewServer()
	// pb.RegisterSearcherServer(grpcserver, &server{})
	// if err := grpcserver.Serve(lis); err != nil {
	// 	logger.Error("failed to set-up application server.")
	// 	panic(err)
	// }
	router := gin.Default()
	router.GET("v1/search", searchHandler)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
