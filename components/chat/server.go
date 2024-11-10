package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
	"github.com/gorilla/websocket"
	"github.com/gin-gonic/gin"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	port         = ":50055"
)

var (
	logger *zap.Logger
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        logger.Error("Error while upgrading:", zap.Error(err))
        return
    }
    defer conn.Close()

    for {
        // Read message from client
        _, msg, err := conn.ReadMessage()
        if err != nil {
            logger.Error("Error while reading message:", zap.Error(err))
            break
        }
        fmt.Printf("Received: %s\n", msg)

        // Echo message back to client
        err = conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            logger.Error("Error while writing message:", zap.Error(err))
            break
        }
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
		log.Println("failed to set-up zap log in chat. \n")
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
	router.GET("v1/chat/send", handleWebSocket)
	router.Run(port)

}

// for goroutin
func exportMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
