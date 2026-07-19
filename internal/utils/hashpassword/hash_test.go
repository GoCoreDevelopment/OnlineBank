package hashpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHashAndCheck(t *testing.T) {
	password := "S3cretPass!"

	hash, err := CreateHash(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash, "hash must not equal the raw password")

	// Правильный пароль проходит проверку.
	assert.NoError(t, CheckValidHash(hash, password))
}

func TestCheckValidHash_WrongPassword(t *testing.T) {
	hash, err := CreateHash("correct-password")
	require.NoError(t, err)

	err = CheckValidHash(hash, "wrong-password")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid password")
}

func TestCreateHash_Unique(t *testing.T) {
	// bcrypt использует соль — два хэша одного пароля должны различаться.
	h1, _ := CreateHash("same")
	h2, _ := CreateHash("same")
	assert.NotEqual(t, h1, h2)
}
