package message

import (
	"encoding/json"
	"messenger/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {

	db := test.DB()
	test.ResetTables(t, db, "message")
	defer db.Close()
	repo := NewRepository(db)
	faker := Faker{}

	// Insert
	msg1 := faker.Message()
	msg1.SenderId = 1
	msg1.RecipientId = 2
	msg1.IsRead = false
	id, err := repo.Create(msg1)
	assert.Nil(t, err)
	t.Logf("New id is %d", id)

	msg2 := faker.Message()
	msg2.SenderId = 2
	msg2.RecipientId = 1
	id, err = repo.Create(msg2)
	assert.Nil(t, err)
	t.Logf("New id is %d", id)

	// Get
	msg, err := repo.Get(id)
	assert.Nil(t, err)
	assert.Equal(t, msg.Id, id)

	// Update
	updatedContent := "[UPDATED_CONTENT]"
	msg.TextContent = updatedContent
	err = repo.Update(msg)
	assert.Nil(t, err)
	msgUpdated, _ := repo.Get(msg.Id)
	assert.Equal(t, updatedContent, msgUpdated.TextContent)

	// Contacts
	contacts, err := repo.Contacts(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(contacts))
	t.Logf("Contact list length: %d", len(contacts))
	_, err = json.Marshal(contacts)
	assert.Nil(t, err)

	// Dialog
	dialog, err := repo.Dialog(1, 2, 0, 5)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(dialog))
	t.Logf("Dialog length: %d", len(dialog))
	data, err := json.Marshal(dialog)
	assert.Nil(t, err)
	t.Logf("Dialog : %s", data)

	// Delete
	err = repo.Delete(id)
	assert.Nil(t, err)

	//CountUnread 
	m := faker.Message()
	m.SenderId = 1
	m.RecipientId = 2
	m.IsRead = false
	_,_  = repo.Create(m)
	m = faker.Message()
	m.SenderId = 1
	m.RecipientId = 2
	m.IsRead = true
	_,_  = repo.Create(m)
	unreadCount, err := repo.CountUnread(2)
	assert.Nil(t, err)
	assert.Equal(t, 2, unreadCount)
	t.Logf("Unread count is: %d", unreadCount)

	// set read
	affected, err := repo.SetMessagesRead(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, unreadCount, affected)

	// unreadCount  unread message count
	unreadCount, err = repo.CountUnread(2)
	assert.Nil(t, err)
	assert.Equal(t, 0, unreadCount)
	t.Logf("Unread count now is: %d", unreadCount)

}
