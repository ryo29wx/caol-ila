package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	col := client.Database("your_database").Collection("your_collection")

}

func (s *server) Recommend(ctx context.Context, req *pb.RecommnedRequest) (*pb.RecommendResponseList, error) {
	userId := pb.user_id
	response, _ := http.Get(fmt.Sprintf("http://localhost:8088/api/recommend/%s?n=10", userId))

	body, _ := ioutil.ReadAll(response.Body)
	var recommendedItems []string
	json.Unmarshal(body, &recommendedItems)

	fmt.Println("Recommended Items:", recommendedItems)
}

type Feedback struct {
	UserId string `json:"UserId"`
	ItemId string `json:"ItemId"`
	Event  string `json:"Event"`
}

func insertFeedback(feedback Feedback) error {
	jsonValue, _ := json.Marshal(feedback)
	_, err := http.Post("http://localhost:8087/insert-feedback", "application/json", bytes.NewBuffer(jsonValue))
	return err
}

func main() {
	// MongoDBから取得したデータをもとにフィードバックを挿入
	feedback := Feedback{
		UserId: "user_id",
		ItemId: "item_id",
		Event:  "like", // または "dislike"
	}
	err := insertFeedback(feedback)
	if err != nil {
		log.Fatal(err)
	}
}

func getRecommendations(userId string, n int) ([]string, error) {
	// Call Gorse's recomend API
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

func main() {
	userId := "your_user_id"
	recommendations, err := getRecommendations(userId, 10)
	if err != nil {
		log.Fatal(err)
	}

	// 推薦結果の処理
	fmt.Println(recommendations)
}
