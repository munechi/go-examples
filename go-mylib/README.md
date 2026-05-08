# 自作パッケージ

## 想定

プログラム開発中にローカルにある自作パッケージを読み込む

## ディレクトリ

```text ln=false
├─ myapp/  ← 開発中のディレクトリ
└─ myutils/  ← 自作パッケージ
```

それぞれのディレクトリに `go.mod` が存在する。

## 読み込み

`myutils/go.mod` に記述されているモジュール名は `github.com/munechi/go-myutils`

`myapp/go.mod` に `replace` 分でローカルのディレクトリを指定する。

```text ln=false
replace github.com/munechi/go-myutils => ../myutils
```

`myapp` のプログラム内で `go-myutils` の関数を使うときは import 文で

```go ln=false
import (
    "github.com/munechi/go-myutils"
)
```

関数を呼ぶときは `パッケージ名.関数名` と書く。

```go ln=false
myutils.Add()
```
