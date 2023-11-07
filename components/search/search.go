package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
	pb "github.com/mockten/mockten_interfaces/searchitem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

/**
*initialize
 */
type server struct {
	pb.UnimplementedSearchItemsServer
}

// Implement SearchItemServer
func (s *server) Search(ctx context.Context, in *pb.Search) (*pb.SearchResponse, error) {
	// logging request log
	logger.Debug(getRequestLog(in.GetProductName(), in.GetSellerName(), in.GetExhibitionDate(), in.GetUpdateDate()))

	// prometheus-exporter server
	searchReqCount.Inc()
	countReqs()

	// Connect DB
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	// close
	defer db.Close()

	productNameForSearch := in.GetProductName()
	sellerNameForSearch := in.GetSellerName()
	exhibitionDateForSearch := in.GetExhibitionDate()
	updateDateForSearch := in.GetUpdateDate()
	categoryForSearch := strconv.Itoa(int(in.GetCategory()))
	rankingFilterForSearch := in.GetRankingFilter()
	pageForSearch := in.GetPage() //int32

	//Get SellerIDs
	sellerIDs := make([]string, 0)
	if len(sellerNameForSearch) > 0 {
		sellerIDs, err = getSellerID(sellerNameForSearch, db)
		if err != nil {
			return nil, err
		}
	}

	// create searchQuery
	searchSQL, countSQL := createSearchSQL(productNameForSearch, sellerIDs, exhibitionDateForSearch, updateDateForSearch, categoryForSearch, rankingFilterForSearch, pageForSearch)

	items, productErr := getProductInfo(db, searchSQL)
	if productErr != nil {
		log.Printf("Search ERROR : %s\n ", productErr)
		return nil, productErr
	}

	total, totalErr := getTotalInfo(db, countSQL)
	if totalErr != nil {
		log.Printf("Search total ERROR : %s\n ", totalErr)
		return nil, totalErr
	}

	log.Printf("Items : %v || Total :  %v \n", items, total)
	searchResCount.Inc()
	countSends()
	return &pb.SearchResponse{TotalNum: total, Response: items}, nil
}

func main() {
	//set-up logging lib.
	var logger *zap.Logger
	var err error

	environment := os.Getenv("GO_ENV")

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
		log.Print("can't implement zap in search component.")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("well done of zap setting in search component.")

	// clean up logger lib
	defer logger.Sync()

	// set-up gRPC receiver
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Error("failed to listen gRPC port in search component.")
		panic(err)
	}

	grpcserver := grpc.NewServer()
	pb.RegisterSearchServer(grpcserver, &server{})
	if err := grpcserver.Serve(lis); err != nil {
		logger.Error("failed to register search server.")
		panic(err)
	}

	// set-up to connect Database.
	db, err := connectDB()
	if err != nil {
		logger.Error("failed to connect database server from search component.")
		panic(err)
	}

	defer db.Close()

	// set-up to connect cache server.
	// TODO

	// set-up to send metrics to metlics server(prometheus by default).
	exportMetrics()
}

var (
	// prometheus metrics name.
	searchReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_request",
		Help: "Total number of requests that have come to search component",
	})

	searchResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_response",
		Help: "Total number of response that send from serch component",
	})
)

func connectDB() (*sql.DB, error) {
	// generally, get environment value from k8s secret resources.
	user := os.Getenv("SECRET_USER")
	pass := os.Getenv("SECRET_PASS")
	sdb := os.Getenv("SECRET_DB")
	table := os.Getenv("SECRET_TABLE")

	mySQLHost := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, sdb, table)
	db, err := sql.Open("mysql", mySQLHost)
	if err = db.Ping(); err != nil {
		logger.Error("failed to connect database server from search component.")
		return nil, err
	}

	return db, nil
}

func getSellerID(sellerName string, db *sql.DB) ([]string, error) {
	sellerIDs := make([]string, 0)
	searchSellerQuery := "SELECT seller_id FROM SELLER_INFO WHERE seller_name LIKE '%" + sellerName + "%';"
	sellerRows, err := db.Query(searchSellerQuery)
	if err != nil {
		log.Printf("DB Seller Query Error: %v | Seller Query is: %v ", err, searchSellerQuery)
		return nil, err
	}
	defer sellerRows.Close()

	var sellerID string
	for sellerRows.Next() {
		if err := sellerRows.Scan(&sellerID); err != nil {
			log.Printf("Search Seller Scan Error: %v", err)
			return nil, err
		}
		sellerIDs = append(sellerIDs, sellerID)
	}
	return sellerIDs, nil
}

