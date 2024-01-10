package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/google/uuid"
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
	n := flag.Int("n", 10, "How many jsons which you want to create")
	// UserProfile構造体のインスタンスを作成し、データを埋めます。
	var result []UserProfile
	for i := 0; i < *n; i++ {
		profile := UserProfile{
			UserID:    generateUUID(),
			UserName:  "山田太郎",
			Sex:       1,
			Age:       20,
			JobTitle:  "hoge",
			Company:   "株式会社サンプル",
			Likes:     []string{"550e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440002"},
			Dislikes:  []string{"550e8400-e29b-41d4-a716-446655440003", "550e8400-e29b-41d4-a716-446655440004"},
			Blocks:    []string{"550e8400-e29b-41d4-a716-446655440005", "550e8400-e29b-41d4-a716-446655440006"},
			MainImage: "http://sample.com/image01",
			ImagePath: []string{"http://sample.com/image02", "http://sample.com/image03"},
			RegistDay: "hoge",
			LastLogin: "fuga",
		}

		result = append(result, profile)
	}

	// JSONにエンコードします。
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// エンコードしたJSONを出力します。
	fmt.Println(string(jsonData))
}

func generateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return "550e8400-e29b-41d4-a716-446655440000"
	}
	return u.String()
}
