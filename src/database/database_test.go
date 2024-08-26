package database

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, _ := hashPassword(password)
	fmt.Println(hashedPassword)
}
