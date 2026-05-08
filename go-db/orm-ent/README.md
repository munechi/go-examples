# Ent

`Ent` は、

> **DBのテーブルを Go の型・コードとして生成する**

思想である （かなりぶっ飛んだ発想）

### インストール

```bash ln=false
go get entgo.io/ent/cmd/ent
go get entgo.io/ent
go get github.com/lib/pq
```

### Schema を書く

例：

`ent/schema/` ディレクトリを作成し、その中に `user.go` ファイルを置く。

```go title:go-db/psql-ent/ent/schema/user.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),

		field.Int("age").
			Positive(),
	}
}
```

SQL でいうと↓

```sql ln=false
CREATE TABLE users (
    id bigserial primary key,
    name text not null,
    age integer not null
);
```

### コード生成

```bash ln=false
go run entgo.io/ent/cmd/ent generate ./ent/schema
```

成功すると `ent/` の中にいくつかディレクトリとファイルが自動生成される。

### テーブル作成

```go title:go-db/orm-ent/create_table/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// migration
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("connected")
}

```

実行：

```bash ln=false
go run orm-ent/create_table/main.go
```

出力例：

```text ln=false
2026/05/08 08:33:49 connected
```

#### 解説

import 文に生成したスキーマ (`go-db/orm-ent/ent`) を指定する。

【注意点】

`ent` はテーブル名を自動的に複数形の名前で作成する（`user` なら `users`。`person` なら `people`）

もちろんテーブル名を固定することも可能。

### INSERT

```go title:go-db/orm-ent/insert/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	u, err := client.User.
		Create().
		SetName("田中").
		SetAge(30).
		Save(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(u.ID)
}
```

実行：

```bash ln=false
go run orm-ent/insert/main.go
```

出力例：

```text ln=false
2026/05/08 08:35:33 1
```

#### 解説

出力例は `log.Println(u.ID)` による出力。ID を表示している。

SQL イメージ：

```sql ln=false
INSERT INTO users (name, age) VALUES ('田中', 30) RETURNING id;
```

### SELECT (1件)

```go title:go-db/orm-ent/select_only/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"
	"go-db/orm-ent/ent/user"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	u, err := client.User.
		Query().
		Where(user.Name("田中")).
		Only(ctx)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(u.ID)
}
```

実行：

```bash ln=false
go run orm-ent/select_only/main.go
```

出力例：

```text ln=false
2026/05/08 08:36:38 1
```

#### 解説

`Where()` でユーザ名が「田中」に一致するものを1件取得し、`log.Println(u.ID)` で ID を出力している。

SQL イメージ：

```sql
SELECT * FROM users WHERE name='田中' LIMIT 2;
```

なせ `LIMIT 1 ではないのか？

- **0件** → NotFound エラー
- **1件** → それを返す
- **2件取れた** → MultipleRows エラー

### SELECT (全件)

```go title:go-db/orm-ent/select_all/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"
	"go-db/orm-ent/ent/user"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	users, err := client.User.
		Query().
		Where(user.Name("田中")).
		All(ctx)

	if err != nil {
		log.Fatal(err)
	}

	for _, u := range users {
		log.Println(u.Name, u.Age)
	}
}
```

実行：

```bash ln=false
go run orm-ent/select_all/main.go
```

出力例：

```text ln=false
2026/05/08 08:37:40 田中 30
```

### UPDATE

```go title:go-db/orm-ent/update/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	err = client.User.
		UpdateOneID(1).
		SetAge(31).
		Exec(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
```

実行：

```bash ln=false
go run orm-ent/update/main.go
```

前出の SELECT ALL を再び実行したときの出力例：

```text ln=false
2026/05/08 08:50:44 田中 31
```

`age` が 30 から 31 に更新されている。


### DELETE

```go title:go-db/orm-ent/delete/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-db/orm-ent/ent"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	err = client.User.
		DeleteOneID(1).
		Exec(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
```

実行：

```bash ln=false
go run orm-ent/delete/main.go
```

前出の SELECT ALL を再び実行しても何も出力されない。
