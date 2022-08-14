package message

import (
	"messenger/internal/entity"
	"time"
)

type Emitter interface {
	Emit(event string, data []byte)
}

type Service struct {
	name       string
	repository Repository
}

func (s *Service) CreateMessage(senderId int, recipientId int, text string) (*entity.ChatMessage, error) {
	msg := entity.ChatMessage{
		TextContent:        text,
		SenderId:           senderId,
		RecipientId:        recipientId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		IsRead:             false,
		Edited:             false,
		DeletedBySender:    false,
		DeletedByRecipient: false,
	}
	id, err := s.repository.Create(msg)
	if err != nil {
		return nil, err
	}
	msg, err = s.repository.Get(id)
	return &msg, err
}

func (s *Service) GetContactList(userId int) error {
	return nil
}

func (s *Service) MessageList(userId int, contactId int, limit int, page int) []entity.ChatMessage {
	return []entity.ChatMessage{}
}
