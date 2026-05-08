# GORM

Go 言語の ORM ライブラリ

## インストール

```bash ln=false
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

## 実行

```bash ln=false
go run orm-gorm/main.go
```

## 出力例

```text
[CREATE Table] name: users
[INSERT] id = 1
[INSERT] id = 2
[SELECT]
id=1 name=田中 age=20
id=2 name=鈴木 age=23
[UPDATE] id = 1 age = 21
[SELECT]
id=1 name=田中 age=21
id=2 name=鈴木 age=23
[DELETE] id = 2
[SELECT]
id=1 name=田中 age=21
[Unscoped SELECT]
id=1 name=田中 age=21 deleted_at=0001-01-01 00:00:00
id=2 name=鈴木 age=23 deleted_at=2026-05-07 16:35:23
[Unscoped DELETE] id = 2
[Unscoped SELECT]
id=1 name=田中 age=21 deleted_at=0001-01-01 00:00:00
```

## GORM の特徴

### メソッドチェーン

GORM は

```go ln=false
db.
	Where("age >= ?", 20).
	Order("age desc").
	Limit(10).
	Find(&users).
	Error
```

こうつなげる思想。
最後に `.Error`
これが作法。

### Find() の戻り値

```go ln=false
*gorm.DB
```

その中に

- Error
- RowsAffected
- Statement

などが入っている。

### 論理削除 (soft delete)

GORM でテーブル定義に `gorm.Model` を埋め込んでいると通常のデータ削除が soft delete になる。

サンプルのテーブル定義はこうなっている。

```go
type User struct {
	gorm.Model
	Name string
	Age  int
}
```

実は `gorm.Model` の中身は

```go
type Model struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
```

こうなっている。

soft delete では `DeletedAt` に日付情報などが入っていて、

```go
db.Find(&users)
```

が、内部的には、

```sql
SELECT * FROM users WHERE deleted_at IS NULL;
```

というふうに処理され、自動的に「削除済み」を隠す。

なぜそういう仕様なのか？
- 復活できる
- 履歴が残る
- 外部キーが壊れにく


soft delete したデータも SELECT したい場合は

```go
db.Unscoped().Find()
```

soft delete せず普通に削除した場合は

```go
db.Unscoped().Delete()
```

というふうに `Unscoped()` を挟む。
