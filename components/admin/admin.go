package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	port         = ":50051"
	keycloakURL  = "http://localhost:8080"
	realm        = "caolila-realm-dev"
	adminUser    = "superadmin"
	adminPass    = "superadmin"
	clientID     = "admin-cli"
	clientSecret = "caolila-client-secret"
)

var user struct {
	Firstname   string `json:"first_name" binding:"required"`
	Lastname    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNum    string `json:"phone_num" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Postcode    string `json:"postcode" binding:"required"`
	Address1    string `json:"address1" binding:"required"`
	Address2    string `json:"address2" binding:"required"`
	Address3    string `json:"address3" binding:"required"`
	Birthday    string `json:"birthday"`
}

var userForDelete struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var (
	createUserReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_user_req",
		Help: "Total number of requests that have come to create-user",
	})

	updateUserReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "update_user_req",
		Help: "Total number of response that send from update-user",
	})

	deleteUserReqCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delete_user_req",
		Help: "Total number of response that send from delete-user",
	})

	logger *zap.Logger
)

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

// Function to get an access token from Keycloak
func getAccessToken() (string, error) {
	data := "client_id=" + clientID + "&client_secret=" + clientSecret + "&username=" + adminUser + "&password=" + adminPass + "&grant_type=password"
	req, err := http.NewRequest("POST", keycloakURL+"/realms/"+realm+"/protocol/openid-connect/token", bytes.NewBuffer([]byte(data)))
	logger.Debug("getAccessToken called: ", zap.String("data", data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result["access_token"].(string), nil
}

func findUserByEmail(email string) ([]map[string]interface{}, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", keycloakURL+"/admin/realms/"+realm+"/users?email="+email, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error")
	}
	logger.Debug(string(b))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to find user: %s", resp.Status)
	}

	var users []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// Function to create a new user in Keycloak
func createUser(c *gin.Context) {
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error:", err)		
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("CreateUser called",
		zap.String("first name", user.Firstname),
		zap.String("last name", user.Lastname),
		zap.String("email", user.Email),
		zap.String("password", user.Password),
		zap.String("postcode", user.Postcode),
		zap.String("address1", user.Address1),
		zap.String("address2", user.Address2),
		zap.String("address3", user.Address3),
		zap.String("phonenum", user.PhoneNum),
		zap.String("birthday", user.Birthday),
	)

	data, err := findUserByEmail(user.Email)
	if err == nil {
		logger.Error("user data already exist.")
		logger.Debug(data[0]["username"].(string))
		return
	}

	token, err := getAccessToken()
	if err != nil {
		logger.Error("failed to get access token.")
		return
	}
	logger.Debug(token)

	username := user.Firstname + "_" + user.Lastname

	user := map[string]interface{}{
		"username":      username,
		"email":         user.Email,
		"enabled":       true,
		"emailVerified": true,
		"firstName":     user.Firstname,
		"lastName":      user.Lastname,
		"credentials": []map[string]interface{}{
			{
				"type":  "password",
				"value": user.Password,
			},
		},
		"attributes": map[string]interface{}{
			"phonenum": user.PhoneNum,
			"postcode": user.Postcode,
			"address1": user.Address1,
			"address2": user.Address2,
			"address3": user.Address3,
			"birthday": user.Birthday,
		},
		"groups": []string{
			"User",
		},
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		logger.Error("error occurd when marshal json in createuser.")
		return
	}

	req, err := http.NewRequest("POST", keycloakURL+"/admin/realms/"+realm+"/users", bytes.NewBuffer(jsonUser))
	if err != nil {
		logger.Error("error occurd when create request URL in createuser.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("error occurd when call create user API of keycloak in createuser.")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error")
	}
	logger.Debug(string(b))

	if resp.StatusCode != http.StatusCreated {
		logger.Error("failed to create new user.")
		return
	}
}

func updateUser(c *gin.Context) {
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("UpdateUser called",
		zap.String("first name", user.Firstname),
		zap.String("last name", user.Lastname),
		zap.String("email", user.Email),
		zap.String("password", user.Password),
		zap.String("postcode", user.Postcode),
		zap.String("address1", user.Address1),
		zap.String("address2", user.Address2),
		zap.String("address3", user.Address3),
		zap.String("phonenum", user.PhoneNum),
		zap.String("birthday", user.Birthday),
	)

	foundUser, err := findUserByEmail(user.Email)
	if err != nil {
		logger.Error("user who is you want to update does not exist.")
		return
	}

	userID := foundUser[0]["id"].(string)
	username := foundUser[0]["username"].(string)
	firstname := foundUser[0]["firstname"].(string)
	lastname := foundUser[0]["lastname"].(string)
	attributes := foundUser[0]["attributes"].(map[string]interface{})
	phonenum := attributes["phonenum"].(string)
	postcode := attributes["postcode"].(string)
	address1 := attributes["address1"].(string)
	address2 := attributes["address2"].(string)
	address3 := attributes["address3"].(string)

	if firstname == user.Firstname {
		firstname = user.Firstname
	}
	if lastname == user.Lastname {
		lastname = user.Lastname
	}
	if firstname == user.Firstname || lastname == user.Lastname {
		username = user.Firstname + "_" + user.Lastname
	}
	if phonenum == user.PhoneNum {
		phonenum = user.PhoneNum
	}
	if postcode == user.Postcode {
		postcode = user.Postcode
	}
	if address1 == user.Address1 {
		address1 = user.Address1
	}
	if address2 == user.Address2 {
		address2 = user.Address2
	}
	if address3 == user.Address3 {
		address3 = user.Address3
	}

	token, err := getAccessToken()
	if err != nil {
		logger.Error("failed to get access token.")
		return
	}
	logger.Debug(token)

	user := map[string]interface{}{
		"username":      username,
		"email":         user.Email,
		"enabled":       true,
		"emailVerified": true,
		"firstName":     firstname,
		"lastName":      lastname,
		"credentials": []map[string]interface{}{
			{
				"type":  "password",
				"value": user.Password,
			},
		},
		"attributes": map[string]interface{}{
			"phonenum": phonenum,
			"postcode": postcode,
			"address1": address1,
			"address2": address2,
			"address3": address3,
		},
		"groups": []string{
			"User",
		},
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		logger.Error("error occurd when marshal json in update user.")
		return
	}

	req, err := http.NewRequest("PUT", keycloakURL+"/admin/realms/"+realm+"/users/"+userID, bytes.NewBuffer(jsonUser))
	if err != nil {
		logger.Error("error occurd when update request URL in update user.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("error occurd when call update user API of keycloak in updateuser.")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error to read response body.")
	}
	logger.Debug(string(b))

	if resp.StatusCode != http.StatusCreated {
		logger.Error("failed to update user.")
		return
	}
}

func deleteUser(c *gin.Context) {
	if err := c.ShouldBindJSON(&userForDelete); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("DeleteUser called",
		zap.String("email", userForDelete.Email),
		zap.String("password", userForDelete.Password),
	)

	foundUser, err := findUserByEmail(userForDelete.Email)
	if err != nil {
		logger.Error("failed to get user data by email.")
		return
	}

	userID := foundUser[0]["id"].(string)

	token, err := getAccessToken()
	if err != nil {
		logger.Error("failed to get access token.")
		return
	}

	logger.Debug(token)

	req, err := http.NewRequest("DELETE", keycloakURL+"/admin/realms/"+realm+"/users/"+userID, nil)
	if err != nil {
		logger.Error("error occurd when delete request URL in deleteuser.")
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("error occurd when call delete user API of keycloak in delete user account.")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error to read response body.")
	}
	logger.Debug(string(b))

	if resp.StatusCode != http.StatusCreated {
		logger.Error("failed to delete user.")
		return
	}
}

func main() {
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
		log.Println("failed to set-up zap log in admin. \n")
		panic(err)
	}

	logger.Debug("this is development environment.")
	logger.Info("success set-up logging function.")

	defer logger.Sync()

	// expose /metrics endpoint for observer(by default Prometheus).
	go exportMetrics()

	// start application
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("v1/admin/account/create", createUser)
	router.PUT("v1/admin/account/update", updateUser)
	router.DELETE("v1/admin/account/delete", deleteUser)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}