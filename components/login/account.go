package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

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
	getMyAccountReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "getmyaccount_request",
		Help: "Total number of requests that have come to getmyaccount query",
	})

	getMyAccountResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "getmyaccount_response",
		Help: "Total number of response that send from getmyaccount query",
	})
	editMyAccountReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "editmyaccount_request",
		Help: "Total number of requests that have come to editmyaccount query",
	})

	editMyAccountResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "editmyaccount_response",
		Help: "Total number of response that send from editmyaccount query",
	})

	createMyAccountReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "createmyaccount_request",
		Help: "Total number of requests that have come to editmyaccount query",
	})

	createMyAccountResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "createmyaccount_response",
		Help: "Total number of response that send from editmyaccount query",
	})

	deleteMyAccountReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "deletemyaccount_request",
		Help: "Total number of requests that have come to editmyaccount query",
	})

	deleteMyAccountResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "deletemyaccount_response",
		Help: "Total number of response that send from editmyaccount query",
	})

	loginMyAccountReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "loginmyaccount_request",
		Help: "Total number of requests that have come to editmyaccount query",
	})

	loginMyAccountResCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "loginmyaccount_response",
		Help: "Total number of response that send from editmyaccount query",
	})

	logger *zap.Logger
	db     *sql.DB
)

type User struct {
	UserID       string    `json:"user_id"`
	FirstName    string    `json:"nick_name"`
	LastName     string    `json:"sex"`
	MailAddress1 string    `json:"mail_address1"`
	MailAddress2 string    `json:"mail_address2"`
	MailAddress3 string    `json:"mail_address3"`
	PhoneNum1    string    `json:"phone_num1"`
	PhoneNum2    string    `json:"phone_num2"`
	PhoneNum3    string    `json:"phone_num3"`
	Address1     string    `json:"address1"`
	Address2     string    `json:"address2"`
	Address3     string    `json:"address3"`
	PostCode     int       `json:"post_code"`
	PayRank      int       `json:"pay_rank"`
	Sex          int       `json:"sex"`
	RegistDay    time.Time `json:"regist_day"`
	Birthday     time.Time `json:"birthday"`
}

type UserData struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

// Manipulate SQL
func getMyAccount(c *gin.Context) {

	userID := c.Query("u")

	if userID == "" {
		logger.Error("User ID is missing when getting user data from SQL.")
		c.JSON(http.StatusNoContent, gin.H{"message": "hoge"})
		return
	}

	// logging request log
	logger.Debug("[getMyAccount] Request log", zap.String("user_id", userID))

	// increment counter
	getMyAccountReqCount.Inc()
	user, err := getMyAccountDataBySQL(userID)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Errorf("no user found with user_id %d", userID),
		})
		return
	}
	// increment counter
	getMyAccountResCount.Inc()

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func getMyAccountDataBySQL(userID string) (*User, error) {
	query := `SELECT * FROM user_info WHERE user_id = $1`
	row := db.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.MailAddress1, &user.MailAddress2, &user.MailAddress3, &user.PhoneNum1, &user.PhoneNum2, &user.PhoneNum3, &user.Address1, &user.Address2, &user.Address3, &user.PostCode, &user.PayRank, &user.Sex, &user.RegistDay, &user.Birthday)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("No user data in MongoDB.")
			return nil, err
		}
		logger.Error("SQL error when user get account data.")
		return nil, err
	}
	return &user, nil
}

