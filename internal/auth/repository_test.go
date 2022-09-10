package auth

import (
	"messenger/internal/entity"
	"messenger/internal/test"
	"testing"
 
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db := test.NewDB(t)
	test.ResetTables(t, &db)
	repo := repository{&db}
	_, err := repo.GetByName("notexist")
	assert.Equal(t, "sql: no rows in result set", err.Error())
	
	// Create
	user := entity.User{
		Username:       "denis",
		HashedPassword: "denispass",
		Email:          "denis@email.com",
	}
	id, err  := repo.Create(user)
	assert.Nil(t, err)

	// Update 
	user.Id = id
	newEmail := "new@email.com"
	user.Email = newEmail
	err = repo.Update(user,"Email")
	assert.Nil(t, err)
	user, err = repo.Get(user.Id)
	assert.Nil(t,err)
	assert.Equal(t,newEmail, user.Email)

}
