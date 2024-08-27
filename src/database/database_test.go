package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, _ := hashPassword(password)
	assert.True(t, checkPasswordHash(password, hashedPassword))
	assert.False(t, checkPasswordHash("invalid", hashedPassword))
}
