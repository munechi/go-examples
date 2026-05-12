# cli

## 標準ライブラリ

標準ライブラリ `os` の `Args` で与えられた引数をスライスを取得できます。

```go
args := os.Args
```

実行：

```bash ln=false
go run Args/main.go foo bar
```

出力例：

```text ln=false
[main.go foo bar]
```

スライスの 0 番目にはプログラム名。それ以降に与えられた引数が格納されています。

build で実行ファイルを作成しても同様です。

ビルド：

```bash ln=false
go build Args/main.go
```

実行：

```bash ln=false
./main foo bar
```

出力例：

```text ln=false
[main foo bar]
```

## Cobra

オプションやヘルプ表示機能などを実装する場合はフレームワーク **Cobra** を使用すると良いです。

利点：

- help が自動作成
- サブコマンドを容易に定義可能
- オプションやフラグを容易に定義可能
- 補完機能

実際にツールとして運用するなら help や補完機能があった方が嬉しいと思います。

### インストール

```bash ln=false
go install github.com/spf13/cobra-cli@latest
```

### プロジェクト作成

```bash ln=false
mkdir Cobra
cd Cobra
cobra-cli init
```

作成されるディレクトリの構成：

```text ln=false
Cobra/
 ├ cmd/
 │  └ root.go
 ├ LICENSE
 └ main.go
```

### サブコマンド追加

```text ln=false
cobra-ini add hello
```

`cmd/hello.go` が生成されます。

実行：

```bash ln=false
go run . hello
```

出力例：

```text ln=false
hello called
```

### ヘルプ表示

```bash ln=false
go run . --help
```

出力例：

```text ln=false
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  Cobra [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  hello       A brief description of your command
  help        Help about any command

Flags:
  -h, --help     help for Cobra
  -t, --toggle   Help message for toggle

Use "Cobra [command] --help" for more information about a command.
```

このようなヘルプが表示されます。

### オプション追加

追加内容：

サブコマンド `hello` に `--name <string>` と `--debug` オプションを追加します。

ソースコード：

`cmd/hello.go` を変更します。

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name") // --name の値を取得
		debug, _ := cmd.Flags().GetBool("debug") // --debug の有無を取得
		fmt.Println("hello", name)
		if debug {
			fmt.Println("debug flag is", debug)
		}
	},
}

func init() {
	helloCmd.Flags().String(
		"name",            // オプション名 (--name)
		"world",           // 初期値
		"name to greet",   // 説明
	)
	helloCmd.Flags().Bool(
		"debug",           // オプション名 (--debug)
		false,             // 初期値
		"debug flag",      // 説明
	)
	rootCmd.AddCommand(helloCmd)
}
```

実行：

```bash ln=false
go run . hello --name Taro --debug
```

出力例：

```text ln=false
hello Taro
debug flag is true
```

オプションなしで実行：

```bash ln=false
go run . hello
```

出力例：

```text ln=false
hello world
```

デフォルト値 (`world`) が適用されています。

ヘルプ出力：

```bash ln=false
go run . hello -h
```

出力例：

```text ln=false
Usage:
  Cobra hello [flags]

Flags:
      --debug         debug flag
  -h, --help          help for hello
      --name string   name to greet (default "world")
```

ヘルプが表示される。

#### よく使う型まとめ

- `Flags().String()`
- `Flags().StringP()`
- `Flags().Int()`
- `Flags().IntP()`
- `Flags().Bool()`
- `Flags().BoolP()`

### 短縮形オプション

`Flags().String()` を `Flags().StringP()` にすると短縮形オプションも扱えるようになります。

```go
	helloCmd.Flags().StringP(
		"name",           // オプション名
        "n",              // 短縮形のオプション
		"world",          // デフォルト値
		"name to greet",  // 説明
	)
```

ヘルプ出力例：

```text ln=false
Usage:
  Cobra hello [flags]

Flags:
      --debug         debug flag
  -h, --help          help for hello
  -n, --name string   name to greet (default "world")
```

`--name` オプションに `-n` が追加されているのが分かります。

### 補完機能

補完のためのスクリプトをシェルに読み込ませる必要があります。

以下は Bash の例：

```bash ln=false
go build
```

その後、とりあえず一時的にスクリプトを読み込ませます。

```bash ln=false
source <(./Cobra completion bash)
```

補完が効くか試してみましょう。

```bash ln=false
./Cobra he<TAB>
```

【注意点】

ビルド時 `-o` で作成する実行ファイル名を指定する場合、ビルド前に `cmd/root.go` の `Use:` に実行ファイル名を定義しておく必要があります。

下記例は実行ファイル名を `myapp` にする場合。

```go
var rootCmd = &cobra.Command{
	Use:   "myapp",
```

上記例では実行中のシェルに一時的に読み込ませていますが、
恒久的に使う場合、個人の `~/.bashrc` でスクリプトを読み込ませたり、
`/usr/share/bash-completion/completions` にスクリプトを格納したりと、運用方法によって変わると思います。
