# Hello world!

フレームワーク `Echo` を使用する場合。

## サーバー起動

```bash ln=false
go run hello-world-echo/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello
```

## 出力

```text ln=false
{"message":"Hello world!"}
```
