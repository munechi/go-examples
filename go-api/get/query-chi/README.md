# クエリパラメータ取得

軽量フレームワーク `Chi` を使用する場合。

`Chi` はルーター回りの面倒を見てくれるだけなのでハンドラー関数の中身は `net/http` と同じ。

## サーバー起動

```bash ln=false
go run get/query-chi/main.go
```

## テスト

```bash ln=false
curl http://localhost:8080/hello?name=Taro
```

## 出力

```text ln=false
{"message":"Hello Taro!"}
```
