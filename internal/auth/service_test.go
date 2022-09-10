package auth

import (
	"messenger/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	db := test.NewDB(t)
	test.ResetTables(t, &db)
	repository := repository{&db}
	userService := NewService(repository)

	userId, err := userService.Register("denis", "pass", "denis@email.com")
	assert.Nil(t, err)
	assert.NotEqual(t, 0, userId)
	id, err := userService.Register("denis", "pass", "denis@email.com")
	assert.NotNil(t, err)
	assert.Equal(t, 0, id)

	sessionKey, err := userService.Login("denis", "pass")
	assert.Nil(t, err)
	userIdFromSession, err := userService.GetUserId(sessionKey)
	assert.Nil(t, err)
	assert.Equal(t, userId, userIdFromSession)

}
