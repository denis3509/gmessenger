package message

import (
	"github.com/stretchr/testify/assert"
	"messenger/internal/test"
	"testing"
)

func TestRepository(t *testing.T) {

	db := test.DB()
	test.ResetTables(t, db, "message")
	defer db.Close()
	repo := NewRepository(db)
	faker := Faker{}

	// insert
	msg1 := faker.Message()
	msg1.SenderId = 1
	msg1.RecipientId = 2
	id, err := repo.Create(msg1)
	assert.Nil(t, err)
	t.Logf("New id is %d", id)

	msg2 := faker.Message()
	msg2.SenderId = 2
	msg2.RecipientId = 1
	id, err = repo.Create(msg2)
	assert.Nil(t, err)
	t.Logf("New id is %d", id)

	// get
	msg, err := repo.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, msg.Id, id)

	// update
	updatedContent := "[UPDATED_CONTENT]"
	msg.TextContent = updatedContent
	err = repo.Update(msg)
	assert.Nil(t, err)
	msgUpdated, _ := repo.Get(msg.Id)
	assert.Equal(t, updatedContent, msgUpdated.TextContent)

	// contacts
	contacts, err := repo.Contacts(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(contacts))
	t.Logf("Contact list length: %d", len(contacts))

	// delete
	err = repo.Delete(id)
	assert.Nil(t, err)

}