func createSearchSQL(productNameForSearch string,
	sellerIDs []string,
	exhibitionDateForSearch string,
	updateDateForSearch string,
	categoryForSearch string,
	rankingFilterForSearch int32,
	pageForSearch int32) (string, string) {
	baseSQL := "SELECT product_id,product_name,seller_name,stock,price,image_path,comment,category FROM PRODUCT_INFO_VIEW WHERE PRODUCT_NAME LIKE '%" + productNameForSearch + "%'"
	baseCountSQL := "SELECT COUNT(*) FROM PRODUCT_INFO WHERE PRODUCT_NAME LIKE '%" + productNameForSearch + "%'"
	// SearchItem all
	if len(sellerIDs) > 0 {
		sellerNameQuery := "("
		for i, sellerID := range sellerIDs {
			if i != 0 {
				sellerNameQuery = sellerNameQuery + ","
			}
			sellerNameQuery = sellerNameQuery + "'" + sellerID + "'"
		}
		sellerNameQuery = sellerNameQuery + ")"
		baseSQL = baseSQL + " AND SELLER_ID IN " + sellerNameQuery
		baseCountSQL = baseCountSQL + " AND SELLER_ID IN " + sellerNameQuery
	}
	if len(exhibitionDateForSearch) > 0 {
		baseSQL = baseSQL + " AND TIME=" + exhibitionDateForSearch
		baseCountSQL = baseCountSQL + " AND TIME=" + exhibitionDateForSearch
	}
	if len(updateDateForSearch) > 0 {
		baseSQL = baseSQL + " AND UPDATE_TIME=" + updateDateForSearch
		baseCountSQL = baseCountSQL + " AND UPDATE_TIME=" + updateDateForSearch
	}
	if categoryForSearch != "99" {
		baseSQL = baseSQL + " AND CATEGORY='" + categoryForSearch + "'"
		baseCountSQL = baseCountSQL + " AND CATEGORY='" + categoryForSearch + "'"
	}

	// attach sortString
	sortStr := CreateSortStr(rankingFilterForSearch)
	baseSQL = baseSQL + sortStr

	pagedata := (10 * pageForSearch) - 10
	baseSQL = baseSQL + " LIMIT " + strconv.Itoa(int(pagedata)) + ",10"
	baseSQL = baseSQL + ";"
	baseCountSQL = baseCountSQL + ";"

	log.Printf("basesQL : " + baseSQL)
	log.Printf("baseCountSQL : " + baseCountSQL)

	return baseSQL, baseCountSQL
}

func CreateSortStr(filter int32) string {
	switch filter {
	case 1:
		return " ORDER BY price DESC"
	case 2:
		return " ORDER BY price ASC"
	case 3:
		return " ORDER BY rate ASC"
	case 4:
		return " ORDER BY update_date ASC"
	case 5:
		return ""
	}
	return ""
}

func getProductInfo(db *sql.DB, searchSQL string) ([]*pb.ResponseResult, error) {

	items := make([]*pb.ResponseResult, 0)
	rows, err := db.Query(searchSQL)
	if err != nil {
		log.Printf("Search Product Query Error: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var responseResult *pb.ResponseResult = &pb.ResponseResult{}
		if err := rows.Scan(&responseResult.ProductId, &responseResult.ProductName, &responseResult.SellerName, &responseResult.Stocks, &responseResult.Price, &responseResult.ImageUrl, &responseResult.Comment, &responseResult.Category); err != nil {
			log.Printf("Search Scan Error: %v | Query is: %v", err, searchSQL)
			return items, err
		}
		items = append(items, responseResult)
		log.Printf("items is : %v", items)
	}

	return items, nil
}

func getTotalInfo(db *sql.DB, countSQL string) (int32, error) {
	var total int32

	rows, err := db.Query(countSQL)
	if err != nil {
		log.Printf("Count Query Error: %v | Query is: %v", err, countSQL)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&total); err != nil {
			log.Printf("Count Scan Error: %v | Query is: %v", err, countSQL)
			return 0, err
		}
	}

	return total, nil
}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
