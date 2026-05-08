# Testify

Python の pytest っぽいフレームワーク

## インストール

```bash ln=false
go get github.com/stretchr/testify
```

## テスト対象

```go:utils.go
package myutils

func Add(a, b int) int {
	return a + b
}
```

## テストコード

```go:utils_test.go
package myutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{1, -1, 0},
		{10, -4, 6},
		{-4, 10, 6},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, Add(tt.a, tt.b))
	}
}
```

## 実行

```bash ln=false
run test .
```

## 出力

```text ln=false
ok      github.com/munechi/myutils      (cached)
```
