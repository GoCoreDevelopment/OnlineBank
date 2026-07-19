package jwtservice

import (
	"testing"

	"api/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newService() JWTService {
	return NewJWTService(&config.Config{SecretKeyJWT: []byte("test-secret-key")})
}

func TestCreateAndCheckJWT(t *testing.T) {
	svc := newService()

	token, err := svc.CreateJWT(42)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	assert.NoError(t, svc.CheckJWT(token))
}

func TestGetIdFromJWT(t *testing.T) {
	svc := newService()

	token, err := svc.CreateJWT(1001)
	require.NoError(t, err)

	id, err := svc.GetIdFromJWT(token)
	require.NoError(t, err)
	assert.Equal(t, 1001, id)
}

func TestCheckJWT_Invalid(t *testing.T) {
	svc := newService()
	assert.Error(t, svc.CheckJWT("garbage.token.value"))
}

func TestCheckJWT_WrongSecret(t *testing.T) {
	// Токен, подписанный другим секретом, не должен проходить проверку.
	signer := NewJWTService(&config.Config{SecretKeyJWT: []byte("secret-A")})
	verifier := NewJWTService(&config.Config{SecretKeyJWT: []byte("secret-B")})

	token, err := signer.CreateJWT(7)
	require.NoError(t, err)

	assert.Error(t, verifier.CheckJWT(token))
}
