package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	n             *int     = flag.Int("n", 10, "How many jsons which you want to create")
	testUserNames []string = []string{"Ryo Kiuchi", "木内 量", "Test User", "Hoge Fuga"}
)

// UserProfile はユーザーのプロファイル情報を表す構造体です。
/*
{"user_id": "string",
 "nick_name": "string",
 "sex": int, //1=man, 2=wemen, 3=other
 "age": int,
 "job_title": "string",
 "company": "string",
 "likes": ["string", "string", ... ],
 "dislikes": ["string", "string", ... ],
 "blocks": ["string", "string", ... ],
 "main_image": "string",
 "image_path": ["string", "string", ... ],
 "regist_day": timestamp,
 "last_login": timestamp}
}
*/
type UserProfile struct {
	UserID    string   `json:"user_id"`
	UserName  string   `json:"user_name"`
	Sex       int      `json:"sex"`
	Age       int      `json:"age"`
	JobTitle  string   `json:"job_title"`
	Company   string   `json:"company"`
	Likes     []string `json:"likes"`
	Dislikes  []string `json:"dislikes"`
	Blocks    []string `json:"blocks"`
	MainImage string   `json:"main_image"`
	ImagePath []string `json:"image_path"`
	RegistDay string   `json:"regist_day"`
	LastLogin string   `json:"last_login"`
}

// 550e8400-e29b-41d4-a716-446655440000
func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	var result []UserProfile
	for i := 0; i < *n; i++ {
		un := rand.Int()
		un = un % len(testUserNames)

		profile := UserProfile{
			UserID:    generateUUID(),
			UserName:  testUserNames[un],
			Sex:       1,
			Age:       20,
			JobTitle:  "hoge",
			Company:   "株式会社サンプル",
			Likes:     []string{generateUUID(), generateUUID()},
			Dislikes:  []string{generateUUID(), generateUUID()},
			Blocks:    []string{generateUUID(), generateUUID()},
			MainImage: "http://sample.com/image01",
			ImagePath: []string{"http://sample.com/image02", "http://sample.com/image03"},
			RegistDay: "hoge",
			LastLogin: "fuga",
		}

		result = append(result, profile)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = os.WriteFile("testData.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

}

func generateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "550e8400-e29b-41d4-a716-446655440000"
	}
	return u.String()
}
