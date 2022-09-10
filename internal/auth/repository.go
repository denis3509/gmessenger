package auth

import (
	"messenger/internal/entity"
	"messenger/pkg/db"

	dbx "github.com/go-ozzo/ozzo-dbx"
)


type Repository interface {
	// Get returns the user with the specified user ID.
	Get(id int) (entity.User, error)
	// GetByNameOrEmail search by email and username
	GetByName(username string) ( entity.User, error)
	// Create saves a new user in the storage.
	Create(user entity.User) (int, error)
	// Delete removes the user with given ID from the storage.
	Delete(id int) error
	// Update 
	Update(user entity.User, attrs ...string) error
}

type repository struct {
	db *db.DB
}


func NewRepository(db.DB) Repository {
	return repository {}
}

// Get returns the user with the specified user ID.
func (r repository) Get(id int) (entity.User, error) {
	var user entity.User
	err := r.db.Select().Model(id, &user)
	return user, err
}

// Create saves a new user in the storage.
func (r repository) Create(user entity.User) (int, error) {
	err := r.db.Model(&user).Insert()
	return user.Id, err
}
func (r repository) Update(user entity.User, attrs ...string) error {
	return r.db.Model(&user).Update(attrs...)
}
// Delete removes the user with given ID from the storage.
func (r repository) Delete(id int) error {
	user, err := r.Get(id)
	if err != nil {
		return err
	}
	return r.db.Model(&user).Delete()
}

func (r repository) GetByName(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Select().
		From("user").
		Where(dbx.HashExp{"username": username}).
		Row(&user.Id, &user.Username,&user.HashedPassword,&user.Email)
	return user, err

}
