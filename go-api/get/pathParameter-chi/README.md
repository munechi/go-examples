# パスパラメータ取得

軽量フレームワーク `Chi` を使用する場合。

## サーバー起動

```bash ln=false
go run get/pathParameter-chi/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello/Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
