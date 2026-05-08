# パスパラメータ取得

フレームワークを使用しない場合。

`strings` パッケージの `trimPrefix()` を使用して `/hello/` から後ろを抜き出す。

## サーバー起動

```bash ln=false
go run get/pathParameter/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello/Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
