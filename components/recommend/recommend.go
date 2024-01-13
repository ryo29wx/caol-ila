package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/ryo29wx/caolila_interfaces/recommend"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	recommendReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "recommend_request",
		Help: "Total number of requests that have come to recommend",
	})

	recommendResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "recommend_response",
		Help: "Total number of response that send from recommend",
	})

	logger     *zap.Logger
	collection *mongo.Collection
	gorseSvc   string
)

const (
	port = ":50052"
)

/*
type RecommnedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Page     int32  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Listsize int32  `protobuf:"varint,3,opt,name=listsize,proto3" json:"listsize,omitempty"`
	Token    string `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
}

type RecommendResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Response []*RecommendResponse `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
}


type RecommendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	NickName  string `protobuf:"bytes,2,opt,name=nick_name,json=nickName,proto3" json:"nick_name,omitempty"`
	Sex       int32  `protobuf:"varint,3,opt,name=sex,proto3" json:"sex,omitempty"`
	Title     string `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Company   string `protobuf:"bytes,5,opt,name=company,proto3" json:"company,omitempty"`
	Like      bool   `protobuf:"varint,6,opt,name=like,proto3" json:"like,omitempty"`
	ImageUrl  string `protobuf:"bytes,7,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	LastLogin string `protobuf:"bytes,8,opt,name=last_login,json=lastLogin,proto3" json:"last_login,omitempty"`
}
*/

type server struct {
	pb.UnimplementedRecommenderServer
}

type UserLike struct {
	Likes []string `bson:"likes"`
}

type User struct {
	UserID    string
	NickName  string
	Sex       int32
	Age       int32
	Title     string
	Company   string
	likes     []string
	dislikes  []string
	blocks    []string
	MainImage string
	ImagePath []string
	RegistDay string
	LastLogin string
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
	userId := req.UserId
	page := req.Page
	size := req.Listsize
	token := req.Token

	logger.Debug("Request log", zap.String("userId", userId), zap.Int32("page", page), zap.Int32("size", size), zap.String("token", token))

	recommendReqCount.Inc()

	var result UserLike
	err := collection.FindOne(context.TODO(), bson.M{
		"user_id": userId,
	}).Decode(&result)
	if err != nil {
		logger.Error("user id is null.")
	}

	// insert data to Gorse API
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
	recommendPath := gorseSvc + "/recommend/"
	resp, err := http.Get(recommendPath + userId + "?n=" + strconv.Itoa(int(page)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var recommendations []string
	json.Unmarshal(body, &recommendations)

	var user User
	response := make([]*pb.RecommendResponse, 0)
	for _, r := range recommendations {
		err = collection.FindOne(context.TODO(), bson.M{
			"user_id": r,
		}).Decode(&user)
		if err != nil {
			logger.Error("user id is null.")
		}

		rUser := pb.RecommendResponse{
			UserId:    user.UserID,
			NickName:  user.NickName,
			Sex:       user.Sex,
			Title:     user.Title,
			Company:   user.Company,
			Like:      true, // 一時的に適当に
			ImageUrl:  user.MainImage,
			LastLogin: user.LastLogin,
		}
		response = append(response, &rUser)
	}

	recommendResCount.Inc()
	return &pb.RecommendResponseList{Response: response}, nil
}

func insertFeedback(feedback Feedback) error {
	jsonValue, _ := json.Marshal(feedback)
	insertPath := gorseSvc + "/insert-feedback"
	_, err := http.Post(insertPath, "application/json", bytes.NewBuffer(jsonValue))
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
		panic(err)
	}
	collection = client.Database("caolila").Collection("user")

	// set-up Gorse API
	gorseSvc = os.Getenv("GORSE_SVC")
	if gorseSvc == "" {
		logger.Error("failed to get gorse dns name from os env.")
		gorseSvc = "localhost:8087"
	}

	// start application
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to set-up port listen with gRPC.")
	}
	grpcserver := grpc.NewServer()
	pb.RegisterRecommenderServer(grpcserver, &server{})
	if err := grpcserver.Serve(lis); err != nil {
		logger.Error("failed to set-up application server.")
		panic(err)
	}

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
