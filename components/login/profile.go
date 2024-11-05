package main

import (
	"context"
	// "io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	// conx       context.Context
)

type User struct {
	UserID       string
	DisplayName  string
	Sex          string
	Age          int
	Title        string
	Company      string
	CompanyEmail string
	Likes        []string
	Dislikes     []string
	Blocks       []string
	MainImage    string
	ImagePath    []string
	RegistDay    time.Time
	LastLogin    time.Time
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func getMyProfile(c *gin.Context) {

	id := c.Query("u")

	if id == "" {
		logger.Error("User ID is missing when getting user data from NoSQL.")
		c.JSON(http.StatusNoContent, gin.H{"message": "missing user id"})
		return
	}

	// logging request log
	logger.Debug("[getMyProfile] Request log", zap.String("user_id", id))

	// increment counter
	getMyProfileReqCount.Inc()

	filter := bson.M{"userid": id}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Not exist this user id
			c.JSON(http.StatusOK, gin.H{"error_message": "can not find user by id"})
			return
		}
		// Unexpected Error
		c.JSON(http.StatusOK, gin.H{"error_message": "can not find user by id"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func createMyProfile(c *gin.Context) {
	displayName := c.PostForm("display_name")
	gender := c.PostForm("gender")
	age := c.PostForm("age")
	title := c.PostForm("title")
	company := c.PostForm("company")
	companyEmail := c.PostForm("company_email")
	description := c.PostForm("description")

	logger.Debug("Request create user profile log",
		zap.String("displayName", displayName),
		zap.String("gender", gender),
		zap.String("age", age),
		zap.String("title", title),
		zap.String("company", company),
		zap.String("companyEmail", companyEmail),
		zap.String("description", description),
	)

	// increment counter
	getMyProfileReqCount.Inc()

	if displayName == "" {
		logger.Error("CreateProfile displayName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}

	// manipulate main image
	mainImage, err := c.FormFile("main_image") // []byte
	if err != nil {
		logger.Error("Something wrong with main_image from client.", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the file data to ./upload（保存先ディレクトリ: ./uploads）
	savePath := "./uploads/" + mainImage.Filename
	if err := c.SaveUploadedFile(mainImage, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the image data"})
		return
	}

	// Save the image file on GCS
	// ctx := context.Background()
	// gcsFilePath := bucketName + "/" + sellerName + "/" + mainImage.Filename
	// logger.Debug("gcs file path:", zap.String("gcsFilePath", gcsFilePath))
	// wc := bucket.Object(gcsFilePath).NewWriter(ctx)
	// mainImageReader, err := FileHeaderToReader(mainImage)

	// if _, err := io.Copy(wc, mainImageReader); err != nil {
	// 	logger.Error("Failed to upload file to GCS:", zap.Error(err))
	// }
	// if err := wc.Close(); err != nil {
	// 	logger.Error("Failed to close GCS writer:", zap.Error(err))
	// 	c.JSON(http.StatusNotFound, "error")

	// }
	ageInt, _ := strconv.Atoi(age)

	user := User{
		UserID:       uuid.New().String(),
		DisplayName:  displayName,
		Sex:          gender,
		Age:		  ageInt,
		Title:        title,
		Company:      company,
		CompanyEmail: companyEmail,
		Likes:        make([]string, 0),
		Dislikes:     make([]string, 0),
		Blocks:       make([]string, 0),
		MainImage:    savePath,
		ImagePath:    make([]string, 0),
		RegistDay:    time.Now(),
		LastLogin:    time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		logger.Error("failed to create profile to mongo", zap.Error(err))
		c.JSON(http.StatusNotFound, "error")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success creating your profile"})
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
	mongoUserDb := os.Getenv("MONGO_USER_DB_NAME")

	if mongoHost == "" {
		logger.Error("does not exist remote mongo host.")
		mongoHost = "127.0.0.1"
	}
	if mongoPort == "" {
		logger.Error("does not exist remote mongo port.")
		mongoPort = "27017"
	}
	if mongoPass == "" {
		logger.Error("does not exist remote mongo password.")
		mongoPass = "password"
	}
	if mongoUser == "" {
		logger.Error("does not exist remote mongo username.")
		mongoUser = "user_info_owner"
	}
	if mongoUserDb == "" {
		logger.Error("does not exist MONGO_USER_DB_NAME.")
		mongoUserDb = "user_info"
	}

	remoteMongoHost := "mongodb://" + mongoUser + ":" + mongoPass + "@" + mongoHost + ":" + mongoPort + "/" + mongoUserDb
	clientOptions := options.Client().ApplyURI(remoteMongoHost)

	client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // 接続の確認
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    logger.Debug("Success to connected")

	mongoUserCollection := os.Getenv("MONGO_USER_COLLECTION_NAME")
	if mongoUserCollection == "" {
		logger.Error("does not exist MONGO_USER_COLLECTION_NAME.")
		mongoUserCollection = "users"
	}

	logger.Debug("mongo user db: " + mongoUserDb + "mongo user collection: " + mongoUserCollection)
	collection = client.Database(mongoUserDb).Collection(mongoUserCollection)

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// expose these API
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("v1/profile", getMyProfile)
	// router.POST("v1/profile/edit", editMyProfile)
	router.POST("v1/profile/create", createMyProfile)
	// router.DELETE("v1/profile/delete", deleteMyProfile)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
