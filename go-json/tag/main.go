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
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

type User2 struct {
	Name string `json:"name"`
	Age  int    `json:"-"`
}

func main() {
	// キー名のみの例
	user0 := User0{Name: "Taro"}

	data0, err := json.Marshal(user0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data0)) // output: {"name":"Taro","age":0}

	// omitempty の例
	user1 := User1{Name: "Taro"}

	data1, err := json.Marshal(user1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data1)) // output: {"name":"Taro"}

	// 無視の例
	user2 := User2{Name: "Taro", Age: 20}

	data2, err := json.Marshal(user2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data2)) // output: {"name":"Taro"}
}
