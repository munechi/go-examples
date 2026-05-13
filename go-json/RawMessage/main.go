package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

type Hoge struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// type: "foo"
type Foo struct {
	Name string `json:"name"`
	Age  *int   `json:"age"`
}

// type: "bar"
type Bar struct {
	Address string `json:"address"`
	Email   string `json:"email"`
}

//go:embed data.json
var data []byte

// Note: サンプルコードなので "embed" を使用しました。

func main() {
	// メモリ確保
	var hoge []Hoge

	// 変換
	err := json.Unmarshal(data, &hoge)
	if err != nil {
		log.Fatal(err)
	}

	for _, h := range hoge {
		switch h.Type {
		case "foo":
			var foo Foo
			err := json.Unmarshal(h.Data, &foo)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(foo.Name, *foo.Age)
		case "bar":
			var bar Bar
			err := json.Unmarshal(h.Data, &bar)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(bar.Address, bar.Email)
		default:
			fmt.Println("Unknown type:", h.Type)
		}
	}
}
