# クエリパラメータ取得

フレームワークを使用しない場合。

## サーバー起動

```bash ln=false
go run get/query/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello?name=Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
