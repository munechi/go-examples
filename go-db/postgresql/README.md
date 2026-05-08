# データベース操作

PostgreSQL を使用した場合について説明する。

## 前準備

- データベース作成
- データベース接続用のユーザ作成

実行ディレクトリに `.env` ファイルを作成してデータベース接続情報を記述する。

```text title:.env
DB_HOST=＜ホスト名＞
DB_PORT=＜ポート番号＞
DB_USER=＜ユーザ名＞
DB_PASSWORD=＜パスワード＞
DB_NAME=＜DB名＞
```

## インストール

PostgreSQL 用ライブラリをインストールする。

```bash ln=false
go get github.com/jackc/pgx/v5
```

## データベース接続

サンプルコード：

`00_connect/main.go`

実行：

```bash ln=false
go run postgresql/00_connect/main.go
```

出力：

```text ln=false
connected!
```

## CRUD

### (1) テーブル作成

サンプルコード：

`01_create_table/main.go`

実行：

```bash ln=false
go run postgresql/01_create_table/main.go
```

出力：

なし

#### 解説

実務では `pgxpool` がよく使われる。

CREATE TABLE 文を `db.Exec()` に渡す。

### (2) データ挿入

サンプルコード：

`02_insert_data/main.go`

実行：

```bash ln=false
go run postgresql/01_create_table/main.go
```

出力：

なし


#### 解説

INSERT 文の VALUES の `$1`, `$2` に該当するデータを `db.Exec()` の3番目以降の引数に渡す。

### (3) データ読み込み

サンプルコード：

`03_select_data/main.go`

実行：

```bash ln=false
go run postgresql/03_select_data/main.go
```

出力：

```text ln=false
1 田中 20
```

#### 解説

SELECT 文で条件を指定してデータを取得し、`for` 文で変数に格納して印字している。

### (4) データ更新

サンプルコード：

`04_update_data/main.go`

実行：

```bash ln=false
go run postgresql/04_update_data/main.go
```

出力例：

```text ln=false
2026/05/08 09:36:55 updated rows: 1
```

#### 解説

UPDATE 文を使用して WHERE 句で `id` を指定し SET 句でデータを代入する。

`tag.RowsAffected()` で実行結果の件数を取得できる。

### (5) データ削除

サンプルコード：

`05_delete_data/main.go`

実行：

```bash ln=false
go run postgresql/05_delete_data/main.go
```

出力例：

```text ln=false
2026/05/08 09:38:31 deleted rows: 1
```

#### 解説

DELETE 文を使って WHERE 句で削除したい `id` を指定する。

WHERE 句を省略すると全て削除される。

### (6) テーブル削除

サンプルコード：

`06_drop_table/main.go`

実行：

```bash ln=false
go run postgresql/06_drop_table/main.go
```

出力：

なし

#### 解説

DROP TABLE 文を使ってテーブルを削除する。
