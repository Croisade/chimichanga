package services

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtAuthService(t *testing.T) {
	jwtService := NewJWTAuthService()
	t.Run("Create Token", func(t *testing.T) {
		got, err := jwtService.CreateToken()

		assert.Nil(t, err)
		assert.Equal(t, len(strings.Split(got, ".")), 3)
	})

	t.Run("Validate Token", func(t *testing.T) {
		token, _ := jwtService.CreateToken()
		got, err := jwtService.ValidateToken(token)

		assert.Nil(t, err)
		assert.Equal(t, len(strings.Split(got.Raw, ".")), 3)
	})

	t.Run("Create Refresh Token", func(t *testing.T) {
		got, err := jwtService.CreateRefreshToken()

		assert.Nil(t, err)
		assert.Equal(t, len(strings.Split(got, ".")), 3)
	})
}
