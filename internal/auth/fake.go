package auth

import (
	"messenger/internal/entity"
	"messenger/pkg/db"
	"testing"

	"github.com/brianvoe/gofakeit"
	"golang.org/x/crypto/bcrypt"
)

type Faker struct {
	repository Repository
}

func NewFaker(DB *db.DB) Faker {
	return Faker{
		repository{DB},
	}
}
func (f *Faker) User(id int) entity.User {
	hashed, _ := bcrypt.GenerateFromPassword([]byte("pass"), 5)
	user := entity.User{
		Username:       gofakeit.Username(),
		HashedPassword: string(hashed),
		Email:          gofakeit.Email(),
	}
	return user
}

func (f *Faker) InsertUsers(t *testing.T, n int) []int {
	if n == 0 {
		n = 5
	}
	ids := make([]int,0)
	for i := 0; i < n; i++ {
		usr := f.User(n + 1)
		id, err := f.repository.Create(usr)
		if err != nil {
			t.Error("auth.Faker.InsertUsers: ",err)
			t.FailNow()
		}
		ids = append(ids,id) 
	}
	return ids
}
