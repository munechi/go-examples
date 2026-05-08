# 設定管理

データベース接続用のユーザ名やパスワード、ホスト名など一か所で読み込んで管理する。

## .env ファイル

実行ディレクトリに `.env` ファイルを置く。

```text title:.env
DB_HOST=＜ホスト名＞
DB_PORT=＜ポート番号＞
DB_USER=＜ユーザ名＞
DB_PASSWORD=＜パスワード＞
DB_NAME=＜DB名＞
```

## 設定ファイル

```go title:config/config.go
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() DBConfig {

	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found")
	}

	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

}

func Load() Config {
	return Config{
		DB: LoadConfig(),
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
	)
}
```

## 読み込み

データベース接続の例：

```go title:main.go
package main

import (
	"context"
	"fmt"
	"go-config/config"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
    // 設定読み込み
	cfg := config.Load()
    dsn := cfg.DB.DSN()

    // データベース接続
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	fmt.Println("connected!")
}
```

実行：

```bash ln=false
go run main.go
```

出力：

```text ln=false
connected!
```

### 解説

`go.mod` のモジュール名は `go-config` にしています。

そのため import で サブディレクトリの `config` を呼ぶときは、

`import "go-config/config"`

になります。

もしも、`go.mod` のモジュール名がデフォルトの `main` のままだと、

`import "main/config"`

になります。
