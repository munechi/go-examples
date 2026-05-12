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
		port, _ := cmd.Flags().GetInt("port")    // --port の値を取得

		fmt.Println("hello", name)

		if debug {
			fmt.Println("debug flag is", debug)
		}

		fmt.Println("port number:", port)
	},
}

func init() {
	helloCmd.Flags().StringP(
		"name",          // オプション名 (--name)
		"n",             // 短縮形
		"world",         // 初期値
		"name to greet", // 説明
	)
	helloCmd.Flags().BoolP(
		"debug",      // オプション名 (--debug)
		"d",          // 短縮形
		false,        // 初期値
		"debug flag", // 説明
	)
	helloCmd.Flags().IntP(
		"port",        // オプション名 (--port)
		"p",           // 短縮形
		8080,          // 初期値
		"Port number", // 説明
	)
	// サブコマンド登録
	rootCmd.AddCommand(helloCmd)
}
