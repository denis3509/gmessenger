package message

import (
	"messenger/internal/entity"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Repository encapsulates the logic to access messages from the data source.
type Repository interface {
	// Get returns the message with the specified message ID.
	Get(id int) (entity.ChatMessage, error)
	// Count returns the number of messages.
	Count() (int, error)
	// Query returns the list of messages with the given offset and limit.
	Query(offset, limit int) ([]entity.ChatMessage, error)
	// Create saves a new message in the storage.
	Create(message entity.ChatMessage) (int, error)
	// Update updates the message with given ID in the storage.
	Update(message entity.ChatMessage, fields... string) error
	// Delete removes the message with given ID from the storage.
	Delete(id int) error
	// Contacts return contact list for user
	Contacts(userId int) ([]entity.Contact, error)
}

// repository persists messages in database
type repository struct {
	db *dbx.DB
	// logger log.Logger

}

// NewRepository creates a new message repository
func NewRepository(db *dbx.DB) Repository {
	return repository{
		db,
		//  logger,
	}
}

// Get reads the message with the specified ID from the database.
func (r repository) Get(id int) (entity.ChatMessage, error) {
	var message entity.ChatMessage
	err := r.db.Select().Model(id, &message)
	return message, err
}

// Create saves a new message record in the database.
// It returns the ID of the newly inserted message record.
func (r repository) Create(message entity.ChatMessage) (int, error) {
	err := r.db.Model(&message).Exclude("created_at", "updated_at").Insert()
	return message.Id, err
}

// Update saves the changes to an message in the database.
func (r repository) Update(message entity.ChatMessage, attrs... string) error {
  return r.db.Model(&message).Update(attrs...)
  
}

// Delete deletes an message with the specified ID from the database.
func (r repository) Delete(id int) error {
	message, err := r.Get(id)
	if err != nil {
		return err
	}
	return r.db.Model(&message).Delete()
}

// Count returns the number of the message records in the database.
func (r repository) Count() (int, error) {
	var count int
	err := r.db.Select("COUNT(*)").From("message").Row(&count)
	return count, err
}

// Query retrieves the message records with the specified offset and limit from the database.
func (r repository) Query(offset, limit int) ([]entity.ChatMessage, error) {
	var messages []entity.ChatMessage
	err := r.db.
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&messages)
	return messages, err
}


// Contacts return contact list
func (r repository) Contacts(userId int  ) ([]entity.Contact, error) {
	q := r.db.NewQuery(contactListSQL)
	q.Bind(dbx.Params{"user": userId}) 
	var contacts []entity.Contact
	err := q.All(&contacts)
	return contacts, err
}

