# パスパラメータ取得

フレームワーク `Echo` を使用する場合。

`Gin` では `gin.H{}` でマップ形式に変換できるけど、`echo` では `map[string]string{}` という標準のマップ記述。

## サーバー起動

```bash ln=false
go run get/pathParameter-echo/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello/Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
