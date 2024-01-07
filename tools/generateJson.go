package main

import (
	"encoding/json"
	"fmt"
)

// UserProfile はユーザーのプロファイル情報を表す構造体です。
type UserProfile struct {
	UserID       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	Title        string   `json:"title"`
	Company      string   `json:"company"`
	LikeUsers    []string `json:"like_users"`
	BlockUsers   []string `json:"block_users"`
	Career       []Career `json:"career"`
	ChatIDList   []string `json:"chat_id_list"`
	Description  string   `json:"description"`
	ImageURLList []string `json:"image_url_list"`
	Confirmation bool     `json:"confirmation"`
	LastLogin    string   `json:"last_login"`
	CreateDate   string   `json:"create_date"`
	Billing      []int    `json:"billing"`
}

// Career はユーザーのキャリア情報を表す構造体です。
type Career struct {
	Position string `json:"役職"`
	Years    string `json:"年数"`
	// その他のキャリア関連情報
}

func main() {
	// UserProfile構造体のインスタンスを作成し、データを埋めます。
	profile := UserProfile{
		UserID:       "user123",
		UserName:     "山田太郎",
		Title:        "エンジニア",
		Company:      "株式会社サンプル",
		LikeUsers:    []string{"user456", "user789"},
		BlockUsers:   []string{"user111", "user222"},
		Career:       []Career{{Position: "マネージャー", Years: "5年"}},
		ChatIDList:   []string{"chat123", "chat456"},
		Description:  "こんにちは、山田太郎です。",
		ImageURLList: []string{"http://example.com/image1.jpg", "http://example.com/image2.jpg"},
		Confirmation: true,
		LastLogin:    "2021-01-01T12:00:00Z",
		CreateDate:   "2020-01-01T00:00:00Z",
		Billing:      []int{1, 2, 3},
	}

	// JSONにエンコードします。
	jsonData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// エンコードしたJSONを出力します。
	fmt.Println(string(jsonData))
}
