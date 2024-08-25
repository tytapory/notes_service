package server

import (
	"fmt"
	"testing"
)

func TestVerify(t *testing.T) {
	str, err := verify("ghdtn      пивет ult ты\nзнаешь я тебя люблю\nvjz vfv тестируем текст д да да урасЖцс")
	fmt.Printf("%s", str)
	if err != nil {
		fmt.Println(err.Error())
	}
}
