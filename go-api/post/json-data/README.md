# POST メソッド (JSON データ)

フレームワークを使用しない場合。

## サーバー起動

```bash ln=false
go run post/json-data/main.go
```

## テスト

```bash ln=false
curl -X POST http://localhost:8080/hello -H "Content-Type: application/json" -d '{"name":"Taro"}'
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
