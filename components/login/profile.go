package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	port = ":50051"
)

var (
	getMyProfileReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "getmyprofile_request",
		Help: "Total number of requests that have come to getmyaccount query",
	})

	getMyProfileResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "getmyprofile_response",
		Help: "Total number of response that send from getmyaccount query",
	})

	logger     *zap.Logger
	collection *mongo.Collection
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

//	Get User Data from SQL and NoSQL.
//		string query = 1;
//		int32 sort = 2;
//		string token = 3;
//	 }
//
// Manipulate SQL
func getMyProfile(c *gin.Context) {

	id := c.Query("u")

	if id == "" {
		logger.Error("User ID is missing when getting user data from NoSQL.")
		c.JSON(http.StatusNoContent, gin.H{"message": "hoge"})
		return
	}

	// logging request log
	logger.Debug("[getMyProfile] Request log", zap.String("user_id", id))

	// increment counter
	getMyProfileReqCount.Inc()

	var users []User
	for _, val := range searchRes.Hits {
		if s, ok := val.(User); ok {
			users = append(users, s)
		} else {
			logger.Error("Value is not of type pb.ResponseResult")
		}
	}

	// increment counter
	searchResCount.Inc()

	for _, user := range users {
		if user.UserID == id {
			c.JSON(http.StatusOK, gin.H{
				user
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error_message": "can not find user by id"
	})
}

func editMyProfile(c *gin.Context) {

	productName := c.PostForm("product_name") // string
	sellerName := c.PostForm("seller_name")   // string
	category := c.PostForm("category")        // number
	price := c.PostForm("price")              // number
	stock := c.PostForm("stock")              // number
	token := c.PostForm("token")              // token

	logger.Debug("Request Add Item log",
		zap.String("productName", productName),
		zap.String("sellerName", sellerName),
		zap.String("category", category),
		zap.String("price", price),
		zap.String("stock", stock),
		zap.String("token", token))

	// increment counter
	addItemReqCount.Inc()

	if productName == "" {
		logger.Error("AddItem productName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if sellerName == "" {
		logger.Error("AddItem sellerName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if category == "" {
		logger.Error("AddItem category parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if price == "" {
		logger.Error("AddItem price parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if stock == "" {
		logger.Error("AddItem stock parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if token == "" {
		logger.Error("AddItem token parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}

	// manipulate main image
	mainImage, err := c.FormFile("file") // []byte
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the image file on GCS
	ctx := context.Background()
	gcsFilePath := bucketName + "/" + sellerName + "/" + mainImage.Filename
	logger.Debug("gcs file path:", zap.String("gcsFilePath", gcsFilePath))
	wc := bucket.Object(gcsFilePath).NewWriter(ctx)
	mainImageReader, err := FileHeaderToReader(mainImage)

	if _, err := io.Copy(wc, mainImageReader); err != nil {
		logger.Error("Failed to upload file to GCS:", zap.Error(err))
	}
	if err := wc.Close(); err != nil {
		logger.Error("Failed to close GCS writer:", zap.Error(err))
		c.JSON(http.StatusNotFound, "error")

	}

	// Save the item data on MongoDB
	numStock, err := strconv.Atoi(stock)
	if err != nil {
		logger.Error("convert error with stock value:", zap.Error(err))
		return
	}

	numCategory, err := strconv.Atoi(category)
	if err != nil {
		logger.Error("convert error with category value:", zap.Error(err))
		return
	}

	numPrice, err := strconv.Atoi(price)
	if err != nil {
		logger.Error("convert error with price value:", zap.Error(err))
		return
	}

	item := Item{ProductId: "hogehoge",
		ProductName: productName,
		SellerName:  sellerName,
		Stocks:      numStock,
		Category:    []int{numCategory},
		Rank:        99999,
		MainImage:   gcsFilePath,
		Summary:     "test item",
		Price:       numPrice,
		RegistDay:   time.Now(),
		LastUpdate:  time.Now(),
	}

	_, err = collection.InsertOne(ctx, item)
	if err != nil {
		logger.Error("failed to add item to mongo", zap.Error(err))
		c.JSON(http.StatusNotFound, "error")
	}

	// increment counter
	addItemResCount.Inc()

	c.JSON(http.StatusOK, "OK")
}

func createMyProfile(c *gin.Context) {
	displayName := c.PostForm("display_name") 
	gender := c.PostForm("gender")   
	age := c.PostForm("age")      
	title := c.PostForm("title")            
	company := c.PostForm("company") 
	companyEmail := c.PostForm("company_email") 

	logger.Debug("Request create user profile log",
		zap.String("displayName", displayName),
		zap.String("gender", gender),
		zap.String("age", age),
		zap.String("title", title),
		zap.String("company", company),
		zap.String("companyEmail", companyEmail))

	// increment counter
	getMyProfileReqCount.Inc()

	if displayName == "" {
		logger.Error("CreateProfile displayName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}

	// manipulate main image
	mainImage, err := c.FormFile("file") // []byte
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("Uploaded File: %+v\n", mainImage.Filename)
	logger.Debug("File Size: %+v\n", mainImage.Size)
	logger.Debug("MIME Header: %+v\n", mainImage.Header)

	// Get the file data to ./upload（保存先ディレクトリ: ./uploads）
	savePath := filepath.Join("uploads", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ファイルの保存に失敗しました"})
		return
	}
	

	// Save the image file on GCS
	ctx := context.Background()
	gcsFilePath := bucketName + "/" + sellerName + "/" + mainImage.Filename
	logger.Debug("gcs file path:", zap.String("gcsFilePath", gcsFilePath))
	wc := bucket.Object(gcsFilePath).NewWriter(ctx)
	mainImageReader, err := FileHeaderToReader(mainImage)

	if _, err := io.Copy(wc, mainImageReader); err != nil {
		logger.Error("Failed to upload file to GCS:", zap.Error(err))
	}
	if err := wc.Close(); err != nil {
		logger.Error("Failed to close GCS writer:", zap.Error(err))
		c.JSON(http.StatusNotFound, "error")

	}

	// アップロード成功メッセージを返す
	c.JSON(http.StatusOK, gin.H{"message": "Success creating your profile"})
}

func deleteMyProfile(c *gin.Context) {

	query := c.Query("q")
	page := c.Query("p")
	token := c.Query("t")

	if query == "" || token == "" || page == "" {
		logger.Error("Search Query parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "There is no users"})
		return
	}

	// logging request log
	logger.Debug("Request log", zap.String("query", query), zap.String("page", page), zap.String("token", token))

	// increment counter
	searchReqCount.Inc()

	searchRes, err := meiliclient.Index("users").Search(query,
		&meilisearch.SearchRequest{
			Limit: 25,
		})
	if err != nil {
		logger.Error("failed to search in some reasons.", zap.Error(err))
		c.JSON(http.StatusNoContent, gin.H{"message": "There is no users"})
		return
	}

	var users []User
	for _, val := range searchRes.Hits {
		if s, ok := val.(User); ok {
			users = append(users, s)
		} else {
			logger.Error("Value is not of type pb.ResponseResult")
		}
	}

	// increment counter
	searchResCount.Inc()

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": page,
	})
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
	collection = client.Database(mongoUserDb).Collection(mongoUserCollection)

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// expose these API
	router := gin.Default()
	router.GET("v1/profile/get", getMyProfile)
	router.POST("v1/profile/edit", editMyProfile)
	router.POST("v1/profile/create", editMyProfile)
	router.DELETE("v1/profile/delete", deleteMyProfile)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
