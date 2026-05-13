package main

import (
	"encoding/json"
	"fmt"
)

// JSON に対応する構造体
type User struct {
	Name string
	Age  int
}

func main() {
	// 変換データの定義
	user := User{
		Name: "Taro",
		Age:  20,
	}

	// 変換
	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	// 変換結果出力
	fmt.Println(string(data)) // Marshal の戻り値は []byte

	// マップ
	map_data := map[string]string{
		"Name": "Jiro",
		"Age":  "19",
	}
	data_map, err := json.Marshal(map_data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data_map))

	// スライス
	slice_data := []User{
		{Name: "Saburo", Age: 18},
		{Name: "Shiro", Age: 17},
	}
	data_slice, err := json.MarshalIndent(slice_data, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data_slice))
}
