# POST メソッド (マルチパートデータ)

フレームワーク `Gin` を使用する場合。

## サーバー起動

```bash ln=false
go run post/multipart-data-gin/main.go
```

## テスト

```bash ln=false
curl -X POST -F "name=Taro" -F "file=@logo.png" http://localhost:8080/upload
```

## 出力

```text ln=false
{"message":"upload successful","name":"Taro"}
```
