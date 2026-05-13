package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// 変換用の構造体
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// 変換する JSON
	data := []byte(`{"name":"Taro","age":20}`)

	// メモリ確保
	var user User

	// 変換
	err := json.Unmarshal(data, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(user.Name, user.Age)

	// マップ(JSON の構造が不明なとき)
	var map_user map[string]any
	err = json.Unmarshal([]byte(`{"name":"Jiro","age":19}`), &map_user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(map_user["name"], map_user["age"]) // 数字は自動的に float64 になる。

	// スライス＋マップ
	var slice_user []map[string]any
	err = json.Unmarshal(
		[]byte(`[{"name":"Saburo","age":18},{"name":"Shiro","age":17}]`),
		&slice_user,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(slice_user[0]["name"], slice_user[0]["age"])
	fmt.Println(slice_user[1]["name"], slice_user[1]["age"])
}
