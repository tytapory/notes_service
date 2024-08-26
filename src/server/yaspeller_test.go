package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate_NormalInput(t *testing.T) {
	str, err := validateAndFixErrors("ghdtn      пивет ult ты\nзнаешь я тебя люблю\nvjz vfv тестируем текст д да да урасЖцс")
	expected := "привет      привет где ты\nзнаешь я тебя люблю\nмоя мам тестируем текст д да да урасЖцс"
	assert.Nil(t, err)
	assert.Equal(t, str, expected)
}

func TestValidate_EmptyInput(t *testing.T) {
	str, err := validateAndFixErrors("")
	expected := ""
	assert.Nil(t, err)
	assert.Equal(t, str, expected)
}

func TestValidate_InvalidCharacters(t *testing.T) {
	str, err := validateAndFixErrors("1234567890 !@#$%^&*()")
	expected := "1234567890 !@#$%^&*()"
	assert.Nil(t, err)
	assert.Equal(t, str, expected)
}