func editMyAccount(c *gin.Context) {
	userID := c.PostForm("user_id")
	firstName := c.PostForm("first_name")       // string
	lastName := c.PostForm("last_name")         // number
	mailAddress1 := c.PostForm("mail_address1") // number
	mailAddress2 := c.PostForm("mail_address2") // number
	mailAddress3 := c.PostForm("mail_address3") // token
	phoneNum1 := c.PostForm("phone_num1")       // token
	phoneNum2 := c.PostForm("phone_num2")       // token
	phoneNum3 := c.PostForm("phone_num3")       // token
	address1 := c.PostForm("address1")          // token
	address2 := c.PostForm("address2")          // token
	address3 := c.PostForm("address3")          // token
	postCode := c.PostForm("post_code")         // token
	payRank := c.PostForm("pay_rank")           // token
	sex := c.PostForm("sex")                    // token
	birthday := c.PostForm("birthday")          // token

	logger.Debug("Request Edit Account data",
		zap.String("userID", userID),
		zap.String("firstName", firstName),
		zap.String("lastName", lastName),
		zap.String("mailAddress1", mailAddress1),
		zap.String("mailAddress2", mailAddress2),
		zap.String("mailAddress3", mailAddress3),
		zap.String("phoneNum1", phoneNum1),
		zap.String("phoneNum2", phoneNum2),
		zap.String("phoneNum3", phoneNum3),
		zap.String("address1", address1),
		zap.String("address2", address2),
		zap.String("address3", address3),
		zap.String("postCode", postCode),
		zap.String("payRank", payRank),
		zap.String("sex", sex),
		zap.String("birthday", birthday))

	// increment counter
	editMyAccountReqCount.Inc()

	if userID == "" {
		logger.Error("userID parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if firstName == "" {
		logger.Error("firstName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if lastName == "" {
		logger.Error("lastName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if mailAddress1 == "" {
		logger.Error("mailAddress1 parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if payRank == "" {
		logger.Error("payRank parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}

	// check the data if it exist
	_, err := getMyAccountDataBySQL(userID)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Errorf("no user found with user_id %d", userID),
		})
		return
	}

	query := `UPDATE user_info SET first_name = $1, last_name = $2, mail_address1 = $3, mail_address2 = $4, mail_address3 = $5, phone_num1 = $6, phone_num2 = $7, phone_num3 = $8, address1 = $9, address2 = $10, address3 = $11, post_code = $12, pay_rank = $13, sex = $14, birthday = $15 WHERE user_id = $16`
	result, err := db.Exec(query, firstName, lastName, mailAddress1, mailAddress2, mailAddress3, phoneNum1, phoneNum2, phoneNum3, address1, address2, address3, postCode, payRank, sex, birthday, userID)
	if err != nil {
		logger.Error("error insert data.")
		c.JSON(http.StatusForbidden, "NG")
		return
	}

	// check the number of affected rows. if it does not exist, return error.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error when checking affected data.")
		c.JSON(http.StatusForbidden, "NG")
		return
	}

	if rowsAffected == 0 {
		logger.Error("No affected data when insert.")
		c.JSON(http.StatusForbidden, "NG")
		return
	}

	// increment counter
	editMyAccountResCount.Inc()
	c.JSON(http.StatusOK, "OK")
}

func createMyAccount(c *gin.Context) {

	firstName := c.PostForm("first_name")       // string
	lastName := c.PostForm("last_name")         // number
	mailAddress1 := c.PostForm("mail_address1") // number
	mailAddress2 := c.PostForm("mail_address2") // number
	mailAddress3 := c.PostForm("mail_address3") // token
	phoneNum1 := c.PostForm("phone_num1")       // token
	phoneNum2 := c.PostForm("phone_num2")       // token
	phoneNum3 := c.PostForm("phone_num3")       // token
	address1 := c.PostForm("address1")          // token
	address2 := c.PostForm("address2")          // token
	address3 := c.PostForm("address3")          // token
	postCode := c.PostForm("post_code")         // token
	payRank := c.PostForm("pay_rank")           // token
	sex := c.PostForm("sex")                    // token
	birthday := c.PostForm("birthday")          // token

	logger.Debug("Request Edit Account data",
		zap.String("firstName", firstName),
		zap.String("lastName", lastName),
		zap.String("mailAddress1", mailAddress1),
		zap.String("mailAddress2", mailAddress2),
		zap.String("mailAddress3", mailAddress3),
		zap.String("phoneNum1", phoneNum1),
		zap.String("phoneNum2", phoneNum2),
		zap.String("phoneNum3", phoneNum3),
		zap.String("address1", address1),
		zap.String("address2", address2),
		zap.String("address3", address3),
		zap.String("postCode", postCode),
		zap.String("payRank", payRank),
		zap.String("sex", sex),
		zap.String("birthday", birthday))

	// increment counter
	createMyAccountReqCount.Inc()

	if firstName == "" {
		logger.Error("firstName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if lastName == "" {
		logger.Error("lastName parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}
	if mailAddress1 == "" {
		logger.Error("mailAddress1 parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Parameter missing"})
		return
	}

	newUUID := uuid.New()

	query := `INSERT INTO user_info (user_id, first_name, mailAddress1, mailAddress2, mailAddress3, phoneNum1, phoneNum2, phoneNum3, address1, address2, address3, post_code, sex, pay_rank, regist_day, birthday) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW())`

	_, err := db.Exec(query, newUUID, firstName, lastName, mailAddress1, mailAddress2, mailAddress3, phoneNum1, phoneNum2, phoneNum3, address1, address2, address3, postCode, payRank, sex, birthday)
	if err != nil {
		logger.Error("error insert data.")
		c.JSON(http.StatusForbidden, "NG")
		return
	}

	// increment counter
	createMyAccountResCount.Inc()
	c.JSON(http.StatusOK, "OK")
}

func deleteMyAccount(c *gin.Context) {

	userID := c.PostForm("userID")             // string
	mailaddress := c.PostForm("mail_address1") // number
	password := c.PostForm("password")         // number

	if userID == "" || mailaddress == "" || password == "" {
		logger.Error("Delete account parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Some mandatory parameter are missing."})
		return
	}
	// logging request log
	logger.Debug("Request log", zap.String("userID", userID), zap.String("mailaddress", mailaddress), zap.String("password", password))

	user, err := getMyAccountDataBySQL(userID)
	if err != nil {
		logger.Error("error on getting data from SQL")
		return
	}

	if user.MailAddress1 != mailaddress {
		logger.Error("different mailaddress")
		return
	}

	if "password" != password {
		logger.Error("different password")
		return
	}

	query := `DELETE FROM user_info WHERE user_id = $1`
	result, err := db.Exec(query, userID)
	if err != nil {
		logger.Error("No user data in MongoDB.")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error("No user data in MongoDB.")
		return
	}

	if rowsAffected == 0 {
		logger.Error("No user data in MongoDB.")
		return
	}

	c.JSON(http.StatusOK, "OK")
	return
}

/*
 *
 */
func loginMyAccount(c *gin.Context) {

	mailaddress := c.PostForm("mail_address1") // number
	password := c.PostForm("password")         // number

	if mailaddress == "" || password == "" {
		logger.Error("Delete account parameter is missing.")
		c.JSON(http.StatusNoContent, gin.H{"message": "Some mandatory parameter are missing."})
		return
	}

	// logging request log
	logger.Debug("Request log", zap.String("mailaddress", mailaddress), zap.String("password", password))

	user, err := getMyAccountDataBySQL(mailaddress)
	if err != nil {
		logger.Error("error on getting data from SQL")
		return
	}

	if user.MailAddress1 != mailaddress {
		logger.Error("different mailaddress")
		return
	}

	if "password" != password {
		logger.Error("different password")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func changeMyPassword(c *gin.Context) {

	userID := c.PostForm("user_id")
	oldPassword := c.PostForm("password")     // number
	newPassword := c.PostForm("new_password") // number
	logger.Debug("Request log", zap.String("userID", userID), zap.String("oldPassword", oldPassword), zap.String("newPassword", newPassword))

	c.JSON(http.StatusOK, "OK")
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

	// set up PostgreSQL
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgInitUser := os.Getenv("PG_INIT_USER")
	pgInitPass := os.Getenv("PG_INIT_PASS")
	pgCaolilaDB := os.Getenv("PG_CAOLILA_DB")
	pgSsl := os.Getenv("PG_SSL")

	if pgHost == "" {
		logger.Error("does not exist PG_HOST.")
		pgHost = "localhost"
	}
	if pgPort == "" {
		logger.Error("does not exist PG_PORT.")
		pgPort = "5432"
	}
	if pgInitUser == "" {
		logger.Error("does not exist PG_INIT_USER.")
		pgInitUser = "power"
	}
	if pgInitPass == "" {
		logger.Error("does not exist PG_INIT_PASS.")
		pgInitPass = "bar"
	}
	if pgCaolilaDB == "" {
		logger.Error("does not exist PG_CAOLILA_DB.")
		pgCaolilaDB = "caolila"
	}
	if pgSsl == "" {
		logger.Error("does not exist PG_SSL.")
		pgSsl = "disable"
	}
	connStr := "host=" + pgHost + " port=" + pgPort + " user=" + pgInitUser + " password=" + pgInitPass + " dbname=" + pgCaolilaDB + " sslmode=" + pgSsl
	logger.Debug("Postgre connection:", connStr)

	// データベースへの接続を開く
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// start application
	router := gin.Default()
	router.GET("v1/account/get", getMyAccount)
	router.POST("v1/account/edit", editMyAccount)
	router.POST("v1/account/create", createMyAccount)
	router.DELETE("v1/account/delete", deleteMyAccount)
	router.POST("v1/account/login", loginMyAccount)
	router.POST("v1/account/changepass", changeMyPassword)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
