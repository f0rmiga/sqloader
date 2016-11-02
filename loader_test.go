package sqloader

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoader(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	sqlfile := path.Dir(filename) + string(filepath.Separator) + "queries_test.sql"

	sqloader, err := NewSQLoader(sqlfile)
	if err != nil {
		panic(err)
	}

	// TODO better tests
	fmt.Println(sqloader.Get("nonExistentQuery"))
	fmt.Println(sqloader.Get("selectUser"))
	fmt.Println(sqloader.Get("listUser"))
}
