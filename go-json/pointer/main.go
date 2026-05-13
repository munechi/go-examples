package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type User0 struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type User1 struct {
	Name string `json:"name"`
	Age  *int   `json:"age"`
}

func main() {
	// 変換する JSON
	data := []byte(`{"name":"Taro","age":0}`)

	// pointer なし
	var user0 User0

	err := json.Unmarshal(data, &user0)
	if err != nil {
		log.Fatal(err)
	}

	if user0.Age == 0 {
		fmt.Println("user.Age is not defined") // age=0 がエラーになる
	} else {
		fmt.Println(user0.Name, user0.Age)
	}

	// pointer あり (age データあり)
	var user1 User1

	err = json.Unmarshal(data, &user1)
	if err != nil {
		log.Fatal(err)
	}

	if user1.Age == nil { // 変数がポインタなので nil で比較できるようになる
		fmt.Println("user.Age is not defined")
	} else {
		fmt.Println(user1.Name, *user1.Age)
	}

	// pointer あり (age データなし)
	data = []byte(`{"name":"Taro"}`)
	var user2 User1

	err = json.Unmarshal(data, &user2)
	if err != nil {
		log.Fatal(err)
	}

	if user2.Age == nil { // 変数がポインタなので nil で比較できるようになる
		fmt.Println("user.Age is not defined")
	} else {
		fmt.Println(user2.Name, *user2.Age)
	}
}
