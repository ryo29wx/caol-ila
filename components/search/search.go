package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	meilisearch "github.com/meilisearch/meilisearch-go"
	pb "github.com/mockten/mockten_interfaces/searchitem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var (
	searchReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_req_total",
		Help: "Total number of requests that have come to search-item",
	})

	searchResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_res_total",
		Help: "Total number of response that send from serch-item",
	})

	logger      *zap.Logger
	meiliclient *meilisearch.Client
)

type server struct {
	pb.UnimplementedSearchItemsServer
}

// Implement SearchItemServer using protocol buffer
func (s *server) SearchItem(ctx context.Context, in *pb.GetSearchItem) (*pb.SearchResponse, error) {
	productNameForSearch := in.GetProductName()
	sellerNameForSearch := in.GetSellerName()
	exhibitionDateForSearch := in.GetExhibitionDate()
	updateDateForSearch := in.GetUpdateDate()
	categoryForSearch := strconv.Itoa(int(in.GetCategory()))
	rankingFilterForSearch := in.GetRankingFilter()
	pageForSearch := in.GetPage() //int32

	// logging request log
	logger.Info("Request log", zap.String("productname", productNameForSearch),
		zap.String("sellername", sellerNameForSearch),
		zap.String("exhibitiondate", exhibitionDateForSearch),
		zap.String("updatedate", updateDateForSearch),
		zap.String("category", categoryForSearch),
		zap.Int32("ranking", rankingFilterForSearch),
		zap.Int32("page", pageForSearch))

	// increment counter
	searchReqCount.Inc()
	products := make([]*pb.ResponseResult, 0)

	searchRes, err := meiliclient.Index("products").Search(productNameForSearch,
		&meilisearch.SearchRequest{
			Limit: 25,
		})
	if err != nil {
		logger.Error("failed to search in some reasons.", zap.Error(err))
		return &pb.SearchResponse{TotalNum: 0, Response: products}, err
	}

	for _, val := range searchRes.Hits {
		if s, ok := val.(*pb.ResponseResult); ok {
			products = append(products, s)
		} else {
			logger.Error("Value is not of type pb.ResponseResult")
		}
	}

	// increment counter
	searchResCount.Inc()

	return &pb.SearchResponse{TotalNum: int32(len(products)), Response: products}, nil
}

func main() {
	// set-up logging environment using zap
	var err error

	environment := os.Getenv("MOCKTEN_ENV")

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
		log.Println("failed to set-up zap log in searchitem. \n")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("success set-up logging function.")

	defer logger.Sync()

	// set-up meilisearch to register products json(documents) to index.
	meiliclient = meilisearch.NewClient(meilisearch.ClientConfig{
		// expect meilisearch sidecar container
		Host:   "http://127.0.0.1:7700",
		APIKey: os.Getenv("MEILISEARCH_MASTERKEY"),
	})

	index := meiliclient.Index("products")

	// If the index 'products' does not exist, Meilisearch creates it when you first add the documents.
	byteValue, err := os.ReadFile("/opt/etc/products.json")
	if err != nil {
		logger.Error("failed to load search json file", zap.Error(err))
		panic(err)
	}

	// decode JSON to struct which is defeined this file.
	var products []*pb.ResponseResult
	json.Unmarshal(byteValue, &products)

	task, err := index.AddDocuments(products)
	if err != nil {
		logger.Error("failed to add document to meilesearch task.", zap.Error(err))
		panic(err)
	}

	logger.Info("success to execute meilisearch", zap.Int64("taskid", task.TaskUID))

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// start application
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to set-up port listen with gRPC.")
	}
	grpcserver := grpc.NewServer()
	pb.RegisterSearchItemsServer(grpcserver, &server{})
	if err := grpcserver.Serve(lis); err != nil {
		logger.Error("failed to set-up application server.")
		panic(err)
	}
}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
