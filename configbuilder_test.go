package configbuilder

import (
	"fmt"
	"testing"
)

type T struct {
	Bool  bool    `csv:"bool"`
	Int   int     `csv:"int"`
	Str   string  `csv:"str"`
	Uint  uint    `csv:"uint;attr=key"`
	Float float32 `csv:"float"`
}

func TestConfigBuilder(t *testing.T) {
	builder := NewConfigBuilder()
	config := builder.MakeConfig(FileType_Csv, "test.csv", (*T)(nil))
	row := config[uint(2)].(*T)

	fmt.Println("struct:", *row)
}
