package message

import (
	"messenger/internal/entity"

	"github.com/brianvoe/gofakeit"
)

type Faker struct{}

func (f *Faker) Message() entity.ChatMessage {
	var msg entity.ChatMessage
	gofakeit.Struct(&msg)
	msg.Id = 0
	msg.TextContent = gofakeit.Sentence(15)
	msg.SenderId = 1
	msg.RecipientId = 2
	return msg
}
