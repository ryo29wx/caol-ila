package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

/*
 * {users:
  [{
  user_id: <string>,
  nick_name: <string>,
  sex: <number>,
  title: <string>,
  company: <string>,
  like: <bool>,
  image_url: <string>
  }]
total: <number>
}
*/

type User struct {
	UserID    string `json:"user_id"`
	NickName  string `json:"nick_name"`
	Sex       int    `json:"sex"`
	Title     string `json:"title"`
	Company   string `json:"company"`
	Like      bool   `json:"like"`
	MainImage string `json:"image_url"`
}

type SearchResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

var mockUsers = []User{
	{
		UserID:    createRandumUUID(),
		NickName:  "sample Ryo Kiuchi",
		Sex:       1,
		Title:     "Engineering Manager",
		Company:   "Splunk",
		Like:      true,
		MainImage: "/images/TestUser1.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "sample Yoshiaki Toyama",
		Sex:       1,
		Title:     "Software Engineer",
		Company:   "Sample Company",
		Like:      true,
		MainImage: "/images/TestUser2.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "sample 四谷　太郎",
		Sex:       2,
		Title:     "部長",
		Company:   "日本電気株式会社",
		Like:      false,
		MainImage: "/images/TestUser3.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "四谷　太郎　卍",
		Sex:       1,
		Title:     "Head of Infrastructure",
		Company:   "LINEヤフー",
		Like:      false,
		MainImage: "/images/TestUser4.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "四谷　太郎　sample",
		Sex:       1,
		Title:     "Head of Infrastructure",
		Company:   "LINEヤフー",
		Like:      false,
		MainImage: "/images/TestUser5.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "sample 吉田 沙織",
		Sex:       1,
		Title:     "Head of Infrastructure",
		Company:   "Sample Campany",
		Like:      false,
		MainImage: "/images/TestUser6.jpeg",
	},
	{
		UserID:    createRandumUUID(),
		NickName:  "Sample User1",
		Sex:       1,
		Title:     "Head of Infrastructure",
		Company:   "Sample Campany",
		Like:      false,
		MainImage: "/images/TestUser7.jpeg",
	},
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	pageStr := r.URL.Query().Get("p")
	fmt.Println("query:" + query + "/pageStr:" + pageStr)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	var responseUsers []User
	for _, user := range mockUsers {
		if strings.Contains(strings.ToLower(user.NickName), strings.ToLower(query)) {
			responseUsers = append(responseUsers, user)
		}
	}

	resp := SearchResponse{
		Users: responseUsers,
		Total: len(responseUsers),
	}

	fmt.Println(resp)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func createRandumUUID() string {
	newUUID := uuid.New()
	uuidString := newUUID.String()
	return uuidString
}

func main() {
	http.HandleFunc("/v1/search", enableCORS(searchHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
