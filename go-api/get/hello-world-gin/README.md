# Hello world!

フレームワーク `Gin` を使用する場合。

### 注目ポイント

`Gin` ではマップデータを `gin.H{}` で生成できる。

## サーバー起動

```bash ln=false
go run hello-world-gin/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello
```

## 出力

```text ln=false
{"message":"Hello world!"}
```
