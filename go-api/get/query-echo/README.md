# クエリパラメータ取得

フレームワーク `Echo` を使用する場合。

## サーバー起動

```bash ln=false
go run get/query-echo/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello?name=Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
