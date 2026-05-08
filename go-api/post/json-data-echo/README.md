# POST メソッド (JSON データ)

フレームワーク `Echo` を使用する場合。

## サーバー起動

```bash ln=false
go run post/json-data-echo/main.go
```

## テスト

```bash ln=false
curl -X POST http://localhost:8080/hello -H "Content-Type: application/json" -d '{"name":"Taro"}'
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
